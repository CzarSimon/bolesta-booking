package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/dbutil"
)

type BookingRepository interface {
	Save(ctx context.Context, booking models.Booking) error
	Find(ctx context.Context, id string) (models.Booking, bool, error)
	FindByFilter(ctx context.Context, f models.BookingFilter) ([]models.Booking, error)
	FindRefsByFilter(ctx context.Context, f models.BookingFilter) ([]models.BookingRef, error)
	Delete(ctx context.Context, id string) error
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepo{
		db:       db,
		userRepo: NewUserRepository(db),
	}
}

type bookingRepo struct {
	db       *sql.DB
	userRepo UserRepository
}

const saveBookingQuery = `
	INSERT INTO booking(id, start_date, end_date, created_at, updated_at, cabin_id, user_id) 
	SELECT ?, ?, ?, ?, ?, ?, ? 
	WHERE NOT EXISTS (
		SELECT 
			1 
		FROM 
			booking 
		WHERE 
			cabin_id = ? 
			AND (
				? BETWEEN start_date AND end_date
				OR ? BETWEEN start_date AND end_date
			)
		LIMIT 1
	);
`

func (r *bookingRepo) Save(ctx context.Context, booking models.Booking) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction%w", err)
	}

	res, err := tx.ExecContext(
		ctx,
		saveBookingQuery,
		booking.ID,
		booking.StartDate,
		booking.EndDate,
		booking.CreatedAt,
		booking.UpdatedAt,
		booking.Cabin.ID,
		booking.User.ID,
		booking.Cabin.ID,  // Match cabin_id:s of existing bookings
		booking.StartDate, // Start date must be after existing booking
		booking.EndDate,   // End date must be before existing booking
	)
	if err != nil {
		dbutil.Rollback(tx)
		return fmt.Errorf("failed to insert %s into booking: %w", booking, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		dbutil.Rollback(tx)
		return fmt.Errorf("failed to get rows affected for booking %s: %w", booking, err)
	}

	if rowsAffected < 1 {
		dbutil.Rollback(tx)
		return httputil.Conflictf("failed to insert %s due to conflicting existing boooking", booking)
	}

	err = tx.Commit()
	if err != nil {
		dbutil.Rollback(tx)
		return fmt.Errorf("failed to commit transaction when inserting %s: %w", booking, err)
	}

	return nil
}

const deleteBookingByIDQuery = "DELETE FROM booking WHERE id = ?"

func (r *bookingRepo) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction%w", err)
	}

	res, err := tx.ExecContext(
		ctx,
		deleteBookingByIDQuery,
		id,
	)
	if err != nil {
		dbutil.Rollback(tx)
		return fmt.Errorf("failed to delete booking(id=%s): %w", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		dbutil.Rollback(tx)
		return fmt.Errorf("failed to get rows affected for booking(id=%s) deletion: %w", id, err)
	}

	if rowsAffected != 1 {
		dbutil.Rollback(tx)
		return fmt.Errorf("failed to delete booking(id=%s) unexpected number of rows affected. Expected 1 got %d", id, rowsAffected)
	}

	err = tx.Commit()
	if err != nil {
		dbutil.Rollback(tx)
		return fmt.Errorf("failed to commit transaction when deleting booking(id=%s): %w", id, err)
	}

	return nil
}

func (r *bookingRepo) Find(ctx context.Context, id string) (models.Booking, bool, error) {
	tx, err := readOnlyTx(ctx, r.db)
	if err != nil {
		return models.Booking{}, false, err
	}
	defer dbutil.Rollback(tx)

	ref, exits, err := findBookingRef(ctx, tx, id)
	if !exits || err != nil {
		return models.Booking{}, exits, err
	}

	cabin, exits, err := findCabin(ctx, tx, ref.CabinID)
	if err != nil {
		return models.Booking{}, exits, err
	}
	if !exits {
		return models.Booking{}, exits, fmt.Errorf("failed to find referenced Cabin(id=%s) from Booking(id=%s)", ref.CabinID, ref.ID)
	}

	user, exits, err := findUser(ctx, tx, ref.UserID)
	if err != nil {
		return models.Booking{}, exits, err
	}
	if !exits {
		return models.Booking{}, exits, fmt.Errorf("failed to find referenced User(id=%s) from Booking(id=%s)", ref.UserID, ref.ID)
	}

	return ref.Booking(cabin, user), true, nil
}

