package friend

import "C"
import (
	"github.com/gin-gonic/gin"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/account"
	"groot/service/account/ctx"
	"groot/service/friend/model"
	"groot/sfw/util"
)


func FriendAsk(c *gin.Context, req *csmsg.CSMsgFriendAskReq, rsp *csmsg.CSMsgFriendAskRsp)(err cserr.ICSErrCodeError){
	uin := ctx.NeedAuthContextGetUin(c)
	checkRet := account.VerifyUserCheckTicketWithUin(uin, req.UserTicket)
	if checkRet == nil {
		return cserr.ErrCheckTicket
	}
	err = DbFriendAsk(uin, checkRet.UinCheck, req.Content, checkRet.Source)
	return err
}
func FriendAnswer(c *gin.Context, req *csmsg.CSMsgFriendAnswerReq, rsp *csmsg.CSMsgFriendAnswerRsp)(err cserr.ICSErrCodeError){
	uin := ctx.NeedAuthContextGetUin(c)
	err = DbFriendAnswer(uin, req.Uin, req.Content, req.AgreeStatus == 1)
	return err
}
func FriendSetBlock(c *gin.Context, req *csmsg.CSMsgFriendSetBlockReq, rsp *csmsg.CSMsgFriendSetBlockRsp)(err cserr.ICSErrCodeError){
	// quest channel
	uin := ctx.NeedAuthContextGetUin(c)
	err = DbFriendSetBlock(uin, req.Uin, req.Block)
	return err
}
func FriendRemove(c *gin.Context, req *csmsg.CSMsgFriendRemoveReq, rsp * csmsg.CSMsgFriendRemoveRsp)(err cserr.ICSErrCodeError){
	// quest channel
	uin := ctx.NeedAuthContextGetUin(c)
	err = DbFriendRemove(uin, req.Uin)
	return err
}
func FriendList(c *gin.Context, req *csmsg.CSMsgFriendListReq , rsp *csmsg.CSMsgFriendListRsp)(err cserr.ICSErrCodeError){
	uin := ctx.NeedAuthContextGetUin(c)

	friendStatList,err := DbFriendGetList(uin,req.Offset, req.Limit, req.TagFilter)
	if util.CheckError(err, "friend get list error !") {
		return err
	}

	settingMap := make(map[string]*model.DbFriendSetting)
	var uinList []string
	for _,fs := range friendStatList {
		uinList = append(uinList, fs.FriendUin)
	}

	DbFriendBatchGetSetting(uin, uinList, settingMap)
	for _,fs := range friendStatList {
		setting := settingMap[fs.FriendUin]
		friendInfo := csmsg.FriendInfo {
			Base: &comm.FriendBase{
				Uin: fs.FriendUin,
				Source: fs.FriendSource,
				FriendType: fs.FriendType,
				Status: fs.FriendStatus,
			},
			UserSetting: & csmsg.FriendUserSetting{},
			UserNote: & csmsg.FriendUserNote{},
		}
		if setting != nil {
			friendInfo.UserNote = model.DbFriendUserSettingToCsUserNote(setting)
			friendInfo.UserSetting = model.DbFriendUserSettingToCsUserSetting(setting)
		}
		rsp.FriendList = append(rsp.FriendList, &friendInfo)
	}

	return nil
}
func FriendUpdateSetting(c *gin.Context, req *csmsg.CSMsgFriendUpdateSettingReq, rsp * csmsg.CSMsgFriendUpdateSettingRsp)(err cserr.ICSErrCodeError){
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	err = DbFriendSettingUpdate(uin, req.Uin, req.UserNote, req.UserSetting)
	return
}
