package worker

import (
	"errors"
	"github.com/karlseguin/ccache/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"os/signal"
	"testing"
	"time"
	"website-analyzer/src/configuration"
	"website-analyzer/src/entity"
	"website-analyzer/src/service"
)

func TestInitAnalyzerWorkerPool(t *testing.T) {

	mockService := new(service.MockAnalyzerService)
	mockService.On("AnalyzeWebsiteContent", mock.Anything).Return(entity.Analysis{}, nil)

	cache := ccache.New(ccache.Configure().MaxSize(10).ItemsToPrune(1))
	cache.Set("123456", entity.Analysis{Id: 123456}, time.Duration(10*time.Second))

	configuration.AppConfig.Worker.BufferSize = 1
	configuration.AppConfig.Worker.PoolSize = 1

	analyzerWorkerObj := NewAnalyzerWorker(mockService, cache, make(chan int64, configuration.AppConfig.Worker.BufferSize), make(chan os.Signal))
	analyzerWorkerObj.InitAnalyzerWorkerPool()

	for i := 0; i < configuration.AppConfig.Worker.PoolSize; i++ {
		jobID := int64(123456)
		analyzerWorkerObj.SendAnalyzeJob(jobID)
	}

	time.Sleep(100 * time.Millisecond)

	interruptSignal := os.Interrupt
	signal.Notify(analyzerWorkerObj.Interrupter, interruptSignal)

	for i := 0; i < configuration.AppConfig.Worker.PoolSize; i++ {
		jobID := int64(123456)
		mockService.AssertCalled(t, "AnalyzeWebsiteContent", entity.Analysis{Id: jobID})
	}
}

func TestAnalyzeWebsiteThread(t *testing.T) {
	mockService := new(service.MockAnalyzerService)
	mockService.On("AnalyzeWebsiteContent", mock.Anything).Return(entity.Analysis{}, nil)

	cache := ccache.New(ccache.Configure().MaxSize(10).ItemsToPrune(1))
	cache.Set("123456", entity.Analysis{Id: 123456}, time.Duration(10*time.Second))
	analyzerWorkerObj := NewAnalyzerWorker(mockService, cache, make(chan int64, configuration.AppConfig.Worker.BufferSize), make(chan os.Signal))

	go func() {
		analyzerWorkerObj.AnalyzeWebsiteThread()
	}()

	jobID := int64(123456)
	analyzerWorkerObj.SendAnalyzeJob(jobID)

	time.Sleep(1000 * time.Millisecond)

	interruptSignal := os.Interrupt
	signal.Notify(analyzerWorkerObj.Interrupter, interruptSignal)

	mockService.AssertCalled(t, "AnalyzeWebsiteContent", entity.Analysis{Id: jobID})
}

func TestAnalyzeWebsiteThreadError(t *testing.T) {
	mockService := new(service.MockAnalyzerService)
	mockService.On("AnalyzeWebsiteContent", mock.Anything).Return(entity.Analysis{}, errors.New(""))

	cache := ccache.New(ccache.Configure().MaxSize(10).ItemsToPrune(1))
	cache.Set("123456", entity.Analysis{Id: 123456}, time.Duration(10*time.Second))
	analyzerWorkerObj := NewAnalyzerWorker(mockService, cache, make(chan int64, configuration.AppConfig.Worker.BufferSize), make(chan os.Signal))

	time.Sleep(100 * time.Millisecond)

	go func() {
		analyzerWorkerObj.AnalyzeWebsiteThread()
	}()

	jobID := int64(123456)
	analyzerWorkerObj.SendAnalyzeJob(jobID)
	time.Sleep(10 * time.Millisecond)

	res := cache.Get("123456")

	interruptSignal := os.Interrupt
	signal.Notify(analyzerWorkerObj.Interrupter, interruptSignal)

	assert.Equal(t, entity.ResponseCodeAnalyzeError, res.Value().(entity.Analysis).Error.Code, "Expected the error code")
}
