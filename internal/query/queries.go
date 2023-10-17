package query

type Queries struct {
	companyRM CompanyReadModel
}

func NewQueries(companyRM CompanyReadModel) *Queries {
	return &Queries{
		companyRM: companyRM,
	}
}
