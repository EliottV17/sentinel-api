package checker

import (
	"context"

	"github.com/EliottV17/sentinel-worker/internal/models"
)

type Monitor = models.Monitor

type Checker interface {
	Check(ctx context.Context, monitor Monitor) (Result, error)
}

type Result struct {
	State          string
	LatencyMs      float64
	StatusCode     *int
	ResponseSample *string
	ErrorMessage   *string
	ExtraData      map[string]any
}