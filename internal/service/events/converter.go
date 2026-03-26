package events

import "github.com/stepanbukhtii/go-blueprint/internal/domain"

func NewEventUserCreatedData(user domain.User) EventUserCreatedData {
	return EventUserCreatedData{
		UserID: user.ID,
	}
}

func NewEventUserUpdatedData(user domain.User) EventUserUpdatedData {
	return EventUserUpdatedData{
		UserID: user.ID,
		Rate:   user.Rate,
	}
}

func NewEventCompanyUpdatedData(company domain.Company) EventCompanyUpdatedData {
	return EventCompanyUpdatedData{
		CompanyID: company.ID,
		Name:      company.Name,
		OwnerID:   company.OwnerID,
		ManagerID: company.ManagerID,
		IsActive:  company.IsActive,
		LogoURL:   company.LogoURL,
	}
}
