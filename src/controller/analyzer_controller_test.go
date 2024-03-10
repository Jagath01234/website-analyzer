package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/karlseguin/ccache/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"website-analyzer/src/entity"
	"website-analyzer/src/entity/response"
	"website-analyzer/src/worker"
)

func TestAnalyzeWebsiteContent(t *testing.T) {
	mockWorker := new(worker.MockAnalyzerWorker)
	cache := ccache.New(ccache.Configure().MaxSize(10).ItemsToPrune(1))
	analyzerController := NewAnalyzerController(cache, mockWorker)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/analyze/basic", analyzerController.AnalyzeWebsiteContent)

	data := url.Values{}
	data.Set("target_url", "http://www.google.com")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/analyze/basic", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	data.Set("target_url", "")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/analyze/basic", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	data.Set("target_url", "invalid-url")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/analyze/basic", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	data.Set("target_url", "http://www.%google.com")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/analyze/basic", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	expectedResponse := response.ErrorResponse{Data: entity.AppError{
		Code:    entity.ResponseCodeInvalidUrl,
		Message: "target_url is invalid",
	}}
	expectedResponseBody, _ := json.MarshalIndent(expectedResponse, "", "    ")
	assert.Equal(t, string(expectedResponseBody), w.Body.String(), "Response body should  with the invalid URL response body")
}

func TestGetAnalyzeJobStatusEmptyJobID(t *testing.T) {

	mockWorker := new(worker.MockAnalyzerWorker)
	cache := ccache.New(ccache.Configure().MaxSize(10).ItemsToPrune(1))
	analyzerController := NewAnalyzerController(cache, mockWorker)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/analyze/status", nil)

	analyzerController.GetAnalyzeJobStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestGetAnalyzeJobStatusJobNotFound(t *testing.T) {
	mockWorker := new(worker.MockAnalyzerWorker)
	cache := ccache.New(ccache.Configure().MaxSize(10).ItemsToPrune(1))
	analyzerController := NewAnalyzerController(cache, mockWorker)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/analyze/status", nil)
	c.Request.URL, _ = url.Parse("?job_id=123")
	analyzerController.GetAnalyzeJobStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestGetAnalyzeJobStatusSuccessful(t *testing.T) {
	mockWorker := new(worker.MockAnalyzerWorker)
	cache := ccache.New(ccache.Configure().MaxSize(10).ItemsToPrune(1))
	analyzerController := NewAnalyzerController(cache, mockWorker)
	mockRes := entity.Analysis{
		Id:        123,
		JobStatus: entity.StatusSuccess,
	}
	cache.Set("123", mockRes, 1*time.Minute)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/analyze/status", nil)
	c.Request.URL, _ = url.Parse("?job_id=123")
	analyzerController.GetAnalyzeJobStatus(c)
	jobCacheItem := cache.Get("123")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, jobCacheItem, "Cache return should be nil")

}
