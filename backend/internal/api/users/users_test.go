package users_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/CzarSimon/bolesta-booking/backend/internal/api/users"
	"github.com/CzarSimon/bolesta-booking/backend/internal/config"
	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/client/rpc"
	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/testutil"
	"github.com/CzarSimon/httputil/timeutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	assert := assert.New(t)
	router, svc := setupRouter(false)
	ctx := context.Background()

	u1 := models.User{
		ID:        id.New(),
		Name:      "Some Name",
		Email:     "mail@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	err := svc.UserRepo.Save(ctx, u1)
	assert.NoError(err)

	u2 := models.User{
		ID:        id.New(),
		Name:      "Other Name",
		Email:     "other@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	err = svc.UserRepo.Save(ctx, u2)
	assert.NoError(err)

	req := testutil.CreateRequest(http.MethodGet, "/v1/users", nil)
	res := testutil.PerformRequest(router, req)
	assert.Equal(http.StatusOK, res.Code)
	var body []models.User
	err = json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)
	assert.Len(body, 2)
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)
	router, svc := setupRouter(false)
	ctx := context.Background()

	u1 := models.User{
		ID:        id.New(),
		Name:      "Some Name",
		Email:     "mail@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	err := svc.UserRepo.Save(ctx, u1)
	assert.NoError(err)

	path := fmt.Sprintf("/v1/users/%s", u1.ID)
	req := testutil.CreateRequest(http.MethodGet, path, nil)
	res := testutil.PerformRequest(router, req)
	assert.Equal(http.StatusOK, res.Code)
	var body models.User
	err = json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)

	u1.Salt = ""
	u1.Password = ""
	assert.Equal(u1, body)

	req = testutil.CreateRequest(http.MethodGet, "/v1/users/does-not-extis", nil)
	res = testutil.PerformRequest(router, req)
	assert.Equal(http.StatusNotFound, res.Code)
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	router, _ := setupRouter(true)

	ur := models.CreateUserRequest{
		Name:     "Some Name",
		Email:    "mail@mail.com",
		Password: id.New(),
	}

	req, _ := rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/users", ur)
	res := testutil.PerformRequest(router, req)
	assert.Equal(http.StatusOK, res.Code)
	var body models.User
	err := json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)

	assert.Equal(ur.Name, body.Name)
	assert.Equal(ur.Email, body.Email)
	assert.Empty(body.Password)

	path := fmt.Sprintf("/v1/users/%s", body.ID)
	req = testutil.CreateRequest(http.MethodGet, path, nil)
	res = testutil.PerformRequest(router, req)
	assert.Equal(http.StatusOK, res.Code)
	var getBody models.User
	err = json.NewDecoder(res.Result().Body).Decode(&getBody)
	assert.NoError(err)

	assert.Equal(body, getBody)
}

func TestCreateUser_notEnabled(t *testing.T) {
	assert := assert.New(t)
	router, _ := setupRouter(false)

	ur := models.CreateUserRequest{
		Name:     "Some Name",
		Email:    "mail@mail.com",
		Password: id.New(),
	}

	req, _ := rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/users", ur)
	res := testutil.PerformRequest(router, req)
	assert.Equal(http.StatusForbidden, res.Code)
}

func TestCreateUser_invalid(t *testing.T) {
	assert := assert.New(t)
	router, _ := setupRouter(true)

	type testCase struct {
		req     models.CreateUserRequest
		status  int
		comment string
	}

	testCases := []testCase{
		{
			req: models.CreateUserRequest{
				Email:    "mail@mail.com",
				Password: id.New(),
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to empty name",
		},
		{
			req: models.CreateUserRequest{
				Name:     "some name",
				Password: id.New(),
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to empty email",
		},
		{
			req: models.CreateUserRequest{
				Name:     "some name",
				Email:    "mail@mail.com",
				Password: "1234567",
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to too short password",
		},
		{
			req: models.CreateUserRequest{
				Name:     "some name",
				Email:    "@mail.com",
				Password: id.New(),
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to too malformed email",
		},
		{
			req: models.CreateUserRequest{
				Name:     "some name",
				Email:    "name@.com",
				Password: id.New(),
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to too malformed email",
		},
		{
			req: models.CreateUserRequest{
				Name:     "some name",
				Email:    "domain.com",
				Password: id.New(),
			},
			status:  http.StatusBadRequest,
			comment: "Should fail due to too malformed email",
		},
	}

	for i, tc := range testCases {
		req, _ := rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/users", tc.req)
		res := testutil.PerformRequest(router, req)
		assert.Equal(tc.status, res.Code, "Test Case #%d failed: %s", i+1, tc.comment)
	}
}

func setupRouter(enableCreateUsers bool) (http.Handler, *service.UserService) {
	db := testutil.InMemoryDB(true, "../../../resources/db/sqlite")
	repo := repository.NewUserRepository(db)
	svc := &service.UserService{
		UserRepo: repo,
	}

	r := httputil.NewRouter("backend", func() error {
		return nil
	})

	users.AttachController(svc, r, config.Config{EnableCreateUsers: enableCreateUsers})
	return r, svc
}
