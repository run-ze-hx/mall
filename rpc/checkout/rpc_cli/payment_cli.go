package rpc_cli

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/payment/paymentservice"
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
