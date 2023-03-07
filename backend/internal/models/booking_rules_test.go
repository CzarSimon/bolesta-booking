package models_test

import (
	"testing"
	"time"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/httputil/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestBookingRules(t *testing.T) {
	assert := assert.New(t)

	r := models.BookingRules{
		MaxBookingLengthDays:  7,
		MaxActiveBookings:     2,
		MustStartWithinMonths: 3,
	}

	type testCase struct {
		err     error
		booking models.Booking
		curr    int
		comment string
	}

	now := timeutil.Now()
	maxLength := time.Hour * 24 * time.Duration(r.MaxBookingLengthDays)

	testCases := []testCase{
		{
			booking: models.Booking{
				StartDate: timeutil.Now(),
				EndDate:   timeutil.Now().Add(maxLength),
			},
			curr:    1,
			err:     nil,
			comment: "Should be allowed",
		},
		{
			booking: models.Booking{
				StartDate: now,
				EndDate:   now.Add(maxLength + 24*time.Hour),
			},
			curr:    1,
			err:     models.ErrBookingToLong,
			comment: "Should fail as booking is too long",
		},
		{
			booking: models.Booking{
				StartDate: now.AddDate(0, int(r.MustStartWithinMonths), -1),
				EndDate:   now.AddDate(0, int(r.MustStartWithinMonths), -1).Add(maxLength),
			},
			curr:    1,
			err:     nil,
			comment: "Should be allowed",
		},
		{
			booking: models.Booking{
				StartDate: now.AddDate(0, int(r.MustStartWithinMonths), -1),
				EndDate:   now.AddDate(0, int(r.MustStartWithinMonths), 0).Add(maxLength),
			},
			curr:    1,
			err:     models.ErrBookingToLong,
			comment: "Should fail as booking is too long",
		},
		{
			booking: models.Booking{
				StartDate: now.AddDate(0, int(r.MustStartWithinMonths), 1),
				EndDate:   now.AddDate(0, int(r.MustStartWithinMonths), 3),
			},
			curr:    1,
			err:     models.ErrBookingToFarInFuture,
			comment: "Should fail as start date is to far in the future",
		},
		{
			booking: models.Booking{
				StartDate: timeutil.Now(),
				EndDate:   timeutil.Now().Add(maxLength),
			},
			curr:    r.MaxActiveBookings,
			err:     models.ErrMaxBookingsExceeded,
			comment: "Should fail as user has to many bookings",
		},
	}

	for i, tc := range testCases {
		err := r.Allowed(tc.booking, tc.curr)
		assert.Equal(tc.err, err, "Test Case #%d failed: %s", i+1, tc.comment)
	}
}
