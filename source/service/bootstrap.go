package service

import (
	"groot/comm/conf"
	"groot/service/feedback"
	"groot/service/notify"
	"groot/service/profile"
	"groot/sfw/app"
	"groot/sfw/middleware"

	"groot/service/account"
	"groot/service/captcha"
	"groot/service/chat"
	"groot/service/file"
	"groot/service/friend"
	"groot/service/record"
	"groot/service/msg"
)

func InitService() error {
	config := conf.GetAppConfig()

	//here register service
	app.AddMiddleWare("logging", middleware.Logging())
	app.AddMiddleWare("request-id", middleware.RequestId())
	app.AddMiddleWare("nocache", middleware.NoCache)
	app.AddMiddleWare("options", middleware.Options)
	app.AddMiddleWare("secure", middleware.Secure)

	account.Initialize()
	profile.Initialize()
	captcha.Initialize()
	notify.Initialize()
	friend.Initialize()
	record.Initialize()
	chat.Initialize()
	msg.Initialize()
	feedback.Initialize()

	if config.RunEnv == "dev" {
		file.Initialize()
	}

	return nil
}
