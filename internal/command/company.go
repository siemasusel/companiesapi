package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/siemasusel/companiesapi/internal/domain"
	"github.com/siemasusel/companiesapi/internal/event"
)

func (c *Commands) CreateCompany(ctx context.Context, name, description string, employeesCount int, registered bool, companyType domain.CompanyType) (uuid.UUID, error) {
	id := uuid.New()

	company, err := domain.NewCompany(id, name, description, employeesCount, registered, companyType)
	if err != nil {
		return uuid.UUID{}, err
	}

	if err = c.companyRepo.CreateCompany(ctx, company); err != nil {
		return uuid.UUID{}, err
	}

	err = c.eventPublisher.PushEvents(
		ctx,
		event.NewCompanyCreatedEvent(company),
	)

	return id, err
}

func (c *Commands) UpdateCompany(ctx context.Context, id uuid.UUID, name, description *string, employeesCount *int, registered *bool, companyType *domain.CompanyType) error {
	err := c.companyRepo.UpdateCompany(
		ctx,
		id,
		func(c domain.Company) (domain.Company, error) {
			if err := c.UpdateCompany(name, description, employeesCount, registered, companyType); err != nil {
				return domain.Company{}, err
			}

			return c, nil
		})
	if err != nil {
		return err
	}

	return c.eventPublisher.PushEvents(
		ctx,
		event.NewCompanyUpdatedEvent(id, name, description, employeesCount, registered, companyType),
	)
}

func (c *Commands) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	if err := c.companyRepo.DeleteCompany(ctx, id); err != nil {
		return err
	}

	return c.eventPublisher.PushEvents(
		ctx,
		event.NewCompanyDeletedEvent(id),
	)
}
