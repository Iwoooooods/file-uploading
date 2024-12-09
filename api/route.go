package api

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
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
	// filename is the combination of fileid and extension
	e.POST("/upload/:userid/:filename", h.uploadFile)
	e.GET("/download/:username/:filename", h.downloadFile)
	e.DELETE("/delete/:username/:filename", h.deleteFile)
}

func (h *Handler) uploadFile(c echo.Context) error {
	ctx := context.Background()

	userid := c.Param("userid")
	filename := c.Param("filename")
	// load env config from dev.env
	basePath := h.cfg.BasePath
	log.Printf("storing files in: %v", basePath)

	uploader, err := localstorage.NewUploader(basePath, userid, h.db)
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

	// Read the entire content
	if _, err := io.Copy(writer, file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to calculate MD5",
		})
	}
	md5Hash := hex.EncodeToString(hasher.Sum(nil))

	// upload file
	err = uploader.UploadFile(ctx, &buffer, filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to upload file",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"fileName": filename,
		"md5Hash":  md5Hash,
		"url":      fmt.Sprintf("%v/%v/%v", h.cfg.ServerHost, userid, filename),
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

func (h *Handler) deleteFile(c echo.Context) error {
	ctx := context.Background()
	username := c.Param("username")
	filename := c.Param("filename")
	basePath := h.cfg.BasePath

	// Create uploader instance to access metadata
	uploader, err := localstorage.NewUploader(basePath, username, h.db)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create uploader")
	}

	filePath := filepath.Join(uploader.BasePath, filename)

	err = uploader.DeleteFile(ctx, filePath)
	if err != nil {
		fmt.Printf("failed to delete file: %v", err)
		return c.String(http.StatusInternalServerError, "failed to delete file")
	}

	fmt.Printf("file deleted: %v", filePath)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) updateFile(c echo.Context) error {

	return nil
}

