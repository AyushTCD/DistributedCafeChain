package account

import (
	"context"
	"errors"
	"loyalty-service/internal/model"

	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service provides methods for account management
type Service struct {
	db *gorm.DB
}

// NewService creates a new account service with a MongoDB collection
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// CreateAccount adds a new loyalty group account to the database along with associating users and allocating points.
func (s *Service) CreateAccount(ctx context.Context, account model.Account, userIds []string, points int) (*model.Account, error) {
	// Generate a new UUID for the account
	accountID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	account.ID = accountID.String()
	account.Points = points // Set initial points for the account

	// Begin a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Attempt to create the account in the database
	if err := tx.WithContext(ctx).Create(&account).Error; err != nil {
		tx.Rollback() // Roll back the transaction on error
		return nil, err
	}

	// Associate users with the account and update their points
	for _, userID := range userIds {
		var user model.User
		// Find the user by their ID
		if err := tx.WithContext(ctx).Where("user_uuid = ?", userID).First(&user).Error; err != nil {
			tx.Rollback() // Roll back the transaction on error
			return nil, err
		}

		// Associate the user with the account
		user.AccountID = &account.ID

		// Update the user's record in the database
		if err := tx.WithContext(ctx).Save(&user).Error; err != nil {
			tx.Rollback() // Roll back the transaction on error
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Return the newly created account
	return &account, nil
}

// GetAccount retrieves an account by its ID
func (s *Service) GetAccount(ctx context.Context, accountID string) (*model.Account, error) {
	var account model.Account
	err := s.db.WithContext(ctx).First(&account, "account_uuid = ?", accountID).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// AddPoints increments points for a user's account based on the transaction amount
func (s *Service) AddPoints(ctx context.Context, userID string, points int) error {
	// Convert the userID string to Int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		// Handle invalid userID format
		return err
	}

	// Begin transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// Check if the user exists
	var user model.User
	if err := tx.First(&user, "id = ?", userIDInt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Update points
	user.Account.Points += points
	if err := tx.Save(&user).Error; err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// SubtractPoints subtracts points from a user's loyalty account
func (s *Service) SubtractPoints(ctx context.Context, userID string, pointsToSubtract int) error {
	if pointsToSubtract <= 0 {
		return errors.New("points to subtract must be positive")
	}

	// Convert the userID string to Int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		// Handle invalid userID format
		return err
	}

	// Begin transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// Check if the user exists
	var user model.User
	if err := tx.First(&user, "id = ?", userIDInt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if the user has enough points
	if user.Account.Points < pointsToSubtract {
		// Update points
		user.Account.Points += pointsToSubtract
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
	} else {
		return errors.New("insufficient points to subtract")
	}

	// Commit transaction
	return tx.Commit().Error
}

// AddUserToAccount adds a user to an account by updating the Account model
func (s *Service) AddUserToAccount(ctx context.Context, userID, accountID string) error {
	// Convert userID and accountID to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	accountIDInt, err := strconv.Atoi(accountID)
	if err != nil {
		return errors.New("invalid account ID format")
	}

	// Find the account
	var account model.Account
	if err := s.db.First(&account, "id = ?", accountIDInt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("account not found")
		}
		return err
	}

	// Find the user
	var user model.User
	if err := s.db.First(&user, "id = ?", userIDInt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Add user to account
	if err := s.db.Model(&account).Association("Users").Append(&user); err != nil {
		return err
	}

	return nil
}
