package tests

import (
	"groot/service/account/model"
	"groot/sfw/db"
	"testing"
)


func TestReflect(t *testing.T) {


	name := sfw_db.GetModelTableNameBase(&model.DbUser{})
	if name == "" {
		t.Errorf("get table name error\n")
		t.Fail()
	}
	t.Logf("get db user tb name base:%s\n", name)
}
