package database

import (
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
)

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	// Create users table
	if err := createUsersTable(db); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create transactions table
	if err := createTransactionsTable(db); err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	// Insert predefined users
	if err := insertPredefinedUsers(db); err != nil {
		return fmt.Errorf("failed to insert predefined users: %w", err)
	}

	return nil
}

func createUsersTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_users_id ON users(id);
	`
	_, err := db.Exec(query)
	return err
}

func createTransactionsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS transactions (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id),
			transaction_id VARCHAR(255) NOT NULL UNIQUE,
			state VARCHAR(10) NOT NULL CHECK (state IN ('win', 'lose')),
			amount DECIMAL(15,2) NOT NULL,
			source_type VARCHAR(20) NOT NULL CHECK (source_type IN ('game', 'server', 'payment')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_transactions_transaction_id ON transactions(transaction_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
	`
	_, err := db.Exec(query)
	return err
}

func insertPredefinedUsers(db *sql.DB) error {
	// Check if users already exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id IN (1, 2, 3)").Scan(&count)
	if err != nil {
		return err
	}

	// If users already exist, skip insertion
	if count >= 3 {
		return nil
	}

	// Insert predefined users with initial balance
	initialBalance := decimal.NewFromFloat(100.00) // Starting with 100.00 balance

	users := []struct {
		id      uint64
		balance decimal.Decimal
	}{
		{1, initialBalance},
		{2, initialBalance},
		{3, initialBalance},
	}

	for _, user := range users {
		query := `
			INSERT INTO users (id, balance) 
			VALUES ($1, $2)
			ON CONFLICT (id) DO NOTHING
		`
		_, err := db.Exec(query, user.id, user.balance)
		if err != nil {
			return fmt.Errorf("failed to insert user %d: %w", user.id, err)
		}
	}

	// Reset the sequence to start from 4 for future auto-generated IDs
	_, err = db.Exec("SELECT setval('users_id_seq', 3, true)")
	if err != nil {
		return fmt.Errorf("failed to reset user sequence: %w", err)
	}

	return nil
}
