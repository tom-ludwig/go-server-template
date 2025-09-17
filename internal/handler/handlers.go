package handler

import (
	"com.tom-ludwig/go-server-template/internal/api"
	"com.tom-ludwig/go-server-template/internal/repository"
)

// compile-time check
var _ api.StrictServerInterface = (*Server)(nil)

type Server struct {
	HealthHandler
	UserHandler
}

func NewServer(queries *repository.Queries) *Server {
	return &Server{
		HealthHandler: HealthHandler{
			queries: queries,
		},
		UserHandler: UserHandler{
			Queries: queries,
		},
	}
}
