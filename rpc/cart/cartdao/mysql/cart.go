package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"mall/model"
)

func AddItem(ctx context.Context, userid uint32, cart *model.Cart) error {
	var row model.Cart
	err := DB.WithContext(ctx).Model(&model.Cart{}).Where(&model.Cart{UserId: userid, ProductId: cart.ProductId}).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果记录不存在，直接插入
			return DB.Create(&model.Cart{
				UserId:    userid,
				ProductId: cart.ProductId,
				Quantity:  cart.Quantity,
			}).Error
		}
		// 其他错误直接返回
		return err
	}
	// 如果记录存在，更新 Quantity
	return DB.WithContext(ctx).Model(&model.Cart{}).Where(&model.Cart{UserId: userid, ProductId: cart.ProductId}).UpdateColumn("quantity", gorm.Expr("quantity+?", cart.Quantity)).Error
}

func GetCartByUserid(ctx context.Context, userid uint32) ([]model.Cart, error) {
	var rows []model.Cart
	err := DB.WithContext(ctx).Model(&model.Cart{}).Where(&model.Cart{UserId: userid}).Find(&rows).Error
	return rows, err
}

func EmptyCart(ctx context.Context, userid uint32) error {
	if userid == 0 {
		return errors.New("userid is required")
	}
	return DB.WithContext(ctx).Where(&model.Cart{UserId: userid}).Delete(&model.Cart{}).Error
}
