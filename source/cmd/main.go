package main

import (
	"groot/service"
	"groot/sfw/app"
	"groot/sfw/util"
)

func main() {
	var err = app.Init()
	if err != nil {
		panic(err)
	}
	err = service.InitService()
	if err != nil {
		panic(err)
	}
	util.CheckError(app.Run(), "main.run")
}
