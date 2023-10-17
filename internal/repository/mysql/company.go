package mysql

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/siemasusel/companiesapi/internal/api"
	"github.com/siemasusel/companiesapi/internal/domain"
	"github.com/siemasusel/companiesapi/internal/utils/apperr"
	"go.uber.org/multierr"
)

const (
	mySQLDeadlockErrorCode        = 1213
	mySQLLockWaitTimeoutErrorCode = 1205
	mySQLDuplicateEntryErrorCode  = 1062
)

type CompanyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *CompanyRepository {
	return &CompanyRepository{
		db: db,
	}
}

func (r *CompanyRepository) CreateCompany(ctx context.Context, company domain.Company) error {
	dbModel, err := marshalCompany(company)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExecContext(
		ctx,
		`INSERT INTO companies(id, name, description, employees_count, registered, type)
		VALUES(:id, :name, :description, :employees_count, :registered, :type)`,
		dbModel)

	if val, ok := errors.Cause(err).(*mysql.MySQLError); ok && (val.Number == mySQLDuplicateEntryErrorCode) {
		return apperr.NewAlreadyExists(err, "company with name '%s' already exists", company.Name)
	}

	return err
}

func (r *CompanyRepository) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM companies WHERE id = ?`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return apperr.NewNotFound(nil, "company '%s' not found", id)
	}

	return nil
}

func (r *CompanyRepository) ReadCompany(ctx context.Context, id uuid.UUID) (api.Company, error) {
	q := `SELECT * FROM companies WHERE id = ? `

	var dbCompany companyModel
	err := r.db.GetContext(ctx, &dbCompany, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return api.Company{}, apperr.NewNotFound(err, "company '%s' not found", id)
	} else if err != nil {
		return api.Company{}, errors.Wrapf(err, "unable to get company '%s' from db", id)
	}

	return companyToApi(dbCompany)
}

func (r *CompanyRepository) UpdateCompany(
	ctx context.Context,
	id uuid.UUID,
	updateFn func(domain.Company) (domain.Company, error),
) error {
	for {
		err := r.updateCompany(ctx, id, updateFn)

		if val, ok := errors.Cause(err).(*mysql.MySQLError); ok && (val.Number == mySQLLockWaitTimeoutErrorCode || val.Number == mySQLDeadlockErrorCode) {
			continue
		}
		return err
	}
}

func (r *CompanyRepository) updateCompany(
	ctx context.Context,
	id uuid.UUID,
	updateFn func(domain.Company) (domain.Company, error),
) (err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	defer func() {
		err = finishTransaction(err, tx)
	}()

	dbCompany, err := getCompanyForUpdateTx(ctx, tx, id)
	if err != nil {
		return err
	}

	company, err := unmarshalCompany(dbCompany)
	if err != nil {
		return err
	}

	company, err = updateFn(company)
	if err != nil {
		return err
	}

	return updateCompanyTx(ctx, tx, company)
}

func getCompanyForUpdateTx(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (companyModel, error) {
	q := `SELECT * FROM companies WHERE id = ? FOR UPDATE`

	var dbCompany companyModel
	err := tx.GetContext(ctx, &dbCompany, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return companyModel{}, apperr.NewNotFound(err, "company '%s' not found", id)
	} else if err != nil {
		return companyModel{}, errors.Wrapf(err, "unable to get company '%s' from db", id)
	}

	return dbCompany, nil
}

func updateCompanyTx(ctx context.Context, tx *sqlx.Tx, company domain.Company) error {
	dbCompany, err := marshalCompany(company)
	if err != nil {
		return errors.Wrap(err, "unable to marshal company")
	}

	_, err = tx.NamedExecContext(
		ctx,
		`UPDATE companies
		SET
			name = :name,
			description = :description,
			employees_count = :employees_count,
			registered = :registered,
			type = :type
		WHERE
			id = :id`,
		dbCompany,
	)
	if err != nil {
		return errors.Wrap(err, "unable to update account")
	}

	return nil
}

func finishTransaction(err error, tx *sqlx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return multierr.Combine(err, rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return errors.Wrap(err, "failed to commit tx")
		}

		return nil
	}
}
