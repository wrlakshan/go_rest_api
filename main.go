package main

import (
	"log"

	"go_rest_api/cmd/api/handlers"
	"go_rest_api/internal/auth"
	authMiddleware "go_rest_api/internal/middleware"
	"go_rest_api/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Database connection setup  
	db, err := storage.NewMySQLConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close() 

	userHandler := handlers.NewUserHandler(db)
	productHandler := handlers.NewProductHandler(db)

	e.GET("/health", handlers.HealthCheckHandler)

	e.POST("/login", auth.Login(db)) 
	e.POST("/signup", userHandler.SignUp)  

	// User routes
	userGroup := e.Group("/users")
	userGroup.Use(authMiddleware.RequireAuth)
	userGroup.GET("", userHandler.GetUsers, authMiddleware.RequireRole("admin"))
	userGroup.POST("", userHandler.CreateUser, authMiddleware.RequireRole("admin"))
	userGroup.GET("/:id", userHandler.GetUser)
	userGroup.PUT("/:id", userHandler.UpdateUser,authMiddleware.RequireRole("admin"))
	userGroup.DELETE("/:id", userHandler.DeleteUser, authMiddleware.RequireRole("admin"))

	// Product routes
	productGroup := e.Group("/products")
	productGroup.Use(authMiddleware.RequireAuth)
	productGroup.GET("", productHandler.GetProducts)
	productGroup.POST("", productHandler.CreateProduct, authMiddleware.RequireRole("admin"))
	productGroup.GET("/:id", productHandler.GetProduct)
	productGroup.PUT("/:id", productHandler.UpdateProduct, authMiddleware.RequireRole("admin"))
	productGroup.DELETE("/:id", productHandler.DeleteProduct, authMiddleware.RequireRole("admin"))

	e.Logger.Fatal(e.Start(":8080"))
}
