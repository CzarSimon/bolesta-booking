package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func readOnlyTx(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to begin read only transaction: %w", err)
	}

	return tx, nil
}

func queryByIDs(ctx context.Context, tx *sql.Tx, baseQuery string, ids []string) (*sql.Rows, error) {
	values := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		values = append(values, id)
	}

	query := baseQuery + createInClause(len(ids))
	rows, err := tx.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to exectute query %s: %w", query, err)
	}

	return rows, nil
}

func createInClause(length int) string {
	values := make([]string, 0, length)
	for i := 0; i < length; i++ {
		values = append(values, "?")
	}

	return fmt.Sprintf(" (%s)", strings.Join(values, ","))
}
