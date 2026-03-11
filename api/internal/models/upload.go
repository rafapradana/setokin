package models

import (
	"time"

	"github.com/google/uuid"
)

// Upload tracks file uploads via MinIO presigned URLs in data.uploads.
type Upload struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FileName   string     `gorm:"not null" json:"file_name"`
	FileType   string     `gorm:"not null" json:"file_type"`
	ObjectKey  string     `gorm:"not null" json:"object_key"`
	EntityType string     `gorm:"not null" json:"entity_type"` // items, suppliers, stock_in_invoice
	EntityID   string     `json:"entity_id,omitempty"`
	Status     string     `gorm:"not null;default:pending" json:"status"` // pending, confirmed, expired
	UploadedBy uuid.UUID  `gorm:"type:uuid;not null" json:"uploaded_by"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()" json:"updated_at"`

	Uploader User `gorm:"foreignKey:UploadedBy" json:"-"`
}

// TableName overrides the default table name.
func (Upload) TableName() string {
	return "data.uploads"
}
