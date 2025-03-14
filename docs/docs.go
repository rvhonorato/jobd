// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/api/get/{id}": {
            "get": {
                "description": "Fetches a job by its ` + "`" + `id` + "`" + ` (provided by the user) with partial content handling",
                "produces": [
                    "application/json"
                ],
                "summary": "Retrieve a job from the queue",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved job",
                        "schema": {
                            "$ref": "#/definitions/jobs.Job"
                        }
                    },
                    "206": {
                        "description": "Partially completed job",
                        "schema": {
                            "$ref": "#/definitions/jobs.Job"
                        }
                    },
                    "404": {
                        "description": "Job not found",
                        "schema": {
                            "$ref": "#/definitions/errors.RestErr"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.RestErr"
                        }
                    }
                }
            }
        },
        "/api/upload": {
            "post": {
                "description": "Upload a payload. ` + "`" + `id` + "`" + ` is a unique user-provided job identificator. The ` + "`" + `input` + "`" + ` field must contain a base64 encoded` + "`" + `.zip` + "`" + ` file with a ` + "`" + `run.sh` + "`" + ` script and the input data. ` + "`" + `slurml` + "`" + ` marks the job for redirection to the ` + "`" + `slurml` + "`" + ` endpoint (wip)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Upload a new job to the queue",
                "parameters": [
                    {
                        "description": "Job to be uploaded",
                        "name": "job",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/jobs.Upload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Job successfully created",
                        "schema": {
                            "$ref": "#/definitions/jobs.Job"
                        }
                    },
                    "400": {
                        "description": "Bad request - validation error",
                        "schema": {
                            "$ref": "#/definitions/errors.RestErr"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.RestErr"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errors.RestErr": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "jobs.Job": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "input": {
                    "type": "string"
                },
                "lastUpdated": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "output": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "slurmID": {
                    "type": "integer"
                },
                "slurml": {
                    "type": "boolean"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "jobs.Upload": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "input": {
                    "type": "string"
                },
                "slurml": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "jobd (Job Daemon) API",
	Description:      "API for managing job queue in jobd application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
