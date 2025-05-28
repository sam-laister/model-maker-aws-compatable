package controller

import (
	"github.com/Soup666/diss-api/model"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
)

type UserAnalyticsController struct {
	userAnalyticsService services.UserAnalyticsService
}

func NewUserAnalyticsController(userAnalyticsService services.UserAnalyticsService) *UserAnalyticsController {
	return &UserAnalyticsController{userAnalyticsService: userAnalyticsService}
}

func (c *UserAnalyticsController) GetAnalytics(ctx *gin.Context) {

	user := ctx.MustGet("user")
	userID := user.(*model.User).Model.ID

	analytics, err := c.userAnalyticsService.GetAnalytics(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"analytics": analytics})
}
