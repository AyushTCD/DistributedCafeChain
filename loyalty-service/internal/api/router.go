package api

import (
	"loyalty-service/internal/account"
	"loyalty-service/internal/invitation"
	"loyalty-service/internal/transaction"
	"loyalty-service/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


// InitializeRouter setups and returns a new instance of *gin.Engine, including all routes and handlers.
func InitializeRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Initialize services
	userService := user.NewService(db)
	accountService := account.NewService(db)
	transactionService := transaction.NewService(db, accountService)
	invitationService := invitation.NewService(db, userService, accountService)

	// Create the handler with services
	handler := NewHandler(userService, transactionService, accountService, invitationService)

	// Setup route handlers
	handler.SetupRoutes(router)

	return router
}
