package auth

import "github.com/golang-jwt/jwt/v5"

// Authenticator defines the interface for managing authentication.
type Authenticator interface {
	GenerateToken(claims jwt.Claims, secret string) (string, error)
	ValidateToken(token, secret string) (*jwt.Token, error)
}
