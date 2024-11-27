package api

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Iwoooooods/fs-upload-go/internal/config"
	"github.com/Iwoooooods/fs-upload-go/internal/localstorage"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	cfg *config.Config
	db  *sql.DB
}

func NewHandler(cfg *config.Config, db *sql.DB) *Handler {
	return &Handler{
		cfg: cfg,
		db:  db,
	}
}

func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/upload/:username/:filename", h.uploadFile)
	e.GET("/download/:username/:filename", h.downloadFile)
}

func (h *Handler) uploadFile(c echo.Context) error {
	ctx := context.Background()

	username := c.Param("username")
	filename := c.Param("filename")
	// load env config from dev.env
	basePath := h.cfg.BasePath
	log.Printf("storing files in: %v", basePath)

	uploader, err := localstorage.NewUploader(basePath, username, h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create uploader",
		})
	}
	log.Printf("current uploader is the user: %v", uploader)

	var buffer bytes.Buffer
	// get file and md5 from request body
	file := c.Request().Body

	// generate md5 hash of the file
	hasher := md5.New()
	writer := io.MultiWriter(&buffer, hasher)

	// Read the entire content through the TeeReader first
	if _, err := io.Copy(writer, file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to calculate MD5",
		})
	}
	md5Hash := hex.EncodeToString(hasher.Sum(nil))

	// check if the md5 of the file already exists
	// return fileId, true if exists, false if not
	exists, err := uploader.CheckFileExists(ctx, md5Hash)
	if err != nil {
		log.Printf("failed to check if file exists: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to check if file exists",
		})
	}
	if exists {
		return c.JSON(http.StatusOK, map[string]string{
			"exists": "true",
		})
	}

	// save file to local storage
	fileId, err := uploader.UploadFile(ctx, &buffer, filename, md5Hash)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to save file")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"fileId": fileId,
	})
}

func (h *Handler) downloadFile(c echo.Context) error {
	// Get username and filename from parameters
	username := c.Param("username")
	filename := c.Param("filename")
	basePath := h.cfg.BasePath

	// Create uploader instance to access metadata
	uploader, err := localstorage.NewUploader(basePath, username, h.db)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create uploader")
	}

	// Get file path
	filePath := filepath.Join(uploader.BasePath, filename)
	log.Printf("getting file from: %v", filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "file not found")
	}

	// Open and return the file
	return c.File(filePath)
}
