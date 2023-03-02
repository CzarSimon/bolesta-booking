package bookings_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/CzarSimon/bolesta-booking/backend/internal/api/bookings"
	"github.com/CzarSimon/bolesta-booking/backend/internal/config"
	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/bolesta-booking/backend/pkg/authutil"
	"github.com/CzarSimon/bolesta-booking/backend/pkg/authutil/authtest"
	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/client/rpc"
	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/CzarSimon/httputil/testutil"
	"github.com/CzarSimon/httputil/timeutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

const correctPassword = "fkalsdjgr3ripåsbiaiogh4åi24h+24fjaga"

func TestGetBooking(t *testing.T) {
	assert := assert.New(t)
	e := newTestEnv()

	user := e.NewUser()
	cabin, found, err := e.cabinRepo.Find(e.ctx, "a4b4f496-767e-423e-9816-83b71e1cfa89") // Bölestastugan
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(cabin)

	booking := models.NewBooking(cabin, user, timeutil.Now(), timeutil.Now().Add(time.Hour))
	err = e.bookingRepo.Save(e.ctx, booking)
	assert.NoError(err)

	authenticator := authtest.NewAuthenticator(e.cfg.JWT)

	path := fmt.Sprintf("/v1/bookings/%s", booking.ID)
	req := testutil.CreateRequest(http.MethodGet, path, nil)
	authenticator.Authenticate(req, id.New(), authutil.UserRole)
	res := testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusOK, res.Code)
	var body models.Booking
	err = json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)

	booking.User.Salt = ""
	booking.User.Password = ""
	assert.Equal(booking, body)

	req = testutil.CreateRequest(http.MethodGet, "/v1/bookings/does-not-extis", nil)
	authenticator.Authenticate(req, id.New(), authutil.UserRole)
	res = testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusNotFound, res.Code)
}

func TestGetBooking_Unauthorized_and_Forbidden(t *testing.T) {
	e := newTestEnv()
	authtest.Test401and403(authtest.TestOpts{
		T:        t,
		Router:   e.router,
		JWTCreds: e.cfg.JWT,
		Method:   http.MethodGet,
		Path:     "/v1/bookings/some-id",
	}, authutil.AnonymousRole)
}

func TestCreateBooking(t *testing.T) {
	assert := assert.New(t)
	e := newTestEnv()

	user := e.NewUser()
	bookingReq := models.BookingRequest{
		CabinID:   "a4b4f496-767e-423e-9816-83b71e1cfa89",
		StartDate: timeutil.Now(),
		EndDate:   timeutil.Now().Add(time.Hour),
	}

	authenticator := authtest.NewAuthenticator(e.cfg.JWT)

	req, _ := rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/bookings", bookingReq)
	authenticator.Authenticate(req, user.ID, authutil.UserRole)
	res := testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusOK, res.Code)

	var body models.Booking
	err := json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)
	assert.Equal(bookingReq.CabinID, body.Cabin.ID)
	assert.Equal(bookingReq.StartDate, body.StartDate)
	assert.Equal(bookingReq.EndDate, body.EndDate)
	user.Password = ""
	user.Salt = ""
	assert.Equal(user, body.User)
}

