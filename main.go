package main

import (
	_ "github.com/go-sql-driver/mysql"
	"golang-restful-api/helper"
)

func main() {
	server := InitializedServer()
	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
