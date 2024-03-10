package response

import (
	"website-analyzer/src/entity"
)

type AnalyzerResponse struct {
	Id string `json:"job_id"`
}

type AnalyzeStatusResponse struct {
	Data AnalysisResponseBody `json:"data"`
}

type AnalysisResponseBody struct {
	Id          string               `json:"Id"`
	TargetUrl   string               `json:"target_url"`
	JobStatus   entity.Status        `json:"job_status"`
	Title       string               `json:"title"`
	HtmlVersion string               `json:"html_version"`
	Headings    []entity.HeadingInfo `json:"headings"`
	Links       entity.LinkInfo      `json:"links"`
	IsLogin     bool                 `json:"is_login"`
	Error       entity.AppError      `json:"error"`
}
