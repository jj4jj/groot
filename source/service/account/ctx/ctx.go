package ctx

import (
	"github.com/gin-gonic/gin"
	"groot/service/account/model"
)

func NeedAuthContextGetUin(c *gin.Context) string {
	cUin, ok := c.Get("auth.UinSigner")
	if ok {
		return cUin.(string)
	} else {
		return ""
	}
}
func NeedAuthContextGetUserTokenState(c *gin.Context) *model.DbUserTokenState {
	cTokenState, ok := c.Get("db.UserTokenState")
	if ok {
		return cTokenState.(*model.DbUserTokenState)
	} else {
		return nil
	}
}

func NeedUserContextGetDbUser(c * gin.Context) * model.DbUser {
	cDbUser, ok := c.Get("db.User")
	if ok {
		return cDbUser.(*model.DbUser)
	} else {
		return nil
	}
}

