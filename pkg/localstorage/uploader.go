package localstorage

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const BUFFER_SIZE = 1024 * 1024 * 4

func NewUploader(serverURL string, username string) (*DefaultUploader, error) {
	// check if user exists, if not create new fold for the user
	basePath := serverURL + "/" + username + "/"
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// give permission to the user to read and write to the folder
		os.MkdirAll(basePath, os.ModePerm)
	}

	metaManager, err := NewMetadataManager(basePath)
	if err != nil {
		log.Printf("failed to create metadata manager: %v", err)
		return nil, err
	}

	return &DefaultUploader{
		ServerURL:   serverURL,
		Username:    username,
		BasePath:    basePath,
		MetaManager: metaManager,
	}, nil
}

func (u *DefaultUploader) UploadFile(src io.Reader, fileName string, md5 string) (fileId string, err error) {
	buffer := make([]byte, BUFFER_SIZE)
	filePath := filepath.Join(u.BasePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.CopyBuffer(file, src, buffer)
	if err != nil {
		return "", err
	}

	newMetadata := make(map[string]FileMetadata)
	newMetadata[fileName] = FileMetadata{
		FileId:   uuid.New().String(),
		FileName: fileName,
		MD5Hash:  md5,
	}
	err = u.MetaManager.SaveMetadata(u.BasePath, newMetadata)
	if err != nil {
		return "", err
	}

	return md5, nil
}

func (u *DefaultUploader) CheckFileExists(fileName string, md5 string) (exists bool, err error) {
	// check if file exists
	if _, err := os.Stat(filepath.Join(u.BasePath, fileName)); os.IsNotExist(err) {
		// if not exists, create new metadata.json and return fileId and false
		newMetadata := make(map[string]FileMetadata)
		newMetadata[fileName] = FileMetadata{
			FileName: fileName,
			MD5Hash:  md5,
		}

		err = u.MetaManager.SaveMetadata(u.BasePath, newMetadata)
		if err != nil {
			log.Printf("failed to save metadata: %v", err)
			return false, err
		}
		return false, nil
	}

	// if exists, check if md5 changes
	// if md5 changes, update metadata.json
	// if md5 does not change, directly return fileId
	metadataFilePath := filepath.Join(u.BasePath, "metadata.json")
	oldMetadata, err := u.MetaManager.LoadMetadata(metadataFilePath)
	if err != nil {
		log.Printf("failed to load metadata: %v", err)
		return false, err
	}

	for _, fileMetadata := range oldMetadata {
		if fileMetadata.MD5Hash == md5 {
			// file exists
			return true, nil
		}
	}

	return false, nil
}
