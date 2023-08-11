package main

import (
	"awesomeProject/database"
	_ "awesomeProject/docs"
	"awesomeProject/handlers"
	"awesomeProject/middlewares"
	"awesomeProject/models"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"log"
	"os"
	"sync"
)

func main() {
	// Load values from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a MySQL database connection using gorm
	dsn := os.Getenv("DB_DSN")
	db, err := database.InitializeDB(dsn)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Get posts for user with ID=7
	posts, err := models.GetPostsForUser(7)
	if err != nil {
		log.Fatal("Error getting posts:", err)
	}

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start a goroutine for each post to fetch comments
	for _, post := range posts {
		wg.Add(1)
		go func(p models.Post) {
			defer wg.Done()
			err := models.InsertPost(db, p)
			if err != nil {
				log.Printf("Error inserting post %d: %s", p.ID, err)
			}
			comments, err := models.GetCommentsForPost(p.ID)
			if err != nil {
				log.Printf("Error getting comments for post %d: %s", p.ID, err)
				return
			}

			// Create a wait group to wait for all comment writing goroutines to finish for this post
			var commentWg sync.WaitGroup

			// Start a goroutine for each comment to write comments to the database
			for _, comment := range comments {
				commentWg.Add(1)
				go func(c models.Comment) {
					defer commentWg.Done()
					err := models.InsertComment(db, c)
					if err != nil {
						log.Printf("Error inserting comment %d for post %d: %s", c.ID, c.PostID, err)
					}
				}(comment)
			}

			// Wait for all comment writing goroutines to finish for this post
			commentWg.Wait()
		}(post)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("All comments and posts have been processed and inserted into the database.")

	// Initialize Echo and router
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Apply JWT middleware for route protection
	apiGroup := e.Group("/api", middlewares.JwtMiddleware)
	apiGroup.POST("/posts", handlers.CreatePostHandler(db))
	apiGroup.POST("/comments", handlers.CreateCommentHandler(db))

	// Routes for posts
	e.GET("/posts", handlers.GetPostsHandler(db))
	e.GET("/posts/:id", handlers.GetPostHandler(db))
	e.PUT("/posts/:id", handlers.UpdatePostHandler(db))
	e.DELETE("/posts/:id", handlers.DeletePostHandler(db))

	// Routes for comments
	e.GET("/comments", handlers.GetCommentsHandler(db))
	e.GET("/comments/:id", handlers.GetCommentHandler(db))
	e.PUT("/comments/:id", handlers.UpdateCommentHandler(db))
	e.DELETE("/comments/:id", handlers.DeleteCommentHandler(db))

	// User registration and login routes
	e.POST("/register", handlers.RegisterHandler(db))
	e.POST("/login", handlers.LoginHandler(db))

	// Google OAuth2 authentication
	e.GET("/auth/google", handlers.HandleGoogleAuth)
	e.GET("/auth/google/callback", handlers.HandleGoogleCallback(db))

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
