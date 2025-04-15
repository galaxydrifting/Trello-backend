package models

import (
	"time"
)

// Board represents a Kanban board
type Board struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Lists     []List         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// List represents a list in a Kanban board
type List struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"not null"`
	BoardID   uint           `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Cards     []Card         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Card represents a card in a Kanban list
type Card struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"not null"`
	Content   string
	ListID    uint           `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}