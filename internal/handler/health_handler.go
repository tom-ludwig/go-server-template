package handler

import (
	"context"

	"com.tom-ludwig/go-server-template/internal/api/health"
	"com.tom-ludwig/go-server-template/internal/repository"
)

// compile-time check
var _ health.StrictServerInterface = (*HealthHandler)(nil)

type HealthHandler struct {
	queries *repository.Queries
}

func NewHealthHandler(queries *repository.Queries) *HealthHandler {
	return &HealthHandler{
		queries: queries,
	}
}

func (s *HealthHandler) GetHealthz(_ context.Context, _ health.GetHealthzRequestObject) (health.GetHealthzResponseObject, error) {
	return health.GetHealthz200JSONResponse{
		Status: "OK",
	}, nil
}

// GetLivez implements health.StrictServerInterface.
func (s *HealthHandler) GetLivez(_ context.Context, _ health.GetLivezRequestObject) (health.GetLivezResponseObject, error) {
	return health.GetLivez200JSONResponse{
		Status: "OK",
	}, nil
}

// GetReadyz implements health.StrictServerInterface.
func (s *HealthHandler) GetReadyz(ctx context.Context, _ health.GetReadyzRequestObject) (health.GetReadyzResponseObject, error) {
	_, err := s.queries.Ping(ctx)
	if err != nil {
		return health.GetReadyz503JSONResponse{
			FailedChecks:      []string{"Database not reachable."},
			SuccessfullChecks: []string{},
		}, nil
	}

	return health.GetReadyz200JSONResponse{
		Status: "OK",
	}, nil
}
