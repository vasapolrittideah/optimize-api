package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator represents a JWT based authenticator.
type JWTAuthenticator struct {
	audience string
	issuer   string
}

// NewJWTAuthenticator creates a new JWTAuthenticator instance.
func NewJWTAuthenticator(audience, issuer string) Authenticator {
	return &JWTAuthenticator{
		audience: audience,
		issuer:   issuer,
	}
}

// GenerateToken generates a JWT token with the given claims and secret.
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ValidateToken validates a JWT token with the given secret.
func (a *JWTAuthenticator) ValidateToken(token, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.audience),
		jwt.WithIssuer(a.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
