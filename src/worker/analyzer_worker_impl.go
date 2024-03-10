package worker

import (
	cache "github.com/karlseguin/ccache/v2"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
	"website-analyzer/src/configuration"
	"website-analyzer/src/entity"
	"website-analyzer/src/service"
)

const logPrefixAnalyzerWorker = "web-analyzer/src/worker/analyzer_worker_impl"

type AnalyzerWorker struct {
	AnalyserService service.AnalyzerServiceInterface
	JobCache        *cache.Cache
	AnalyzerChannel chan int64
	Interrupter     chan os.Signal
}

func NewAnalyzerWorker(AnalyserService service.AnalyzerServiceInterface, jobCache *cache.Cache, AnalyzerChannel chan int64, interrupter chan os.Signal) *AnalyzerWorker {
	return &AnalyzerWorker{
		AnalyserService: AnalyserService,
		JobCache:        jobCache,
		AnalyzerChannel: AnalyzerChannel,
		Interrupter:     interrupter,
	}
}

func (a AnalyzerWorker) InitAnalyzerWorkerPool() {

	for i := 1; i <= configuration.AppConfig.Worker.PoolSize; i++ {
		go a.AnalyzeWebsiteThread()
	}
}

func (a AnalyzerWorker) AnalyzeWebsiteThread() {
	signal.Notify(a.Interrupter, os.Interrupt)
AnalyzerWorkerLoop:
	for {
		select {
		case jobId := <-a.AnalyzerChannel:
			jobIdStr := strconv.FormatInt(jobId, 10)
			job := a.JobCache.Get(jobIdStr)
			if job != nil {
				res, err := a.AnalyserService.AnalyzeWebsiteContent(job.Value().(entity.Analysis))
				if err != nil {
					log.Printf("%v,%v,%v,%v,%v", "ERROR", logPrefixAnalyzerWorker, "error analyzing URL for job_id:", jobIdStr, err.Error())
					res.Error = entity.AppError{
						Code:    entity.ResponseCodeAnalyzeError,
						Message: err.Error(),
					}
				}
				a.JobCache.Set(jobIdStr, res, time.Duration(configuration.AppConfig.Cache.ExpiryTimeSecs)*time.Second)
			}
		case <-a.Interrupter:
			break AnalyzerWorkerLoop
		}
	}
}

func (a AnalyzerWorker) SendAnalyzeJob(jobId int64) {
	a.AnalyzerChannel <- jobId
}
