// package routes

// import (
// 	"cineverse/config"
// 	"cineverse/controllers"
// 	"cineverse/middlewares"

// 	"github.com/gin-gonic/gin"
// )

// func SetupRouter() *gin.Engine {
// 	r := gin.Default()

// 	// load HTML templates (templates folder)
// 	r.LoadHTMLGlob("templates/*")

// 	api := r.Group("/api")
// 	{
// 		// Auth
// 		api.POST("/register", controllers.Register(config.DB))
// 		api.POST("/login", controllers.Login(config.DB))

// 		// Public
// 		api.GET("/movies", controllers.GetMovies(config.DB))
// 		api.GET("/movies/:id", controllers.GetMovieDetails(config.DB))
// 		api.GET("/movies/:id/shows", controllers.GetShowsByMovie(config.DB))
// 	}

// 	// Protected user routes
// 	user := r.Group("/api/user").Use(middlewares.AuthMiddleware())
// 	{
// 		user.POST("/book", controllers.BookTickets(config.DB))
// 		user.GET("/mybookings", controllers.GetUserBookings(config.DB))
// 	}

// 	// Admin routes: require auth + admin
// 	admin := r.Group("/admin").Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
// 	{
// 		admin.POST("/movies", controllers.AdminAddMovie(config.DB))
// 		admin.POST("/shows", controllers.AdminAddShow(config.DB))
// 		admin.GET("/bookings", controllers.AdminListBookings(config.DB))
// 		admin.PATCH("/bookings/:id", controllers.AdminUpdateBookingStatus(config.DB))
// 		admin.GET("/dashboard", controllers.AdminDashboard(config.DB))

// 		// Movies
// 		// admin.GET("/movies", controllers.AdminGetMovies)
// 		// admin.POST("/movies", controllers.AdminAddMovie)
// 		// admin.DELETE("/movies/:id", controllers.AdminDeleteMovie)

// 		// // Shows
// 		// admin.GET("/shows", controllers.AdminGetShows)
// 		// admin.POST("/shows", controllers.AdminAddShow)
// 		// admin.DELETE("/shows/:id", controllers.AdminDeleteShow)
// 	}

// 	// public HTML routes (for minimal template front-end)
// 	r.GET("/", func(c *gin.Context) {
// 		c.HTML(200, "public_movies.html", gin.H{})
// 	})
// 	r.GET("/login", func(c *gin.Context) {
// 		c.HTML(200, "login.html", gin.H{})
// 	})
// 	r.GET("/register", func(c *gin.Context) {
// 		c.HTML(200, "register.html", gin.H{})
// 	})

// 	admin.GET("/dashboard/html", func(c *gin.Context) {
// 		c.HTML(200, "admin_dashboard.html", gin.H{})
// 	})

// 	r.GET("/movie/:id", func(c *gin.Context) {
// 		c.HTML(200, "public_movie_details.html", gin.H{"movie_id": c.Param("id")})
// 	})
// 	r.NoRoute(func(c *gin.Context) {
// 		c.JSON(404, gin.H{"error": "page not found"})
// 	})

// 	return r
// }

package routes

import (
	"cineverse/config"
	"cineverse/controllers"
	"cineverse/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Load all HTML templates
	r.LoadHTMLGlob("templates/*")

	// -------------------------------
	// 🔹 Public API Routes
	// -------------------------------
	api := r.Group("/api")
	{
		// Auth routes
		api.POST("/register", controllers.Register(config.DB))
		api.POST("/login", controllers.Login(config.DB))

		// Public movie routes
		api.GET("/movies", controllers.GetMovies(config.DB))
		api.GET("/movies/:id", controllers.GetMovieDetails(config.DB))
		api.GET("/movies/:id/shows", controllers.GetShowsByMovie(config.DB))
	}

	// -------------------------------
	// 🔹 Protected User Routes (Require Login)
	// -------------------------------
	user := r.Group("/api/user").Use(middlewares.AuthMiddleware())
	{
		user.POST("/book", controllers.BookTickets(config.DB))
		user.GET("/mybookings", controllers.GetUserBookings(config.DB))
	}

	// -------------------------------
	// 🔹 Admin Routes (Require Admin Access)
	// -------------------------------
	admin := r.Group("/api/admin").Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.POST("/movies", controllers.AdminAddMovie(config.DB))
		admin.POST("/shows", controllers.AdminAddShow(config.DB))
		admin.GET("/bookings", controllers.AdminListBookings(config.DB))
		admin.PATCH("/bookings/:id", controllers.AdminUpdateBookingStatus(config.DB))
		admin.GET("/dashboard", controllers.AdminDashboard(config.DB))
	}

	// -------------------------------
	// 🔹 Public HTML Pages
	// -------------------------------
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "public_movies.html", gin.H{})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{})
	})

	r.GET("/movie/:id", func(c *gin.Context) {
		c.HTML(200, "public_movie_details.html", gin.H{
			"movie_id": c.Param("id"),
		})
	})

	// -------------------------------
	// 🔹 Admin HTML Page (Protected)
	// -------------------------------
	r.GET("/admin/dashboard", func(c *gin.Context) {
		c.HTML(200, "admin_dashboard.html", gin.H{})
	})

	// -------------------------------
	// 🔹 Fallback for Unknown Routes
	// -------------------------------
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "page not found"})
	})

	return r
}
