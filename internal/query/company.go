package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/siemasusel/companiesapi/internal/api"
)

type CompanyReadModel interface {
	ReadCompany(ctx context.Context, id uuid.UUID) (api.Company, error)
}

func (q *Queries) ReadCompany(ctx context.Context, id uuid.UUID) (api.Company, error) {
	return q.companyRM.ReadCompany(ctx, id)
}
