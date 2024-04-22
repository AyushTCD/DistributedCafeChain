package transaction

import (
	"context"
	"loyalty-service/internal/account"
	"loyalty-service/internal/model"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service provides methods to interact with transaction data.
type Service struct {
	db         *gorm.DB
	accountSvc *account.Service
}

// NewService creates a new transaction service.
func NewService(db *gorm.DB, accountSvc *account.Service) *Service {
	return &Service{
		db:         db,
		accountSvc: accountSvc,
	}
}

func (s *Service) ProcessTransaction(ctx context.Context, transaction model.Transaction, usePoints bool) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		transactionID, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		transaction.ID = transactionID.String()

		var account model.Account
		err = s.db.First(&account, "account_uuid = ?", transaction.AccountID).Error
		if err != nil {
			return err
		}

		var pointsChange int

		if usePoints {
			pointsRequired := int(math.Ceil(transaction.Amount)) * 10

			if pointsRequired <= account.Points {
				pointsChange = -pointsRequired
			} else {
				pointsChange = -account.Points
			}
		} else {
			pointsChange = int(math.Floor(transaction.Amount))
		}

		account.Points += pointsChange
		transaction.PointsEarned = pointsChange

		err = tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		err = tx.Save(&account).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// // GetTransactionsByUserID retrieves transactions for a specific user from the database
// func (s *Service) GetTransactionsByUserID(ctx context.Context, userID string) ([]model.Transaction, error) {
// 	var transactions []model.Transaction

// 	err := s.db.WithContext(ctx).Where("account_uuid = ?", userID).Find(&transactions).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return transactions, nil
// }
