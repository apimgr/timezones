package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AdminCredentials represents admin user credentials
type AdminCredentials struct {
	ID           int
	Username     string
	PasswordHash string
	Token        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateAdminUser creates a new admin user with hashed password and token
func CreateAdminUser(username, password string) (*AdminCredentials, error) {
	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate random token
	token, err := generateToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Insert into database
	result, err := DB.Exec(`
		INSERT INTO admin_users (username, password_hash, token)
		VALUES (?, ?, ?)
	`, username, string(passwordHash), token)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	id, _ := result.LastInsertId()

	return &AdminCredentials{
		ID:           int(id),
		Username:     username,
		PasswordHash: string(passwordHash),
		Token:        token,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// ValidatePassword checks if the provided password matches the stored hash
func ValidatePassword(username, password string) (*AdminCredentials, error) {
	var creds AdminCredentials

	err := DB.QueryRow(`
		SELECT id, username, password_hash, token, created_at, updated_at
		FROM admin_users
		WHERE username = ?
	`, username).Scan(&creds.ID, &creds.Username, &creds.PasswordHash, &creds.Token, &creds.CreatedAt, &creds.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid username or password")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return &creds, nil
}

// ValidateToken checks if the provided token is valid
func ValidateToken(token string) (*AdminCredentials, error) {
	var creds AdminCredentials

	err := DB.QueryRow(`
		SELECT id, username, password_hash, token, created_at, updated_at
		FROM admin_users
		WHERE token = ?
	`, token).Scan(&creds.ID, &creds.Username, &creds.PasswordHash, &creds.Token, &creds.CreatedAt, &creds.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid token")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &creds, nil
}

// GetAdminUser retrieves admin user by username
func GetAdminUser(username string) (*AdminCredentials, error) {
	var creds AdminCredentials

	err := DB.QueryRow(`
		SELECT id, username, password_hash, token, created_at, updated_at
		FROM admin_users
		WHERE username = ?
	`, username).Scan(&creds.ID, &creds.Username, &creds.PasswordHash, &creds.Token, &creds.CreatedAt, &creds.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &creds, nil
}

// AdminUserExists checks if any admin user exists
func AdminUserExists() (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM admin_users").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// generateToken generates a random hex token
func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
