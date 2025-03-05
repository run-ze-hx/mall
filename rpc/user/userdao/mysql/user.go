package mysql

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mall/model"
)

func CreateUser(user *model.User) error {
	return DB.Create(user).Error
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found") // 用户邮箱不存在
		}
		return nil, result.Error // 其他错误
	}
	return &user, nil
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
