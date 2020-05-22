package tests

import (
	"fmt"
	"groot/proto/comm"
	"groot/proto/csmsg"
	"groot/tests/csuser"
	"testing"
)

func FriendAll(t *testing.T) {

	csuser1 := csuser.GetUser("15102991436")
	csuser2 := csuser.GetUser("15102991437")
	//friend 1

	//search user 1

	req1 := csmsg.CSMsgFriendAskReq {
		Tel:  "15102991436",
	}
	var rsp1 csmsg.CSMsgAccountLoginOrRegisterGetSmsCodeRsp
	err := MainServer.Call("account.login_or_register_get_sms_code", &req1, &rsp1)
	if err != nil {
		t.Error(err)
	}
	//ok
	fmt.Println(rsp1.String())
	req2 := csmsg.CSMsgAccountLoginOrRegisterReq{
		AuthType:  comm.UserAuthType_USER_AUTH_SMS_CODE,
		AuthId:       "15102991436",
		AuthCode: rsp1.Ticket,
		Ticket:         rsp1.Ticket,
		LoginDevice: &csmsg.CSUserLoginDevice{
			DeviceId: "1234",
			DeviceName: "iphone-x",
			DeviceType: "",
			Imei:"XXX",
			LoginScene: 1,
			Language:"zh_CN",
		},
	}
	var rsp2 csmsg.CSMsgAccountLoginOrRegisterRsp
	err = MainServer.Call("account.login_or_register", &req2, &rsp2)
	if err != nil {
		t.Error(err)
	}
	//ok
	fmt.Println(rsp2.String())

	//get key
	MainServer.SetJwt(rsp2.AuthInfo.Token)
	MainServer.SetBase(rsp2.Uin, req2.LoginDevice.Imei, req2.LoginDevice.DeviceType)
	//logout
	req3 := csmsg.CSMsgAccountLogoutReq {
	}
	var rsp3 csmsg.CSMsgAccountLogoutRsp
	err = MainServer.Call("account.logout", &req3, &rsp3)
	if err != nil {
		t.Error(err)
	}
	//ok
	fmt.Println(rsp3.String())

	//log out again
	err = MainServer.Call("account.logout", &req3, &rsp3)
}

