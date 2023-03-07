package models

import (
	"fmt"
	"time"

	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/timeutil"
)

// BookingRuleViloation type of vialotion
type BookingRuleViloation int

const (
	bookingOk BookingRuleViloation = iota
	ErrInvalidBooking
	ErrBookingToLong
	ErrBookingToFarInFuture
	ErrMaxBookingsExceeded
)

func (v BookingRuleViloation) Error() string {
	switch v {
	case ErrInvalidBooking:
		return "ErrInvalidBooking"
	case ErrBookingToLong:
		return "ErrBookingToLong"
	case ErrBookingToFarInFuture:
		return "ErrBookingToFarInFuture"
	case ErrMaxBookingsExceeded:
		return "ErrMaxBookingsExceeded"
	default:
		return "Unkown"
	}
}

// Booking booking of a cabin by a user with a specificed start and end date
type Booking struct {
	ID        string    `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Cabin     Cabin     `json:"cabin"`
	User      User      `json:"user"`
}

func (b Booking) String() string {
	return fmt.Sprintf("Booking(id=%s, startDate=%v, endDate=%v, createdAt=%v, updatedAt=%v, cabin=%s, user=%s)", b.ID, b.StartDate, b.EndDate, b.CreatedAt, b.UpdatedAt, b.Cabin, b.User)
}

func NewBooking(cabin Cabin, user User, startDate, endDate time.Time) Booking {
	return Booking{
		ID:        id.New(),
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
		Cabin:     cabin,
		User:      user,
	}
}

// BookingRequest a request to book a cabin by a user during a specified date
type BookingRequest struct {
	CabinID   string    `json:"cabinId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	UserID    string    `json:"userId"`
	DryRun    bool
}

func (r BookingRequest) Valid() error {
	if r.CabinID == "" {
		return ErrInvalidBooking
	}

	if r.UserID == "" {
		return ErrInvalidBooking
	}

	if r.StartDate.After(r.EndDate) {
		return ErrInvalidBooking
	}

	if r.EndDate.Before(timeutil.Now()) {
		return ErrInvalidBooking
	}

	return nil
}

func (r BookingRequest) String() string {
	return fmt.Sprintf("BookingRequest(cabinId=%s, userId=%s, startDate=%v, endDate=%v, dryRun=%v)", r.CabinID, r.UserID, r.StartDate, r.EndDate, r.DryRun)
}

// BookingRules rules governing if a booking is allowed
type BookingRules struct {
	MaxBookingLengthDays  int
	MaxActiveBookings     int
	MustStartWithinMonths int
}

func (r BookingRules) Allowed(b Booking, currentBookings int) error {
	maxEndDate := b.StartDate.Add(r.maxLength())
	if b.EndDate.After(maxEndDate) {
		return ErrBookingToLong
	}

	maxStartDate := timeutil.Now().AddDate(0, r.MustStartWithinMonths, 0)
	if b.StartDate.After(maxStartDate) {
		return ErrBookingToFarInFuture
	}

	if currentBookings >= r.MaxActiveBookings {
		return ErrMaxBookingsExceeded
	}

	return nil
}

func (r BookingRules) maxLength() time.Duration {
	return time.Hour * 24 * time.Duration(r.MaxBookingLengthDays)
}

func DefaultBookingRules() BookingRules {
	return BookingRules{
		MaxBookingLengthDays:  7,
		MaxActiveBookings:     2,
		MustStartWithinMonths: 3,
	}
}

// DeleteBookingRequest request to delete a booking
type DeleteBookingRequest struct {
	BookingID string
	UserID    string
}

func (r DeleteBookingRequest) String() string {
	return fmt.Sprintf("DeleteBookingRequest(bookingId=%s, userId=%s)", r.BookingID, r.UserID)
}

// BookingFilter filter instruction to recive list of bookings
type BookingFilter struct {
	CabinID string `json:"cabinId,omitempty"`
	UserID  string `json:"userId,omitempty"`
}

func (f BookingFilter) String() string {
	return fmt.Sprintf("BookingFilter(cabinId=%s, userId=%s)", f.CabinID, f.UserID)
}

// BookingRef representation of booking with only references
type BookingRef struct {
	ID        string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	CabinID   string
	UserID    string
}

func (b BookingRef) Booking(cabin Cabin, user User) Booking {
	return Booking{
		ID:        b.ID,
		StartDate: b.StartDate,
		EndDate:   b.EndDate,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		Cabin:     cabin,
		User:      user,
	}
}

func (b BookingRef) String() string {
	return fmt.Sprintf("BookingRef(id=%s, startDate=%v, endDate=%v, createdAt=%v, updatedAt=%v, cabinId=%s, userId=%s)", b.ID, b.StartDate, b.EndDate, b.CreatedAt, b.UpdatedAt, b.CabinID, b.UserID)
}
