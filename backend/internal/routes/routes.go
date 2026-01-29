package routes

import (
	"website-dummy/internal/controllers"
	"website-dummy/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// Initialize controllers
	userController := controllers.NewUserController()
	healthController := controllers.NewHealthController()
	authController := controllers.NewAuthController()
	propertyController := controllers.NewPropertyController()
	contactRequestController := controllers.NewContactRequestController()
	rentalController := controllers.NewRentalController()
	articleController := controllers.NewArticleController()
	testimonialController := controllers.NewTestimonialController()
	sitemapController := controllers.NewSitemapController()

	// Public routes (no authentication required)
	r.GET("/api/health", healthController.HealthCheck)
	r.GET("/sitemap.xml", sitemapController.GetSitemap)
	r.POST("/api/auth/login", authController.Login)
	r.GET("/api/auth/logout", authController.Logout)
	r.GET("/api/featured-properties", propertyController.GetFeaturedProperties)
	r.GET("/api/v1/properties/initial", propertyController.GetInitialProperties)
	r.POST("/api/v1/properties", propertyController.QueryProperties)
	r.GET("/api/v1/properties", propertyController.QueryProperties)
	r.GET("/api/neighborhoods/:neighborhood", propertyController.GetNeighborhoodProperties)
	
	// Contact requests (public)
	r.POST("/api/contact-requests", contactRequestController.CreateContactRequest)
	
	// Articles (public endpoints)
	r.GET("/api/articles", articleController.GetArticles)
	r.GET("/api/articles/:id", articleController.GetArticle)
	
	// Testimonials (public endpoints)
	r.GET("/api/testimonials", testimonialController.GetTestimonials)
	
	// Rentals (public endpoints)
	r.GET("/api/rentals/public", rentalController.GetPublicRentals)
	r.GET("/api/rentals/public/with-properties", rentalController.GetPublicRentalsWithProperties)

	// Protected API routes (require authentication)
	api := r.Group("/api")
	api.Use(middleware.Auth0Middleware())
	{
		// User profile
		api.GET("/profile", authController.Profile)
		api.PUT("/auth/change-password", authController.ChangePassword)
		
		// Admin routes
		admin := api.Group("/admin")
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Admin dashboard", "user": c.MustGet("user")})
			})
		}
		
		// Contact requests (protected)
		api.GET("/contact-requests", contactRequestController.GetContactRequests)
		api.GET("/contact-requests/:id", contactRequestController.GetContactRequest)
		api.DELETE("/contact-requests/:id", contactRequestController.DeleteContactRequest)
		
		// Rentals (protected)
		api.GET("/rentals", rentalController.GetRentals)
		api.POST("/rentals", rentalController.CreateRental)
		api.PUT("/rentals/:id", rentalController.UpdateRental)
		api.PUT("/rentals/orders", rentalController.BulkUpdateRentalOrders)
		api.DELETE("/rentals/:id", rentalController.DeleteRental)
		api.GET("/rentals/search", rentalController.SearchProperties)
		
		// Articles (protected admin endpoints)
		api.GET("/articles/admin/list", articleController.GetArticles)
		api.GET("/articles/admin/:id", articleController.GetArticle)
		api.POST("/articles", articleController.CreateArticle)
		api.PUT("/articles/:id", articleController.UpdateArticle)
		api.DELETE("/articles/:id", articleController.DeleteArticle)
		api.POST("/articles/upload-image", articleController.UploadImage)
		api.GET("/articles/:id/files", articleController.GetArticleFiles)
		api.DELETE("/articles/:id/files", articleController.RemoveImage)
		api.DELETE("/articles/images/*path", articleController.DeleteImage)
		
		// Testimonials (protected admin endpoints)
		api.GET("/admin/testimonials", testimonialController.GetTestimonialsAdmin)
		api.GET("/admin/testimonials/:id", testimonialController.GetTestimonial)
		api.POST("/admin/testimonials", testimonialController.CreateTestimonial)
		api.PUT("/admin/testimonials/:id", testimonialController.UpdateTestimonial)
		api.DELETE("/admin/testimonials/:id", testimonialController.DeleteTestimonial)
	}

	// Legacy routes (keeping for compatibility)
	r.GET("/api/name", userController.GetName)
}
