package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User represents a user in the authentication system.
type User struct {
	ID                        bson.ObjectID `bson:"_id,omitempty"`
	FullName                  string        `bson:"full_name"`
	Email                     string        `bson:"email"`
	PasswordHash              string        `bson:"password_hash"`
	Verified                  bool          `bson:"verified"`
	VerificationCode          string        `bson:"verification_code"`
	VerificationCodeExpiresAt time.Time     `bson:"verification_code_expires_at"`
	CreatedAt                 time.Time     `bson:"created_at"`
	UpdatedAt                 time.Time     `bson:"updated_at"`
}

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id string, params UpdateUserParams) (*User, error)
	DeleteUser(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context, params FilterUsersParams) ([]*User, error)
}

// UpdateUserParams defines the optional parameters for updating a user.
// Only the fields that are not nil will be updated.
type UpdateUserParams struct {
	Email        *string
	FullName     *string
	PasswordHash *string
}

// FilterUsersParams defines the parameters for filtering and paginating users.
type FilterUsersParams struct {
	Email    *string
	Verified *bool
	Limit    uint64
	Offset   uint64
	SortBy   *string
	SortDesc bool
}
