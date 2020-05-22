package tests

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"groot/proto/csmsg"
	"testing"
)

func f(pb proto.Message) {
	fmt.Println(pb)
}

func TestProtoCv(t *testing.T) {
	rsp := csmsg.CSMsgAccountLoginOrRegisterReq{

	}

	var pnil interface{}
	f(pnil.(proto.Message))
	pnil = &rsp
	f(nil)
	f(pnil.(proto.Message))
	//fmt.Println(nil.(proto.Message))
	//fmt.Println(rsp.(proto.Message))
	//fmt.Println((&rsp).(proto.Message))
}
