package paymentdao

import "mall/rpc/payment/paymentdao/mysql"

func Init() {
	mysql.Init()
}
