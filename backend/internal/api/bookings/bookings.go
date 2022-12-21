package bookings

import (
	"net/http"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/service"
	"github.com/CzarSimon/httputil"
	"github.com/gin-gonic/gin"
)

// Controller http handler for cabins
type controller struct {
	svc *service.BookingService
}

// AttachController attaches a controller to the specified route group.
func AttachController(svc *service.BookingService, r gin.IRouter) {
	controller := &controller{svc: svc}
	g := r.Group("/v1/bookings")

	g.POST("", controller.createBooking)
	g.GET("", controller.listBookings)
	g.GET("/:id", controller.getBooking)
}

func (h *controller) getBooking(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	booking, err := h.svc.GetBooking(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *controller) listBookings(c *gin.Context) {
	ctx := c.Request.Context()
	f := parseBookingFilter(c)

	bookings, err := h.svc.ListBookings(ctx, f)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *controller) createBooking(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := parseBookingRequest(c)
	if err != nil {
		c.Error(err)
		return
	}

	booking, err := h.svc.CreateBooking(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, booking)
}

func parseBookingRequest(c *gin.Context) (models.BookingRequest, error) {
	var body models.BookingRequest
	err := c.BindJSON(&body)
	if err != nil {
		err = httputil.BadRequestf("failed to parse request body. %w", err)
		return models.BookingRequest{}, err
	}

	err = body.Valid()
	if err != nil {
		err = httputil.BadRequestf("invalid %s. %w", body, err)
		return models.BookingRequest{}, err
	}

	return body, nil
}

func parseBookingFilter(c *gin.Context) models.BookingFilter {
	return models.BookingFilter{
		CabinID: c.Query("cabinId"),
		UserID:  c.Query("userId"),
	}
}
