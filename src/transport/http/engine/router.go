package engine

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
	"time"
	"website-analyzer/docs"
	"website-analyzer/src/configuration"
	"website-analyzer/src/controller"
	"website-analyzer/src/metrics"
)

type HttpEngine struct {
	analyzerController controller.AnalyzerControllerInterface
}

func NewHttpEngine(analyzerController controller.AnalyzerControllerInterface) *HttpEngine {
	return &HttpEngine{
		analyzerController: analyzerController,
	}
}

func (e HttpEngine) Engine() *gin.Engine {
	engine := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("*")
	config.AddAllowMethods("*")
	engine.Use(
		cors.New(config),
		gin.Logger(),
		gin.Recovery(),
		uuidAttacherMiddleware(),
		metricsCollectorMiddleware(),
	)

	engine.GET("/ping", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	analyzerGroup := engine.Group("/analyze")
	{
		analyzerGroup.POST("/basic", func(c *gin.Context) {
			e.analyzerController.AnalyzeWebsiteContent(c)
		})

		analyzerGroup.GET("/status", func(c *gin.Context) {
			e.analyzerController.GetAnalyzeJobStatus(c)
		})
	}

	if configuration.AppConfig.ApiDocs.IsEnabled {
		docs.SwaggerInfo.BasePath = "/"
		engine.GET("/doc/ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return engine
}

func uuidAttacherMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New().String()
		ctx := context.WithValue(c.Request.Context(), "uuid", uuid)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func metricsCollectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		metrics.HttpRequestSummary.WithLabelValues(c.Request.Method, c.Request.URL.Path, strconv.Itoa(c.Writer.Status())).Observe(float64(time.Since(t)))
	}
}
