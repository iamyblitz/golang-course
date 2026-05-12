package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRepositoryHandler(log *slog.Logger, repository *usecase.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repoURL := r.URL.Query().Get("url")
		if repoURL == "" {
			writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{
				Error: "url query parameter is required",
			})
			return
		}

		parsedURL, err := url.ParseRequestURI(repoURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{
				Error: "invalid repository url",
			})
			return
		}

		info, err := repository.GetInfo(r.Context(), repoURL)
		if err != nil {
			log.Error("failed to get repository info", "error", err)

			httpStatus := httpStatusFromError(err)
			writeJSON(w, httpStatus, dto.ErrorResponse{
				Error: errorMessage(httpStatus),
			})
			return
		}

		writeJSON(w, http.StatusOK, dto.RepositoryInfoResponse{
			FullName:    info.FullName,
			Description: info.Description,
			Stars:       info.Stars,
			Forks:       info.Forks,
			CreatedAt:   info.CreatedAt,
		})
	}
}

func NewSubscriptionsInfoHandler(log *slog.Logger, repository *usecase.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repositories, err := repository.GetSubscriptionsInfo(r.Context())
		if err != nil {
			log.Error("failed to get subscriptions info", "error", err)

			httpStatus := httpStatusFromError(err)
			writeJSON(w, httpStatus, dto.ErrorResponse{
				Error: errorMessage(httpStatus),
			})
			return
		}

		response := dto.RepositoriesInfoResponse{
			Repositories: make([]dto.RepositoryInfoResponse, 0, len(repositories)),
		}
		for _, info := range repositories {
			response.Repositories = append(response.Repositories, dto.RepositoryInfoResponse{
				FullName:    info.FullName,
				Description: info.Description,
				Stars:       info.Stars,
				Forks:       info.Forks,
				CreatedAt:   info.CreatedAt,
			})
		}

		writeJSON(w, http.StatusOK, response)
	}
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func httpStatusFromError(err error) int {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func errorMessage(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad request"
	case http.StatusNotFound:
		return "repository not found"
	case http.StatusServiceUnavailable:
		return "service unavailable"
	default:
		return "failed to get repository info"
	}
}
