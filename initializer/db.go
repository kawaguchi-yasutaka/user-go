package initializer

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"user-go/config"
	inframysql "user-go/infra/mysql"
)

func InitDbConnection(conf config.Config) *gorm.DB {
	dsn := getDbDsn(conf)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&inframysql.User{})
	db.AutoMigrate(&inframysql.UserAuthentication{})
	db.AutoMigrate(&inframysql.UserRemember{})
}

func getDbDsn(conf config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		conf.DB.UserName,
		conf.DB.UserPassword,
		conf.DB.Address,
		conf.DB.Port,
		conf.DB.Table,
	)
}
