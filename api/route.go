package api

import (
	"log"
	"net/http"

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
	basePath := "."
	uploader, err := localstorage.NewUploader(basePath, username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create uploader")
	}
	log.Printf("current uploader is the user: %v", uploader)

	// get file and md5 from request body
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "file is required")
	}
	md5 := c.FormValue("md5")
	// check if the md5 of the file already exists
	// return fileId, true if exists, false if not
	exists, err := uploader.CheckFileExists(filename, md5)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to check if file exists")
	}
	if exists {
		return c.JSON(http.StatusOK, map[string]string{
			"exists": "true",
		})
	}

	// save file to local storage
	err = uploader.SaveFile(file, filename, md5)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to save file")
	}

	return c.String(http.StatusOK, "pong")
}

func downloadFile(c echo.Context) error {

}
