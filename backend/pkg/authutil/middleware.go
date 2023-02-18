package authutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/gin-gonic/gin"
)

const principalKey = "bolesta-booking:backend:pkg:auth:principalKey"

// GetPrincipal returns the authenticated user if exists.
func GetPrincipal(c *gin.Context) (jwt.User, bool) {
	val, ok := c.Get(principalKey)
	if !ok {
		return jwt.User{}, false
	}

	user, ok := val.(jwt.User)
	return user, ok
}

// MustGetPrincipal returns the authenticated user or an error if none exists.
func MustGetPrincipal(c *gin.Context) (jwt.User, error) {
	principal, ok := GetPrincipal(c)
	if !ok {
		return jwt.User{}, fmt.Errorf("failed to parse prinipal from authenticated request")
	}

	return principal, nil
}

type Middleware struct {
	verifier jwt.Verifier
}

func NewMiddleware(creds jwt.Credentials) *Middleware {
	return &Middleware{
		verifier: jwt.NewVerifier(creds, time.Minute),
	}
}

func (m *Middleware) Secure(p Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := m.extractUserFromRequest(c)
		if err != nil {
			c.AbortWithStatusJSON(err.Status, err)
		}

		c.Set(principalKey, user)
		if !isAllowed(user, p) {
			err := httputil.Forbiddenf("forbidden, lacks permission: %s", p)
			c.AbortWithStatusJSON(err.Status, err)
		}

		c.Next()
	}
}

func (m *Middleware) extractUserFromRequest(c *gin.Context) (jwt.User, *httputil.Error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return jwt.User{}, httputil.Unauthorizedf("missing Authorization header")
	}

	token := strings.Replace(header, "Bearer ", "", 1)
	user, err := m.verifier.Verify(token)
	if err != nil {
		return jwt.User{}, httputil.UnauthorizedError(err)
	}

	return user, nil
}

func isAllowed(user jwt.User, p Permission) bool {
	for _, roleName := range user.Roles {
		r, ok := getRole(roleName)
		if !ok {
			continue
		}

		allowed := r.Allowed(p)
		if allowed {
			return true
		}
	}

	return false
}
