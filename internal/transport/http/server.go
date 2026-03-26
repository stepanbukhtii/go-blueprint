package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/rest"
	easymiddleware "github.com/stepanbukhtii/easy-tools/rest/middleware"
	"github.com/stepanbukhtii/easy-tools/rest/swagger"

	"github.com/stepanbukhtii/go-blueprint/docs"
	"github.com/stepanbukhtii/go-blueprint/internal/app"
	"github.com/stepanbukhtii/go-blueprint/internal/config"
)

//	@title			Blueprint Swagger Example API
//	@version		1.0
//	@description	This is a sample server

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

type Server struct {
	config        config.Config
	services      *app.Services
	router        *gin.Engine
	jwtMiddleware *easymiddleware.JWTAuth
}

func NewServer(app *app.App) (*http.Server, error) {
	gin.SetMode(gin.ReleaseMode)

	server := &Server{
		config:   app.Config,
		services: app.Services,
		//		router:   rest.NewRouter(app.Config.API, app.Config.Service),
	}

	router := rest.NewRouter(app.Config.API, app.Config.Service)

	if err := swagger.RegisterSwagger(router, docs.FilesFS); err != nil {
		return nil, err
	}

	if err := server.registerJWTMiddleware(); err != nil {
		return nil, err
	}

	server.registerMiddlewares(router)

	server.registerHandlers(router)

	return rest.NewServer(app.Config.API, router), nil
}
