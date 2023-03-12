package auth

import (
	"net/http"
	"strings"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/httputil"
	"github.com/gin-gonic/gin"
)

// Controller http handler for cabins
type controller struct {
	svc *service.AuthService
}

// AttachController attaches a controller to the specified route group.
func AttachController(svc *service.AuthService, r gin.IRouter) {
	controller := &controller{
		svc: svc,
	}
	g := r.Group("/v1")

	g.POST("/login", controller.login)
}

func (h *controller) login(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := parseCreateUserRequest(c)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := h.svc.Authenticate(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func parseCreateUserRequest(c *gin.Context) (models.LoginRequest, error) {
	var body models.LoginRequest
	err := c.BindJSON(&body)
	if err != nil {
		err = httputil.BadRequestf("failed to parse request body. %w", err)
		return models.LoginRequest{}, err
	}

	body.Email = strings.ToLower(body.Email)
	return body, nil
}
