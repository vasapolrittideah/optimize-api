package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Session represents an authentication user session with access and refresh tokens.
type Session struct {
	ID                    bson.ObjectID `bson:"_id,omitempty"`
	UserID                string        `bson:"user_id"`
	AccessToken           string        `bson:"access_token"`
	RefreshToken          string        `bson:"refresh_token"`
	AccessTokenExpiresAt  time.Time     `bson:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time     `bson:"refresh_token_expires_at"`
	IPAddress             *string       `bson:"ip_address"`
	UserAgent             *string       `bson:"user_agent"`
	CreatedAt             time.Time     `bson:"created_at"`
	UpdatedAt             time.Time     `bson:"updated_at"`
}

// SessionRepository defines the interface for session-related database operations.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) (*Session, error)
	GetSessionByUserID(ctx context.Context, userID string) (*Session, error)
	UpdateTokens(ctx context.Context, id string, params UpdateTokensParams) (*Session, error)
}

// UpdateTokensParams defines the parameters for updating session tokens.
type UpdateTokensParams struct {
	AccessToken           string    `bson:"access_token"`
	RefreshToken          string    `bson:"refresh_token"`
	AccessTokenExpiresAt  time.Time `bson:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `bson:"refresh_token_expires_at"`
}
