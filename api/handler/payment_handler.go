package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/payment"
	"mall/kitex_gen/payment/paymentservice"
	"mall/utils"
	"net/http"
)

var PaymentCli paymentservice.Client

// PaymentCliInit 创建 payment 服务的 Client，并连接到 etcd
func PaymentCliInit(r discovery.Resolver) {
	pc, err := paymentservice.NewClient("payment",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	PaymentCli = pc
}

// Charge 处理支付请求
func Charge(c context.Context, ctx *app.RequestContext) {
	var req payment.ChargeReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Code: 4005001,
			Msg:  err.Error(),
		})
		return
	}

	resp, err := PaymentCli.Charge(c, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Code: 5005001,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Code: 200,
		Msg:  "Charge successful",
		Data: resp,
	})
}
