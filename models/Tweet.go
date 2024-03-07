package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tweet struct {
	ID        string `gorm:"primaryKey"`
	Text      string
	ImageUrl  string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (t *Tweet) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
