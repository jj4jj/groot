package account

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"groot/proto/cserr"
	"groot/service/account/crud"
	"groot/service/account/model"
	"groot/sfw/app"
	"groot/sfw/crypto/token"
	"groot/sfw/util"
	"time"
)

//add context:db.UserTokenState
//add context:auth.UinSigner
func AccountNeedAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token
		ctx, e := token.ParseRequest(c)
		if e != nil {
			app.HttpRpcReply(c, cserr.ErrAccountTokenInvalid, nil)
			c.Abort()
			return
		}
		//verify expired
		tokenState := crud.DbUserTokenGet(ctx.UIN,ctx.DeviceID)
		if tokenState == nil {
			log.Errorf(nil, "token pass but not found device:%v", ctx)
			app.HttpRpcReply(c, cserr.ErrAccountTokenInvalid, nil)
			c.Abort()
			return
		}
		timeNow := time.Now().Unix()
		if timeNow > tokenState.JwtExpiredTime || timeNow < tokenState.JwtIssueTime {
			log.Errorf(nil, "token is expired time")
			app.HttpRpcReply(c, cserr.ErrAccountTokenInvalid, nil)
			c.Abort()
			return
		}
		c.Set("db.UserTokenState", tokenState)
		c.Set("auth.UinSigner", tokenState.Uin)

		/*
		// Read the body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		//Restore the io.ReadCloser to it's original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		//Check the jwt and baserequest
		request := & server.Request{}
		if e := proto.Unmarshal(bodyBytes, request); e != nil {
			app.HttpRpcReply(c, cserr.ErrProto, nil)
			c.Abort()
			return
		}

		//校验当前请求参数是否和签名一致，若不一致则返回失败。
		if strings.Compare(ctx.UIN, request.BaseRequest.UIN) != 0 ||
			strings.Compare(ctx.DeviceType, request.BaseRequest.DeviceType) != 0 ||
			strings.Compare(ctx.DeviceID, request.BaseRequest.DeviceID) != 0 ||
			ctx.ClientVersion != request.BaseRequest.ClientVersion {

			log.Warnf("Not Match Current ClientVersion is %u,ctx ClientVersion is %u ",
				request.BaseRequest.ClientVersion, ctx.ClientVersion)
			log.Warnf("Not Match Current DeviceID is %s,ctx DeviceID is %s ",
				request.BaseRequest.DeviceID, ctx.DeviceID)
			log.Warnf("Not Match Current DeviceType is %s,ctx DeviceType is %s ",
				request.BaseRequest.DeviceType, ctx.DeviceType)
			log.Warnf("Not Match Current UIN is %s,ctx UIN is %s ", request.BaseRequest.UIN, ctx.UIN)
			app.HttpRpcReply(c, cserr.ErrTokenVerify, nil)
			c.Abort()
			return
		}
		*/
		c.Next()
	}
}

func NeedUserContextGetDbUser(c *gin.Context) *model.DbUser {
	cDbUser,ok := c.Get("db.User")
	if ok {
		return cDbUser.(*model.DbUser)
	} else {
		return nil
	}
}

//add context:db.User
func AccountNeedUserMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		Uin, ok := c.Get("auth.UinSigner")
		if !ok || Uin == "" {
			log.Errorf(nil, "get auth.UinSigner fail , need_auth middleware order ?")
			app.HttpRpcReply(c, cserr.ErrAccountUserNotExist, nil)
			return
		}
		dbuser, err := crud.DbUserGetByUin(Uin.(string))
		if !util.ErrOK(err) || dbuser == nil {
			app.HttpRpcReply(c, cserr.ErrAccountUserNotExist, nil)
			return
		}
		c.Set("db.User", dbuser)
		c.Next()
	}
}
