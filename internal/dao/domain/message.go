package domain

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Phone   string
	Content string `gorm:"type:varchar(160);not null"`
	IsSent  bool
}
