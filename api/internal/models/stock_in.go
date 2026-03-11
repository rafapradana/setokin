package models

import (
	"time"

	"github.com/google/uuid"
)

// StockIn represents a stock in transaction in data.stock_in.
type StockIn struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ItemID        uuid.UUID  `gorm:"type:uuid;not null" json:"item_id"`
	BatchID       uuid.UUID  `gorm:"type:uuid;not null" json:"batch_id"`
	SupplierID    *uuid.UUID `gorm:"type:uuid" json:"supplier_id,omitempty"`
	Quantity      float64    `gorm:"type:numeric(10,3);not null" json:"quantity"`
	PurchaseDate  time.Time  `gorm:"type:date;not null" json:"purchase_date"`
	PurchasePrice *float64   `gorm:"type:numeric(15,2)" json:"purchase_price,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"created_at"`

	Item     Item      `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	Batch    Batch     `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
	Supplier *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Creator  User      `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}

// TableName overrides the default table name.
func (StockIn) TableName() string {
	return "data.stock_in"
}
