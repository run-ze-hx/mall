package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"mall/api/handler"
	"mall/api/router"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	// 创建 etcd 注册中心
	r, err := etcd.NewEtcdResolver([]string{"localhost:2379"})
	if err != nil {
		log.Fatal(err)
	}

	//创建 auth 服务的 Client，并连接到 etcd
	handler.AuthCliInit(r)
	//创建 user 服务的 Client，并连接到 etcd
	handler.UserCliInit(r)
	//创建 product 服务的 Client，并连接到 etcd
	handler.ProductCliInit(r)
	//创建 cart 服务的 Client，并连接到 etcd
	handler.CartCliInit(r)
	//创建 payment 服务的 Client，并连接到 etcd
	handler.PaymentCliInit(r)
	//创建 checkout 服务的 Client，并连接到 etcd
	handler.CheckoutCliInit(r)
	//创建 order 服务的 Client，并连接到 etcd
	handler.OrderCliInit(r)

	// 启动 HTTP 服务器来暴露 pprof 接口
	go func() {
		log.Println("Starting pprof server on :6060")
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	hz := server.Default(server.WithHostPorts("localhost:8889"))
	router.Register(hz, handler.AuthCli)

	if err = hz.Run(); err != nil {
		log.Fatal(err)
	}

}
