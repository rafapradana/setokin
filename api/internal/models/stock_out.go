package models

import (
	"time"

	"github.com/google/uuid"
)

// StockOut represents a stock out transaction in data.stock_out.
type StockOut struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ItemID    uuid.UUID `gorm:"type:uuid;not null" json:"item_id"`
	Quantity  float64   `gorm:"type:numeric(10,3);not null" json:"quantity"`
	UsageDate time.Time `gorm:"type:date;not null" json:"usage_date"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`

	Item    Item             `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	Creator User             `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	Details []StockOutDetail `gorm:"foreignKey:StockOutID" json:"details,omitempty"`
}

// TableName overrides the default table name.
func (StockOut) TableName() string {
	return "data.stock_out"
}

// StockOutDetail represents FEFO batch deduction details in data.stock_out_details.
type StockOutDetail struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StockOutID   uuid.UUID `gorm:"type:uuid;not null" json:"stock_out_id"`
	BatchID      uuid.UUID `gorm:"type:uuid;not null" json:"batch_id"`
	QuantityUsed float64   `gorm:"type:numeric(10,3);not null" json:"quantity_used"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`

	Batch Batch `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
}

// TableName overrides the default table name.
func (StockOutDetail) TableName() string {
	return "data.stock_out_details"
}
