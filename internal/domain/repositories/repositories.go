package repositories

import (
	"context"

	"transaction-service/internal/domain/entities"

	"github.com/shopspring/decimal"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByID(ctx context.Context, userID uint64) (*entities.User, error)
	UpdateBalance(ctx context.Context, userID uint64, newBalance decimal.Decimal) error
	Create(ctx context.Context, user *entities.User) error
}

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entities.Transaction) error
	ExistsByTransactionID(ctx context.Context, transactionID string) (bool, error)
	GetByUserID(ctx context.Context, userID uint64) ([]*entities.Transaction, error)
}
