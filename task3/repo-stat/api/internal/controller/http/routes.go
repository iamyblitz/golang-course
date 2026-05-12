package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, repository *usecase.Repository) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewRepositoryHandler(log, repository))
	mux.Handle("GET /swagger/index.html", NewSwaggerIndexHandler())
	mux.Handle("GET /swagger/openapi.json", NewOpenAPIHandler())
}
