package rpc_cli

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/product/productcatalogservice"
)

var ProductCli productcatalogservice.Client

// ProductCliInit 创建 product 服务的 Client，并连接到 etcd
func ProductCliInit(r discovery.Resolver) {
	pc, err := productcatalogservice.NewClient("product",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	ProductCli = pc
}
