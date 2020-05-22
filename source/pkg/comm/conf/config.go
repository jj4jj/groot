package conf

import (
	"fmt"
	"github.com/lexkong/log"
	"encoding/json"
)


type BrokerConnx struct {
	Name         string
	Enable       bool
	Backend      string
	ConnxEnv     string
	MaxIdleConnx int
	MaxOpenConnx int
}


type EventDispatcherCfg struct {
	BrokerName    string
	DispatcherNum int
}


type DBConnection struct {
	Name         string
	Enable       bool
	Backend      string
	ConnxEnv     string
	MaxIdleConnx int
	MaxOpenConnx int
	Desc 		 string
}



type ServerListener struct {
	Name        string
	Bind        string
	TlsCertFile string
	TlsKeyFile  string
}


type RunDevEnvConfig struct {
	UploadFileDir 	string
}


type SvcUseDbConfig struct {
	SplitDbNum		uint32
	SplitTbNum		uint32
	DbUseMap		map[string]string
}


type AccountSvcConfig  struct {
	UseDb SvcUseDbConfig
}

type CaptchaSvcConfig struct {
	UseDb SvcUseDbConfig
}

type FriendSvcConfig  struct {
	UseDb SvcUseDbConfig
}


type ChatSvcConfig  struct {
	UseDb     SvcUseDbConfig
	MsgSvcUrl string
}


type MsgSvcConfig  struct {
	EnableSplitDbBegIdx 		uint32
	EnableSplitDbEndIdx 		uint32
}

type ProfileSvcConfig  struct {
	UseDb SvcUseDbConfig
}


type NotifySvcConfig  struct {
	UseDb SvcUseDbConfig
}


type AuthTokenConfig struct {
	SessionSignKeySecret	string
	JwtSignKeySecret		string
	FriendSignKeySecret		string
	FileSignKeySecret		string
	SmsCodeSignKeySecret	string
}

//定义AppConfig 基础结构,用于嵌套直接获取
type AppConfig struct {
	LogConf     	  log.PassLagerCfg `yaml:"log"`
	Listeners         []ServerListener `yaml:"listeners"`
	DbConnxs          []DBConnection   `yaml:"dbconnxs"`
	AuthRsaPrivateKey string
	RunEnv            string //dev/prod/test
	Brokers           []BrokerConnx
	Event             EventDispatcherCfg
	TaskBrokerName	  string
	DevEnv			  RunDevEnvConfig 	//when run env == dev,test,prod (it's valid)
	Account 		  AccountSvcConfig
	Captcha 		  CaptchaSvcConfig
	Friend 			  FriendSvcConfig
	Chat 			  ChatSvcConfig
	Msg 			  MsgSvcConfig
	Profile 		  ProfileSvcConfig
	AuthToken 		  AuthTokenConfig
	Notify 			  NotifySvcConfig
}

type IosApnsConfig struct  {
	DevApnsP12	string
}

func AuxGenDbUseMapAllInOneConnx(dbNameBase string, dbNum int, connxName string) map[string]string {
	config := map[string]string {dbNameBase:connxName}
	for i:=0; i < dbNum; i++ {
		config[fmt.Sprintf("%s_%d", dbNameBase, i)] = connxName
	}
	return config
}
func AuxMergeUseMap(mapList... map[string]string) map[string]string {
	mp := make(map[string]string)
	for i := range mapList {
		for k,v := range mapList[i] {
			mp[k] = v
		}
	}
	return mp
}

func AuxGenDbUseMapRoundRobinDbConnx(dbNameBase string, dbNum int, connxNameBase string, connxNum int) map[string]string {
	config := map[string]string {dbNameBase:connxNameBase}
	if connxNum <= 1 {
		return AuxGenDbUseMapAllInOneConnx(dbNameBase, dbNum, connxNameBase)
	}
	for i:=0; i < dbNum; i++ {
		config[fmt.Sprintf("%s_%d", dbNameBase, i)] = fmt.Sprintf("%s-%d", connxNameBase, i%connxNum)
	}
	return config
}


func (c AppConfig) String() string {
	b,_ := json.Marshal(&c)
	return string(b)
}


