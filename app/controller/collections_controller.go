package controller

import (
	"net/http"
	"strconv"

	"github.com/Soup666/diss-api/model"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
)

type CollectionsController struct {
	collectionsService services.CollectionsService
}

func NewCollectionsController(collectionsService services.CollectionsService) *CollectionsController {
	return &CollectionsController{collectionsService: collectionsService}
}

func (c *CollectionsController) CreateCollection(ctx *gin.Context) {
	collection := &model.Collection{}

	if err := ctx.ShouldBindJSON(collection); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user := ctx.MustGet("user").(*model.User)
	collection.UserID = user.Model.ID

	err := c.collectionsService.CreateCollection(collection)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"collection": collection})
}

func (c *CollectionsController) GetCollection(ctx *gin.Context) {
	collectionIDParam := ctx.Param("collectionID")
	collectionID, err := strconv.Atoi(collectionIDParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	collection, err := c.collectionsService.GetCollection(uint(collectionID))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"collection": collection})
}

func (c *CollectionsController) GetCollections(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	collections, err := c.collectionsService.GetCollections(user.Model.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"collections": collections})
}

func (c *CollectionsController) ArchiveCollection(ctx *gin.Context) {
	collectionIDParam := ctx.Param("collectionID")
	collectionID, err := strconv.Atoi(collectionIDParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	if err := c.collectionsService.ArchiveCollection(uint(collectionID)); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *CollectionsController) SaveCollection(ctx *gin.Context) {
	collection := &model.Collection{}

	if err := ctx.ShouldBindJSON(collection); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user := ctx.MustGet("user").(*model.User)
	collection.UserID = user.Model.ID

	err := c.collectionsService.SaveCollection(collection)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"collection": collection})
}
