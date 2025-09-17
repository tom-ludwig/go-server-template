package handler

import (
	"context"

	"com.tom-ludwig/go-server-template/internal/api"
	"com.tom-ludwig/go-server-template/internal/repository"
)

type HealthHandler struct {
	queries *repository.Queries
}

func (s *HealthHandler) GetHealthz(ctx context.Context, request api.GetHealthzRequestObject) (api.GetHealthzResponseObject, error) {
	return api.GetHealthz200JSONResponse{
		Status: "OK",
	}, nil
}

// GetLivez implements api.StrictServerInterface.
func (s *HealthHandler) GetLivez(ctx context.Context, request api.GetLivezRequestObject) (api.GetLivezResponseObject, error) {
	return api.GetLivez200JSONResponse{
		Status: "OK",
	}, nil
}

// GetReadyz implements api.StrictServerInterface.
func (s *HealthHandler) GetReadyz(ctx context.Context, request api.GetReadyzRequestObject) (api.GetReadyzResponseObject, error) {
	_, err := s.queries.Ping(ctx)
	if err != nil {
		return api.GetReadyz503JSONResponse{
			FailedChecks:      []string{"Database not reachable."},
			SuccessfullChecks: []string{},
		}, nil
	}

	return api.GetReadyz200JSONResponse{
		Status: "OK",
	}, nil
}
