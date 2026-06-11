package v1

import (
	"construction_transport_server/api/rest/v1/delivery"
	"construction_transport_server/api/rest/v1/middleware"
	"construction_transport_server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authHandler *delivery.AuthHandler, vehicleHandler *delivery.VehicleHandler, bookingHandler *delivery.BookingHandler, jobHandler *delivery.JobHandler) {
	// public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// protected routes
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(utils.JWTManager))
	{
		// profile (to be implemented)
		api.GET("/profile", profileHandler.Get)
		api.PUT("/profile", profileHandler.Update)

		// vehicle (transporter only)
		vehicles := api.Group("/vehicles")
		vehicles.Use(middleware.RoleMiddleware("TRANSPORTER"))
		{
			vehicles.POST("/", vehicleHandler.Create)
			vehicles.GET("/", vehicleHandler.List)
			vehicles.GET("/:id", vehicleHandler.Get)
			vehicles.PUT("/:id", vehicleHandler.Update)
			vehicles.DELETE("/:id", vehicleHandler.Delete)
		}

		// bookings (customer)
		bookings := api.Group("/bookings")
		bookings.Use(middleware.RoleMiddleware("USER"))
		{
			bookings.POST("/", bookingHandler.Create)
			bookings.GET("/", bookingHandler.ListCustomer)
			bookings.GET("/:id", bookingHandler.Get)
			bookings.POST("/:id/cancel", bookingHandler.Cancel)
		}

		// jobs (transporter)
		jobs := api.Group("/jobs")
		jobs.Use(middleware.RoleMiddleware("TRANSPORTER"))
		{
			jobs.GET("/", jobHandler.ListMyJobs)
			jobs.PATCH("/:id/status", jobHandler.UpdateStatus)
		}
	}
}
