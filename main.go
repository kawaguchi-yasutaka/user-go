package main

import (
	"user-go/config"
	"user-go/initializer"
	"user-go/web"
)

func main() {
	conf := config.SetConfig()
	db := initializer.InitDbConnection(conf)
	infra := initializer.NewInfra(conf)
	repository := initializer.InitRepository(db)
	service := initializer.NewService(infra, repository)
	initializer.Migrate(db)
	web.Init(service)
}
