package worker

import (
	"github.com/stretchr/testify/mock"
)

type MockAnalyzerWorker struct {
	mock.Mock
}

func (m *MockAnalyzerWorker) InitAnalyzerWorkerPool() {

}
func (m *MockAnalyzerWorker) AnalyzeWebsiteThread() {

}
func (m *MockAnalyzerWorker) SendAnalyzeJob(jobId int64) {

}
