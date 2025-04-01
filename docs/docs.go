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
        "/api/v1/apps": {
            "post": {
                "description": "Add a new app to the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "apps"
                ],
                "summary": "Add a new app",
                "parameters": [
                    {
                        "description": "App object to be added",
                        "name": "app",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.App"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "App added successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/apps/:name": {
            "delete": {
                "description": "Delete all reviews for a given app name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reviews"
                ],
                "summary": "Delete reviews for an app",
                "parameters": [
                    {
                        "type": "string",
                        "description": "App name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reviews deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/reviews": {
            "post": {
                "description": "Add a new review to the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reviews"
                ],
                "summary": "Add a new review",
                "parameters": [
                    {
                        "description": "Review object to be added",
                        "name": "review",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Review"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Review added successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/reviews/:name": {
            "delete": {
                "description": "Delete all reviews for a given app name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reviews"
                ],
                "summary": "Delete reviews for an app",
                "parameters": [
                    {
                        "type": "string",
                        "description": "App name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reviews deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.JSONResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.App": {
            "type": "object",
            "properties": {
                "androidVer": {
                    "type": "string"
                },
                "category": {
                    "type": "string"
                },
                "contentRating": {
                    "type": "string"
                },
                "currentVer": {
                    "type": "string"
                },
                "genres": {
                    "type": "string"
                },
                "installs": {
                    "type": "string"
                },
                "lastUpdated": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "string"
                },
                "rating": {
                    "type": "number"
                },
                "reviews": {
                    "type": "integer"
                },
                "size": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models.Review": {
            "type": "object",
            "properties": {
                "app": {
                    "type": "string"
                },
                "sentiment": {
                    "type": "string"
                },
                "sentimentPolarity": {
                    "type": "number"
                },
                "sentimentSubjectivity": {
                    "type": "number"
                },
                "translatedReview": {
                    "type": "string"
                }
            }
        },
        "utils.JSONResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "status": {
                    "type": "string"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
