package repository

import (
	"time"

	"gorm.io/gorm"
)

type ContaModel struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	Saldo     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Operacao struct {
	gorm.Model

	Type          int
	ContaSourceID *uint
	ContaSource   *ContaModel
	ContaTargetID *uint
	ContaTarget   *ContaModel
	Value         int
}

const (
	Addition = iota
	Autumn
	Winter
	Spring
)

// 0 transferencia
// 1 depósito
// 2 saque

// ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
