package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"groot/comm/constk"
	"strconv"
	//"github.com/golang/protobuf/proto"
	"github.com/teris-io/shortid"
	"math/rand"
	"time"
)

var (
	UtilRandom *rand.Rand
)

func Init() {
	//初始化常用组件
	UtilRandom = rand.New(rand.NewSource(time.Now().UnixNano()))

}

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GinGetRequestID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	return ""
}

func GetMsgShortStr(msg interface{}) string {
	ret := fmt.Sprintf("%v", msg)
	tlen := len(ret)
	if tlen > 128 {
		return ret[:128] + " ...(" + strconv.Itoa(tlen) + ")"
	} else {
		return ret
	}
}

func GetBrokerTopicName(topicType constk.BrokerTopicType, args ...interface{}) string {
	var ret = string(topicType)
	for _, v := range args {
		ret = ret + ":" + fmt.Sprint(v)
	}
	return ret
}
