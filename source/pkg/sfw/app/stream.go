package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lexkong/log"
	"groot/pkg/util"
	"net"
	"net/http"
	"reflect"
	"sync/atomic"
	"time"
	"errors"
	"github.com/golang/protobuf/proto"
)

type (
	ServiceStreamCtx struct {
		logic             ServiceStreamLogic
		middlewares       []string
		reqType           reflect.Type
		heartBeatInterval int64
	}
	ServiceStream interface {
		Close()
		Send(response proto.Message) error
	}
	ServiceStreamCloseReason int
	ServiceStreamLogic interface {
		OnOpen(s ServiceStream, c *gin.Context) bool
		OnMsg(s ServiceStream, request proto.Message)
		OnClose(s ServiceStream, r ServiceStreamCloseReason)
	}
	ServiceWebsocketStream struct {
		ws 		* websocket.Conn
		send 	chan []byte
		close 	chan ServiceStreamCloseReason
		lastActiveTime time.Time
		reqInst 	proto.Message
		closed		int32
	}
)
const (
	StreamClosedActive				ServiceStreamCloseReason = 0
	StreamClosedByPeer				ServiceStreamCloseReason = 1
	StreamClosedWriteError			ServiceStreamCloseReason = 2
	StreamClosedHeartBeatTimeout	ServiceStreamCloseReason = 3
	StreamClosedOpenFail			ServiceStreamCloseReason = 4
	StreamClosedReadError			ServiceStreamCloseReason = 5
	StreamClosedByMsgParse			ServiceStreamCloseReason = 6
	StreamClosedException			ServiceStreamCloseReason = 7
)

const (
	SVC_MAX_STREAM_READ_BUFF_SIZE	int64 =	1024*512

)

var (
	ErrStreamClosed = errors.New("stream is closed")
	wsHttpUpgrader = websocket.Upgrader {
		ReadBufferSize:  1024 * 1024,
		WriteBufferSize: 1024 * 1024,
		CheckOrigin: func(r *http.Request) bool {
			//ignore origin check
			return true
		},
	}
)



func(stream *ServiceWebsocketStream) String() string {
	return fmt.Sprintf("stream:%p", stream.ws)
}

func (stream *ServiceWebsocketStream) Send(rsp proto.Message) error {
	b, e:= proto.Marshal(rsp)
	if util.CheckError(e,"ws:%v marshal rsp marshal error", stream) {
		return e
	}
	stream.send <- b
	return nil
}
func (stream *ServiceWebsocketStream) Close()  {
	stream.closeStream(StreamClosedActive)
}

func (stream *ServiceWebsocketStream) closeStream(reason ServiceStreamCloseReason)  {
	if atomic.CompareAndSwapInt32(&stream.closed, 0, 1){
		stream.close <- reason
	}
}

func (stream *ServiceWebsocketStream) closeStreamAndDispatch(streamCtx *ServiceStreamCtx, reason ServiceStreamCloseReason)  {
	if atomic.CompareAndSwapInt32(&stream.closed, 0, 1){
		util.CheckError(stream.ws.Close(),"stream:%s close reason:%d", stream, reason)
		streamCtx.logic.OnClose(stream, reason)
		stream.ws = nil
	}
}

