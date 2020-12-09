package database

import (
	"database/sql"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

const nilVersion int = -1

type migration struct {
	sql string
}

type migrator struct {
	db              *sql.DB
	migrationsTable string

	migrations []migration
}

func Migrate(db *sql.DB) error {
	m := migrator{
		db:              db,
		migrationsTable: "migrations",

		// only append!
		migrations: []migration{
			{
				// language=sql
				sql: `CREATE TABLE auth (player_id int64 NOT NULL, token VARCHAR(32) NOT NULL)`,
			},
		},
	}

	return m.migrate()
}

func (m *migrator) migrate() error {
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	curVersion := m.version()

	if err := m.up(curVersion); err != nil {
		return fmt.Errorf("up, current version '%d': %w", curVersion, err)
	}

	return nil
}

func (m *migrator) createMigrationsTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (version uint64);
  		CREATE UNIQUE INDEX IF NOT EXISTS version_unique ON %s (version);
	`, m.migrationsTable, m.migrationsTable)

	if _, err := m.db.Exec(query); err != nil {
		return err
	}

	return nil
}

func (m *migrator) version() (version int) {
	query := fmt.Sprintf("SELECT version FROM %s LIMIT 1", m.migrationsTable)

	if err := m.db.QueryRow(query).Scan(&version); err != nil {
		return nilVersion
	}

	return
}

func (m *migrator) up(curVersion int) error {
	for version, migration := range m.migrations {
		if version <= curVersion {
			continue
		}

		if err := m.executeQuery(migration.sql); err != nil {
			return fmt.Errorf("execute error, version '%d', query: '%s': %w",
				version,
				migration.sql,
				err)
		}

		if err := m.setVersion(version); err != nil {
			return fmt.Errorf("set version '%d': %w",
				version,
				err)
		}
	}

	return nil
}

func (m *migrator) executeQuery(query string) error {
	tx, err := m.db.Begin()

	if err != nil {
		return fmt.Errorf("transaction start failed: %w", err)
	}

	if _, err := tx.Exec(query); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = multierror.Append(err, errRollback)
		}

		return fmt.Errorf("execute error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (m *migrator) setVersion(version int) error {
	tx, err := m.db.Begin()

	if err != nil {
		return fmt.Errorf("transaction start failed: %w", err)
	}

	// delete old...
	queryDelete := fmt.Sprintf("DELETE FROM %s", m.migrationsTable)

	if _, err := tx.Exec(queryDelete); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = multierror.Append(err, errRollback)
		}

		return fmt.Errorf("execute error, query '%s': %w", queryDelete, err)
	}

	// insert new...
	queryInsert := fmt.Sprintf(`INSERT INTO %s (version) VALUES (?)`, m.migrationsTable)

	if _, err := tx.Exec(queryInsert, version); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = multierror.Append(err, errRollback)
		}

		return fmt.Errorf("execute error, query '%s': %w", queryDelete, err)
	}

	// commit...
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
