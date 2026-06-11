package delivery

import (
    "construction_transport_server/api/rest/v1/dto"
    "construction_transport_server/internal/vehicle/domain"
    "construction_transport_server/internal/vehicle/usecase"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type VehicleHandler struct {
    uc *usecase.VehicleUsecase
}

func NewVehicleHandler(uc *usecase.VehicleUsecase) *VehicleHandler {
    return &VehicleHandler{uc: uc}
}

func (h *VehicleHandler) Create(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    role := c.GetString("role")
    if role != "TRANSPORTER" {
        SendError(c, http.StatusForbidden, "only transporters can add vehicles")
        return
    }

    var input domain.CreateVehicleInput
    if err := c.ShouldBindJSON(&input); err != nil {
        SendError(c, http.StatusBadRequest, "invalid input")
        return
    }
    vehicle, err := h.uc.Create(c.Request.Context(), userID, input)
    if err != nil {
        SendError(c, http.StatusInternalServerError, err.Error())
        return
    }
    SendResponse(c, http.StatusCreated, "vehicle created", vehicle)
}

func (h *VehicleHandler) List(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    vehicles, err := h.uc.ListMyVehicles(c.Request.Context(), userID)
    if err != nil {
        SendError(c, http.StatusInternalServerError, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "vehicles fetched", vehicles)
}

func (h *VehicleHandler) Get(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        SendError(c, http.StatusBadRequest, "invalid vehicle id")
        return
    }
    v, err := h.uc.GetByID(c.Request.Context(), id, userID)
    if err != nil {
        SendError(c, http.StatusNotFound, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "vehicle fetched", v)
}

func (h *VehicleHandler) Update(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        SendError(c, http.StatusBadRequest, "invalid id")
        return
    }
    var input domain.UpdateVehicleInput
    if err := c.ShouldBindJSON(&input); err != nil {
        SendError(c, http.StatusBadRequest, "invalid input")
        return
    }
    if err := h.uc.Update(c.Request.Context(), id, userID, input); err != nil {
        SendError(c, http.StatusInternalServerError, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "vehicle updated", nil)
}

func (h *VehicleHandler) Delete(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        SendError(c, http.StatusBadRequest, "invalid id")
        return
    }
    if err := h.uc.Delete(c.Request.Context(), id, userID); err != nil {
        SendError(c, http.StatusInternalServerError, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "vehicle deleted", nil)
}