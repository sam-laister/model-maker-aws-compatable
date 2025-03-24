package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	repositories "github.com/Soup666/diss-api/repository"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
)

type VisionController struct {
	VisionService  services.VisionService
	TaskRepository repositories.TaskRepository
	TaskService    services.TaskService
}

func NewVisionController(visionService services.VisionService, taskRepository repositories.TaskRepository, taskService services.TaskService) *VisionController {
	return &VisionController{VisionService: visionService, TaskRepository: taskRepository, TaskService: taskService}
}

func (c *VisionController) AnalyzeTask(ctx *gin.Context) {

	taskIdParam := ctx.Param("taskID")
	taskId, err := strconv.Atoi(taskIdParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := c.TaskRepository.GetTaskByID(uint(taskId))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.TaskService.FullyLoadTask(task)

	if (len(task.Images)) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No images found for the task"})
		return
	}

	result, err := c.VisionService.AnalyseImage(fmt.Sprintf("./uploads/%d/%s", task.Id, task.Images[0].Filename), "")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to analyze the image"})
		return
	}

	c.TaskService.UpdateMeta(task, "ai-description", result)

	ctx.JSON(http.StatusOK, gin.H{"message": result})
}

func (c *VisionController) AnalyzeImage(ctx *gin.Context) {

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	f, err := os.CreateTemp("", "sample")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	fmt.Println("Temp file name:", f.Name())

	_, err = io.Copy(f, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	image := fmt.Sprintf("/%s", f.Name())

	defer os.Remove(f.Name())

	result, err := c.VisionService.AnalyseImage(image, "")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to analyze the image"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": result})
}
