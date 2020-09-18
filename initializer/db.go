package initializer

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	inframysql "user-go/infra/mysql"
)

func InitDbConnection() *gorm.DB {
	dsn := "master:gagagigu123@tcp(127.0.0.1:3306)/usergo?charset=utf8&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&inframysql.User{})
	db.AutoMigrate(&inframysql.UserAuthentication{})
}
