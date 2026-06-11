package delivery

import (
    "construction_transport_server/internal/account/usecase"
    "net/http"
    "github.com/gin-gonic/gin"
)

type AccountHandler struct {
    uc *usecase.AccountUsecase
}

func NewAccountHandler(uc *usecase.AccountUsecase) *AccountHandler {
    return &AccountHandler{uc: uc}
}

func (h *AccountHandler) GetProfile(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    profile, err := h.uc.GetProfile(c.Request.Context(), userID)
    if err != nil {
        SendError(c, http.StatusInternalServerError, err.Error())
        return
    }
    if profile == nil {
        SendResponse(c, http.StatusOK, "profile not completed", gin.H{})
        return
    }
    SendResponse(c, http.StatusOK, "profile fetched", profile)
}

func (h *AccountHandler) UpdateProfile(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    var req struct {
        Name         string `json:"full_name"`
        PhoneNumber  string `json:"phone_number"`
        ProfileImage string `json:"profile_image"`
        Location     string `json:"location"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        SendError(c, http.StatusBadRequest, "invalid input")
        return
    }
    profile, err := h.uc.CreateOrUpdateProfile(c.Request.Context(), userID, 
        req.Name, req.PhoneNumber, req.ProfileImage, req.Location)
    if err != nil {
        SendError(c, http.StatusInternalServerError, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "profile updated", profile)
}