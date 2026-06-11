package delivery

import (
	"construction_transport_server/api/rest/v1/dto"
	"construction_transport_server/internal/booking/domain"
	"construction_transport_server/internal/booking/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	uc *usecase.BookingUsecase
}

func NewBookingHandler(uc *usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{uc: uc}
}

func (h *BookingHandler) Create(c *gin.Context) {
	userID := c.GetInt64("auth_id")
	role := c.GetString("role")
	if role != "USER" {
		SendError(c, http.StatusForbidden, "only customers can create bookings")
		return
	}
	var input domain.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		SendError(c, http.StatusBadRequest, "invalid input")
		return
	}
	booking, err := h.uc.CreateBooking(c.Request.Context(), userID, input)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(c, http.StatusCreated, "booking created", booking)
}

func (h *BookingHandler) Get(c *gin.Context) {
	userID := c.GetInt64("auth_id")
	role := c.GetString("role")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendError(c, http.StatusBadRequest, "invalid booking id")
		return
	}
	booking, err := h.uc.GetBooking(c.Request.Context(), id, userID, role)
	if err != nil {
		SendError(c, http.StatusNotFound, err.Error())
		return
	}
	SendResponse(c, http.StatusOK, "booking fetched", booking)
}

func (h *BookingHandler) ListCustomer(c *gin.Context) {
	userID := c.GetInt64("auth_id")
	bookings, err := h.uc.ListCustomerBookings(c.Request.Context(), userID)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(c, http.StatusOK, "bookings fetched", bookings)
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	userID := c.GetInt64("auth_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.uc.CancelBooking(c.Request.Context(), id, userID); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	SendResponse(c, http.StatusOK, "booking cancelled", nil)
}
