package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/http/handlers/request"
)

type Auth struct {
	service domain.AuthService
}

func NewAuth(service domain.AuthService) *Auth {
	return &Auth{
		service: service,
	}
}

// Login godoc
//
//	@Summary	List of users
//	@Tags		auth
//	@Param		request	body		request.Login	true	"Login"
//	@Success	200		{object}	string
//	@Router		/auth/login [post]
func (h *Auth) Login(c *gin.Context) {
	var req request.Login
	if !api.ParseRequest(c, &req) {
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.ToDomain())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, token)
}