func (stream *ServiceWebsocketStream)ReadMsgQueueAndSend(streamCtx *ServiceStreamCtx){
	//dispatch goroutine
	stream.ws.SetPingHandler(func(appData string) error {
		log.Debugf("ping msg recv will add read and write dead line")
		r1:=stream.ws.SetReadDeadline(time.Now().Add(time.Duration(2 * streamCtx.heartBeatInterval)*time.Second))
		r2:=stream.ws.SetWriteDeadline(time.Now().Add(time.Duration(2 * streamCtx.heartBeatInterval)*time.Second))
		if util.CheckError(r1,"") {
			return r1
		}
		if util.CheckError(r2,""){
			return r2
		}
		return nil
	})
	intervalCheck := time.Second*time.Duration(streamCtx.heartBeatInterval)
	ticker := time.NewTicker(intervalCheck)
	defer func(){
		defer ticker.Stop()
		close(stream.close)
		close(stream.send)
		stream.closeStreamAndDispatch(streamCtx, StreamClosedException)
	}()
	for {
		select {
		case <- ticker.C:
			expiredTime := stream.lastActiveTime.Add(intervalCheck+intervalCheck)
			if !time.Now().Before(expiredTime)  {
				//expired
				log.Errorf(nil, "stream:%s heart-beat timeout with last check time:%s",
					stream, stream.lastActiveTime.Format("2006-01-02T15:04:05"))
				stream.closeStreamAndDispatch(streamCtx, StreamClosedHeartBeatTimeout)
				return
			}
		case msg:= <- stream.send:
			stream.lastActiveTime = time.Now()
			err := stream.ws.WriteMessage(websocket.BinaryMessage, msg)
			if util.CheckError(err,"stream:%s write msg error will streamClosed", stream) {
				stream.closeStreamAndDispatch(streamCtx, StreamClosedWriteError)
				return
			}
		case closeReason:= <-stream.close:
			log.Debugf("stream:%s recv close reason:%d", stream, closeReason)
			stream.closeStreamAndDispatch(streamCtx, closeReason)
			return
		}
	}
}
func (stream *ServiceWebsocketStream)RecvMsgAndDispatch(streamCtx *ServiceStreamCtx){
	//read msg goroutine
	defer util.CatchExcpetion(nil,"stream:%s read goroutine", stream)
	stream.ws.SetReadLimit(SVC_MAX_STREAM_READ_BUFF_SIZE)
	for {
		mtype, message, err := stream.ws.ReadMessage()
		if err != nil {
			if stream.closed == 1 {
				log.Debugf("stream:%s has been closed", stream)
				return
			}
			if nerr,ok := err.(net.Error); ok && nerr.Timeout(){
				stream.close <- StreamClosedByPeer
				return
			}
			if websocket.IsUnexpectedCloseError(err) {
				log.Errorf(err,"stream:%s read close error", stream)
				stream.close <- StreamClosedReadError
				return
			}
			log.Errorf(err, "stream:%s read message error !", stream)
			return
		} else if mtype == websocket.BinaryMessage {
			stream.lastActiveTime = time.Now()
			if stream.reqInst != nil {
				if !util.CheckError(proto.Unmarshal(message, stream.reqInst),
					"stream:%s proto msg unmarshal error ", stream) {
					streamCtx.logic.OnMsg(stream, stream.reqInst)
				} else {
					stream.close <- StreamClosedByMsgParse
					return
				}
			}
		} else if mtype == websocket.TextMessage {
			stream.lastActiveTime = time.Now()
			log.Debugf("recv stream:%s text msg:[%s] len:%d", stream, string(message), len(message))
		} else {
			stream.lastActiveTime = time.Now()
			log.Warnf("stream:%s recv unkown type:%d msg len:%d", stream, mtype, len(message))
		}
	}
}

func createStreamWrapper(streamCtx *ServiceStreamCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		//转换为wsocket执行stream logic
		ws, err := wsHttpUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Errorf(err,"ws upgrade error")
			c.String(http.StatusUpgradeRequired, "The incorrect stream protocol\n")
			c.Abort()
			return
		}
		defer util.CatchExcpetion(nil, "stream:%s catch a exception")
		stream := &ServiceWebsocketStream {
			lastActiveTime: time.Now(),
			reqInst: nil,
			closed:0,
		}
		if streamCtx.reqType != reflect.TypeOf(nil) {
			stream.reqInst = reflect.New(streamCtx.reqType).Interface().(proto.Message)
		}
		if streamCtx.logic.OnOpen(stream, c) == false {
			log.Debugf("stream open fail")
			c.String(http.StatusForbidden, "The stream auth fail\n")
			c.Abort()
			if err = ws.Close(); err != nil {
				log.Errorf(err, "stream close error")
			}
			streamCtx.logic.OnClose(stream, StreamClosedOpenFail)
			return
		}
		//for send queue
		stream.ws = ws
		stream.send = make(chan []byte, 32)
		stream.close = make(chan ServiceStreamCloseReason)
		ws.SetCloseHandler(func(code int, text string) error {
			log.Warnf("stream:%s closed code:%d desc:%s ", stream, code, text)
			stream.closeStream(StreamClosedByPeer)
			return nil
		})
		go stream.RecvMsgAndDispatch(streamCtx)
		go stream.ReadMsgQueueAndSend(streamCtx)
	}
}