const findBookingRefByIDQuery = `
	SELECT 
		id, 
		start_date, 
		end_date, 
		created_at, 
		updated_at, 
		cabin_id, 
		user_id
	FROM 
		booking
	WHERE
		id = ?`

func findBookingRef(ctx context.Context, tx *sql.Tx, id string) (models.BookingRef, bool, error) {
	var b models.BookingRef
	err := tx.QueryRowContext(ctx, findBookingRefByIDQuery, id).Scan(&b.ID, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt, &b.CabinID, &b.UserID)
	if err == sql.ErrNoRows {
		return models.BookingRef{}, false, nil
	}

	if err != nil {
		return models.BookingRef{}, false, fmt.Errorf("failed to query Booking(id=%s): %w", id, err)
	}

	return b, true, nil
}

func (r *bookingRepo) FindByFilter(ctx context.Context, f models.BookingFilter) ([]models.Booking, error) {
	tx, err := readOnlyTx(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer dbutil.Rollback(tx)

	refs, err := findBookingRefsByFilter(ctx, tx, f)
	if err != nil {
		return nil, err
	}

	cabins, err := findCabinsByIDs(ctx, tx, mapRefsToUniqueCabinIDs(refs))
	if err != nil {
		return nil, err
	}

	users, err := findUsersByIDs(ctx, tx, mapRefsToUniqueUserIDs(refs))
	if err != nil {
		return nil, err
	}

	return mapRefsToBookings(refs, cabins, users)
}

func (r *bookingRepo) FindRefsByFilter(ctx context.Context, f models.BookingFilter) ([]models.BookingRef, error) {
	tx, err := readOnlyTx(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer dbutil.Rollback(tx)

	refs, err := findBookingRefsByFilter(ctx, tx, f)
	if err != nil {
		return nil, err
	}

	return refs, nil
}

func findBookingRefsByFilter(ctx context.Context, tx *sql.Tx, f models.BookingFilter) ([]models.BookingRef, error) {
	query, values := createBookingRefsFilterQuery(f)
	rows, err := tx.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings by %s", f)
	}
	defer rows.Close()

	refs := make([]models.BookingRef, 0)
	var b models.BookingRef
	for rows.Next() {
		err = rows.Scan(&b.ID, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt, &b.CabinID, &b.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}

		refs = append(refs, b)
	}

	return refs, nil
}

func createBookingRefsFilterQuery(f models.BookingFilter) (string, []interface{}) {
	values := make([]interface{}, 0)
	clauses := make([]string, 0)
	if f.CabinID != "" {
		values = append(values, f.CabinID)
		clauses = append(clauses, "cabin_id = ?")
	}

	if f.UserID != "" {
		values = append(values, f.UserID)
		clauses = append(clauses, "user_id = ?")
	}

	baseQuery := `
		SELECT 
			id, 
			start_date, 
			end_date, 
			created_at, 
			updated_at, 
			cabin_id, 
			user_id
		FROM 
			booking`

	if len(values) == 0 {
		return baseQuery, values
	}

	whereClause := strings.Join(clauses, " AND ")
	return fmt.Sprintf("%s WHERE %s", baseQuery, whereClause), values
}

func mapRefsToBookings(refs []models.BookingRef, cabins map[string]models.Cabin, users map[string]models.User) ([]models.Booking, error) {
	bookings := make([]models.Booking, 0, len(refs))
	for _, ref := range refs {
		cabin, ok := cabins[ref.CabinID]
		if !ok {
			return nil, fmt.Errorf("Could not find Cabin(id=%s) referenced in Booking(id=%s)", ref.CabinID, ref.ID)
		}

		user, ok := users[ref.UserID]
		if !ok {
			return nil, fmt.Errorf("Could not find User(id=%s) referenced in Booking(id=%s)", ref.UserID, ref.ID)
		}

		bookings = append(bookings, ref.Booking(cabin, user))
	}

	return bookings, nil
}

func mapRefsToUniqueCabinIDs(refs []models.BookingRef) []string {
	ids := make([]string, 0)
	unique := make(map[string]bool)

	for _, ref := range refs {
		id := ref.CabinID
		_, present := unique[id]
		if !present {
			unique[id] = true
			ids = append(ids, id)
		}
	}

	return ids
}

func mapRefsToUniqueUserIDs(refs []models.BookingRef) []string {
	ids := make([]string, 0)
	unique := make(map[string]bool)

	for _, ref := range refs {
		id := ref.UserID
		_, present := unique[id]
		if !present {
			unique[id] = true
			ids = append(ids, id)
		}
	}

	return ids
}
