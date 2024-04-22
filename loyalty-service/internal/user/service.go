package user

import (
	"context"
	"loyalty-service/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service provides methods to interact with user data.
type Service struct {
	db *gorm.DB
}

// NewService creates a new user service.
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// CreateUser creates a new user in the database.
func (s *Service) CreateUser(ctx context.Context, u model.User) (*model.User, error) {
	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	u.ID = userID.String()

	// Hash the user's password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashedPassword)

	// Create user in the database
	if err := s.db.WithContext(ctx).Create(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserByID retrieves a user by their ID from the database.
func (s *Service) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	return &user, s.db.WithContext(ctx).First(&user, "user_uuid = ?", userID).Error
}

// // GetUserByEmail retrieves a user by their email from the database.
// func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
// 	var user model.User
// 	err := s.db.WithContext(ctx).Where("email_address = ?", email).First(&user).Error
// 	return &user, err
// }
