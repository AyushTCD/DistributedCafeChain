package api

import (
	"log"
	"net/http"
	"time"

	"loyalty-service/internal/account"
	"loyalty-service/internal/invitation"
	"loyalty-service/internal/transaction"
	"loyalty-service/internal/user"

	"loyalty-service/internal/model"

	"github.com/gin-gonic/gin"
)

// Handler struct centralizes dependencies for HTTP handlers.
type Handler struct {
	userService        *user.Service
	transactionService *transaction.Service
	accountService     *account.Service
	invitationService  *invitation.Service
}

// NewHandler is the constructor for Handler.
func NewHandler(userSvc *user.Service, transactionSvc *transaction.Service, accountSvc *account.Service, invitationSvc *invitation.Service) *Handler {
	return &Handler{
		userService:        userSvc,
		transactionService: transactionSvc,
		accountService:     accountSvc,
		invitationService:  invitationSvc,
	}
}

// SetupRoutes defines all application's routes.
func (h *Handler) SetupRoutes(router *gin.Engine) {
	// User account management
	router.POST("/users", h.RegisterUser) // Register a new user
	router.GET("/users/:id", h.GetUser)   // Retrieve user details
	// router.PUT("/users/:id", h.UpdateUser) // Update user details

	// Managing loyalty-card accounts (Linking family and friends)
	router.POST("/loyalty-accounts", h.CreateLoyaltyAccount) // Create a new loyalty account
	// router.PUT("/loyalty-accounts/:id", h.AddUserToLoyaltyAccount)  // Add a user to an existing loyalty account
	router.GET("/loyalty-accounts/:id", h.GetLoyaltyAccountDetails) // Get details of a loyalty account

	// Transaction history
	router.POST("/transactions", h.ProcessTransaction) // Log a new transaction
	// router.GET("/users/:id/transactions", h.GetUserTransactions) // Retrieve a user's transaction history

	// Invitations
	router.POST("invitations/create", h.CreateInvitation)
	router.POST("/invitations/accept", h.AcceptInvitation)
	router.POST("/invitations/decline", h.DeclineInvitation)
}

