package api

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Iwoooooods/fs-upload-go/pkg/localstorage"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group) {
	e.GET("/ping", ping)
	e.POST("/upload/:username/:filename", uploadFile)
	e.GET("/download/:username/:filename", downloadFile)
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func uploadFile(c echo.Context) error {
	username := c.Param("username")
	filename := c.Param("filename")
	basePath := "./private/files"
	uploader, err := localstorage.NewUploader(basePath, username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create uploader")
	}
	log.Printf("current uploader is the user: %v", uploader)

	// get file and md5 from request body
	file := c.Request().Body

	// generate md5 hash of the file
	hasher := md5.New()
	tee := io.TeeReader(file, hasher)
	md5 := hex.EncodeToString(hasher.Sum(nil))

	// check if the md5 of the file already exists
	// return fileId, true if exists, false if not
	exists, err := uploader.CheckFileExists(filename, md5)
	if err != nil {
		log.Printf("failed to check if file exists: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "file exists",
		})
	}
	if exists {
		return c.JSON(http.StatusOK, map[string]string{
			"exists": "true",
		})
	}

	// save file to local storage
	fileId, err := uploader.UploadFile(tee, filename, md5)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to save file")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"fileId": fileId,
	})
}

func downloadFile(c echo.Context) error {
	// Get username and filename from parameters
	username := c.Param("username")
	filename := c.Param("filename")
	basePath := "./private/files"

	// Create uploader instance to access metadata
	uploader, err := localstorage.NewUploader(basePath, username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create uploader")
	}

	// Get file path
	filePath := filepath.Join(uploader.BasePath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "file not found")
	}

	// Open and return the file
	return c.File(filePath)
}
