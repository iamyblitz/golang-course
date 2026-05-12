package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewCreateSubscriptionHandler(log *slog.Logger, subscription *usecase.Subscription) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request dto.SubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
			return
		}

		created, err := subscription.Create(r.Context(), request.Owner, request.Repo)
		if err != nil {
			log.Error("failed to create subscription", "error", err)
			writeSubscriptionError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, subscriptionToResponse(created))
	}
}

func NewDeleteSubscriptionHandler(log *slog.Logger, subscription *usecase.Subscription) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")

		if err := subscription.Delete(r.Context(), owner, repo); err != nil {
			log.Error("failed to delete subscription", "error", err)
			writeSubscriptionError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func NewListSubscriptionsHandler(log *slog.Logger, subscription *usecase.Subscription) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptions, err := subscription.List(r.Context())
		if err != nil {
			log.Error("failed to list subscriptions", "error", err)
			writeSubscriptionError(w, err)
			return
		}

		response := dto.SubscriptionsResponse{
			Subscriptions: make([]dto.SubscriptionResponse, 0, len(subscriptions)),
		}
		for _, subscription := range subscriptions {
			response.Subscriptions = append(response.Subscriptions, subscriptionToResponse(subscription))
		}

		writeJSON(w, http.StatusOK, response)
	}
}

func subscriptionToResponse(subscription domain.Subscription) dto.SubscriptionResponse {
	return dto.SubscriptionResponse{
		Owner: subscription.Owner,
		Repo:  subscription.Repo,
	}
}

func writeSubscriptionError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to process subscription"})
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: st.Message()})
	case codes.NotFound:
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: st.Message()})
	case codes.AlreadyExists:
		writeJSON(w, http.StatusConflict, dto.ErrorResponse{Error: st.Message()})
	case codes.Unavailable:
		writeJSON(w, http.StatusServiceUnavailable, dto.ErrorResponse{Error: "service unavailable"})
	default:
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to process subscription"})
	}
}
