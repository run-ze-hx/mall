package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/order"
	"mall/kitex_gen/order/orderservice"
	"mall/utils"
	"net/http"
)

var OrderCli orderservice.Client

// OrderCliInit 创建 payment 服务的 Client，并连接到 etcd
func OrderCliInit(r discovery.Resolver) {
	oc, err := orderservice.NewClient("order",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	OrderCli = oc
}

// PlaceOrder 创建订单
func PlaceOrder(c context.Context, ctx *app.RequestContext) {
	var req order.PlaceOrderReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Code: 4004002,
			Msg:  err.Error(),
		})
		return
	}

	resp, err := OrderCli.PlaceOrder(c, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Code: 5004002,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Code: 200,
		Msg:  "Order placed",
		Data: resp,
	})
}

// MarkOrderPaid 标记订单为已支付
func MarkOrderPaid(c context.Context, ctx *app.RequestContext) {
	var req order.MarkOrderPaidReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Code: 4004003,
			Msg:  err.Error(),
		})
		return
	}

	resp, err := OrderCli.MarkOrderPaid(c, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Code: 5004003,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Code: 200,
		Msg:  "Order marked as paid",
		Data: resp,
	})
}

// ListOrder 获取订单列表
func ListOrder(c context.Context, ctx *app.RequestContext) {
	var req order.ListOrderReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Code: 4004004,
			Msg:  err.Error(),
		})
		return
	}

	resp, err := OrderCli.ListOrder(c, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Code: 5004004,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Code: 200,
		Msg:  "Order list retrieved",
		Data: resp,
	})
}
