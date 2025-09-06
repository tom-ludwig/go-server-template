package handler

import (
	"com.tom-ludwig/go-server-template/internal/api"
	"com.tom-ludwig/go-server-template/internal/repository"
)

// compile-time check
var _ api.StrictServerInterface = (*Server)(nil)

type Server struct {
	UserHandler
}

func NewServer(queries *repository.Queries) *Server {
	return &Server{
		UserHandler: UserHandler{
			Queries: queries,
		},
	}
}
