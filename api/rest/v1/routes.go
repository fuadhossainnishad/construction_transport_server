package v1

import (
	"construction_transport_server/api/rest/v1/delivery"
	"construction_transport_server/api/rest/v1/middleware"
	"construction_transport_server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authHandler *delivery.AuthHandler,
	profileHandler *delivery.AccountHandler,
	vehicleHandler *delivery.VehicleHandler,
	bookingHandler *delivery.BookingHandler,
	jobHandler *delivery.JobHandler,
	jwtManager *utils.JWTManager,
) {
	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// Protected routes
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Profile (any authenticated user)
		api.GET("/profile", accountHandler.GetProfile)
		api.PUT("/profile", accountHandler.UpdateProfile)

		// Vehicles (transporter only)
		vehicles := api.Group("/vehicles")
		vehicles.Use(middleware.RoleMiddleware("TRANSPORTER"))
		{
			vehicles.POST("/", vehicleHandler.Create)
			vehicles.GET("/", vehicleHandler.List)
			vehicles.GET("/:id", vehicleHandler.Get)
			vehicles.PUT("/:id", vehicleHandler.Update)
			vehicles.DELETE("/:id", vehicleHandler.Delete)
		}

		// Bookings (customer only)
		bookings := api.Group("/bookings")
		bookings.Use(middleware.RoleMiddleware("USER"))
		{
			bookings.POST("/", bookingHandler.Create)
			bookings.GET("/", bookingHandler.ListCustomer)
			bookings.GET("/:id", bookingHandler.Get)
			bookings.POST("/:id/cancel", bookingHandler.Cancel)
		}

		// Jobs (transporter only)
		jobs := api.Group("/jobs")
		jobs.Use(middleware.RoleMiddleware("TRANSPORTER"))
		{
			jobs.GET("/", jobHandler.ListMyJobs)
			jobs.GET("/:id", jobHandler.GetJobDetails)
			jobs.PATCH("/:id/status", jobHandler.UpdateStatus)
		}
	}
}
