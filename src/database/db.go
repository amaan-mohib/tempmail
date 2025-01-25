package database

import (
	"context"
	"fmt"
	"os"
	"tempgalias/src/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/scan"
	"github.com/stephenafamo/scan/pgxscan"
)

var DB *pgxpool.Pool
var DBCtx context.Context

func SetupDatabase() {
	DBCtx = context.Background()
	dbpool, err := pgxpool.New(DBCtx, config.Config.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	DB = dbpool

	ping := DB.Ping(context.Background())
	if ping != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
}

func ScanOne[T any](mapper scan.Mapper[T], sql string, args ...any) (T, error) {
	return pgxscan.One(DBCtx, DB, mapper, sql, args...)
}

func ScanAll[T any](mapper scan.Mapper[T], sql string, args ...any) ([]T, error) {
	return pgxscan.All(DBCtx, DB, mapper, sql, args...)
}

func Query(sql string, args ...any) (pgconn.CommandTag, error) {
	return DB.Exec(DBCtx, sql, args...)
}

/*
If the data is already in a [][]any

	rows := [][]any{
		{"John", "Smith", int32(36)},
		{"Jane", "Doe", int32(29)},
	}
	copyCount, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"people"},
		[]string{"first_name", "last_name", "age"},
		pgx.CopyFromRows(rows),
	)

If the data already have a typed array using CopyFromSlice can be more convenient

	rows := []User{
		{"John", "Smith", 36},
		{"Jane", "Doe", 29},
	}
	copyCount, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"people"},
		[]string{"first_name", "last_name", "age"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
				return []any{rows[i].FirstName, rows[i].LastName, rows[i].Age}, nil
		}),
	)
*/
func InsertMultiple(tableName string, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return DB.CopyFrom(DBCtx, pgx.Identifier{tableName}, columnNames, rowSrc)
}
