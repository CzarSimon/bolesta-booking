package authtest

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/CzarSimon/httputil/testutil"
	"github.com/stretchr/testify/assert"
)

type Authenticator struct {
	issuer jwt.Issuer
}

func NewAuthenticator(creds jwt.Credentials) Authenticator {
	return Authenticator{
		issuer: jwt.NewIssuer(creds),
	}
}

func (a Authenticator) Authenticate(req *http.Request, id, role string) {
	token, err := a.issuer.Issue(jwt.User{
		ID:    id,
		Roles: []string{role},
	}, time.Hour)

	if err != nil {
		log.Fatalf("failed to genereate token: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}

func Authenticate(req *http.Request, id, role string, creds jwt.Credentials) {
	NewAuthenticator(creds).Authenticate(req, id, role)
}

type TestOpts struct {
	T        *testing.T
	Router   http.Handler
	JWTCreds jwt.Credentials
	Method   string
	Path     string
}

func Test401and403(opts TestOpts, forbiddenRoles ...string) {
	assert := assert.New(opts.T)

	req := testutil.CreateRequest(opts.Method, opts.Path, nil)
	res := testutil.PerformRequest(opts.Router, req)
	assert.Equal(http.StatusUnauthorized, res.Code)

	authenticator := NewAuthenticator(opts.JWTCreds)

	for _, role := range forbiddenRoles {
		authenticator.Authenticate(req, id.New(), role)
		res = testutil.PerformRequest(opts.Router, req)
		assert.Equal(http.StatusForbidden, res.Code)
	}
}
