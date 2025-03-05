package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	checkout "mall/kitex_gen/checkout/checkoutservice"
	"mall/rpc/checkout/rpc_cli"
	"net"
)

func main() {

	// 创建 etcd 注册中心,调用其他微服务
	r, err := etcd.NewEtcdResolver([]string{"localhost:2379"})
	if err != nil {
		log.Fatal(err)
	}

	//启动调用到的微服务
	rpc_cli.CartCliInit(r)
	rpc_cli.PaymentCliInit(r)
	rpc_cli.ProductCliInit(r)
	rpc_cli.OrderCliInit(r)

	etcdRegistry, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"}) // r不应重复使用。
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8883")
	svr := checkout.NewServer(new(CheckoutServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "checkout"}),
		server.WithRegistry(etcdRegistry),
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
