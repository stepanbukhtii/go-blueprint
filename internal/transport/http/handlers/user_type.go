package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/http/handlers/request"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/http/handlers/response"
)

type UserType struct {
	service domain.UserTypeService
}

func NewUserType(service domain.UserTypeService) *UserType {
	return &UserType{
		service: service,
	}
}

// List godoc
//
//	@Summary	List of user types
//	@Tags		user-types
//	@Success	200	{object}	[]response.UserType
//	@Router		/user-types [get]
func (h *UserType) List(c *gin.Context) {
	userTypes, err := h.service.List(c.Request.Context())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewUserTypes(userTypes))
}

// Create godoc
//
//	@Summary	Create user type
//	@Tags		user-types
//	@Param		request	body		request.CreateUserType	true	"Create user type"
//	@Success	200		{object}	response.UserType
//	@Router		/user-types [post]
func (h *UserType) Create(c *gin.Context) {
	var req request.CreateUserType
	if !api.ParseRequest(c, &req) {
		return
	}

	userType, err := h.service.Create(c.Request.Context(), req.ToDomain())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewUserType(userType))
}

// Get godoc
//
//	@Summary	Get user
//	@Tags		user-types
//	@Param		user_type_code	path		string	true	"User type code"
//	@Success	200				{object}	response.UserType
//	@Router		/user-types/{user_type_code} [get]
func (h *UserType) Get(c *gin.Context) {
	var req request.UserTypeCodeURI
	if !api.ParseRequest(c, &req) {
		return
	}

	userType, err := h.service.Get(c.Request.Context(), req.UserTypeCode)
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewUserType(userType))
}

// Update godoc
//
//	@Summary	Update user
//	@Tags		user-types
//	@Param		user_type_code	path		string					true	"User type code"
//	@Param		request			body		request.UpdateUserType	true	"Update user type"
//	@Success	200				{object}	response.UserType
//	@Router		/user-types/{user_type_code} [patch]
func (h *UserType) Update(c *gin.Context) {
	var req request.UpdateUserType
	if !api.ParseRequest(c, &req) {
		return
	}

	userType, err := h.service.Update(c.Request.Context(), req.ToDomain())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewUserType(userType))
}

// Delete godoc
//
//	@Summary	Delete user type
//	@Tags		user-types
//	@Param		request	body		request.UpdateUser	true	"Delete user type"
//	@Success	200		{object}	api.Response
//	@Router		/user-types/{user_type_code} [delete]
func (h *UserType) Delete(c *gin.Context) {
	var req request.UserTypeCodeURI
	if !api.ParseRequest(c, &req) {
		return
	}

	if err := h.service.Delete(c.Request.Context(), req.UserTypeCode); err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondOK(c)
}
