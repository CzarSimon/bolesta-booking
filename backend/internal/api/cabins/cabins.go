package cabins

import (
	"net/http"

	"github.com/CzarSimon/bolesta-booking/backend/internal/config"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/bolesta-booking/backend/pkg/authutil"
	"github.com/gin-gonic/gin"
)

// Controller http handler for cabins
type controller struct {
	svc *service.CabinService
}

// AttachController attaches a controller to the specified route group.
func AttachController(svc *service.CabinService, r gin.IRouter, cfg config.Config) {
	controller := &controller{svc: svc}
	authz := authutil.NewMiddleware(cfg.JWT)

	g := r.Group("/v1/cabins").Use(authz.Secure(authutil.ReadCabin))

	g.GET("", controller.GetCabins)
	g.GET("/:id", controller.GetCabin)
}

func (h *controller) GetCabin(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	cabin, err := h.svc.GetCabin(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cabin)
}

func (h *controller) GetCabins(c *gin.Context) {
	ctx := c.Request.Context()

	cabins, err := h.svc.GetCabins(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cabins)
}
