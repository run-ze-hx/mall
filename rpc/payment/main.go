package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	payment "mall/kitex_gen/payment/paymentservice"
	"mall/rpc/payment/paymentdao"
	"net"
)

func main() {
	//数据库启动
	paymentdao.Init()

	//启动服务注册中心
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"}) // r不应重复使用。
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8885")
	svr := payment.NewServer(new(PaymentServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "payment"}),
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
