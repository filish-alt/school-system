package db

import (
	"context"
	"database/sql"
	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)"
	return sql.Open("sqlite", dsn)
}

func ExecBatch(ctx context.Context, db *sql.DB, script string) error {
	_, err := db.ExecContext(ctx, script)
	return err
}

