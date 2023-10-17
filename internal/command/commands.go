package command

import (
	"github.com/google/uuid"
	"github.com/siemasusel/companiesapi/internal/domain"
	"github.com/siemasusel/companiesapi/internal/event"
)

type Commands struct {
	companyRepo    domain.CompanyRepository
	eventPublisher event.Publisher
}

type UUIDGenerator interface {
	New() uuid.UUID
}

func NewCommands(companyRepo domain.CompanyRepository, eventPublisher event.Publisher) *Commands {
	return &Commands{
		companyRepo:    companyRepo,
		eventPublisher: eventPublisher,
	}
}
