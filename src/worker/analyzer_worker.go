package worker

type AnalyzerWorkerInterface interface {
	InitAnalyzerWorkerPool()
	AnalyzeWebsiteThread()
	SendAnalyzeJob(jobId int64)
}
