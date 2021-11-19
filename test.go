package main

import "github.com/leapig/tecsun/helper"

func main() {
	helper.RegisterHelper()
	helper.RegisterMysql()
	helper.RegisterRedis()
	helper.Daemon()
}
