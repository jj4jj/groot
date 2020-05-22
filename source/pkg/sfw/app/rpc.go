package app

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/lexkong/log"
	"io/ioutil"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/pkg/util"
	"net/http"
	"reflect"
	"time"
)

type (
	ServiceRpcHandler   interface{} //func(c *gin.Context, req Request,rsp Response) (e cserr.ICSErrCodeError)
	ServiceRpcCtx struct {
		handler     ServiceRpcHandler
		middlewares []string
		reqType     reflect.Type
		rspType     reflect.Type
		caller 		reflect.Value
	}
)

func checkServiceRpcHandlerType(handler ServiceRpcHandler) bool {
	ft := reflect.TypeOf(handler)
	if !(ft.Kind() == reflect.Func && ft.NumIn() == 3 && ft.NumOut() == 1) {
		log.Errorf(nil,
			"rpc handler type is error need:[func(*app.Context,req proto.Message,rsp Response)cserr.ICSErrCodeError]")
		return false
	}
	ctxType := ft.In(0)
	reqType := ft.In(1)
	rspType := ft.In(2)
	errType := ft.Out(0)
	protoMessageInterfaceType := reflect.TypeOf(new(proto.Message)).Elem()
	errInterfaceType := reflect.TypeOf(new(cserr.ICSErrCodeError)).Elem()
	if !(reqType.Implements(protoMessageInterfaceType) && rspType.Implements(protoMessageInterfaceType) &&
		errType.Implements(errInterfaceType) && ctxType == reflect.TypeOf(&gin.Context{})) {
		log.Errorf(nil,
			"rpc handler type is error need:[func(*gin.Context,req proto.Message,rsp Response)cserr.ICSErrCodeError]")
		return false
	}
	return true
}


func createRpcWrapper(rpc *ServiceRpcCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//异常发生.统一返回内部错误
				log.Warnf("call rpc internal error:%v", err)
				HttpRpcReply(c, cserr.ErrInternal, nil)
			}
		}()
		//获取请求结构(base+req) byte
		var bodyBytes []byte
		bodyBytes, e := ioutil.ReadAll(c.Request.Body)
		if e != nil {
			log.Error("Request Body is nil...please check the progress", e)
			HttpRpcReply(c, cserr.ErrProto, nil)
			return
		}
		//取参数
		req := csmsg.CSMsgHttpFrame{}
		if e := proto.Unmarshal(bodyBytes, &req); e != nil {
			log.Errorf(e, "unpack frame a error")
			HttpRpcReply(c, cserr.ErrProto, nil)
			return
		}

		//
		c.Set("cs.req.frame.head", req.ReqHead)
		//req.ReqHead.Cookie
		csReq := reflect.New(rpc.reqType)
		if e := proto.Unmarshal(req.Payload, csReq.Interface().(proto.Message)); e != nil {
			log.Errorf(e, "unpack frame a error")
			HttpRpcReply(c, cserr.ErrProto, nil)
			return
		}
		log.Debugf("rpc path:%s|req:%v", c.FullPath(), util.GetMsgShortStr(csReq))
		//msg
		csRsp := reflect.New(rpc.rspType)
		in := []reflect.Value{reflect.ValueOf(c), csReq, csRsp}
		out := rpc.caller.Call(in)
		var err cserr.ICSErrCodeError
		if out[0].Interface() != nil {
			err = out[0].Interface().(cserr.ICSErrCodeError)
		}
		//err := rpc.routeHandler(c, csReq, csRsp)
		if util.ErrOK(err) {
			rsp := csRsp.Interface().(proto.Message)
			log.Debugf("rpc path:%s|success|rsp:%s", c.FullPath(), util.GetMsgShortStr(rsp))
			HttpRpcReply(c, err, rsp)
		} else {
			log.Debugf("rpc path:%s|fail|error:%d", c.FullPath(), util.ErrCode(err))
			HttpRpcReply(c, err, nil)
		}
	}
}


//legacy
func HttpRpcReply(c *gin.Context, err cserr.ICSErrCodeError, pb proto.Message) {
	reqHead, rok := c.Get("cs.req.frame.head")
	if !rok  {
		log.Errorf(err,"error get cs context for rpc reply ")
		return
	}
	code := util.ErrCode(err)
	msg := util.ErrStr(err)
	head := reqHead.(*csmsg.CSMsgReqHead)
	csRsp := &csmsg.CSMsgHttpFrame {
		FrameType: csmsg.CSFrameType_CS_FRAME_RSP,
		RspHead : &csmsg.CSMsgRspHead {
			Cmd: csmsg.CSMsgCmd(int32(head.Cmd) + 1),
			Cookie: head.Cookie,
			Time: time.Now().Unix(),
			Code: cserr.CSErrCode(code),
			Msg: msg,
			Seq: head.Seq,
		},
	}

	if pb != nil {
		pbData, e := proto.Marshal(pb)
		if e != nil {
			log.Errorf(e, "Proto Marshal ResponseData ErrStr")
			pbData = []byte{}
		}
		csRsp.Payload = pbData
	}
	responseData, e := proto.Marshal(csRsp)
	if e != nil {
		log.Error("Proto Marshal ResponseData ErrStr", e)
		c.Data(http.StatusInternalServerError, "application/octet-stream",  nil)
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", responseData)
}
