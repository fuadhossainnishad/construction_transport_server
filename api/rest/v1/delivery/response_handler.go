package delivery

import (
	"construction_transport_server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode,
		utils.ResponseSuccess(statusCode, message, data),
	)
}

func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode,
		utils.ResponseError(utils.AppError{
			StatusCode: statusCode,
			Message:    message,
		}),
	)
}
