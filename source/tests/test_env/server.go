package test_env_test

import (
	"fmt"
	sfw_test "groot/sfw/test"
	"testing"
)

var (
	ServerRoot = "http://127.0.0.1:8080"
	MainServer *sfw_test.AppServer
)

func init() {
	fmt.Println("testing init env create server with:", ServerRoot)
	MainServer = sfw_test.NewServer(ServerRoot, "main")
}


func TestError(err error, t *testing.T){
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
