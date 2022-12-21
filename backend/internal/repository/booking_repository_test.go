package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/testutil"
	"github.com/CzarSimon/httputil/timeutil"
	"github.com/stretchr/testify/assert"
)

func Test_bookingRepo_Save(t *testing.T) {
	assert := assert.New(t)
	db := testutil.InMemoryDB(true, "../../resources/db/sqlite")
	userRepo := repository.NewUserRepository(db)
	cabinRepo := repository.NewCabinRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	ctx := context.Background()

	user := models.User{
		ID:        id.New(),
		Name:      "Some Name",
		Email:     "mail@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	err := userRepo.Save(ctx, user)
	assert.NoError(err)

	c1, found, err := cabinRepo.Find(ctx, "a4b4f496-767e-423e-9816-83b71e1cfa89") // Bölestastugan
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(c1)

	c2, found, err := cabinRepo.Find(ctx, "63e71fef-0037-451f-b731-27249c0164d9") // Gulhuset
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(c1)

	var rowCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM booking").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(0, rowCount)

	startDate := timeutil.Now()
	endDate := timeutil.Now().Add(24 * time.Hour)
	booking := models.NewBooking(c1, user, startDate, endDate)
	err = bookingRepo.Save(ctx, booking) // Should insert booking of c1
	assert.NoError(err)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM booking").Scan(&rowCount)
	assert.NoError(err) // Should insert booking of c1
	assert.Equal(1, rowCount)

	booking.Cabin = c2
	err = bookingRepo.Save(ctx, booking)
	assert.Error(err) // Should fail on primary key conflict

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM booking").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(1, rowCount)

	booking.ID = id.New()
	booking.StartDate = startDate.Add(time.Hour)
	booking.EndDate = endDate.Add(-1 * time.Hour)
	err = bookingRepo.Save(ctx, booking)
	assert.NoError(err) // Should insert booking of c2

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM booking").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(2, rowCount)

	booking.ID = id.New()
	booking.StartDate = startDate.Add(time.Hour)
	booking.EndDate = endDate.Add(time.Hour)
	err = bookingRepo.Save(ctx, booking)
	assert.Error(err) // Should not insert booking of c2 as start date is before end_date of existing booking

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM booking").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(2, rowCount)

	booking.ID = id.New()
	booking.StartDate = startDate.Add(-48 * time.Hour)
	booking.EndDate = startDate.Add(5 * time.Hour)
	err = bookingRepo.Save(ctx, booking)
	assert.Error(err) // Should not insert booking of c2 as end date is after start_date of existing booking

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM booking").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(2, rowCount)

	booking = models.NewBooking(c1, user, endDate.Add(time.Hour), endDate.Add(25*time.Hour))
	err = bookingRepo.Save(ctx, booking)
	assert.NoError(err)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM booking").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(3, rowCount)

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	booking.ID = id.New()
	err = bookingRepo.Save(ctx, booking)
	assert.Error(err)
	assert.True(errors.Is(err, context.Canceled))
}

func Test_bookingRepo_Find(t *testing.T) {
	assert := assert.New(t)
	db := testutil.InMemoryDB(true, "../../resources/db/sqlite")
	userRepo := repository.NewUserRepository(db)
	cabinRepo := repository.NewCabinRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	ctx := context.Background()

	user := models.User{
		ID:        id.New(),
		Name:      "Some Name",
		Email:     "mail@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	err := userRepo.Save(ctx, user)
	assert.NoError(err)

	c1, found, err := cabinRepo.Find(ctx, "a4b4f496-767e-423e-9816-83b71e1cfa89") // Bölestastugan
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(c1)

	c2, found, err := cabinRepo.Find(ctx, "63e71fef-0037-451f-b731-27249c0164d9") // Gulhuset
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(c2)

	b1 := models.NewBooking(c1, user, timeutil.Now(), timeutil.Now().Add(24*time.Hour))
	err = bookingRepo.Save(ctx, b1)
	assert.NoError(err)

	b2 := models.NewBooking(c2, user, timeutil.Now(), timeutil.Now().Add(24*time.Hour))
	err = bookingRepo.Save(ctx, b2)
	assert.NoError(err)

	actual, found, err := bookingRepo.Find(ctx, b1.ID)
	assert.NoError(err)
	assert.True(found)
	assert.Equal(b1, actual)

	actual, found, err = bookingRepo.Find(ctx, id.New())
	assert.NoError(err)
	assert.False(found)
	assert.Empty(actual)

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, _, err = bookingRepo.Find(ctx, b1.ID)
	assert.Error(err)
	assert.True(errors.Is(err, context.Canceled))
}

func Test_bookingRepo_FindByFilter(t *testing.T) {
	assert := assert.New(t)
	db := testutil.InMemoryDB(true, "../../resources/db/sqlite")
	userRepo := repository.NewUserRepository(db)
	userSvc := &service.UserService{UserRepo: userRepo}
	cabinRepo := repository.NewCabinRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	ctx := context.Background()

	user1, err := userSvc.CreateUser(ctx, models.CreateUserRequest{
		Name:     "Some Name",
		Email:    "some.name@mail.com",
		Password: "some-password",
	})
	assert.NoError(err)
	user2, err := userSvc.CreateUser(ctx, models.CreateUserRequest{
		Name:     "Other Name",
		Email:    "other.name@mail.com",
		Password: "other-password",
	})
	assert.NoError(err)

	cabin1, _, _ := cabinRepo.Find(ctx, "a4b4f496-767e-423e-9816-83b71e1cfa89") // Bölestastugan
	assert.NotEmpty(cabin1)

	cabin2, _, _ := cabinRepo.Find(ctx, "63e71fef-0037-451f-b731-27249c0164d9") // Gulhuset
	assert.NotEmpty(cabin2)

	now := timeutil.Now()
	inAWeek := now.Add(7 * 24 * time.Hour)
	inTwoWeeks := now.Add(2 * 7 * 24 * time.Hour)
	inThreeWeeks := now.Add(3 * 7 * 24 * time.Hour)

	bookings := []models.Booking{
		models.NewBooking(cabin1, user1, now, inAWeek),
		models.NewBooking(cabin1, user2, inTwoWeeks, inThreeWeeks),
		models.NewBooking(cabin2, user1, inAWeek, inTwoWeeks),
	}

	for _, booking := range bookings {
		err = bookingRepo.Save(ctx, booking)
		assert.NoError(err)
	}

	res, err := bookingRepo.FindByFilter(ctx, models.BookingFilter{})
	assert.NoError(err)
	assert.Len(res, 3)
	assertBookings(t, res, bookings...)

	res, err = bookingRepo.FindByFilter(ctx, models.BookingFilter{
		CabinID: cabin1.ID,
	})
	assert.NoError(err)
	assert.Len(res, 2)
	assertBookings(t, res, bookings[0], bookings[1])

	res, err = bookingRepo.FindByFilter(ctx, models.BookingFilter{
		UserID: user1.ID,
	})
	assert.NoError(err)
	assert.Len(res, 2)
	assertBookings(t, res, bookings[0], bookings[2])

	res, err = bookingRepo.FindByFilter(ctx, models.BookingFilter{
		CabinID: cabin1.ID,
		UserID:  user1.ID,
	})
	assert.NoError(err)
	assert.Len(res, 1)
	assertBookings(t, res, bookings[0])

	res, err = bookingRepo.FindByFilter(ctx, models.BookingFilter{
		CabinID: cabin2.ID,
		UserID:  user1.ID,
	})
	assert.NoError(err)
	assert.Len(res, 1)
	assertBookings(t, res, bookings[2])

	res, err = bookingRepo.FindByFilter(ctx, models.BookingFilter{
		CabinID: cabin1.ID,
		UserID:  user2.ID,
	})
	assert.NoError(err)
	assert.Len(res, 1)
	assertBookings(t, res, bookings[1])
}

func assertBookings(t *testing.T, actual []models.Booking, expected ...models.Booking) {
	assert := assert.New(t)

	actualMap := make(map[string]models.Booking)
	for _, b := range actual {
		actualMap[b.ID] = b
	}

	for i, booking := range expected {
		found, ok := actualMap[booking.ID]
		assert.True(ok, "Should have found booking #%d with id=%s", i+1, booking.ID)
		assert.Equal(booking, found)
	}
}
