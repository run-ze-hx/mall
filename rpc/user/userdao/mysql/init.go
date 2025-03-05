package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mall/model"
)

var dsn = "root:qwe@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True"
var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	//自动迁移
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Println(" table users AutoMigrate failed")
	}

}
