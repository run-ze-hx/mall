package userdao

import "mall/rpc/user/userdao/mysql"

func Init() {
	mysql.Init()
}
