package models

import (
	"fmt"
	"time"

	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/timeutil"
)

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
}

func (r BookingRequest) Valid() error {
	if r.CabinID == "" {
		return fmt.Errorf("CabinID cannot be empty")
	}

	if r.UserID == "" {
		return fmt.Errorf("UserID cannot be empty")
	}

	if r.StartDate.After(r.EndDate) {
		return fmt.Errorf("StartDate must be before EndDate")
	}

	if r.EndDate.Before(timeutil.Now()) {
		return fmt.Errorf("EndDate must be in the future")
	}

	return nil
}

func (r BookingRequest) String() string {
	return fmt.Sprintf("BookingRequest(cabinId=%s, userId=%s, startDate=%v, endDate=%v)", r.CabinID, r.UserID, r.StartDate, r.EndDate)
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
