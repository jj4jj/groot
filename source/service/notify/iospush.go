package notify

import (
	"encoding/hex"
	"github.com/gogo/protobuf/proto"
	"github.com/lexkong/log"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"groot/comm/conf"
	"groot/proto/comm"
	"groot/proto/csmsg"
	"groot/service/account/crud"
	crud2 "groot/service/profile/crud"
	"groot/sfw/db"
)

//var apnsToken *token.Token = nil
type IOSApnsPushedData struct {
	Title        string
	Desc         string
	Scene        string
	SendUserName string
	MessageID    string
	BookID       string
	NodeID       string
}

func GetPushedDataFromCsUserEvent(title string, event *comm.CSUserEvent) *IOSApnsPushedData {
	//查询Setting
	setting := crud2.DbUserSettingGetByUin(event.Uin)
	if setting == nil {
		return nil
	}
	if event.EvtType != comm.UserEventType_USR_EVT_FRIEND_MSG_NEW {
		return nil
	}

	param := csmsg.NotifyUserEventFriendMsgNew{}
	if e := proto.Unmarshal(event.EvtParam, &param); e != nil {
		log.Errorf(e, "unmarshal msg event param fail")
		return nil
	}

	msgType := param.MsgType
	desc := ""
	switch msgType {
	case comm.MsgType_MSG_TYPE_TEXT:
		desc = "发来了一个消息"
	case comm.MsgType_MSG_TYPE_REPLY:
		desc = "回复了你"
	case comm.MsgType_MSG_TYPE_CUSTOM_EMOTION:
		desc = "发来了一个表情"
	case comm.MsgType_MSG_TYPE_LOCATION:
		desc = "发来了一个位置"
	case comm.MsgType_MSG_TYPE_FILE:
		desc = "发来了一个文件"
	case comm.MsgType_MSG_TYPE_SHARE_URL:
		desc = "发来了一个网址"
	case comm.MsgType_MSG_TYPE_APPMSG:
		desc = "发来了一个应用消息"
	case comm.MsgType_MSG_TYPE_VOICE:
		desc = "发来了一段语音"
	case comm.MsgType_MSG_TYPE_VOIP:
		desc = "发来了一个电话"
	case comm.MsgType_MSG_TYPE_VIDEO:
		desc = "发来了一段视频"
	case comm.MsgType_MSG_TYPE_USERCARD:
		desc = "发来了一个名片"
	case comm.MsgType_MSG_TYPE_MEDIAS:
		desc = "发来了一个文件"
	}

	settingV := setting.GetJsonUserSetting()
	if settingV.PushHideDetail {
		desc = "发来了一条消息"
	}

	iOSApnsPushedData := &IOSApnsPushedData{
		Title:        title,
		Desc:         desc,
		Scene:        "MSG",
		SendUserName:  param.FriendUin,
		MessageID:     event.EventId,
	}
	return iOSApnsPushedData
}

func SendApnsPushToClient(Uin string, pushedData *IOSApnsPushedData) {
	log.Debugf("the Uin %s is not online", Uin)
	//获取用户的所有设备调用推送接口
	deviceList := crud.DbUserLoginDeviceGetList(Uin, comm.DeviceOsType_OS_TYPE_IOS)
	for i := range deviceList {
		deviceInfo := &deviceList[i]
		log.Infof("devide : %s DeviceToken length:%d,Device Is Online:%b unreadcount:%d",
			deviceInfo.DeviceID, len(deviceInfo.DevicePushToken), deviceInfo.IsOnline, deviceInfo.UnreadCount)
		if deviceInfo.IsOnline && len(deviceInfo.DevicePushToken) > 0 {
			deviceInfo.UnreadCount++
			go IosPushToUser([]byte(deviceInfo.DevicePushToken), deviceInfo.UnreadCount, Uin, pushedData)
			log.Debugf("deviceInfo.UnreadCount is %d", deviceInfo.UnreadCount)
			crud.DbUserLoginDeviceUpdate(deviceInfo.Uin, deviceInfo.DeviceID, sfw_db.UpdateFieldsMap{"unread_count": deviceInfo.UnreadCount})
		}
	}
}

var client *apns2.Client = nil
func createClient() bool {
	config := conf.GetAppConfig()
	certPath := config.IosApns.DevApnsP12
	cert, e := certificate.FromP12File(certPath, "")
	if e != nil {
		log.Infof("has a error:%v", e)
		return false
	}
	client = apns2.NewClient(cert).Development()
	if client != nil {
		return true
	} else {
		log.Errorf(nil, "create apns client fail !")
	}
	return false
}

func IosPushToUser(DevicePushToken []byte, unreadCount uint32, UIN string, pushedData *IOSApnsPushedData) {
	if client == nil {
		if !createClient() {
			return
		}
	}
	log.Info("start to push user.")
	notification := &apns2.Notification{}
	notification.DeviceToken = hex.EncodeToString(DevicePushToken)
	notification.Topic = "cn.xinyuetech.joychat"
	alertMap := map[string]string{}
	log.Infof("body is %s", pushedData.Desc)
	alertMap["subtitle"] = pushedData.Title
	alertMap["body"] = pushedData.Desc
	// userinfoMap := map[string]string{}
	// userinfoMap["Scene"] = pushedData.Scene
	// userinfoMap["MessageID"] = pushedData.MessageID
	// userinfoMap["SenderUserName"] = pushedData.SendUserName
	notification.Payload = payload.NewPayload().Badge(int(unreadCount)).Sound("default").Alert(alertMap).Custom("Scene", pushedData.Scene).Custom("MessageID", pushedData.MessageID).Custom("SenderUserName", pushedData.SendUserName).Custom("BookID", pushedData.BookID).Custom("NodeID", pushedData.NodeID)
	// Custom("SenderUserName", pushedData.SendUserName).Custom("MessageID", pushedData.MessageID).Custom("Scene", pushedData.Scene)

	res, e := client.Push(notification)
	if e != nil {
		log.Errorf(e, "ios push fail")
		return
	}
	if res.Sent() {
		log.Debugf("ios push success")
		return
	}
}