// Register a new user
func (h *Handler) RegisterUser(c *gin.Context) {
	var newUser model.User

	// Parse the JSON request body into newUser
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Basic validation (you might want to expand this)
	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Call your service's createUser function and retrieve the newly created user
	createdUser, err := h.userService.CreateUser(c.Request.Context(), newUser)
	if err != nil {
		// Handle specific error cases here as needed
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Prepare the response, excluding the password field
	userResponse := struct {
		ID           string    `json:"id"`
		Name         string    `json:"name"`
		Email        string    `json:"email"`
		Phone        string    `json:"phone,omitempty"` // Include other fields as necessary
		CreationDate time.Time `json:"creationDate"`
	}{
		ID:           createdUser.ID,
		Name:         createdUser.Name,
		Email:        createdUser.Email,
		Phone:        createdUser.Phone,
		CreationDate: createdUser.CreationDate,
	}

	// Return success message with created user object
	c.JSON(http.StatusCreated, userResponse)
}

// // Fetch transactions for a user.
// func (h *Handler) GetUserTransactions(c *gin.Context) {
// 	// Placeholder implementation. Extract user ID from path, validate it, call transactionService to fetch transactions, return the result.
// 	userID := c.Param("id")
// 	c.JSON(http.StatusOK, gin.H{"userID": userID, "transactions": "TODO: Need to implement"})
// }

// // Logs a new transaction to a user's account.
// func (h *Handler) LogTransaction(c *gin.Context) {
// 	// Implementation goes here.
// 	c.JSON(http.StatusCreated, gin.H{"message": "TODO: Need to implement"})
// }

// GetUser handles fetching user details by ID.
func (h *Handler) GetUser(c *gin.Context) {
	// Extract the user ID from the URL path
	userID := c.Param("id")

	// Use the userService to fetch the user by their ID
	user, err := h.userService.GetUserByID(c, userID)
	if err != nil {
		// If the user is not found or there's another error, return an appropriate response
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// If the user is found, return their details
	c.JSON(http.StatusOK, gin.H{
		"ID":           user.ID,
		"AccountID":    user.AccountID,
		"Name":         user.Name,
		"Email":        user.Email,
		"Phone":        user.Phone,
		"CreationDate": user.CreationDate,
	})
}

// // Update user details
// func (h *Handler) UpdateUser(c *gin.Context) {
// 	// Implementation goes here.
// 	c.JSON(http.StatusCreated, gin.H{"message": "TODO: Need to implement"})
// }

// // Add a user to an existing loyalty account
// func (h *Handler) AddUserToLoyaltyAccount(c *gin.Context) {
// 	// Implementation goes here.
// 	c.JSON(http.StatusCreated, gin.H{"message": "TODO: Need to implement"})
// }

// CreateLoyaltyAccount handles the creation of a new loyalty account, associating it with users and setting initial points.
func (h *Handler) CreateLoyaltyAccount(c *gin.Context) {
	var request struct {
		UserIDs []string `json:"userIds"` // Array of user IDs to associate with the account
		Points  int      `json:"points"`  // Initial points to assign to the account
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	acc := model.Account{
		Points: request.Points,
	}

	createdAccount, err := h.accountService.CreateAccount(c.Request.Context(), acc, request.UserIDs, request.Points)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdAccount)
}

// Get details of a loyalty account
func (h *Handler) GetLoyaltyAccountDetails(c *gin.Context) {
	accountID := c.Param("id")
	acc, err := h.accountService.GetAccount(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(200, &acc)
}

// Process a new transaction and either add points or use points based on the transaction details.
func (h *Handler) ProcessTransaction(c *gin.Context) {
	var trans model.Transaction

	if err := c.ShouldBindJSON(&trans); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	usePoints := c.Query("usePoints") == "true"

	if err := h.transactionService.ProcessTransaction(c.Request.Context(), trans, usePoints); err != nil {
		log.Printf("Error processing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction processed successfully", "id": trans.ID})
}

func (h *Handler) CreateInvitation(c *gin.Context) {
	var req struct {
		Email     string `json:"email"`     // Email of the person being invited
		InviterID string `json:"inviterID"` // ID of the person sending the invitation
		AccountID string `json:"accountID"` // ID of the account/group to which the invitee is being invited
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// `CreateInvitation` equires the email of the invitee, the inviterID, and the accountID.
	CreatedInvitation, err := h.invitationService.CreateInvitation(c.Request.Context(), req.Email, req.InviterID, req.AccountID)
	if err != nil {
		// Handle errors as appropriate for your application
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invitation"})
		return
	}

	// Prepare the response, excluding the password field
	invitationResponse := struct {
		InvitationUUID string    `json:"invitationID"`
		Email          string    `json:"email"`
		AccountUUID    string    `json:"accountID"`
		InviterUUID    string    `json:"inviterID"`
		Token          string    `json:"token"`
		CreationDate   time.Time `json:"creationDate"`
		ExpirationDate time.Time `json:"expirationDate"`
		Status         string    `json:"status"`
	}{
		InvitationUUID: CreatedInvitation.InvitationUUID,
		Email:          CreatedInvitation.Email,
		AccountUUID:    CreatedInvitation.AccountUUID,
		InviterUUID:    CreatedInvitation.InviterUUID,
		Token:          CreatedInvitation.Token,
		CreationDate:   CreatedInvitation.CreationDate,
		ExpirationDate: CreatedInvitation.ExpirationDate,
		Status:         CreatedInvitation.Status,
	}

	c.JSON(http.StatusCreated, invitationResponse)
}

func (h *Handler) AcceptInvitation(c *gin.Context) {
	var req struct {
		Token        string `json:"token"`
		InviteeEmail string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err := h.invitationService.AcceptInvitation(c.Request.Context(), req.Token, req.InviteeEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation accepted successfully"})
}

func (h *Handler) DeclineInvitation(c *gin.Context) {
	var req struct {
		Token        string `json:"token"`
		InviteeEmail string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err := h.invitationService.DeclineInvitation(c.Request.Context(), req.Token, req.InviteeEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation declined successfully"})
}
