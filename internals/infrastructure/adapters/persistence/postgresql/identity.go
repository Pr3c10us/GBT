package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/logger"
	"go.uber.org/zap"
	"time"
)

type IdentityRepositoryPG struct {
	db     *sql.DB
	logger logger.Logger
}

func NewIdentityRepositoryPG(db *sql.DB, logger logger.Logger) identity.Repository {
	return &IdentityRepositoryPG{db: db, logger: logger}
}

// CreateUser adds a user to database
func (repository *IdentityRepositoryPG) CreateUser(user *identity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := repository.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		repository.logger.LogWithFields("error", "Failed to start transaction", zap.Error(err))
		return err
	}

	var statement *sql.Stmt
	statement, err = tx.PrepareContext(ctx, `INSERT INTO users (username, password, security_question, security_answer) VALUES ($1,$2,$3,$4)`)
	if err != nil {
		_ = tx.Rollback()
		repository.logger.LogWithFields("error", "Failed to prepare sql statement", zap.Error(err))
		return err
	}

	_, err = statement.ExecContext(ctx, user.Username, user.Password, user.SecurityQuestion, user.SecurityAnswer)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		repository.logger.LogWithFields("error", "Failed to commit transaction", zap.Error(err))
		return err
	}
	return nil
}

func (repository *IdentityRepositoryPG) GetUser(username string) (*identity.User, error) {
	query :=
		`SELECT id,username,password,security_question,security_answer FROM users where username=$1`
	stmt, err := repository.db.Prepare(query)
	if err != nil {
		repository.logger.LogWithFields("error", "Failed to prepare sql query", zap.Error(err))
		return &identity.User{}, err
	}
	defer stmt.Close()

	var user identity.User
	switch err := stmt.QueryRow(username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.SecurityQuestion,
		&user.SecurityAnswer,
	); {
	case errors.Is(err, sql.ErrNoRows):
		var userNotFoundErr = errors.New(fmt.Sprintf("user with username '%v' does not exit", username))
		return &identity.User{}, appError.NotFound(userNotFoundErr)
	case err == nil:
		return &user, nil
	default:
		return &identity.User{}, err
	}
}
