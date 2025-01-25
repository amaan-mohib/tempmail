package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInit, downInit)
}

func upInit(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx,
		`create table if not exists aliases (
			id SERIAL primary key,
			alias varchar(255) NOT NULL,
			created_at_ts BIGINT NOT NULL,
			expiry_at BIGINT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "create unique index idx_alias_unique on aliases(alias)")
	if err != nil {
		return err
	}
	return err
}

func downInit(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "drop index if exists idx_alias_unique")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "drop table if exists aliases")
	if err != nil {
		return err
	}
	return err
}
