package users

import (
	"net/http"

	"github.com/CzarSimon/bolesta-booking/backend/internal/config"
	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/bolesta-booking/backend/pkg/authutil"
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
	authz := authutil.NewMiddleware(cfg.JWT)
	g := r.Group("/v1/users")

	g.GET("", authz.Secure(authutil.ReadUser), controller.getUsers)
	g.POST("", controller.createUser)
	g.GET("/:id", authz.Secure(authutil.ReadUser), controller.getUser)
	g.PUT("/:id/password", authz.Secure(authutil.UpdateUser), controller.changePassword)
}

func (h *controller) getUser(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *controller) getUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.svc.GetUsers(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *controller) createUser(c *gin.Context) {
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

func (h *controller) changePassword(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := parseChangePasswordRequest(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.svc.ChangePassword(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	httputil.SendOK(c)
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

func parseChangePasswordRequest(c *gin.Context) (models.ChangePasswordRequest, error) {
	principal, err := authutil.MustGetPrincipal(c)
	if err != nil {
		return models.ChangePasswordRequest{}, err
	}

	var body models.ChangePasswordRequest
	err = c.BindJSON(&body)
	if err != nil {
		err = httputil.BadRequestf("failed to parse request body. %w", err)
		return models.ChangePasswordRequest{}, err
	}
	body.UserID = principal.ID

	err = body.Valid()
	if err != nil {
		err = httputil.BadRequestf("invalid %s. %w", body, err)
		return models.ChangePasswordRequest{}, err
	}

	return body, nil
}
