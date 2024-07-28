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
        "/message": {
            "post": {
                "description": "Создает новое сообщение и сохраняет его в базе данных",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "Создать сообщение",
                "parameters": [
                    {
                        "description": "Сообщение",
                        "name": "message",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Message_Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Сообщение успешно создано",
                        "schema": {
                            "$ref": "#/definitions/models.Message_Request"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/stats": {
            "get": {
                "description": "Get statistics from the database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stats"
                ],
                "summary": "Get statistics",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit the number of results",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Statistics retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/models.Stats"
                        }
                    },
                    "500": {
                        "description": "Failed to get or encode stats",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Message_Request": {
            "type": "object",
            "properties": {
                "text": {
                    "type": "string"
                }
            }
        },
        "models.Stats": {
            "type": "object",
            "properties": {
                "failed_messages": {
                    "description": "Количество сообщений, которые не удалось обработать",
                    "type": "integer"
                },
                "pending_messages": {
                    "description": "Количество сообщений в ожидании",
                    "type": "integer"
                },
                "processed_messages": {
                    "description": "Количество обработанных сообщений",
                    "type": "integer"
                },
                "total_messages": {
                    "description": "Общее количество сообщений",
                    "type": "integer"
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
