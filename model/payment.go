package model

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	UserId        uint32    `json:"user_id"`
	OrderId       string    `json:"order_id"`
	TransactionId string    `json:"transaction_id"`
	Amount        float32   `json:"amount"`
	PayAt         time.Time `json:"pay_at"`
}

func (p Payment) TableName() string {
	return "payment"
}
