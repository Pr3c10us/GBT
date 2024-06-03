package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/logger"
	"go.uber.org/zap"
)

type DebtorRepositoryPG struct {
	db     *sql.DB
	logger logger.Logger
}

func NewDebtorRepositoryPG(db *sql.DB, logger logger.Logger) debtor.Repository {
	return &DebtorRepositoryPG{db: db, logger: logger}
}

func (repository *DebtorRepositoryPG) AddDebtor(debtor *debtor.Debtor) error {
	statement, err := repository.db.Prepare(`INSERT INTO debtors ( name, phone_number, user_id) VALUES ($1,$2,$3)`)
	if err != nil {
		repository.logger.LogWithFields("error", "Failed to prepare sql statement", zap.Error(err))
		return err
	}

	_, err = statement.Exec(debtor.Name, debtor.PhoneNumber, debtor.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *DebtorRepositoryPG) RemoveDebtor(id string) error {
	statement, err := repository.db.Prepare(`DELETE FROM debtors WHERE id = $1`)
	if err != nil {
		repository.logger.LogWithFields("error", "Failed to prepare sql statement", zap.Error(err))
		return err
	}

	var queryResult sql.Result
	queryResult, err = statement.Exec(id)
	if err != nil {
		return err
	}
	rowsAffected, err := queryResult.RowsAffected()
	if err != nil {
		repository.logger.LogWithFields("error", "Failed to get affected rows", zap.Error(err))
		return err
	}
	if rowsAffected == 0 {
		return appError.BadRequest(errors.New("no debtor deleted: the specified debtor id does not exist"))
	}

	return nil
}

func (repository *DebtorRepositoryPG) GetDebtorByID(id string) (*debtor.Debtor, error) {
	query := `SELECT id, name, phone_number, user_id FROM debtors WHERE id=$1;`
	row := repository.db.QueryRow(query, id)
	var d debtor.Debtor
	switch err := row.Scan(&d.ID, &d.Name, &d.PhoneNumber, &d.UserID); {
	case errors.Is(err, sql.ErrNoRows):
		var notFoundErr = errors.New(fmt.Sprintf("debtor with id '%v' does not exit", id))
		return &debtor.Debtor{}, appError.NotFound(notFoundErr)
	case err == nil:
		return &d, nil
	default:
		return &debtor.Debtor{}, err
	}
}
