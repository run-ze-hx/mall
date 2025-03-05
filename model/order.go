package model

import "gorm.io/gorm"

type Consignee struct {
	Email         string
	StreetAddress string
	City          string
	State         string
	Country       string
	Zipcode       int32
}

type Order struct {
	gorm.Model
	OrderId      string      `gorm:"type:varchar(100);uniqueIndex"`
	UserId       uint32      `gorm:"type:int(11)"`
	UserCurrency string      `gorm:"type:varchar(10)"`
	Consignee    Consignee   `gorm:"embedded"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"`
	IsPaid       bool        `gorm:"type:tinyint(1)"`
}

func (Order) TableName() string {
	return "order"
}
