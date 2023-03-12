package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/CzarSimon/bolesta-booking/backend/internal/api/auth"
	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/client/rpc"
	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/CzarSimon/httputil/testutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var jwtCreds = jwt.Credentials{
	Issuer: "bolesta-booking",
	Secret: id.New(),
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	router, userSvc := setupRouter(true)
	ctx := context.Background()

	ur := models.CreateUserRequest{
		Name:     "Some Name",
		Email:    "mail@mail.com",
		Password: id.New(),
	}

	user, err := userSvc.CreateUser(ctx, ur)
	assert.NoError(err)

	lr := models.LoginRequest{
		Email:    user.Email,
		Password: ur.Password,
	}

	req, _ := rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/login", lr)
	res := testutil.PerformRequest(router, req)
	assert.Equal(http.StatusOK, res.Code)

	var body models.AuthenticatedResponse
	err = json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)

	assert.Equal(user.ID, body.User.ID)
	assert.Equal(ur.Name, body.User.Name)
	assert.Equal(ur.Email, body.User.Email)
	assert.Empty(body.User.Password)

	jwtUser, err := jwt.NewVerifier(jwtCreds, time.Minute).Verify(body.Token)
	assert.NoError(err)
	assert.Equal(user.ID, jwtUser.ID)
	assert.True(jwtUser.HasRole(models.UserRole))
	assert.Contains(jwtUser.Roles, models.UserRole)

	lr.Email = "Mail@mail.com"
	req, _ = rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/login", lr)
	res = testutil.PerformRequest(router, req)
	assert.Equal(http.StatusOK, res.Code)

	lr.Password = "wrong-password"
	req, _ = rpc.NewClient(time.Second).CreateRequest(http.MethodPost, "/v1/login", lr)
	res = testutil.PerformRequest(router, req)
	assert.Equal(http.StatusUnauthorized, res.Code)
}

func setupRouter(enableCreateUsers bool) (http.Handler, *service.UserService) {
	db := testutil.InMemoryDB(true, "../../../resources/db/sqlite")
	repo := repository.NewUserRepository(db)
	authSvc := &service.AuthService{
		UserRepo: repo,
		Issuer:   jwt.NewIssuer(jwtCreds),
	}

	userSvc := &service.UserService{
		UserRepo: repo,
	}

	r := httputil.NewRouter("backend", func() error {
		return nil
	})

	auth.AttachController(authSvc, r)
	return r, userSvc
}
