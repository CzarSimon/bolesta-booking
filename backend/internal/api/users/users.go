package users

import (
	"net/http"

	"github.com/CzarSimon/bolesta-booking/backend/internal/config"
	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/httputil"
	"github.com/gin-gonic/gin"
)

// Controller http handler for cabins
type controller struct {
	svc               *service.UserService
	enableCreateUsers bool
}

// AttachController attaches a controller to the specified route group.
func AttachController(svc *service.UserService, r gin.IRouter, cfg config.Config) {
	controller := &controller{
		svc:               svc,
		enableCreateUsers: cfg.EnableCreateUsers,
	}
	g := r.Group("/v1/users")

	g.GET("", controller.GetUsers)
	g.POST("", controller.CreateUser)
	g.GET("/:id", controller.GetUser)
}

func (h *controller) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *controller) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.svc.GetUsers(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *controller) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	if !h.enableCreateUsers {
		c.Error(httputil.Forbiddenf("Creating users is not enabled"))
		return
	}

	req, err := parseCreateUserRequest(c)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := h.svc.CreateUser(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func parseCreateUserRequest(c *gin.Context) (models.CreateUserRequest, error) {
	var body models.CreateUserRequest
	err := c.BindJSON(&body)
	if err != nil {
		err = httputil.BadRequestf("failed to parse request body. %w", err)
		return models.CreateUserRequest{}, err
	}

	err = body.Valid()
	if err != nil {
		err = httputil.BadRequestf("invalid %s. %w", body, err)
		return models.CreateUserRequest{}, err
	}

	return body, nil
}
