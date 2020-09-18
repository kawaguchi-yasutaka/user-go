package initializer

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	inframysql "user-go/infra/mysql"
)

func InitDbConnection() *gorm.DB {
	dsn := getDbDsn()
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

func getDbDsn() string {
	err := godotenv.Load(".env.local")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		os.Getenv("DB_USER_NAME"),
		os.Getenv("DB_USER_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)
}
