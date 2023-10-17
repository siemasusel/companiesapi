package mysql

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/siemasusel/companiesapi/internal/api"
	"github.com/siemasusel/companiesapi/internal/domain"
)

type companyModel struct {
	ID             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	Description    *string   `db:"description"`
	EmployeesCount int       `db:"employees_count"`
	Registered     bool      `db:"registered"`
	Type           string    `db:"type"`
}

func marshalCompany(company domain.Company) (companyModel, error) {
	dbCompany := companyModel{
		ID:             company.ID,
		Name:           company.Name,
		EmployeesCount: company.EmployeesCount,
		Registered:     company.Registered,
	}
	if company.Description != "" {
		dbCompany.Description = &company.Description
	}

	ct, err := api.MarshalCompanyType(company.Type)
	if err != nil {
		return companyModel{}, err
	}
	dbCompany.Type = string(ct)

	return dbCompany, nil
}

func unmarshalCompany(dbCompany companyModel) (domain.Company, error) {
	company := domain.Company{
		ID:             dbCompany.ID,
		Name:           dbCompany.Name,
		EmployeesCount: dbCompany.EmployeesCount,
		Registered:     dbCompany.Registered,
	}

	if dbCompany.Description != nil {
		company.Description = *dbCompany.Description
	}

	ct, err := api.UnmarshalCompanyType(api.CompanyType(dbCompany.Type))
	if err != nil {
		return domain.Company{}, err
	}
	company.Type = ct

	return company, nil
}

func companyToApi(dbCompany companyModel) (api.Company, error) {
	return api.Company{
		Id:             dbCompany.ID.String(),
		Name:           dbCompany.Name,
		Description:    dbCompany.Description,
		EmployeesCount: dbCompany.EmployeesCount,
		Registered:     dbCompany.Registered,
		Type:           api.CompanyType(dbCompany.Type),
	}, nil
}

func parseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, errors.Wrapf(err, "unable to parse id '%s'", s)
	}

	return id, nil
}
