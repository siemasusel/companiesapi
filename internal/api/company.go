package api

import (
	"github.com/pkg/errors"
	"github.com/siemasusel/companiesapi/internal/domain"
)

func MarshalCompanyType(ct domain.CompanyType) (CompanyType, error) {
	switch ct {
	case domain.Corporations:
		return Corporations, nil
	case domain.NonProfit:
		return NonProfit, nil
	case domain.Cooperative:
		return Cooperative, nil
	case domain.SoleProprietorship:
		return SoleProprietorship, nil
	}

	return "", errors.Errorf("unable to marshal company type %d", ct)
}

func UnmarshalCompanyType(ct CompanyType) (domain.CompanyType, error) {
	switch ct {
	case Corporations:
		return domain.Corporations, nil
	case NonProfit:
		return domain.NonProfit, nil
	case Cooperative:
		return domain.Cooperative, nil
	case SoleProprietorship:
		return domain.SoleProprietorship, nil
	}

	return 0, errors.Errorf("uanble to unmarshal company type '%s'", ct)
}
