package controller

import (
	"net/http"
	"strconv"

	"github.com/Soup666/diss-api/model"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
)

type ReportsController struct {
	reportsService services.ReportsService
}

func NewReportsController(reportsService services.ReportsService) *ReportsController {
	return &ReportsController{reportsService: reportsService}
}

func (c *ReportsController) GetReports(ctx *gin.Context) {

	user := ctx.MustGet("user").(*model.User)

	reports, err := c.reportsService.GetReports(uint(user.Model.ID))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"reports": reports})
}

func (c *ReportsController) GetReportByID(ctx *gin.Context) {

	reportIDParam := ctx.Param("reportID")
	reportID, err := strconv.Atoi(reportIDParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	report, err := c.reportsService.GetReport(uint(reportID))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"report": report})
}

func (c *ReportsController) CreateReport(ctx *gin.Context) {
	report := &model.Report{}

	if err := ctx.ShouldBindJSON(report); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user := ctx.MustGet("user").(*model.User)
	report.UserID = user.Model.ID

	err := c.reportsService.CreateReport(report)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"report": report})
}

func (c *ReportsController) ArchiveReport(ctx *gin.Context) {
	reportIDParam := ctx.Param("reportID")
	reportID, err := strconv.Atoi(reportIDParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	if err := c.reportsService.ArchiveReport(uint(reportID)); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report archived successfully"})
}

func (c *ReportsController) SaveReport(ctx *gin.Context) {
	report := &model.Report{}
	if err := ctx.ShouldBindJSON(report); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.reportsService.SaveReport(report); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"report": report})
}
