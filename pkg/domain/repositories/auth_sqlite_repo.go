package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/qkveri/player_core/pkg/domain"
)

type authSqliteRepo struct {
	db *sql.DB
}

func NewAuthSqliteRepo(db *sql.DB) *authSqliteRepo {
	return &authSqliteRepo{
		db: db,
	}
}

func (a *authSqliteRepo) Set(ctx context.Context, auth *domain.Auth) error {
	tx, err := a.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("transaction start failed: %w", err)
	}

	// delete old...
	// language=sql
	queryDelete := "DELETE FROM auth"

	if _, err := tx.ExecContext(ctx, queryDelete); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = multierror.Append(err, errRollback)
		}

		return fmt.Errorf("execute error, query '%s': %w",
			queryDelete, err)
	}

	// insert new...
	// language=sql
	queryInsert := "INSERT INTO auth (player_id, token) VALUES (?, ?)"

	if _, err := tx.ExecContext(ctx, queryInsert,
		auth.PlayerID,
		auth.Token,
	); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = multierror.Append(err, errRollback)
		}

		return fmt.Errorf("execute error: %w, query '%s', playerID: %d, token: %s",
			err, queryDelete, auth.PlayerID, auth.Token)
	}

	// commit...
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (a *authSqliteRepo) Get(ctx context.Context) (*domain.Auth, error) {
	// language=sql
	query := "SELECT player_id, token FROM auth LIMIT 1"

	stmt, err := a.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("db.PrepareContext: %w", err)
	}

	defer func() {
		_ = stmt.Close()
	}()

	auth := &domain.Auth{}

	if err := stmt.QueryRowContext(ctx).Scan(
		&auth.PlayerID,
		&auth.Token,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return auth, nil
}
