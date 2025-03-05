package handler

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/auth/authservice"
)

var AuthCli authservice.Client

// AuthCliInit 创建 auth 服务的 Client，并连接到 etcd
func AuthCliInit(r discovery.Resolver) {

	ac, err := authservice.NewClient("auth",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err != nil {
		log.Fatal(err)
	}
	AuthCli = ac
}
