package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/comm/constk"
	"groot/pkg/broker"
	db "groot/pkg/db"
	sfw_event "groot/pkg/event"
	"groot/pkg/task"
	"groot/pkg/util"
	"net/http"
	"sync"
)

var (
	wg sync.WaitGroup
)

func Init() error {
	//解析命令行
	parseCmdLine()

	//处理命令行TODO
	configPath := GetConfigPath()

	//解析配置文件初始化基础日志和配置组件
	if e := conf.LoadConfig(configPath); e != nil {
		panic(e)
	}

	//fmt.Println("TODO: check cmd line routeHandler and execute")
	//初始化DB配置组件
	var config = conf.GetAppConfig()
	log.Infof("load config info:%v", config)

	//fmt.Println("set gin mode:", config.GinMode)
	if config.RunEnv != "prod"  {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	for i, dbcx := range config.DbConnxs {
		if dbcx.Enable == false {
			continue
		}
		if e := db.AddDb(dbcx.Name, dbcx.Backend, dbcx.ConnxEnv, dbcx.MaxIdleConnx, dbcx.MaxOpenConnx); e != nil {
			log.Errorf(e, "db connection init error ith:%d name:%s error", i, dbcx.Name)
			panic(e)
		}
	}

	for i, brkcx := range config.Brokers {
		if brkcx.Enable == false {
			continue
		}
		if e := broker.AddBroker(brkcx.Name, brkcx.Backend, brkcx.ConnxEnv); e != nil {
			log.Errorf(e, "broker connection init error ith:%d name:%s error", i, brkcx.Name)
			panic(e)
		}

	}

	InitUtils(config)

	//fmt.Println("TODO: Read Config File and Set Config Options")
	//加载初始化注册的应用服务

	//fmt.Println("TODO: Init some components from config options")

	return nil
}

func StartServer(svrConf *conf.ServerListener) {
	log.Infof("Start One Server:%s Listen:%s ...", svrConf.Name, svrConf.Bind)
	cert := svrConf.TlsCertFile
	key := svrConf.TlsKeyFile
	oneRouter := GetServerRouter(svrConf.Name)
	if oneRouter == nil {
		log.Errorf(nil, "server name:%s router is not registerd !", svrConf.Name)
		return
	}
	var err error
	if cert != "" && key != "" {
		err = http.ListenAndServeTLS(svrConf.Bind, cert, key, oneRouter)
	} else {
		err = http.ListenAndServe(svrConf.Bind, oneRouter)
	}
	if err != nil {
		log.Errorf(err, "server listen run error ")
	}
}

func StartServiceTcsWorker(svrConf *conf.ServerListener) {
	svr := GetOrNewServer(svrConf.Name)
	for _,svc := range svr.mpService {
		if len(svc.mpTcsMethod) > 0 {
			//service tcs topic
			topic := util.GetBrokerTopicName(constk.BRK_TOPIC_CALL_SERVICE,svc.name)
			log.Infof("tcs service:%s worker num:%d is listening(%s) ...", svc.name, svc.tcsWorkerNum, topic)
			go task.RunWorker(topic, svc.tcsWorkerNum, createTcsWorkerHandler(svc.mpTcsMethod))
		}
	}
}

func SyncDb() error {
	config := conf.GetAppConfig()
	for i,_ := range config.Listeners {
		if err := InitSyncDb(config.Listeners[i].Name); err != nil {
			log.Errorf(err, "server:%s sync db error", config.Listeners[i].Name)
			return err
		}
	}
	return nil
}

func Run() error {
	if *cmd_bp_syncdb {
		//sync db
		return SyncDb()
	}


	//running
	config := conf.GetAppConfig()
	for i,_ := range config.Listeners {
		if err := InitServices(config.Listeners[i].Name); err != nil {
			log.Errorf(err, "init server:%s error", config.Listeners[i].Name)
			return err
		}
	}
	//start run server
	wg.Add(len(config.Listeners))
	for i, _ := range config.Listeners {
		go func(svr *conf.ServerListener) {
			defer wg.Done()
			StartServiceTcsWorker(svr)
			StartServer(svr)
		}(&config.Listeners[i])
	}

	//start dispatcher
	wg.Add(1)
	sfw_event.SetEventBroker(config.Event.BrokerName)
	go func() {
		defer wg.Done()
		sfw_event.StartEventDispatcher(config.Event.DispatcherNum)
	}()

	wg.Wait()
	return nil
}
