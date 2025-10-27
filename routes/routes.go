package routes

import (
	"cineverse/config"
	"cineverse/controllers"
	"cineverse/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// load HTML templates (templates folder)
	r.LoadHTMLGlob("templates/*")

	api := r.Group("/api")
	{
		// Auth
		api.POST("/register", controllers.Register(config.DB))
		api.POST("/login", controllers.Login(config.DB))

		// Public
		api.GET("/movies", controllers.GetMovies(config.DB))
		api.GET("/movies/:id", controllers.GetMovieDetails(config.DB))
		api.GET("/movies/:id/shows", controllers.GetShowsByMovie(config.DB))
	}

	// Protected user routes
	user := r.Group("/api").Use(middlewares.AuthMiddleware())
	{
		user.POST("/book", controllers.BookTickets(config.DB))
		user.GET("/mybookings", controllers.GetUserBookings(config.DB))
	}

	// Admin routes: require auth + admin
	admin := r.Group("/admin").Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.POST("/movies", controllers.AdminAddMovie(config.DB))
		admin.POST("/shows", controllers.AdminAddShow(config.DB))
		admin.GET("/bookings", controllers.AdminListBookings(config.DB))
		admin.PATCH("/bookings/:id", controllers.AdminUpdateBookingStatus(config.DB))
		admin.GET("/dashboard", controllers.AdminDashboard(config.DB))
	}

	// public HTML routes (for minimal template front-end)
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "public_movies.html", gin.H{})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{})
	})

	return r
}
