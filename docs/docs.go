// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/analyze/basic": {
            "post": {
                "description": "Push a job to analyze basic information ona website",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "analyzer"
                ],
                "summary": "Analyze website content",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Target URL for website content analysis",
                        "name": "target_url",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.AnalyzerResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/analyze/status": {
            "get": {
                "description": "Analysis result of the website analysis",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    ""
                ],
                "summary": "Website analysis status",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Job ID for analysis status",
                        "name": "job_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.AnalyzeStatusResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.AppError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.HeadingInfo": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "level": {
                    "type": "string"
                }
            }
        },
        "entity.LinkInfo": {
            "type": "object",
            "properties": {
                "external_links": {
                    "type": "integer"
                },
                "inaccessible_links": {
                    "type": "integer"
                },
                "internal_links": {
                    "type": "integer"
                }
            }
        },
        "entity.Status": {
            "type": "integer",
            "enum": [
                0,
                1,
                2
            ],
            "x-enum-varnames": [
                "StatusPending",
                "StatusSuccess",
                "StatusFail"
            ]
        },
        "response.AnalysisResponseBody": {
            "type": "object",
            "properties": {
                "Id": {
                    "type": "string"
                },
                "error": {
                    "$ref": "#/definitions/entity.AppError"
                },
                "headings": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.HeadingInfo"
                    }
                },
                "html_version": {
                    "type": "string"
                },
                "is_login": {
                    "type": "boolean"
                },
                "job_status": {
                    "$ref": "#/definitions/entity.Status"
                },
                "links": {
                    "$ref": "#/definitions/entity.LinkInfo"
                },
                "target_url": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "response.AnalyzeStatusResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/response.AnalysisResponseBody"
                }
            }
        },
        "response.AnalyzerResponse": {
            "type": "object",
            "properties": {
                "job_id": {
                    "type": "string"
                }
            }
        },
        "response.ErrorResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/entity.AppError"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
