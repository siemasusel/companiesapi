package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/siemasusel/companiesapi/internal/utils/apperr"
)

type CompanyRepository interface {
	CreateCompany(ctx context.Context, c Company) error
	UpdateCompany(ctx context.Context, id uuid.UUID, f func(Company) (Company, error)) error
	DeleteCompany(ctx context.Context, id uuid.UUID) error
}

type Company struct {
	ID             uuid.UUID
	Name           string
	Description    string
	EmployeesCount int
	Registered     bool
	Type           CompanyType
}

func (c Company) Validate() error {
	if !c.Type.IsValid() {
		return apperr.NewInvalidArgument(nil, "invalid company type")
	}

	if len(c.Name) == 0 {
		return apperr.NewInvalidArgument(nil, "name is requried")
	}
	if len(c.Name) > 15 {
		return apperr.NewInvalidArgument(nil, "name cannot be longer than 15, got %d", len(c.Name))
	}

	if c.EmployeesCount < 1 {
		return apperr.NewInvalidArgument(nil, "invalid employees count %d, must be greater than 0", c.EmployeesCount)
	}

	if len(c.Description) > 0 && len(c.Description) > 3000 {
		return apperr.NewInvalidArgument(nil, "description cannot be longer than 3000, got %d", len(c.Description))
	}

	return nil
}

type CompanyType uint8

const (
	Corporations CompanyType = iota
	NonProfit
	Cooperative
	SoleProprietorship

	companyTypeCount
)

func (c CompanyType) IsValid() bool {
	return c < companyTypeCount
}

func NewCompany(id uuid.UUID, name, description string, employeesCount int, registered bool, companyType CompanyType) (Company, error) {
	c := Company{
		ID:             id,
		Name:           name,
		Description:    description,
		EmployeesCount: employeesCount,
		Registered:     registered,
		Type:           companyType,
	}
	if err := c.Validate(); err != nil {
		return Company{}, err
	}

	return c, nil
}

func (c *Company) UpdateCompany(name, description *string, employeesCount *int, registered *bool, companyType *CompanyType) error {
	if name != nil {
		c.Name = *name
	}
	if description != nil {
		c.Description = *description
	}
	if employeesCount != nil {
		c.EmployeesCount = *employeesCount
	}
	if registered != nil {
		c.Registered = *registered
	}
	if companyType != nil {
		c.Type = *companyType
	}

	return c.Validate()
}
