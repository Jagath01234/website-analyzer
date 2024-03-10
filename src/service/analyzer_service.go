package service

import (
	"website-analyzer/src/entity"
)

type AnalyzerServiceInterface interface {
	AnalyzeWebsiteContent(request entity.Analysis) (entity.Analysis, error)
}
