package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
)

func AddRoutes(
	mux *http.ServeMux,
	log *slog.Logger,
	ping *usecase.Ping,
	repository *usecase.Repository,
	subscription *usecase.Subscription,
) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewRepositoryHandler(log, repository))
	mux.Handle("GET /subscriptions/info", NewSubscriptionsInfoHandler(log, repository))
	mux.Handle("POST /subscriptions", NewCreateSubscriptionHandler(log, subscription))
	mux.Handle("DELETE /subscriptions/{owner}/{repo}", NewDeleteSubscriptionHandler(log, subscription))
	mux.Handle("GET /subscriptions", NewListSubscriptionsHandler(log, subscription))
	mux.Handle("GET /swagger/index.html", NewSwaggerIndexHandler())
	mux.Handle("GET /swagger/openapi.json", NewOpenAPIHandler())
}
