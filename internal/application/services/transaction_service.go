package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"transaction-service/internal/domain/entities"
	"transaction-service/internal/domain/repositories"

	"github.com/shopspring/decimal"
)

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrDuplicateTransaction    = errors.New("transaction already processed")
	ErrInvalidAmount           = errors.New("invalid amount")
	ErrInvalidTransactionState = errors.New("invalid transaction state")
	ErrInvalidSourceType       = errors.New("invalid source type")
)

// TransactionService handles transaction business logic
type TransactionService struct {
	userRepo        repositories.UserRepository
	transactionRepo repositories.TransactionRepository
}

// NewTransactionService creates a new TransactionService
func NewTransactionService(
	userRepo repositories.UserRepository,
	transactionRepo repositories.TransactionRepository,
) *TransactionService {
	return &TransactionService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

// ProcessTransaction processes a new transaction
func (s *TransactionService) ProcessTransaction(
	ctx context.Context,
	userID uint64,
	req entities.TransactionRequest,
	sourceType entities.SourceType,
) error {
	// Validating a source type
	if !sourceType.IsValid() {
		return ErrInvalidSourceType
	}

	// Checking for duplicate transactions
	exists, err := s.transactionRepo.ExistsByTransactionID(ctx, req.TransactionID)
	if err != nil {
		return fmt.Errorf("failed to check transaction existence: %w", err)
	}
	if exists {
		return ErrDuplicateTransaction
	}

	// Parse and validate the amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return ErrInvalidAmount
	}
	if amount.IsNegative() || amount.IsZero() {
		return ErrInvalidAmount
	}

	// Validating transaction state
	state := entities.TransactionState(req.State)
	if !state.IsValid() {
		return ErrInvalidTransactionState
	}

	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Calculate new balance
	var newBalance decimal.Decimal
	switch state {
	case entities.StateWin:
		newBalance = user.Balance.Add(amount)
	case entities.StateLose:
		newBalance = user.Balance.Sub(amount)
		if newBalance.IsNegative() {
			return ErrInsufficientFunds
		}
	}

	// Create transaction record
	transaction := &entities.Transaction{
		UserID:        userID,
		TransactionID: req.TransactionID,
		State:         state,
		Amount:        amount,
		SourceType:    sourceType,
		CreatedAt:     time.Now(),
	}

	// Save the transaction
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update user balance
	if err := s.userRepo.UpdateBalance(ctx, userID, newBalance); err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	return nil
}

// GetUserBalance retrieves the current user balance
func (s *TransactionService) GetUserBalance(
	ctx context.Context,
	userID uint64,
) (*entities.BalanceResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &entities.BalanceResponse{
		UserID:  user.ID,
		Balance: user.Balance.StringFixed(2),
	}, nil
}
