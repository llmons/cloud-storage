package dal

import (
	"cloud-storage/biz/dal/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dsn = "root:123456@(127.0.0.1:3306)/cloud_storage?charset=utf8mb4&parseTime=True&loc=Local"

func Init() {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	query.SetDefault(db)
}
