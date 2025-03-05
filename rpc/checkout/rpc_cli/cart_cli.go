package rpc_cli

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/cart/cartservice"
)

var CartCli cartservice.Client

// CartCliInit 创建 cart 服务的 Client，并连接到 etcd
func CartCliInit(r discovery.Resolver) {
	cc, err := cartservice.NewClient("cart",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	CartCli = cc
}
