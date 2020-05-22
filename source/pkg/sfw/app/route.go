package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"groot/proto/cserr"
	"groot/pkg/util"
	"net/http"
	"reflect"
)
type (
	ServiceRouteHandler interface{}	//func(c *Context) (rsp interface{}, e cserr.ICSErrCodeError)
	ServiceMethodCtx struct {
		method       string
		routeHandler ServiceRouteHandler
		rawhandler   gin.HandlerFunc
		middlewares  []string
		caller 		 reflect.Value
	}
	RouteResponse struct {
		Code    int32         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func checkServiceRouteHandlerType(handler ServiceRpcHandler) bool {
	ft := reflect.TypeOf(handler)
	if !(ft.Kind() == reflect.Func && ft.NumIn() == 1 && ft.NumOut() == 2) {
		log.Errorf(nil,
			"route handler type is error need:[func(*gin.Context)(rsp Response,err cserr.ICSErrCodeError)]")
		return false
	}
	ctxType := ft.In(0)
	errType := ft.Out(1)
	errInterfaceType := reflect.TypeOf(new(cserr.ICSErrCodeError)).Elem()
	if !(errType.Implements(errInterfaceType) && ctxType == reflect.TypeOf(&gin.Context{})) {
		log.Errorf(nil,
			"route handler type is error need:[func(*gin.Context)(rsp Response,err cserr.ICSErrCodeError)]")
		return false
	}
	return true
}

func createRouteWrapper(meth *ServiceMethodCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//异常发生.统一返回内部错误
				log.Warnf("route internal error panic:%v", err)
				response := RouteResponse{
					Code:    int32(cserr.CSErrCode_CS_ERR_INTERNAL),
					Message: cserr.CSErrCode_CS_ERR_INTERNAL.String(),
				}
				c.JSON(http.StatusOK, response)
			}
		}()
		if meth.routeHandler != nil {
			in := []reflect.Value{reflect.ValueOf(c)}
			out := meth.caller.Call(in)
			//rsp,eno := meth.routeHandler(c)
			rsp := out[0].Interface()
			eno := out[1].Interface().(cserr.ICSErrCodeError)
			code, message := util.ErrCode(eno), util.ErrStr(eno)
			response := RouteResponse{
				Code:    int32(code),
				Data:    rsp,
				Message: message,
			}
			c.JSON(http.StatusOK, response)
		} else if meth.rawhandler != nil {
			meth.rawhandler(c)
		} else {
			log.Errorf(nil, "path:%s routeHandler is nil !", c.Request.URL.Path)
			response := RouteResponse{
				Code:    int32(cserr.CSErrCode_CS_ERR_ACCOUNT),
				Message: cserr.CSErrCode_CS_ERR_INTERNAL.String(),
			}
			c.JSON(http.StatusNotImplemented, response)
		}
	}
}

