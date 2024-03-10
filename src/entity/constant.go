package entity

type Status int

const (
	StatusPending Status = iota
	StatusSuccess
	StatusFail
)

const (
	ResponseCodeEmptyRequestParam = 1001
	ResponseCodeInvalidUrl        = 1003
	ResponseCodeJobNotFound       = 1004
	ResponseCodeAnalyzeError      = 1005
)
