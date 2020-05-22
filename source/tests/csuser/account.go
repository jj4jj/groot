package csuser

import (
	"groot/proto/csmsg"
	"groot/sfw/util"
)

//user ticket
func (user *CSUser)SearchTel(tel string) string {
	searchReq := &csmsg.CSMsgAccountSearchUserReq{}
	searchReq.Search = tel
	rsp := &csmsg.CSMsgAccountSearchUserRsp{}
	if util.CheckError(user.server.Call("account.search", searchReq, rsp),""){
		return ""
	}
	return rsp.UserTicket
}
