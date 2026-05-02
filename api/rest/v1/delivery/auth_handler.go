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
