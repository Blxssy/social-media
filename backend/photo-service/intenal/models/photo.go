package models

import (
	"gorm.io/gorm"
	"time"
)

type Photo struct {
	gorm.Model
	ID            uint64
	UserID        uint64
	ImageURL      string
	LikesCount    uint32
	CommentsCount uint32
	Comments      []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
