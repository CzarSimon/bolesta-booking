package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/httputil/dbutil"
	"github.com/CzarSimon/httputil/timeutil"
)

type UserRepository interface {
	Save(ctx context.Context, user models.User) error
	Find(ctx context.Context, id string) (models.User, bool, error)
	FindAll(ctx context.Context) ([]models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, bool, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

type userRepo struct {
	db *sql.DB
}

func (r *userRepo) Save(ctx context.Context, user models.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction%w", err)
	}

	existing, found, err := findUser(ctx, tx, user.ID)
	if err != nil {
		dbutil.Rollback(tx)
		return err
	}

	if found {
		err = updateUser(ctx, tx, user, existing)
	} else {
		err = saveNewUser(ctx, tx, user)
	}

	if err != nil {
		dbutil.Rollback(tx)
		return err
	}

	err = tx.Commit()
	if err != nil {
		dbutil.Rollback(tx)
		return fmt.Errorf("failed to commit transaction when saving user(id=%s): %w", user.ID, err)
	}

	return nil
}

const saveNewUserQuery = `
	INSERT INTO user_account(id, name, email, password, salt, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)
`

func saveNewUser(ctx context.Context, tx *sql.Tx, user models.User) error {
	_, err := tx.ExecContext(ctx, saveNewUserQuery, user.ID, user.Name, user.Email, user.Password, user.Salt, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert %s into user_account: %w", user, err)
	}

	return nil
}

const updateUserQuery = `
	UPDATE
		user_account
	SET
		email = ?,
		password = ?,
		salt = ?,
		updated_at = ?
	WHERE
		id = ?
`

func updateUser(ctx context.Context, tx *sql.Tx, update, existing models.User) error {
	_, err := tx.ExecContext(ctx, updateUserQuery, update.Email, update.Password, update.Salt, timeutil.Now(), existing.ID)
	if err != nil {
		return fmt.Errorf("failed to update user_account. Existing=%s, Update=%s: %w", existing, update, err)
	}

	return nil
}

const findUserByIDQuery = `
	SELECT
		id,
		name,
		email,
		password,
		salt,
		created_at,
		updated_at
	FROM
		user_account
	WHERE
		id = ?
`

func (r *userRepo) Find(ctx context.Context, id string) (models.User, bool, error) {
	tx, err := readOnlyTx(ctx, r.db)
	if err != nil {
		return models.User{}, false, err
	}
	defer dbutil.Rollback(tx)

	return findUser(ctx, tx, id)
}

func findUser(ctx context.Context, tx *sql.Tx, id string) (models.User, bool, error) {
	var u models.User
	err := tx.QueryRowContext(ctx, findUserByIDQuery, id).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Salt, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return models.User{}, false, nil
	}

	if err != nil {
		return models.User{}, false, fmt.Errorf("failed to query User(id=%s): %w", id, err)
	}

	return u, true, nil
}

const findAllUsersQuery = `
	SELECT
		id,
		name,
		email,
		password,
		salt,
		created_at,
		updated_at
	FROM
		user_account
`

func (r *userRepo) FindAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, findAllUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query all users: %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	var u models.User
	for rows.Next() {
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Password,
			&u.Salt,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, u)
	}

	return users, nil
}

const findUserByEmailQuery = `
	SELECT
		id,
		name,
		email,
		password,
		salt,
		created_at,
		updated_at
	FROM
		user_account
	WHERE
		email = ?
`

func (r *userRepo) FindByEmail(ctx context.Context, email string) (models.User, bool, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, findUserByEmailQuery, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Salt, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return models.User{}, false, nil
	}

	if err != nil {
		return models.User{}, false, fmt.Errorf("failed to query be email: %w", err)
	}

	return u, true, nil
}

const findUsersByIDsQuery = `
	SELECT
		id,
		name,
		email,
		password,
		salt,
		created_at,
		updated_at
	FROM
		user_account
	WHERE 
		id IN`

func findUsersByIDs(ctx context.Context, tx *sql.Tx, ids []string) (map[string]models.User, error) {
	rows, err := queryByIDs(ctx, tx, findUsersByIDsQuery, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(map[string]models.User)
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Salt, &u.CreatedAt, &u.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users[u.ID] = u
	}

	return users, nil
}
