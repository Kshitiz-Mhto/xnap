package logs

import "time"

type Log struct {
	ID                uint      `gorm:"primary_key"`
	Action            string    `gorm:"not null"`
	Command           string    `gorm:"not null"`
	Status            string    `gorm:"not null"`
	ErrorMessage      string    `gorm:"type:text"`
	UserName          string    `gorm:"index"`
	ExecutionDuration float64   `gorm:"type:double;default:0"`
	CreatedAt         time.Time `gorm:"default:current_timestamp"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp"`
}

func (Log) TableName() string {
	return "logs"
}
