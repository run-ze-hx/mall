package mysql

import (
	"context"
	"mall/model"
)

func CreatePayment(ctx context.Context, p *model.Payment) error {
	return DB.WithContext(ctx).Model(&model.Payment{}).Create(p).Error
}
