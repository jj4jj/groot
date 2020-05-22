package csuser

import (
	"fmt"
	"groot/proto/comm"
	"groot/proto/csmsg"
	sfw_test "groot/sfw/test"
	"groot/sfw/util"
	"time"
)

type (
	CSUser struct {
		uin string
		auth * csmsg.CSLoginUserAuthInfo
		server	*sfw_test.AppServer
		stream  *sfw_test.ServiceStream
		friend_list []*csmsg.FriendInfo
		record_book_list []*csmsg.RecordBookBrief
		record_node_list []*csmsg.RecordNodeInfo
		moment_node_list []*csmsg.RecordNodeInfo
		check_session_msg_seq uint64
	}
)
var (
	serverRoot = "http://127.0.0.1:8080"
)


func GetUser(tel string) (user *CSUser) {
	server := sfw_test.NewServer(serverRoot, "main")
	user = nil
	defer util.CatchExcpetion(nil,"get user")
	req1 := csmsg.CSMsgAccountLoginOrRegisterGetSmsCodeReq {
		Tel:  tel,
	}
	var rsp1 csmsg.CSMsgAccountLoginOrRegisterGetSmsCodeRsp
	err := server.Call("account.login_or_register_get_sms_code", &req1, &rsp1)
	if err != nil {
		panic(err)
	}
	//ok
	deviceId := "device_id-1234"
	fmt.Println(rsp1.String())
	req2 := csmsg.CSMsgAccountLoginOrRegisterReq{
		AuthType:  comm.UserAuthType_USER_AUTH_SMS_CODE,
		AuthId:       tel,
		AuthCode: rsp1.Ticket,
		Ticket:         rsp1.Ticket,
		LoginDevice: &csmsg.CSUserLoginDevice{
			DeviceId: deviceId,
			DeviceName: "iphone-x",
			DeviceType: "",
			Imei:"XXX",
			LoginScene: 1,
			Language:"zh_CN",
			PushToken:tel,
		},
	}
	var rsp2 csmsg.CSMsgAccountLoginOrRegisterRsp
	err = server.Call("account.login_or_register", &req2, &rsp2)
	if err != nil {
		panic(err)
	}

	//
	server.SetBase(rsp2.Uin, deviceId, "device_type-123")
	server.SetJwt(rsp2.AuthInfo.Token)
	fmt.Println("uin:%s expired time:%s", rsp2.Uin, time.Unix(rsp2.AuthInfo.ExpiredTimeStamp,0).String())

	//ok
	user = & CSUser{
		uin : rsp2.Uin,
		auth: rsp2.AuthInfo,
		stream: server.StreamDial("notify.push"),
	}
	return user
}
