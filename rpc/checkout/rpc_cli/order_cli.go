package rpc_cli

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/order/orderservice"
)

var OrderCli orderservice.Client

// OrderCliInit 创建 order 服务的 Client，并连接到 etcd
func OrderCliInit(r discovery.Resolver) {
	oc, err := orderservice.NewClient("order",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	OrderCli = oc
}
