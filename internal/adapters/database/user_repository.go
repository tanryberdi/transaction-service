package database

import (
	"context"
	"database/sql"
	"fmt"

	"transaction-service/internal/domain/entities"

	"github.com/shopspring/decimal"
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, userID uint64) (*entities.User, error) {
	query := "SELECT id, balance FROM users WHERE id = $1"

	var user entities.User
	var balanceStr string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &balanceStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	balance, err := decimal.NewFromString(balanceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse balance: %w", err)
	}
	user.Balance = balance

	return &user, nil
}

// UpdateBalance updates the user's balance
func (r *UserRepository) UpdateBalance(ctx context.Context, userID uint64, newBalance decimal.Decimal) error {
	query := "UPDATE users SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"

	result, err := r.db.ExecContext(ctx, query, newBalance, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := "INSERT INTO users (balance) VALUES ($1) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, user.Balance).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
