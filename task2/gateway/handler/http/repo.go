package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/iamyblitz/golang-course/task2/proto"
)

type repoGetter interface {
	GetRepo(ctx context.Context, owner, repo string) (*pb.RepoResponse, error)
}

type Handler struct {
	repoGetter repoGetter
	timeout    time.Duration
}

func New(repoGetter repoGetter, timeout time.Duration) *Handler {
	return &Handler{
		repoGetter: repoGetter,
		timeout:    timeout,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/repos/", h.handleGetRepo)
	mux.HandleFunc("/swagger/openapi.yaml", h.handleOpenAPI)
	mux.HandleFunc("/swagger/", h.handleSwaggerUI)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) handleGetRepo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/repos/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.Error(w, "invalid path, expected /repos/{owner}/{repo}", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.repoGetter.GetRepo(ctx, parts[0], parts[1])
	if err != nil {
		h.writeGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"name":        resp.Name,
		"description": resp.Description,
		"stars":       resp.Stars,
		"forks":       resp.Forks,
		"created_at":  resp.CreatedAt,
	}); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
	}
}

func (h *Handler) writeGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, "collector unavailable", http.StatusBadGateway)
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		http.Error(w, st.Message(), http.StatusBadRequest)
	case codes.NotFound:
		http.Error(w, st.Message(), http.StatusNotFound)
	case codes.Unavailable:
		http.Error(w, "collector unavailable", http.StatusBadGateway)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (h *Handler) handleOpenAPI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write([]byte(openAPISpec))
}

func (h *Handler) handleSwaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerHTML))
}

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Task2 Swagger</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: '/swagger/openapi.yaml',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>`

const openAPISpec = `openapi: 3.0.3
info:
  title: GitHub Repository Info API
  version: 1.0.0
servers:
  - url: http://localhost:8080
paths:
  /health:
    get:
      summary: Health check
      responses:
        '200':
          description: Service is healthy
  /repos/{owner}/{repo}:
    get:
      summary: Get repository information
      parameters:
        - in: path
          name: owner
          required: true
          schema:
            type: string
        - in: path
          name: repo
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Repository information
          content:
            application/json:
              schema:
                type: object
                properties:
                  name:
                    type: string
                  description:
                    type: string
                    nullable: true
                  stars:
                    type: integer
                  forks:
                    type: integer
                  created_at:
                    type: string
        '400':
          description: Invalid request
        '404':
          description: Repository not found
        '500':
          description: Internal error
        '502':
          description: Collector unavailable
`
