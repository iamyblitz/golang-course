package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := ping.Execute(r.Context())

		response := dto.PingResponse{
			Status: status.Status,
		}
		httpStatus := http.StatusOK
		for _, service := range status.Services {
			response.Services = append(response.Services, dto.PingService{
				Name:   service.Name,
				Status: string(service.Status),
			})
			if service.Status == domain.PingStatusDown {
				httpStatus = http.StatusServiceUnavailable
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
