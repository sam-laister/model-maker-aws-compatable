package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthController is the controller for handling authentication requests
type ObjectController struct{}

func NewObjectController() *ObjectController {
	return &ObjectController{}
}

func (c *ObjectController) GetObject(ctx *gin.Context) {
	taskId := ctx.Param("taskID")

	// Construct the full file path
	filePath := fmt.Sprintf("objects/%s/%s", taskId, "mvs/final_model.glb")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	// Serve the file
	ctx.File(filePath)
}
