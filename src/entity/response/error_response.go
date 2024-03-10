package response

import "website-analyzer/src/entity"

type ErrorResponse struct {
	Data entity.AppError `json:"data"`
}
