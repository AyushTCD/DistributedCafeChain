package invitation

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"loyalty-service/internal/account"
	"loyalty-service/internal/model"
	"loyalty-service/internal/user"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	userSvc    *user.Service
	accountSvc *account.Service
}

func NewService(db *gorm.DB, userSvc *user.Service, accountSvc *account.Service) *Service {
	return &Service{
		db:         db,
		userSvc:    userSvc,
		accountSvc: accountSvc,
	}
}

// GetInvitationByToken looks up an invitation by its token
func (s *Service) GetUserByInvite(ctx context.Context, token string) (*model.User, error) {
	var user model.User

	err := s.db.WithContext(ctx).Where("invite_code = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) CreateInvitation(ctx context.Context, email, inviterID, accountID string) (*model.Invitation, error) {
	// Check if the inviter is part of the specified account
	var inviter model.User
	result := s.db.WithContext(ctx).Where("user_uuid = ? AND account_uuid = ?", inviterID, accountID).First(&inviter)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to verify inviter: %w", result.Error)
	}

	invitationId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// Generate a unique token for the invitation
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invitation token: %w", err)
	}

	// Set the expiration date for the invitation (e.g., 48 hours from now)
	expirationDate := time.Now().Add(48 * time.Hour)

	// Create the invitation record
	invitation := model.Invitation{
		InvitationUUID: invitationId.String(),
		Email:          email,
		InviterUUID:    inviterID,
		AccountUUID:    accountID,
		Token:          token,
		CreationDate:   time.Now(),
		ExpirationDate: expirationDate,
		Status:         "pending",
	}

	result = s.db.WithContext(ctx).Create(&invitation)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", result.Error)
	}

	return &invitation, nil
}
func generateToken() (string, error) {
	// Simple token generation
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 16)
	for i := range b {
		val, err := rand.Int(rand.Reader, big.NewInt(52))
		if err != nil {
			return "", err
		}

		b[i] = letters[val.Int64()]
	}

	return string(b), nil
}

func (s *Service) AcceptInvitation(ctx context.Context, token string, email string) error {
	// Find the invitation by token and ensure it's valid
	var invitation model.Invitation
	err := s.db.WithContext(ctx).Where("token = ? AND email = ?", token, email).First(&invitation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invitation not found or does not match email")
		}
		return err
	}

	// Ensure the invitation is still valid (not expired and status is pending)
	if invitation.Status != "pending" || invitation.ExpirationDate.Before(time.Now()) {
		return errors.New("invitation is not valid or has expired")
	}

	// Find user by email and update their account_uuid to the one in the invitation
	var user model.User
	err = s.db.WithContext(ctx).Where("email_address = ?", email).First(&user).Error
	if err != nil {
		return err
	}

	user.AccountID = &invitation.AccountUUID

	// Update user account id and save the user
	if err = s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user's account: %w", err)
	}

	// Mark invitation as accepted
	invitation.Status = "accepted"
	if err := s.db.WithContext(ctx).Save(&invitation).Error; err != nil {
		return fmt.Errorf("failed to update invitation status: %w", err)
	}

	return nil
}

func (s *Service) DeclineInvitation(ctx context.Context, token string, email string) error {
	// Find the invitation by token and ensure it matches the email
	var invitation model.Invitation
	err := s.db.WithContext(ctx).Where("token = ? AND email = ?", token, email).First(&invitation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invitation not found or does not match email")
		}
		return err
	}

	// Ensure the invitation is still valid (not expired)
	if invitation.ExpirationDate.Before(time.Now()) {
		return errors.New("invitation has expired")
	}

	// Mark invitation as declined
	invitation.Status = "declined"
	if err := s.db.WithContext(ctx).Save(&invitation).Error; err != nil {
		return fmt.Errorf("failed to update invitation status to declined: %w", err)
	}

	return nil
}
