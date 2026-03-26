package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/http/handlers/request"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/http/handlers/response"
)

type User struct {
	service domain.UserService
}

func NewUser(service domain.UserService) *User {
	return &User{
		service: service,
	}
}

// List godoc
//
//	@Summary	List of users
//	@Tags		users
//	@Param		request	query		api.Pagination	true	"List of users"
//	@Success	200		{object}	[]response.User
//	@Router		/users [get]
func (h *User) List(c *gin.Context) {
	var req api.Pagination
	if !api.ParseRequest(c, &req) {
		return
	}

	users, total, err := h.service.ListPaginate(c.Request.Context(), req)
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondDataPages(c, response.NewUsers(users), req.Pages(total))
}

// Create godoc
//
//	@Summary	Create user
//	@Tags		users
//	@Param		request	body		request.CreateUser	true	"Create user"
//	@Success	200		{object}	response.User
//	@Router		/users [post]
func (h *User) Create(c *gin.Context) {
	var req request.CreateUser
	if !api.ParseRequest(c, &req) {
		return
	}

	user, err := h.service.Create(c.Request.Context(), req.ToDomain())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewUser(user))
}

// Get godoc
//
//	@Summary	Get user
//	@Tags		users
//	@Param		user_id	path		string	true	"User ID"
//	@Success	200		{object}	response.User
//	@Router		/users/{user_id} [get]
func (h *User) Get(c *gin.Context) {
	var req request.UserURI
	if !api.ParseRequest(c, &req) {
		return
	}

	user, err := h.service.Get(c.Request.Context(), req.UserID)
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewUser(user))
}

// Update godoc
//
//	@Summary	Update user
//	@Tags		users
//	@Param		user_id	path		string				true	"User ID"
//	@Param		request	body		request.UpdateUser	true	"Update user"
//	@Success	200		{object}	response.User
//	@Router		/users/{user_id} [patch]
func (h *User) Update(c *gin.Context) {
	var req request.UpdateUser
	if !api.ParseRequest(c, &req) {
		return
	}

	user, err := h.service.Update(c.Request.Context(), req.ToDomain())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewUser(user))
}

// Delete godoc
//
//	@Summary	Delete user
//	@Tags		users
//	@Param		request	body		request.UpdateUser	true	"Delete user"
//	@Success	200		{object}	api.Response
//	@Router		/users/{user_id} [delete]
func (h *User) Delete(c *gin.Context) {
	var req request.UserURI
	if !api.ParseRequest(c, &req) {
		return
	}

	if err := h.service.Delete(c.Request.Context(), req.UserID); err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondOK(c)
}
