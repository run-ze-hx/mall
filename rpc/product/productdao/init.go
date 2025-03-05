package productdao

import (
	"mall/rpc/product/productdao/mysql"
	"mall/rpc/product/productdao/redis"
)

func Init() {
	mysql.Init()
	redis.InitRedis()
}
