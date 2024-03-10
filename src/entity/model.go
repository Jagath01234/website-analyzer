package entity

import (
	"time"
)

type Analysis struct {
	Id          int64         `json:"Id"`
	TargetUrl   string        `json:"target_url"`
	JobStatus   Status        `json:"job_status"`
	Title       string        `json:"title"`
	HtmlVersion string        `json:"html_version"`
	Headings    []HeadingInfo `json:"headings"`
	Links       LinkInfo      `json:"links"`
	IsLogin     bool          `json:"is_login"`
	Error       AppError      `json:"error"`
	Timestamp   time.Time     `json:"-"`
}

type HeadingInfo struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}
type LinkInfo struct {
	InternalLinks     int `json:"internal_links"`
	ExternalLinks     int `json:"external_links"`
	InaccessibleLinks int `json:"inaccessible_links"`
}
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
