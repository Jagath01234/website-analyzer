package controller

import "github.com/gin-gonic/gin"

type AnalyzerControllerInterface interface {
	AnalyzeWebsiteContent(c *gin.Context)
	GetAnalyzeJobStatus(c *gin.Context)
}
