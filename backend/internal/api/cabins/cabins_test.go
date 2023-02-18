package cabins_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CzarSimon/bolesta-booking/backend/internal/api/cabins"
	"github.com/CzarSimon/bolesta-booking/backend/internal/config"
	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/bolesta-booking/backend/pkg/authutil"
	"github.com/CzarSimon/bolesta-booking/backend/pkg/authutil/authtest"
	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/CzarSimon/httputil/testutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestGetCabin(t *testing.T) {
	assert := assert.New(t)
	router, cfg := setupRouter()

	req := testutil.CreateRequest(http.MethodGet, "/v1/cabins/63e71fef-0037-451f-b731-27249c0164d9", nil)
	authtest.Authenticate(req, id.New(), authutil.UserRole, cfg.JWT)
	res := testutil.PerformRequest(router, req)

	assert.Equal(http.StatusOK, res.Code)
	var body models.Cabin
	err := json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)
	assert.Equal("63e71fef-0037-451f-b731-27249c0164d9", body.ID)
	assert.Equal("Gulhuset", body.Name)
	assert.NotEmpty(body.CreatedAt)
	assert.NotEmpty(body.UpdatedAt)

	req = testutil.CreateRequest(http.MethodGet, "/v1/cabins/does-not-extis", nil)
	authtest.Authenticate(req, id.New(), authutil.UserRole, cfg.JWT)
	res = testutil.PerformRequest(router, req)

	assert.Equal(http.StatusNotFound, res.Code)
}

func TestGetCabin_Unauthorized_and_Forbidden(t *testing.T) {
	router, cfg := setupRouter()
	authtest.Test401and403(authtest.TestOpts{
		T:        t,
		Router:   router,
		JWTCreds: cfg.JWT,
		Method:   http.MethodGet,
		Path:     "/v1/cabins/some-cabin-id",
	}, authutil.AnonymousRole)
}

func TestGetCabins(t *testing.T) {
	assert := assert.New(t)
	router, cfg := setupRouter()

	req := testutil.CreateRequest(http.MethodGet, "/v1/cabins", nil)
	authtest.Authenticate(req, id.New(), authutil.UserRole, cfg.JWT)
	res := testutil.PerformRequest(router, req)

	assert.Equal(http.StatusOK, res.Code)
	var body []models.Cabin
	err := json.NewDecoder(res.Result().Body).Decode(&body)
	assert.NoError(err)
	assert.Len(body, 3)
}

func TestGetCabins_Unauthorized_and_Forbidden(t *testing.T) {
	router, cfg := setupRouter()
	authtest.Test401and403(authtest.TestOpts{
		T:        t,
		Router:   router,
		JWTCreds: cfg.JWT,
		Method:   http.MethodGet,
		Path:     "/v1/cabins",
	}, authutil.AnonymousRole)
}

func setupRouter() (http.Handler, config.Config) {
	db := testutil.InMemoryDB(true, "../../../resources/db/sqlite")
	repo := repository.NewCabinRepository(db)
	svc := &service.CabinService{
		CabinRepo: repo,
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

	cabins.AttachController(svc, r, cfg)
	return r, cfg
}
