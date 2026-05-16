package v1

import (
	"construction_transport_server/api/rest/v1/delivery"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authHandler *delivery.AuthHandler) {
	auth := router.Group("/auth")
	auth.POST("/register", authHandler.Register)
}
