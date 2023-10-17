package event

import (
	"github.com/google/uuid"
	"github.com/siemasusel/companiesapi/internal/api"
	"github.com/siemasusel/companiesapi/internal/domain"
)

type CompanyCreatedPayload struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description,omitempty"`
	EmployeesCount int             `json:"employees_count"`
	Registered     bool            `json:"registered"`
	Type           api.CompanyType `json:"type"`
}

func NewCompanyCreatedEvent(company domain.Company) Event {
	payload := CompanyCreatedPayload{
		ID:             company.ID,
		Name:           company.Name,
		Description:    company.Description,
		EmployeesCount: company.EmployeesCount,
		Registered:     company.Registered,
	}

	payload.Type, _ = api.MarshalCompanyType(company.Type)

	return Event{
		AggregateID:   company.ID,
		AggregateName: AggregateNameCompany,
		Type:          EventTypeCompanyCreated,
		Payload:       payload,
	}
}

type CompanyUpdatedPayload struct {
	Name           *string          `json:"name"`
	Description    *string          `json:"description,omitempty"`
	EmployeesCount *int             `json:"employees_count"`
	Registered     *bool            `json:"registered"`
	Type           *api.CompanyType `json:"type"`
}

func NewCompanyUpdatedEvent(id uuid.UUID, name, description *string, employeesCount *int, registered *bool, companyType *domain.CompanyType) Event {
	payload := CompanyUpdatedPayload{
		Name:           name,
		Description:    description,
		EmployeesCount: employeesCount,
		Registered:     registered,
	}

	if companyType != nil {
		ct, _ := api.MarshalCompanyType(*companyType)
		payload.Type = &ct
	}

	return Event{
		AggregateID:   id,
		AggregateName: AggregateNameCompany,
		Type:          EventTypeCompanyUpdated,
		Payload:       payload,
	}
}

type CompanyDeletedPayload struct {
}

func NewCompanyDeletedEvent(id uuid.UUID) Event {
	return Event{
		AggregateID:   id,
		AggregateName: AggregateNameCompany,
		Type:          EventTypeCompanyUpdated,
		Payload:       CompanyDeletedPayload{},
	}
}
