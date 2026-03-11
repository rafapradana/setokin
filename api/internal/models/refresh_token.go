package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a JWT refresh token in data.refresh_tokens.
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	TokenHash string     `gorm:"not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the default table name.
func (RefreshToken) TableName() string {
	return "data.refresh_tokens"
}

// IsExpired returns true if the token has expired.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsRevoked returns true if the token has been revoked.
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}
