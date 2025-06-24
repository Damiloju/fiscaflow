package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for user data operations
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]User, error)
	CreateSession(ctx context.Context, session *UserSession) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*UserSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new user
func (r *repository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *repository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}

// List retrieves a list of users with pagination
func (r *repository) List(ctx context.Context, offset, limit int) ([]User, error) {
	var users []User
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// CreateSession creates a new user session
func (r *repository) CreateSession(ctx context.Context, session *UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByRefreshToken retrieves a session by refresh token
func (r *repository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*UserSession, error) {
	var session UserSession
	err := r.db.WithContext(ctx).Where("refresh_token = ? AND revoked_at IS NULL", refreshToken).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

// RevokeSession revokes a specific session
func (r *repository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&UserSession{}).Where("id = ?", sessionID).Update("revoked_at", gorm.Expr("NOW()")).Error
}

// RevokeAllUserSessions revokes all sessions for a user
func (r *repository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&UserSession{}).Where("user_id = ?", userID).Update("revoked_at", gorm.Expr("NOW()")).Error
}

// Custom errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
)
