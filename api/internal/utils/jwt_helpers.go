package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MustParseRefreshExpiry extracts the expiry time from a refresh token.
// Returns a safe default if parsing fails.
func MustParseRefreshExpiry(tokenString, secret string) time.Time {
	claims, err := ValidateRefreshToken(tokenString, secret)
	if err != nil || claims.ExpiresAt == nil {
		return time.Now().Add(7 * 24 * time.Hour) // fallback: 7 days
	}
	return claims.ExpiresAt.Time
}

// ExtractTokenID extracts the JTI (token ID) from a refresh token.
func ExtractTokenID(tokenString, secret string) string {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return ""
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims.ID
	}
	return ""
}
