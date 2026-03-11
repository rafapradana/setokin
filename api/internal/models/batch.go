package models

import (
	"time"

	"github.com/google/uuid"
)

// Batch represents a stock batch in data.batches.
type Batch struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ItemID            uuid.UUID `gorm:"type:uuid;not null" json:"item_id"`
	BatchNumber       string    `gorm:"not null" json:"batch_number"`
	InitialQuantity   float64   `gorm:"type:numeric(10,3);not null" json:"initial_quantity"`
	RemainingQuantity float64   `gorm:"type:numeric(10,3);not null" json:"remaining_quantity"`
	ExpiryDate        time.Time `gorm:"type:date;not null" json:"expiry_date"`
	IsDepleted        bool      `gorm:"default:false" json:"is_depleted"`
	CreatedAt         time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time `gorm:"not null;default:now()" json:"updated_at"`

	Item Item `gorm:"foreignKey:ItemID" json:"item,omitempty"`
}

// TableName overrides the default table name.
func (Batch) TableName() string {
	return "data.batches"
}

// Status returns the batch status based on expiry date.
func (b *Batch) Status() string {
	now := time.Now().Truncate(24 * time.Hour)
	expiry := b.ExpiryDate.Truncate(24 * time.Hour)

	if expiry.Before(now) {
		return "expired"
	}
	if expiry.Sub(now) <= 3*24*time.Hour {
		return "expiring_soon"
	}
	return "good"
}

// DaysUntilExpiry returns the number of days until the batch expires.
func (b *Batch) DaysUntilExpiry() int {
	now := time.Now().Truncate(24 * time.Hour)
	expiry := b.ExpiryDate.Truncate(24 * time.Hour)
	return int(expiry.Sub(now).Hours() / 24)
}
