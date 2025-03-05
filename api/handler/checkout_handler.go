package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/checkout"
	"mall/kitex_gen/checkout/checkoutservice"
	"mall/kitex_gen/payment"
	"mall/utils"
	"net/http"
)

// CheckOutReq 获取请求参数
type CheckOutReq struct {
	UserId     uint32                  `json:"user_id,omitempty"`
	Firstname  string                  `json:"firstname,omitempty"`
	Lastname   string                  `json:"lastname,omitempty"`
	Email      string                  `json:"email,omitempty"`
	Address    *checkout.Address       `json:"address,omitempty"`
	CreditCard *payment.CreditCardInfo `json:"credit_card,omitempty"`
}

// CheckOutResp 返回请求结果
type CheckOutResp struct {
	OrderId       string `json:"order_id,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
}

var CheckoutCli checkoutservice.Client

// CheckoutCliInit 创建 checkout 服务的 Client，并连接到 etcd
func CheckoutCliInit(r discovery.Resolver) {
	cc, err := checkoutservice.NewClient("checkout",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	CheckoutCli = cc
}

// Checkout 服务
func Checkout(c context.Context, ctx *app.RequestContext) {
	//获取参数
	var req CheckOutReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 4002001, Msg: "请求参数获取异常:" + err.Error()})
		return
	}
	if req.Firstname == "" || req.Lastname == "" || req.Email == "" || req.Address == nil || req.CreditCard == nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 4002002, Msg: "参数不能为空"})
		return
	}

	var resp *checkout.CheckoutResp
	resp, err = CheckoutCli.Checkout(c, &checkout.CheckoutReq{
		UserId:     req.UserId,
		Firstname:  req.Firstname,
		Lastname:   req.Lastname,
		Email:      req.Email,
		Address:    req.Address,
		CreditCard: req.CreditCard,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5002003, Msg: "checkout fail:" + err.Error()})
		return
	}
	if resp == nil || resp.OrderId == "" || resp.TransactionId == "" {
		ctx.JSON(http.StatusOK, utils.Response{Code: 5002004, Msg: "checkout 微服务 resp 为空"})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Data: CheckOutResp{
			OrderId:       resp.OrderId,
			TransactionId: resp.TransactionId,
		},
		Msg: "成功",
	})
}
