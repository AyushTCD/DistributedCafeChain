package main

import (
	"log"
	"loyalty-service/internal/account"
	"loyalty-service/internal/api"
	"loyalty-service/internal/invitation"
	"loyalty-service/internal/transaction"
	"loyalty-service/internal/user"
	"loyalty-service/pkg/db"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
)

type ConfigRegion []string
type Config map[string]ConfigRegion

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cfgData, err := os.ReadFile(cwd + "/loyalty-service.toml")
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = toml.Unmarshal(cfgData, &cfg)
	if err != nil {
		panic(err)
	}

	// Connect to MySQL
	database, err := db.Connect(cfg["default"])
	if err != nil {
		panic(err)
	}

	// Initialize services with the database
	userService := user.NewService(database)
	accountService := account.NewService(database)
	transactionService := transaction.NewService(database, accountService)
	invitationService := invitation.NewService(database, userService, accountService)

	// Set up Gin router and routes
	router := gin.Default()

	// Initialize the handler with the services
	handler := api.NewHandler(userService, transactionService, accountService, invitationService)

	// Setup routes using the handler
	handler.SetupRoutes(router)

	// Start the HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
