package controller

import (
	"fmt"
	"net/http"
	"os"

	services "github.com/Soup666/diss-api/services"

	"github.com/gin-gonic/gin"
)

// AuthController is the controller for handling authentication requests
type ObjectController struct {
	authService *services.AuthService
}

func NewObjectController(authService *services.AuthService) *ObjectController {
	return &ObjectController{authService}
}

func (c *ObjectController) GetObject(ctx *gin.Context) {
	filename := ctx.Param("filename")
	taskId := ctx.Param("taskID")

	// Construct the full file path
	filePath := fmt.Sprintf("objects/task-%s/%s", taskId, filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	// Serve the file
	ctx.File(filePath)
}
