package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email            string     `json:"email" gorm:"unique;not null"`
	PasswordHash     string     `json:"-" gorm:"not null"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Phone            string     `json:"phone"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	Timezone         string     `json:"timezone" gorm:"default:'UTC'"`
	Locale           string     `json:"locale" gorm:"default:'en-US'"`
	Role             UserRole   `json:"role" gorm:"default:'user'"`
	Status           UserStatus `json:"status" gorm:"default:'active'"`
	EmailVerified    bool       `json:"email_verified" gorm:"default:false"`
	PhoneVerified    bool       `json:"phone_verified" gorm:"default:false"`
	TwoFactorEnabled bool       `json:"two_factor_enabled" gorm:"default:false"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleUser         UserRole = "user"
	UserRolePremium      UserRole = "premium"
	UserRoleAdmin        UserRole = "admin"
	UserRoleFamilyOwner  UserRole = "family_owner"
	UserRoleFamilyMember UserRole = "family_member"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// UserSession represents a user session
type UserSession struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	RefreshToken    string     `json:"refresh_token" gorm:"unique;not null"`
	AccessTokenHash string     `json:"access_token_hash"`
	DeviceInfo      string     `json:"device_info" gorm:"type:jsonb"`
	IPAddress       string     `json:"ip_address"`
	UserAgent       string     `json:"user_agent"`
	ExpiresAt       time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt       *time.Time `json:"revoked_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Timezone  string `json:"timezone"`
	Locale    string `json:"locale"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Phone       string     `json:"phone"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Timezone    string     `json:"timezone"`
	Locale      string     `json:"locale"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Phone         string     `json:"phone"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	Timezone      string     `json:"timezone"`
	Locale        string     `json:"locale"`
	Role          UserRole   `json:"role"`
	Status        UserStatus `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	PhoneVerified bool       `json:"phone_verified"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// TableName specifies the table name for UserSession
func (UserSession) TableName() string {
	return "user_sessions"
}
