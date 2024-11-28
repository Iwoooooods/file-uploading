package localstorage

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Iwoooooods/fs-upload-go/internal/repositories"
	"github.com/Iwoooooods/fs-upload-go/internal/services"
)

const BUFFER_SIZE = 1024 * 1024 * 4

type Uploader interface {
	UploadFile(filePath string, md5 string, metaService services.MetaService) (fileId string, err error)
}

type DefaultUploader struct {
	ServerURL   string
	Username    string
	BasePath    string
	MetaService services.MetaService
}

func NewUploader(serverURL string, username string, db *sql.DB) (*DefaultUploader, error) {
	// check if user exists, if not create new fold for the user
	basePath := serverURL + "/" + username + "/"
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// give permission to the user to read and write to the folder
		os.MkdirAll(basePath, 0755)
	}

	metaRepo := repositories.NewMetaRepositorySQLite(db)
	metaService := services.NewMetaService(metaRepo)

	return &DefaultUploader{
		ServerURL:   serverURL,
		Username:    username,
		BasePath:    basePath,
		MetaService: metaService,
	}, nil
}

func (u *DefaultUploader) UploadFile(ctx context.Context, src io.Reader, fileName string) error {

	buffer := make([]byte, BUFFER_SIZE)
	filePath := filepath.Join(u.BasePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return err
	}
	defer file.Close()

	_, err = io.CopyBuffer(file, src, buffer)
	if err != nil {
		log.Printf("failed to copy file: %v", err)
		return err
	}

	return nil
}

func (u *DefaultUploader) CheckFileExists(ctx context.Context, md5 string) (exists bool, err error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	metadata, err := u.MetaService.GetMetadataByMD5(timeoutCtx, md5)
	if err == sql.ErrNoRows {
		// file does not exist
		return false, nil
	}
	if err != nil {
		log.Printf("failed to load metadata: %v", err)
		return false, err
	}
	if metadata.MD5Hash == md5 {
		return true, nil
	}

	return false, nil
}
