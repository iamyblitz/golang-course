package http

import (
	"net/http"
)

const swaggerIndexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Repo Stat API Swagger</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/swagger/openapi.json",
        dom_id: "#swagger-ui"
      });
    };
  </script>
</body>
</html>`

const openAPIJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Repo Stat API",
    "version": "1.0.0",
    "description": "API Gateway for repository statistics services"
  },
  "paths": {
    "/api/ping": {
      "get": {
        "summary": "Check service statuses",
        "responses": {
          "200": {
            "description": "Processor and subscriber are available",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/PingResponse"
                }
              }
            }
          },
          "503": {
            "description": "At least one service is unavailable",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/PingResponse"
                }
              }
            }
          }
        }
      }
    },
    "/api/repositories/info": {
      "get": {
        "summary": "Get GitHub repository information",
        "parameters": [
          {
            "name": "url",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "https://github.com/golang/go"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Repository information",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/RepositoryInfoResponse"
                }
              }
            }
          },
          "400": {
            "description": "Invalid or missing repository URL",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "404": {
            "description": "Repository was not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "503": {
            "description": "Downstream service is unavailable",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    },
    "/subscriptions": {
      "post": {
        "summary": "Create repository subscription",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SubscriptionRequest"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Subscription created",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SubscriptionResponse"
                }
              }
            }
          },
          "400": {
            "description": "Invalid request",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "404": {
            "description": "Repository was not found on GitHub",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "409": {
            "description": "Subscription already exists",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      },
      "get": {
        "summary": "List repository subscriptions",
        "responses": {
          "200": {
            "description": "Subscriptions list",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SubscriptionsResponse"
                }
              }
            }
          }
        }
      }
    },
    "/subscriptions/{owner}/{repo}": {
      "delete": {
        "summary": "Delete repository subscription",
        "parameters": [
          {
            "name": "owner",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "repo",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Subscription deleted"
          },
          "404": {
            "description": "Subscription was not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    },
    "/subscriptions/info": {
      "get": {
        "summary": "Get GitHub information for subscribed repositories",
        "responses": {
          "200": {
            "description": "Repository information for subscriptions",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/RepositoriesInfoResponse"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "PingResponse": {
        "type": "object",
        "properties": {
          "status": {
            "type": "string",
            "example": "ok"
          },
          "services": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/PingService"
            }
          }
        }
      },
      "PingService": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "example": "processor"
          },
          "status": {
            "type": "string",
            "example": "up"
          }
        }
      },
      "RepositoryInfoResponse": {
        "type": "object",
        "properties": {
          "full_name": {
            "type": "string",
            "example": "golang/go"
          },
          "description": {
            "type": "string",
            "example": "The Go programming language"
          },
          "stars": {
            "type": "integer",
            "format": "int64",
            "example": 123456
          },
          "forks": {
            "type": "integer",
            "format": "int64",
            "example": 12345
          },
          "created_at": {
            "type": "string",
            "example": "2009-11-10T23:00:00Z"
          }
        }
      },
      "RepositoriesInfoResponse": {
        "type": "object",
        "properties": {
          "repositories": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/RepositoryInfoResponse"
            }
          }
        }
      },
      "SubscriptionRequest": {
        "type": "object",
        "required": ["owner", "repo"],
        "properties": {
          "owner": {
            "type": "string",
            "example": "golang"
          },
          "repo": {
            "type": "string",
            "example": "go"
          }
        }
      },
      "SubscriptionResponse": {
        "type": "object",
        "properties": {
          "owner": {
            "type": "string",
            "example": "golang"
          },
          "repo": {
            "type": "string",
            "example": "go"
          }
        }
      },
      "SubscriptionsResponse": {
        "type": "object",
        "properties": {
          "subscriptions": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/SubscriptionResponse"
            }
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {
            "type": "string",
            "example": "internal error"
          }
        }
      }
    }
  }
}`

func NewSwaggerIndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(swaggerIndexHTML))
	}
}

func NewOpenAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(openAPIJSON))
	}
}