func TestCreateBooking_invalid(t *testing.T) {
	assert := assert.New(t)
	e := newTestEnv()

	user := e.NewUser()

	type testCase struct {
		req     models.BookingRequest
		status  int
		comment string
	}

	testCases := []testCase{
		{
			req: models.BookingRequest{
				CabinID:   "no-such-cabin-id",
				StartDate: timeutil.Now(),
				EndDate:   timeutil.Now().Add(time.Hour),
				UserID:    user.ID,
			},
			status:  http.StatusPreconditionRequired,
			comment: "Should fail due to invalid cabin id",
		},
		{
			req: models.BookingRequest{
				CabinID:   "a4b4f496-767e-423e-9816-83b71e1cfa89",
				StartDate: timeutil.Now(),
				EndDate:   timeutil.Now().Add(time.Hour),
				UserID:    id.New(),
			},
			status:  http.StatusPreconditionRequired,
			comment: "Should fail due to invalid user id",
		},
		{
			req: models.BookingRequest{
				CabinID:   "a4b4f496-767e-423e-9816-83b71e1cfa89",
				StartDate: timeutil.Now().Add(time.Hour),
				EndDate:   timeutil.Now(),
				UserID:    user.ID,
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to startDate > endDate",
		},
		{
			req: models.BookingRequest{
				CabinID:   "",
				StartDate: timeutil.Now(),
				EndDate:   timeutil.Now().Add(time.Hour),
				UserID:    user.ID,
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to empty cabinId",
		},
		{
			req: models.BookingRequest{
				CabinID:   "a4b4f496-767e-423e-9816-83b71e1cfa89",
				StartDate: timeutil.Now().Add(-2 * time.Hour),
				EndDate:   timeutil.Now().Add(-1 * time.Hour),
				UserID:    user.ID,
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to end date being in the past",
		},
	}

	authenticator := authtest.NewAuthenticator(e.cfg.JWT)

	for i, tc := range testCases {
		userID := tc.req.UserID
		tc.req.UserID = ""
		req, _ := rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/bookings", tc.req)
		authenticator.Authenticate(req, userID, authutil.UserRole)
		res := testutil.PerformRequest(e.router, req)
		assert.Equal(tc.status, res.Code, "Test Case #%d failed: %s", i+1, tc.comment)
	}
}

func TestListBookings(t *testing.T) {
	assert := assert.New(t)
	e := newTestEnv()

	user1 := e.NewUser()
	user2 := e.NewUser()

	bookingReqs := []models.BookingRequest{
		{
			CabinID:   "a4b4f496-767e-423e-9816-83b71e1cfa89",
			StartDate: timeutil.Now(),
			EndDate:   timeutil.Now().Add(time.Hour),
			UserID:    user1.ID,
		},
		{
			CabinID:   "a4b4f496-767e-423e-9816-83b71e1cfa89",
			StartDate: timeutil.Now().Add(2 * time.Hour),
			EndDate:   timeutil.Now().Add(3 * time.Hour),
			UserID:    user2.ID,
		},
		{
			CabinID:   "63e71fef-0037-451f-b731-27249c0164d9",
			StartDate: timeutil.Now().Add(2 * time.Hour),
			EndDate:   timeutil.Now().Add(3 * time.Hour),
			UserID:    user1.ID,
		},
	}

	authenticator := authtest.NewAuthenticator(e.cfg.JWT)

	for i, br := range bookingReqs {
		userID := br.UserID
		br.UserID = ""
		req, _ := rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/bookings", br)
		authenticator.Authenticate(req, userID, authutil.UserRole)
		res := testutil.PerformRequest(e.router, req)
		assert.Equal(http.StatusOK, res.Code, "Failed to create booking #%d", i+1)
	}

	path := "/v1/bookings?cabinId=a4b4f496-767e-423e-9816-83b71e1cfa89"
	req := testutil.CreateRequest(http.MethodGet, path, nil)
	authenticator.Authenticate(req, id.New(), authutil.UserRole)
	res := testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusOK, res.Code)
	var body []models.Booking
	err := json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)
	assert.Len(body, 2)

	path = "/v1/bookings"
	req = testutil.CreateRequest(http.MethodGet, path, nil)
	authenticator.Authenticate(req, id.New(), authutil.UserRole)
	res = testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusOK, res.Code)
	err = json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)
	assert.Len(body, 3)

	path = fmt.Sprintf("/v1/bookings?userId=%s", user2.ID)
	req = testutil.CreateRequest(http.MethodGet, path, nil)
	authenticator.Authenticate(req, id.New(), authutil.UserRole)
	res = testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusOK, res.Code)
	err = json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)
	assert.Len(body, 1)

	path = fmt.Sprintf("/v1/bookings?userId=%s&cabinId=2aa15162-2443-48f1-9b8f-6314f90faf9a", user2.ID)
	req = testutil.CreateRequest(http.MethodGet, path, nil)
	authenticator.Authenticate(req, id.New(), authutil.UserRole)
	res = testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusOK, res.Code)
	err = json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)
	assert.Len(body, 0)
}

func TestListBookings_Unauthorized_and_Forbidden(t *testing.T) {
	e := newTestEnv()
	authtest.Test401and403(authtest.TestOpts{
		T:        t,
		Router:   e.router,
		JWTCreds: e.cfg.JWT,
		Method:   http.MethodGet,
		Path:     "/v1/bookings",
	}, authutil.AnonymousRole)
}

