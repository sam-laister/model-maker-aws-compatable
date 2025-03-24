package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"

	"github.com/gin-gonic/gin"
)

// AuthController is the controller for handling authentication requests
type UploadController struct {
}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (c *UploadController) UploadFile(ctx *gin.Context) {

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	// Save the file
	savePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	// Save file metadata in the database
	image := models.AppFile{
		Filename: header.Filename,
		Url:      fmt.Sprintf("/%s", savePath),
	}
	database.DB.Create(&image)

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "image": image})
}

func (c *UploadController) GetFile(ctx *gin.Context) {
	taskId := ctx.Param("taskId")
	filename := ctx.Param("filename")

	// Construct the full file path
	filePath := fmt.Sprintf("uploads/%s/%s", taskId, filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Image not found", "path": filePath})
		return
	}

	// Serve the file
	ctx.File(filePath)
}

func (c *UploadController) GetObject(ctx *gin.Context) {
	filename := ctx.Param("filename")

	// Construct the full file path
	filePath := filepath.Join("objects", filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	// Serve the file
	ctx.File(filePath)
}
