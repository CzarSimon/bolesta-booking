package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/httputil/dbutil"
)

type CabinRepository interface {
	Find(ctx context.Context, id string) (models.Cabin, bool, error)
	FindAll(ctx context.Context) ([]models.Cabin, error)
}

func NewCabinRepository(db *sql.DB) CabinRepository {
	return &cabinRepo{
		db: db,
	}
}

type cabinRepo struct {
	db *sql.DB
}

const findCabinByIDQuery = `
	SELECT 
		id, 
		name, 
		created_at, 
		updated_at 
	FROM 
		cabin 
	WHERE 
		id = ?`

func (r *cabinRepo) Find(ctx context.Context, id string) (models.Cabin, bool, error) {
	tx, err := readOnlyTx(ctx, r.db)
	if err != nil {
		return models.Cabin{}, false, err
	}
	defer dbutil.Rollback(tx)

	return findCabin(ctx, tx, id)
}

func findCabin(ctx context.Context, tx *sql.Tx, id string) (models.Cabin, bool, error) {
	var c models.Cabin
	err := tx.QueryRowContext(ctx, findCabinByIDQuery, id).Scan(
		&c.ID,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Cabin{}, false, nil
	}
	if err != nil {
		return models.Cabin{}, false, fmt.Errorf("failed to query for Cabin(id=%s). Error: %w", id, err)
	}

	return c, true, nil
}

const findAllCabinsQuery = `
	SELECT 
		id, 
		name, 
		created_at, 
		updated_at 
	FROM 
		cabin`

func (r *cabinRepo) FindAll(ctx context.Context) ([]models.Cabin, error) {
	rows, err := r.db.QueryContext(ctx, findAllCabinsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query all cabins: %w", err)
	}
	defer rows.Close()

	cabins := make([]models.Cabin, 0)
	var c models.Cabin
	for rows.Next() {
		err = rows.Scan(
			&c.ID,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan cabin: %w", err)
		}

		cabins = append(cabins, c)
	}

	return cabins, nil
}

const findCabinsByIDsQuery = `
	SELECT 
		id, 
		name, 
		created_at, 
		updated_at 
	FROM 
		cabin
	WHERE 
		id IN`

func findCabinsByIDs(ctx context.Context, tx *sql.Tx, ids []string) (map[string]models.Cabin, error) {
	rows, err := queryByIDs(ctx, tx, findCabinsByIDsQuery, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cabins := make(map[string]models.Cabin)
	for rows.Next() {
		var c models.Cabin
		err = rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan cabin: %w", err)
		}

		cabins[c.ID] = c
	}

	return cabins, nil
}
