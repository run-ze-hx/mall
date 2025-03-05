package main

import (
	"context"
	"errors"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"mall/kitex_gen/cart"
	checkout "mall/kitex_gen/checkout"
	"mall/kitex_gen/order"
	"mall/kitex_gen/payment"
	"mall/kitex_gen/product"
	"mall/rpc/checkout/rpc_cli"
	"strconv"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct{}

// Checkout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) Checkout(ctx context.Context, req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	//拿到购物车
	var response *cart.GetCartResp
	response, err = rpc_cli.CartCli.GetCart(ctx, &cart.GetCartReq{UserId: req.UserId})
	if err != nil {
		return nil, kerrors.NewBizStatusError(500301, err.Error())
	}
	if response == nil || response.Cart.Items == nil {
		return nil, kerrors.NewBizStatusError(500302, "cart not exist")
	}

	//拿到商品价格和购买数量
	var total float32
	var orderItems []*order.OrderItem
	for _, item := range response.Cart.Items {
		resp, Err := rpc_cli.ProductCli.GetProduct(ctx, &product.GetProductReq{Id: item.ProductId})
		if Err != nil {
			return nil, kerrors.NewBizStatusError(500303, Err.Error())
		}
		if resp.Product == nil {
			continue
		}
		orderItems = append(orderItems, &order.OrderItem{
			Item: item,
			Cost: resp.Product.Price,
		})
		total += resp.Product.Price * float32(item.Quantity) //计算价格
	}

	//下单

	var zipcode int
	zipcode, err = strconv.Atoi(req.Address.ZipCode)
	if err != nil {
		return nil, errors.New("zipcode 类型转换失败")
	}
	var orderId string
	var orderresp *order.PlaceOrderResp
	orderresp, err = rpc_cli.OrderCli.PlaceOrder(ctx, &order.PlaceOrderReq{
		UserId: req.UserId,
		Email:  req.Email,
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			Country:       req.Address.Country,
			State:         req.Address.State,
			ZipCode:       int32(zipcode),
		},
		OrderItems: orderItems,
	})
	if err != nil {
		return nil, kerrors.NewBizStatusError(500304, err.Error())
	}
	if orderresp != nil {
		if orderresp.Order != nil {
			orderId = orderresp.Order.OrderId
		} else {
			return nil, errors.New("order is nil")
		}
	} else {
		return nil, errors.New("order response is nil ")
	}

	//支付
	var Payresp *payment.ChargeResp
	Payresp, err = rpc_cli.PaymentCli.Charge(ctx, &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
		},
	})
	if err != nil {
		return nil, kerrors.NewBizStatusError(500305, err.Error())
	}

	//修改订单信息为已支付
	_, err = rpc_cli.OrderCli.MarkOrderPaid(ctx, &order.MarkOrderPaidReq{OrderId: orderId})
	if err != nil {
		return nil, kerrors.NewBizStatusError(500306, err.Error())
	}

	//支付完清空购物车
	_, err = rpc_cli.CartCli.EmptyCart(ctx, &cart.EmptyCartReq{UserId: req.UserId})
	if err != nil {
		return nil, kerrors.NewBizStatusError(500307, err.Error())
	}

	checkoutresp := &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: Payresp.TransactionId,
	}

	return checkoutresp, nil
}
