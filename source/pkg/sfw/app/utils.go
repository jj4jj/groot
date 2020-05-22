package app

import (
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/pkg/task"
	"groot/pkg/util"
	"math/rand"
	"time"
)

var (
	AppRandom *rand.Rand
)

//初始化常用组件
func InitUtils(config *conf.AppConfig) {
	AppRandom = rand.New(rand.NewSource(time.Now().UnixNano()))

	util.Init()

	e := task.Init(config.TaskBrokerName)
	if e != nil {
		log.Errorf(e, "init task broker:%s error !", config.TaskBrokerName)
		panic(e)
	}

}
