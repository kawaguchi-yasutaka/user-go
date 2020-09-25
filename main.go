package main

import (
	"user-go/config"
	"user-go/initializer"
)

func main() {
	conf := config.SetConfig()
	db := initializer.InitDbConnection(conf)
	initializer.Migrate(db)
}
