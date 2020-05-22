package tests

import (
	"fmt"
	"testing"
	. "groot/tests/test_env"
	"time"
)

func TestPing(t *testing.T) {
	//for testing
	//var rsp server.LogoutRequest
	for i:=0; i < 1; i++ {
		stream := MainServer.StreamDial("notify.ping")
		if stream == nil {
			t.Errorf("error connect")
			t.FailNow()
		}
		fmt.Println("recv one msg:")
	}
	time.Sleep(3*time.Second)
}

func TestNotify(t *testing.T) {
	r,e := MainServer.Get("noitify.test_broadcast")
	TestError(e, t)
	fmt.Println(r)
}
