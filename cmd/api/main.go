package main

import (
	"github.com/Zosimov735/go-chat/internal/app"
	"github.com/Zosimov735/go-chat/internal/db"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	// Database initialization
	db, err := db.Initialize()
	if err != nil {
		log.Fatal("Could not set up database: ", err)
		return
	}
	defer db.Close()

	// Setting up the application with dependencies
	appHandler := app.NewHandler(db)
	appHandler.RegisterRoutes(router)

	// Starting the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
