package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/siemasusel/companiesapi/internal/domain"
)

func (c *Commands) CreateCompany(ctx context.Context, name, description string, employeesCount int, registered bool, companyType domain.CompanyType) (uuid.UUID, error) {
	id := uuid.New()

	company, err := domain.NewCompany(id, name, description, employeesCount, registered, companyType)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = c.companyRepo.CreateCompany(ctx, company)

	return id, err
}

func (c *Commands) UpdateCompany(ctx context.Context, id uuid.UUID, name, description *string, employeesCount *int, registered *bool, companyType *domain.CompanyType) error {
	return c.companyRepo.UpdateCompany(
		ctx,
		id,
		func(c domain.Company) (domain.Company, error) {
			if err := c.UpdateCompany(name, description, employeesCount, registered, companyType); err != nil {
				return domain.Company{}, err
			}

			return c, nil
		})
}

func (c *Commands) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	return c.companyRepo.DeleteCompany(ctx, id)
}
