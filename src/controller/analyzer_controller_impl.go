package controller

import (
	"github.com/gin-gonic/gin"
	cache "github.com/karlseguin/ccache/v2"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"website-analyzer/src/configuration"
	"website-analyzer/src/entity"
	"website-analyzer/src/entity/response"
	"website-analyzer/src/util"
	"website-analyzer/src/worker"
)

const logPrefixAnalyzerController = "web-analyzer/src/controller/analyzer_controller_impl"

type AnalyzerController struct {
	jobCache *cache.Cache
	worker   worker.AnalyzerWorkerInterface
}

func NewAnalyzerController(jobCache *cache.Cache, worker worker.AnalyzerWorkerInterface) *AnalyzerController {
	return &AnalyzerController{
		jobCache: jobCache,
		worker:   worker,
	}
}

// @Summary Analyze website content
// @Description Push a job to analyze basic information ona website
// @Tags analyzer
// @Param target_url query string true "Target URL for website content analysis"
// @Produce  json
// @Success 200 {object} response.AnalyzerResponse  // Response containing the job ID
// @Failure 400 {object} response.ErrorResponse  // Error response if target_url is missing
// @Failure 500 {object} response.ErrorResponse  // Error response for internal server error
// @Router /analyze/basic [post]
func (a *AnalyzerController) AnalyzeWebsiteContent(c *gin.Context) {

	ctx := c.Request.Context()
	targetUrl := strings.TrimSpace(c.PostForm("target_url"))
	if targetUrl == "" {
		log.Printf("%v,%v,%v,%v", "ERROR", logPrefixAnalyzerController, ctx, "target_url parameter empty")
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponse{Data: entity.AppError{
			Code:    entity.ResponseCodeEmptyRequestParam,
			Message: "target_url parameter is required",
		},
		})
		return
	}
	if !util.IsValidURL(targetUrl) {
		log.Printf("%v,%v,%v,%v,%v", "ERROR", logPrefixAnalyzerController, ctx, "invalid target_url:", targetUrl)
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponse{Data: entity.AppError{
			Code:    entity.ResponseCodeInvalidUrl,
			Message: "target_url is invalid",
		},
		})
		return
	}
	id := time.Now().UnixNano()
	req := entity.Analysis{
		Id:        id,
		TargetUrl: targetUrl,
		JobStatus: entity.StatusPending,
		Timestamp: time.Now(),
	}
	a.jobCache.Set(strconv.FormatInt(id, 10), req, time.Duration(configuration.AppConfig.Cache.ExpiryTimeSecs)*time.Second)
	a.worker.SendAnalyzeJob(id)
	idString := strconv.FormatInt(id, 10)
	c.IndentedJSON(http.StatusOK, response.AnalyzerResponse{Id: idString})
}

// @Summary Website analysis status
// @Description Analysis result of the website analysis
// @Tags
// @Produce  json
// @Success 200 {object} response.AnalyzeStatusResponse
// @Router /analyze/status [get]
// @Param job_id query int true "Job ID for analysis status"
func (a *AnalyzerController) GetAnalyzeJobStatus(c *gin.Context) {
	ctx := c.Request.Context()
	jobIdStr := c.Query("job_id")
	if len(jobIdStr) == 0 {
		log.Printf("%v,%v,%v,%v", "ERROR", logPrefixAnalyzerController, ctx, "job_id parameter empty")
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponse{Data: entity.AppError{
			Code:    entity.ResponseCodeEmptyRequestParam,
			Message: "job_id is empty",
		},
		})

	} else {
		jobMapJob := a.jobCache.Get(jobIdStr)
		if jobMapJob == nil {
			log.Printf("%v,%v,%v,%v,%v", "ERROR", logPrefixAnalyzerController, ctx, jobIdStr, "job_not found")
			c.IndentedJSON(http.StatusBadRequest, response.ErrorResponse{Data: entity.AppError{
				Code:    entity.ResponseCodeJobNotFound,
				Message: "Job Not found",
			},
			})
			return
		}

		job := jobMapJob.Value().(entity.Analysis)
		resp := response.AnalyzeStatusResponse{
			Data: response.AnalysisResponseBody{
				Id:          jobIdStr,
				TargetUrl:   job.TargetUrl,
				JobStatus:   job.JobStatus,
				Title:       job.Title,
				HtmlVersion: job.HtmlVersion,
				Headings:    job.Headings,
				Links:       job.Links,
				IsLogin:     job.IsLogin,
				Error:       job.Error,
			},
		}
		if job.JobStatus == entity.StatusFail || job.JobStatus == entity.StatusSuccess {
			a.jobCache.Delete(jobIdStr)
		}
		c.IndentedJSON(http.StatusOK, resp)
	}
}
