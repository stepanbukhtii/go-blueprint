package http

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/http/handlers"
)

func (s *Server) registerHandlers(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")

	s.registerLoginHandlers(apiV1)
	s.registerUserHandlers(apiV1)
	s.registerUserTypeHandlers(apiV1)
	s.registerCompanyHandlers(apiV1)
}

func (s *Server) registerLoginHandlers(r *gin.RouterGroup) {
	h := handlers.NewAuth(s.services.Auth)

	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
	}
}

func (s *Server) registerUserHandlers(r *gin.RouterGroup) {
	h := handlers.NewUser(s.services.User)

	users := r.Group("/users")
	{
		users.GET("", h.List)
		users.POST("", h.Create)

		userID := users.Group("/:user_id")
		{
			userID.GET("", h.Get)
			userID.PATCH("", h.Update)
			userID.DELETE("", h.Delete)
		}
	}
}

func (s *Server) registerUserTypeHandlers(r *gin.RouterGroup) {
	h := handlers.NewUser(s.services.User)

	users := r.Group("/user-types")
	{
		users.GET("", h.List)
		users.POST("", h.Create)

		userID := users.Group("/:user_type_code")
		{
			userID.GET("", h.Get)
			userID.PATCH("", h.Update)
			userID.DELETE("", h.Delete)
		}
	}
}

func (s *Server) registerCompanyHandlers(r *gin.RouterGroup) {
	h := handlers.NewCompany(s.services.Company)

	companies := r.Group("/companies", s.jwtMiddleware.Auth)
	{
		companies.GET("", h.List)
		companies.POST("", h.Create)

		companyID := companies.Group("/:company_id")
		{
			companyID.GET("", h.Get)
			companyID.PATCH("", h.Update, s.jwtMiddleware.AuthRole(domain.RoleAdmin))
			companyID.DELETE("", h.Delete, s.jwtMiddleware.AuthRole(domain.RoleAdmin))
		}

		owner := companies.Group("/owner")
		{
			owner.GET("", h.GetCompanyByOwner)
		}

		multiple := companies.Group("/multiple")
		{
			multiple.POST("", h.CreateMultiple)
		}
	}
}
