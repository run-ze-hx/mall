package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mall/kitex_gen/cart"
	order "mall/kitex_gen/order"
	"mall/model"
	"mall/rpc/order/orderdao/mysql"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// PlaceOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	if len(req.OrderItems) == 0 {
		return nil, kerrors.NewBizStatusError(500401, "items is empty")
	}

	var orderId string
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		orderid, _ := uuid.NewUUID()

		o := &model.Order{
			OrderId:      orderid.String(),
			UserId:       req.UserId,
			UserCurrency: req.UserCurrency,
			Consignee: model.Consignee{
				Email: req.Email,
			},
		}
		if req.Address != nil {
			a := req.Address
			o.Consignee.StreetAddress = a.StreetAddress
			o.Consignee.City = a.City
			o.Consignee.State = a.State
			o.Consignee.Country = a.Country
		}
		if err := tx.Create(o).Error; err != nil {
			return err
		}

		var items []model.OrderItem
		for _, v := range req.OrderItems {
			items = append(items, model.OrderItem{
				OrderIdRefer: orderid.String(),
				ProductId:    v.Item.ProductId,
				Quantity:     v.Item.Quantity,
				Cost:         v.Cost,
			})
		}
		if err := tx.Create(items).Error; err != nil {
			return err
		}

		// 将订单 ID 赋值给外部变量
		orderId = orderid.String()
		return nil
	})

	if err != nil {
		return nil, kerrors.NewBizStatusError(500402, err.Error())
	}

	// 在事务外构建响应
	resp = &order.PlaceOrderResp{
		Order: &order.OrderResult{
			OrderId: orderId,
		},
	}

	return resp, nil
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	var orderlist []*model.Order
	orderlist, err = mysql.ListOrder(ctx, req.UserId)
	if err != nil {
		return nil, kerrors.NewBizStatusError(500403, err.Error())
	}
	var orders []*order.Order
	for _, v := range orderlist {

		//每个订单的商品信息
		var items []*order.OrderItem
		for _, oi := range v.OrderItems {
			items = append(items, &order.OrderItem{
				Item: &cart.CartItem{
					ProductId: oi.ProductId,
					Quantity:  oi.Quantity,
				},
				Cost: oi.Cost,
			})
		}

		//订单信息
		orders = append(orders, &order.Order{
			OrderId:      v.OrderId,
			UserId:       v.UserId,
			UserCurrency: v.UserCurrency,
			Email:        v.Consignee.Email,
			Address: &order.Address{
				StreetAddress: v.Consignee.StreetAddress,
				Country:       v.Consignee.Country,
				State:         v.Consignee.State,
				City:          v.Consignee.City,
				ZipCode:       v.Consignee.Zipcode,
			},
			OrderItems: items,
		})
	}
	return &order.ListOrderResp{Orders: orders}, nil
}

// MarkOrderPaid implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {

	orderId := req.OrderId

	// 查找订单
	var order *model.Order
	order, err = mysql.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, kerrors.NewBizStatusError(500404, orderId)
	}

	// 检查订单是否已标记为已支付
	if order.IsPaid {
		return nil, kerrors.NewBizStatusError(500405, "Order already marked as paid")
	}

	// 标记订单为已支付
	order.IsPaid = true

	// 保存更新到数据库
	err = mysql.UpdateOrder(ctx, order)
	if err != nil {
		return nil, kerrors.NewBizStatusError(500406, err.Error())
	}

	// 返回响应
	return nil, nil
}
