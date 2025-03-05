package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	cart "mall/kitex_gen/cart"
	"mall/model"
	"mall/rpc/cart/cartdao/mysql"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct{}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {

	ct := &model.Cart{UserId: req.GetUserId(),
		ProductId: req.GetItem().ProductId,
		Quantity:  req.GetItem().Quantity}

	err = mysql.AddItem(ctx, req.UserId, ct)
	if err != nil {
		return nil, kerrors.NewBizStatusError(500201, err.Error())
	}

	return &cart.AddItemResp{}, nil
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	var list []model.Cart
	list, err = mysql.GetCartByUserid(ctx, req.GetUserId())
	if err != nil {
		return nil, kerrors.NewBizStatusError(500202, err.Error())
	}

	var items []*cart.CartItem
	for _, item := range list {
		items = append(items, &cart.CartItem{ProductId: item.ProductId, Quantity: item.Quantity})
	}

	c := &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: req.GetUserId(),
			Items:  items,
		},
	}

	return c, nil
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	err = mysql.EmptyCart(ctx, req.GetUserId())
	if err != nil {
		return nil, kerrors.NewBizStatusError(500203, err.Error())
	}
	return &cart.EmptyCartResp{}, nil
}
