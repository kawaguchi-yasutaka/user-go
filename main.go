package main

import (
	"fmt"
	"user-go/initializer"
)

func main() {
	fmt.Println("hello world")
	db := initializer.InitDbConnection()
	initializer.Migrate(db)
}
