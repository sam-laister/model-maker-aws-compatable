package controller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"
	services "github.com/Soup666/diss-api/services"

	"github.com/gin-gonic/gin"
)

// AuthController is the controller for handling authentication requests
type UploadController struct {
	authService *services.AuthService
}

func NewUploadController(authService *services.AuthService) *UploadController {
	return &UploadController{authService}
}

func (c *UploadController) UploadFile(ctx *gin.Context) {

	// Extract API key from request header
	apiKey := ctx.GetHeader("Authorization")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is missing"})
		return
	}

	// Remove "Bearer " if present
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")

	_, err := c.authService.FireAuth.VerifyIDToken(context.Background(), apiKey)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

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
	image := models.Image{
		Filename: header.Filename,
		Url:      fmt.Sprintf("/%s", savePath),
	}
	database.DB.Create(&image)

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "image": image})
}

func (c *UploadController) GetFile(ctx *gin.Context) {
	filename := ctx.Param("filename")

	// Construct the full file path
	filePath := filepath.Join("uploads", filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
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
