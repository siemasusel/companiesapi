package command

import (
	"github.com/google/uuid"
	"github.com/siemasusel/companiesapi/internal/domain"
)

type Commands struct {
	companyRepo domain.CompanyRepository
}

type UUIDGenerator interface {
	New() uuid.UUID
}

func NewCommands(companyRepo domain.CompanyRepository) *Commands {
	return &Commands{
		companyRepo: companyRepo,
	}
}
