package local

import (
	"time"

	"github.com/google/uuid"
)

type Backup struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	FileName   string    `gorm:"not null"`
	SourcePath string    `gorm:"not null"`
	BackupPath string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"default:current_timestamp"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp;onUpdate:current_timestamp"`
}

func (Backup) TableName() string {
	return "backups"
}
