package database

import (
	"context"
	"database/sql"
	"fmt"

	"transaction-service/internal/domain/entities"

	"github.com/shopspring/decimal"
)

// TransactionRepository implements the transaction repository interface
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(ctx context.Context, transaction *entities.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, transaction_id, state, amount, source_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		transaction.UserID,
		transaction.TransactionID,
		transaction.State,
		transaction.Amount,
		transaction.SourceType,
		transaction.CreatedAt,
	).Scan(&transaction.ID)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// ExistsByTransactionID checks if a transaction with the given ID exists
func (r *TransactionRepository) ExistsByTransactionID(ctx context.Context, transactionID string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM transactions WHERE transaction_id = $1)"

	var exists bool
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check transaction existence: %w", err)
	}

	return exists, nil
}

// GetByUserID retrieves all transactions for a user
func (r *TransactionRepository) GetByUserID(ctx context.Context, userID uint64) ([]*entities.Transaction, error) {
	query := `
		SELECT id, user_id, transaction_id, state, amount, source_type, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		var transaction entities.Transaction
		var amountStr string

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.TransactionID,
			&transaction.State,
			&amountStr,
			&transaction.SourceType,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = amount

		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}
