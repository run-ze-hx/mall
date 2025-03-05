package mysql

import (
	"context"
	"mall/model"
)

func ListOrder(c context.Context, userid uint32) (orders []*model.Order, err error) {
	err = DB.WithContext(c).Where("user_id = ?", userid).Preload("OrderItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return
}
func GetOrderById(ctx context.Context, orderId string) (*model.Order, error) {
	var order model.Order
	if err := DB.Model(&order).Where("order_id = ?", orderId).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func UpdateOrder(ctx context.Context, order *model.Order) error {
	return DB.Save(order).Error
}
