package model

import "time"

//对应users表

type User struct {
	Id           int32  `gorm:"primaryKey;autoIncrement:auto"`
	Email        string `gorm:"uniqueIndex;size:255"`         //用户邮箱
	PasswordHash string `gorm:"column:passwordHash;not null"` //加密后的密码
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
