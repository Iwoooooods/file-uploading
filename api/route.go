package api

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Iwoooooods/fs-upload-go/pkg/blobstorage"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group) {
	e.GET("/ping", ping)
	e.POST("/upload/:filename", uploadFile)
	e.GET("/download/:downloadPath", downloadFile)
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func uploadFile(c echo.Context) error {
	fileName := c.Param("filename")
	fs := blobstorage.NewLocalFS("./tmp/")
	if fileName == "" {
		return c.String(http.StatusBadRequest, "filename is required")
	}
	file, err := fs.Upload(context.Background(), c.Request().Body, fileName)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, file)
}

func downloadFile(c echo.Context) error {
	downloadPath := c.Param("downloadPath")
	log.Println("ext of file: ", downloadPath)
	ext := strings.Split(downloadPath, ".")[1]

	fs := blobstorage.NewLocalFS("./tmp/")
	reader, err := fs.Read(context.Background(), downloadPath)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	contentType := "application/octet-stream"

	switch ext {
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	case "gif":
		contentType = "image/gif"
	case "webp":
		contentType = "image/webp"
	}

	c.Response().Header().Set("Content-Type", contentType)

	return c.Stream(http.StatusOK, contentType, reader)
}
