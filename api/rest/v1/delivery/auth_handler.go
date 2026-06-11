package delivery

import (
	"construction_transport_server/api/rest/v1/dto"
	"construction_transport_server/internal/auth/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	register_usecase *usecase.RegisteredUsecase
	login_usecase    *usecase.LoginUseCase
}

func NewAuthHandler(register_usecase *usecase.RegisteredUsecase, login_usecase *usecase.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		register_usecase: register_usecase,
		login_usecase:    login_usecase,
	}
}

func (handler *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	err := handler.register_usecase.Execute(ctx, usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(c, http.StatusOK, "OTP_SENT", gin.H{
		"message": "OTP has been sent to your email, please verify your account",
	})

}

func (handler *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        SendError(c, http.StatusBadRequest, "invalid request")
        return
    }
    resp, err := handler.login_usecase.Execute(c.Request.Context(), usecase.LoginInput{
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil {
        SendError(c, http.StatusUnauthorized, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "login success", resp)
}

func (handler *AuthHandler) VerifyOTP(c *gin.Context) {
    var req struct {
        Email string `json:"email"`
        OTP   string `json:"otp"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        SendError(c, http.StatusBadRequest, "invalid request")
        return
    }
    err := handler.otpService.VerifyOTP(c.Request.Context(), req.Email, req.OTP)
    if err != nil {
        SendError(c, http.StatusBadRequest, err.Error())
        return
    }
    // mark email verified in DB
    if err := handler.authRepo.VerifyEmail(c.Request.Context(), req.Email); err != nil {
        SendError(c, http.StatusInternalServerError, "failed to verify email")
        return
    }
    SendResponse(c, http.StatusOK, "email verified", nil)
}

// ForgotPassword: generate OTP and send
func (handler *AuthHandler) ForgotPassword(c *gin.Context) {
    var req struct { Email string `json:"email"` }
    if err := c.ShouldBindJSON(&req); err != nil {
        SendError(c, http.StatusBadRequest, "invalid email")
        return
    }
    // check user exists
    _, err := handler.authRepo.GetAuth(c.Request.Context(), req.Email)
    if err != nil {
        SendError(c, http.StatusNotFound, "user not found")
        return
    }
    if err := handler.otpService.GenerateAndSendOTP(c.Request.Context(), req.Email); err != nil {
        SendError(c, http.StatusInternalServerError, "could not send OTP")
        return
    }
    SendResponse(c, http.StatusOK, "password reset OTP sent", nil)
}

// ResetPassword: verify OTP + update password
func (handler *AuthHandler) ResetPassword(c *gin.Context) {
    var req struct {
        Email       string `json:"email"`
        OTP         string `json:"otp"`
        NewPassword string `json:"new_password"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        SendError(c, http.StatusBadRequest, "invalid input")
        return
    }
    if err := handler.otpService.VerifyOTP(c.Request.Context(), req.Email, req.OTP); err != nil {
        SendError(c, http.StatusBadRequest, "invalid OTP")
        return
    }
    hashed, err := handler.hashFunc.Hash(req.NewPassword)
    if err != nil {
        SendError(c, http.StatusInternalServerError, "could not hash password")
        return
    }
    if err := handler.authRepo.UpdatePassword(c.Request.Context(), req.Email, hashed); err != nil {
        SendError(c, http.StatusInternalServerError, "failed to update password")
        return
    }
    SendResponse(c, http.StatusOK, "password reset successful", nil)
}

// Refresh token endpoint
func (handler *AuthHandler) Refresh(c *gin.Context) {
    var req struct { RefreshToken string `json:"refresh_token"` }
    if err := c.ShouldBindJSON(&req); err != nil {
        SendError(c, http.StatusBadRequest, "missing refresh token")
        return
    }
    accessToken, newRefreshToken, err := handler.refreshUsecase.Refresh(c.Request.Context(), req.RefreshToken)
    if err != nil {
        SendError(c, http.StatusUnauthorized, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "token refreshed", gin.H{
        "access_token":  accessToken,
        "refresh_token": newRefreshToken,
    })
}