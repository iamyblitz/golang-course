package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}

type ServiceStatus struct {
	Name   string
	Status domain.PingStatus
}

type PingResult struct {
	Status   string
	Services []ServiceStatus
}

type Ping struct {
	processor  Pinger
	subscriber Pinger
}

func NewPing(processor Pinger, subscriber Pinger) *Ping {
	return &Ping{
		processor:  processor,
		subscriber: subscriber,
	}
}

func (u *Ping) Execute(ctx context.Context) PingResult {
	services := []ServiceStatus{
		{
			Name:   "processor",
			Status: u.processor.Ping(ctx),
		},
		{
			Name:   "subscriber",
			Status: u.subscriber.Ping(ctx),
		},
	}

	status := "ok"
	for _, service := range services {
		if service.Status == domain.PingStatusDown {
			status = "degraded"
			break
		}
	}

	return PingResult{
		Status:   status,
		Services: services,
	}
}
