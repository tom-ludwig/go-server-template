package routes

import (
	"log"

	"com.tom-ludwig/go-server-template/internal/api"
	"com.tom-ludwig/go-server-template/internal/handler"
	"com.tom-ludwig/go-server-template/internal/middleware"
	"com.tom-ludwig/go-server-template/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/gin-middleware"
)

func NewRouter(queries *repository.Queries) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders)

	swagger, _ := api.GetSwagger()
	r.Use(ginmiddleware.OapiRequestValidator(swagger))

	err := r.SetTrustedProxies([]string{
		"10.0.0.0/8",     // Internal Kubernetes networking
		"172.16.0.0/12",  // Docker default bridge network
		"192.168.0.0/16", // on-premise
	})
	if err != nil {
		log.Fatalf("Error while setting trusted proxies: %s", err)
	}

	server := handler.NewServer(queries)
	strictserver := api.NewStrictHandler(server, nil)
	api.RegisterHandlers(r, strictserver)

	return r
}
