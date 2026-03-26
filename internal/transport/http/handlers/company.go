package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/rest/api"

	"github.com/stepanbukhtii/go-blueprint/internal/domain"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/http/handlers/request"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/http/handlers/response"
)

type Company struct {
	service domain.CompanyService
}

func NewCompany(service domain.CompanyService) *Company {
	return &Company{
		service: service,
	}
}

// List godoc
//
//	@Summary	List of companies
//	@Tags		companies
//	@Security	ApiKeyAuth
//	@Param		request	query		api.Pagination	true	"List of companies"
//	@Success	200		{object}	[]response.Company
//	@Router		/companies [get]
func (h *Company) List(c *gin.Context) {
	var req api.Pagination
	if !api.ParseRequest(c, &req) {
		return
	}

	companies, total, err := h.service.ListPaginate(c.Request.Context(), req)
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondDataPages(c, response.NewCompanies(companies), req.Pages(total))
}

// Create godoc
//
//	@Summary	Create company
//	@Tags		companies
//	@Security	ApiKeyAuth
//	@Param		request	body		request.CreateCompany	true	"Create company"
//	@Success	200		{object}	response.Company
//	@Router		/companies [post]
func (h *Company) Create(c *gin.Context) {
	var req request.CreateCompany
	if !api.ParseRequest(c, &req) {
		return
	}

	company, err := h.service.Create(c.Request.Context(), req.ToDomain())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewCompany(company))
}

// CreateMultiple godoc
//
//	@Summary	Create multiple company
//	@Tags		companies
//	@Security	ApiKeyAuth
//	@Param		request	body		[]request.CreateCompany	true	"Create multiple company"
//	@Success	200		{object}	[]response.Company
//	@Router		/companies/multiple [post]
func (h *Company) CreateMultiple(c *gin.Context) {
	var req request.CreateCompanyMultiple
	if !api.ParseRequest(c, &req) {
		return
	}

	companies, err := h.service.CreateMultiple(c.Request.Context(), req.ToDomain())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewCompanies(companies))
}

// Get godoc
//
//	@Summary	Get company
//	@Tags		companies
//	@Security	ApiKeyAuth
//	@Param		company_id	path		string	true	"Company ID"
//	@Success	200			{object}	response.Company
//	@Router		/companies/{company_id} [get]
func (h *Company) Get(c *gin.Context) {
	var req request.CompanyURI
	if !api.ParseRequest(c, &req) {
		return
	}

	company, err := h.service.Get(c.Request.Context(), req.CompanyID)
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewCompany(company))
}

// GetCompanyByOwner godoc
//
//	@Summary	Get companies by owner id
//	@Tags		companies
//	@Security	ApiKeyAuth
//	@Success	200	{object}	[]response.Company
//	@Router		/companies/owner [get]
func (h *Company) GetCompanyByOwner(c *gin.Context) {
	companies, err := h.service.GetCompanyByOwner(c.Request.Context())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewCompanies(companies))
}

// Update godoc
//
//	@Summary	Update company
//	@Tags		companies
//	@Security	ApiKeyAuth
//	@Param		request	body		request.UpdateCompany	true	"Update company"
//	@Success	200		{object}	response.Company
//	@Router		/companies/{company_id} [post]
func (h *Company) Update(c *gin.Context) {
	var req request.UpdateCompany
	if !api.ParseRequest(c, &req) {
		return
	}

	company, err := h.service.Update(c.Request.Context(), req.ToDomain())
	if err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondData(c, response.NewCompany(company))
}

// Delete godoc
//
//	@Summary	Delete company
//	@Tags		companies
//	@Security	ApiKeyAuth
//	@Param		request	body		request.UpdateUser	true	"Delete company"
//	@Success	200		{object}	any
//	@Router		/companies/{company_id} [delete]
func (h *Company) Delete(c *gin.Context) {
	var req request.CompanyURI
	if !api.ParseRequest(c, &req) {
		return
	}

	if err := h.service.Delete(c.Request.Context(), req.CompanyID); err != nil {
		api.ServeError(c, err)
		return
	}

	api.RespondOK(c)
}