func TestDeleteBooking(t *testing.T) {
	assert := assert.New(t)
	e := newTestEnv()

	user := e.NewUser()

	booking := models.Booking{
		ID:        id.New(),
		StartDate: timeutil.Now(),
		EndDate:   timeutil.Now().Add(time.Hour),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
		Cabin: models.Cabin{
			ID: "a4b4f496-767e-423e-9816-83b71e1cfa89",
		},
		User: user,
	}

	err := e.bookingRepo.Save(e.ctx, booking)
	assert.NoError(err)
	_, found, err := e.bookingRepo.Find(e.ctx, booking.ID)
	assert.NoError(err)
	assert.True(found)

	authenticator := authtest.NewAuthenticator(e.cfg.JWT)
	req := testutil.CreateRequest(http.MethodDelete, fmt.Sprintf("/v1/bookings/%s", booking.ID), nil)
	authenticator.Authenticate(req, user.ID, authutil.UserRole)
	res := testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusOK, res.Code)
	_, found, err = e.bookingRepo.Find(e.ctx, booking.ID)
	assert.NoError(err)
	assert.False(found)
}

func TestDeleteBooking_wrongUser(t *testing.T) {
	assert := assert.New(t)
	e := newTestEnv()

	user := e.NewUser()

	booking := models.Booking{
		ID:        id.New(),
		StartDate: timeutil.Now(),
		EndDate:   timeutil.Now().Add(time.Hour),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
		Cabin: models.Cabin{
			ID: "a4b4f496-767e-423e-9816-83b71e1cfa89",
		},
		User: user,
	}

	err := e.bookingRepo.Save(e.ctx, booking)
	assert.NoError(err)
	_, found, err := e.bookingRepo.Find(e.ctx, booking.ID)
	assert.NoError(err)
	assert.True(found)

	otherUser := e.NewUser()
	authenticator := authtest.NewAuthenticator(e.cfg.JWT)
	req := testutil.CreateRequest(http.MethodDelete, fmt.Sprintf("/v1/bookings/%s", booking.ID), nil)
	authenticator.Authenticate(req, otherUser.ID, authutil.UserRole)
	res := testutil.PerformRequest(e.router, req)

	assert.Equal(http.StatusForbidden, res.Code)
	_, found, err = e.bookingRepo.Find(e.ctx, booking.ID)
	assert.NoError(err)
	assert.True(found)
}

func TestDeleteBooking_Unauthorized_and_Forbidden(t *testing.T) {
	e := newTestEnv()
	authtest.Test401and403(authtest.TestOpts{
		T:        t,
		Router:   e.router,
		JWTCreds: e.cfg.JWT,
		Method:   http.MethodDelete,
		Path:     "/v1/bookings/some-id",
	}, authutil.AnonymousRole, authutil.AdminRole)
}

type testEnv struct {
	router      http.Handler
	cabinRepo   repository.CabinRepository
	userRepo    repository.UserRepository
	bookingRepo repository.BookingRepository
	svc         *service.BookingService
	cfg         config.Config
	ctx         context.Context
}

func (e testEnv) NewUser() models.User {
	name := id.New()
	req := models.CreateUserRequest{
		Name:     name,
		Email:    fmt.Sprintf("%s@mail.com", name),
		Password: correctPassword,
	}

	svc := &service.UserService{UserRepo: e.userRepo}
	user, err := svc.CreateUser(e.ctx, req)
	if err != nil {
		log.Fatalf("failed to insert %s: %v", user, err)
	}

	return user
}

func newTestEnv() testEnv {
	db := testutil.InMemoryDB(true, "../../../resources/db/sqlite")

	cabinRepo := repository.NewCabinRepository(db)
	userRepo := repository.NewUserRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	svc := &service.BookingService{
		BookingRepo: bookingRepo,
		CabinRepo:   cabinRepo,
		UserRepo:    userRepo,
	}

	r := httputil.NewRouter("backend", func() error {
		return nil
	})

	cfg := config.Config{
		JWT: jwt.Credentials{
			Issuer: "bolesta-booking/backend",
			Secret: id.New(),
		},
	}

	bookings.AttachController(svc, r, cfg)
	return testEnv{
		router:      r,
		cabinRepo:   cabinRepo,
		userRepo:    userRepo,
		bookingRepo: bookingRepo,
		svc:         svc,
		cfg:         cfg,
		ctx:         context.Background(),
	}
}