func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		LogConf: log.PassLagerCfg {
			Writers:        "file,stdout",
			LoggerLevel:    "DEBUG",
			LoggerFile:     "../logs/jcs.log",
			LogFormatText:  true, //true with plain text false with json
			RollingPolicy:  "size",
			LogRotateDate:  1,  //enable with date policy
			LogRotateSize:  10, //MB
			LogBackupCount: 10, //Latest N
		},
		Listeners: []ServerListener{
			{
				Name: "main",
				Bind: ":8080",
			},
		},
		DbConnxs: []DBConnection{
			{
				Name:         "dev",
				Enable:       true,
				Backend:      "sqlite3",
				ConnxEnv:     "../db/dev.db",
				MaxIdleConnx: 120,
				MaxOpenConnx: 1024,
			},
			{
				Name:         "test",
				Enable:       false,
				Backend:      "mysql",
				ConnxEnv:     "test:test123@/127.0.01:3306/jcs?charset=utf8&parseTime=True&loc=Local",
				MaxIdleConnx: 120,
				MaxOpenConnx: 1024,
			},
			{
				Name:         "prod-account",
				Enable:       false,
				Backend:      "mysql",
				ConnxEnv:     "test:test123@/127.0.01:3306/jcs?charset=utf8&parseTime=True&loc=Local",
				MaxIdleConnx: 120,
				MaxOpenConnx: 1024,
			},
			{
				Name:         "prod-msg",
				Enable:       false,
				Backend:      "mysql",
				ConnxEnv:     "test:test123@/127.0.01:3306/jcs?charset=utf8&parseTime=True&loc=Local",
				MaxIdleConnx: 120,
				MaxOpenConnx: 1024,
			},
			{
				Name:         "prod-record",
				Enable:       false,
				Backend:      "mysql",
				ConnxEnv:     "test:test123@/127.0.01:3306/jcs?charset=utf8&parseTime=True&loc=Local",
				MaxIdleConnx: 120,
				MaxOpenConnx: 1024,
			},
			{
				Name:         "prod-misc",
				Enable:       false,
				Backend:      "mysql",
				ConnxEnv:     "test:test123@/127.0.01:3306/jcs?charset=utf8&parseTime=True&loc=Local",
				MaxIdleConnx: 120,
				MaxOpenConnx: 1024,
			},
		},
		AuthRsaPrivateKey: "../keys/rsa_pri.pem",
		RunEnv:            "dev",
		Brokers: []BrokerConnx{
			{
				Name:         "dev",
				Enable:       true,
				Backend:      "redis",
				ConnxEnv:     "127.0.0.1:6379",
				MaxIdleConnx: 120,
				MaxOpenConnx: 1024,
			},
		},
		Event: EventDispatcherCfg{
			BrokerName:    "dev",
			DispatcherNum: 2,
		},
		TaskBrokerName: "dev",
		DevEnv:RunDevEnvConfig {
			UploadFileDir: "../upload",
		},
		Account: AccountSvcConfig {
			UseDb: SvcUseDbConfig{
				SplitDbNum: 1,
				SplitTbNum: 1,
				DbUseMap:   AuxGenDbUseMapAllInOneConnx("db_account", 1, "prod-account"),
			},
		},
		Captcha: CaptchaSvcConfig {
			UseDb: SvcUseDbConfig{
				SplitDbNum: 1,
				SplitTbNum: 1,
				DbUseMap:   AuxGenDbUseMapAllInOneConnx("db_captcha", 1, "prod-account"),
			},
		},
		Friend: FriendSvcConfig {
			UseDb: SvcUseDbConfig{
				SplitDbNum: 1,
				SplitTbNum: 1,
				DbUseMap:   AuxGenDbUseMapAllInOneConnx("db_friend", 1, "prod-account"),
			},
		},
		Chat: ChatSvcConfig {
			UseDb: SvcUseDbConfig{
				SplitDbNum: 1,
				SplitTbNum: 1,
				DbUseMap:   AuxGenDbUseMapAllInOneConnx("db_chat", 1, "prod-account"),
			},
			MsgSvcUrl: "127.0.0.1:8080",
		},
		Msg: MsgSvcConfig {
			EnableSplitDbBegIdx:0,
			EnableSplitDbEndIdx:0,
		},
		AuthToken: AuthTokenConfig{
			SessionSignKeySecret: "1234567890abcdef",
			JwtSignKeySecret: "1234567890abcdef",
			SmsCodeSignKeySecret: "1234567890abcdef",
		},
		Profile: ProfileSvcConfig{
			UseDb: SvcUseDbConfig{
				SplitDbNum: 1,
				SplitTbNum: 1,
				DbUseMap:   AuxGenDbUseMapAllInOneConnx("db_profile", 1, "prod-account"),

			},
		},
		Notify: NotifySvcConfig {
			UseDb: SvcUseDbConfig{
				SplitDbNum: 1,
				SplitTbNum: 1,
				DbUseMap:   AuxGenDbUseMapAllInOneConnx("db_notify", 1, "prod-account"),
			},
		},
	}
}
