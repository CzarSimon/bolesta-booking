package service

import (
	"context"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/httputil"
)

type BookingService struct {
	BookingRepo repository.BookingRepository
	UserRepo    repository.UserRepository
	CabinRepo   repository.CabinRepository
}

func (s *BookingService) GetBooking(ctx context.Context, id string) (models.Booking, error) {
	booking, exists, err := s.BookingRepo.Find(ctx, id)
	if err != nil {
		return models.Booking{}, err
	}

	if !exists {
		return models.Booking{}, httputil.NotFoundError(err)
	}

	return booking, nil
}

func (s *BookingService) ListBookings(ctx context.Context, f models.BookingFilter) ([]models.Booking, error) {
	return s.BookingRepo.FindByFilter(ctx, f)
}

func (s *BookingService) CreateBooking(ctx context.Context, req models.BookingRequest) (models.Booking, error) {
	cabin, err := s.mustGetCabin(ctx, req.CabinID)
	if err != nil {
		return models.Booking{}, err
	}

	user, err := s.mustGetUser(ctx, req.UserID)
	if err != nil {
		return models.Booking{}, err
	}

	booking := models.NewBooking(cabin, user, req.StartDate, req.EndDate)

	err = s.BookingRepo.Save(ctx, booking)
	if err != nil {
		return models.Booking{}, err
	}

	return booking, nil
}

func (s *BookingService) DeleteBooking(ctx context.Context, req models.DeleteBookingRequest) error {
	_, err := s.mustGetUser(ctx, req.UserID)
	if err != nil {
		return err
	}

	booking, err := s.mustGetBooking(ctx, req.BookingID)
	if err != nil {
		return err
	}

	if booking.User.ID != req.UserID {
		return httputil.Forbiddenf("User(id=%s) does not own Booking(id=%s)", req.UserID, req.BookingID)
	}

	return s.BookingRepo.Delete(ctx, req.BookingID)
}

func (s *BookingService) mustGetCabin(ctx context.Context, cabinID string) (models.Cabin, error) {
	cabin, exists, err := s.CabinRepo.Find(ctx, cabinID)
	if err != nil {
		return models.Cabin{}, err
	}

	if !exists {
		return models.Cabin{}, httputil.PreconditionRequiredError(err)
	}

	return cabin, nil
}

func (s *BookingService) mustGetUser(ctx context.Context, cabinID string) (models.User, error) {
	user, exists, err := s.UserRepo.Find(ctx, cabinID)
	if err != nil {
		return models.User{}, err
	}

	if !exists {
		return models.User{}, httputil.PreconditionRequiredError(err)
	}

	return user, nil
}

func (s *BookingService) mustGetBooking(ctx context.Context, bookingID string) (models.Booking, error) {
	booking, exists, err := s.BookingRepo.Find(ctx, bookingID)
	if err != nil {
		return models.Booking{}, err
	}

	if !exists {
		return models.Booking{}, httputil.PreconditionRequiredError(err)
	}

	return booking, nil
}
