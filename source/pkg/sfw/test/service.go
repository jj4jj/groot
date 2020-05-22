package sfw_test

import (
	"bytes"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"io/ioutil"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"net/http"
	"os"
	"strings"
	"golang.org/x/net/websocket"
	"time"
)

type (
	AppServer struct {
		Name string
		Addr string
		JwtTicket string
		BaseCtx		csmsg.AuthContext
	}
)

func NewServer(addr, name string) *AppServer {
	return &AppServer{
		Name: name,
		Addr: addr,
	}
}
func (s *AppServer) SetJwt(ticket string){
	s.JwtTicket = ticket
}
func (s *AppServer) SetBase(Uin,DeviceId,DeviceType string) {
	s.BaseCtx.DeviceId = DeviceId
	s.BaseCtx.DeviceType = DeviceType
	s.BaseCtx.Uin = Uin
}

type ServiceStream struct {
	ws	*websocket.Conn
}
func (stream *ServiceStream) Send(request proto.Message) error {
	b,e := proto.Marshal(request)
	if e != nil {
		return e
	}
	_,r := stream.ws.Write(b)
	return r
}
func (stream *ServiceStream) Recv(response proto.Message) error {
	buffer := make([]byte, 1024*256)
	_,e := stream.ws.Read(buffer)
	if e != nil {
		return e
	}
	e = proto.Unmarshal(buffer, response)
	return e
}

func (stream *ServiceStream) Close() {
	stream.ws.Close()
}



func (s *AppServer) StreamDial(path string) *ServiceStream {
	url := s.Addr + "/" + strings.ReplaceAll(path, ".", "/")
	url = strings.ReplaceAll(url,"http://","ws://")
	ws,err := websocket.Dial(url, "ws",  s.Addr)
	if err !=  nil {
		fmt.Fprintf(os.Stderr,"connect websocket server err:%v path:%s  !\n", err.Error(), url)
		return nil
	}
	return & ServiceStream{
		ws: ws,
	}
}

func (s *AppServer)Get(path string) (string, error) {
	url := s.Addr + "/" + strings.ReplaceAll(path, ".", "/")
	hreq, he := http.NewRequest("GET", url, nil)
	if he != nil {
		fmt.Fprintf(os.Stderr, "http request path:%s req error:%v\n", url,  he)
		return "",he
	}
	resp, he2 := http.DefaultClient.Do(hreq)
	if he2 != nil {
		fmt.Fprintf(os.Stderr, "http client do request error:%v\n", he2)
		return "",he2
	}
	respBody, e := ioutil.ReadAll(resp.Body)
	return string(respBody),e
}

func (s *AppServer)Post(path string, body []byte) (string, error) {
	url := s.Addr + "/" + strings.ReplaceAll(path, ".", "/")
	hreq, he := http.NewRequest("POST", url, bytes.NewReader(body))
	if he != nil {
		fmt.Fprintf(os.Stderr, "http request path:%s req len:%d error:%v\n", url, body, he)
		return "",he
	}
	resp, he2 := http.DefaultClient.Do(hreq)
	if he2 != nil {
		fmt.Fprintf(os.Stderr, "http client do request error:%v\n", he2)
		return "",he2
	}
	respBody, e := ioutil.ReadAll(resp.Body)
	return string(respBody),e
}

func (s *AppServer) Call(path string, req proto.Message, rsp proto.Message) error {
	b, e := proto.Marshal(req)
	if e != nil {
		fmt.Fprintf(os.Stderr, "call path:%s request mashal1 error:%v\n", path, e)
		return e
	}
	url := s.Addr + "/" + strings.ReplaceAll(path, ".", "/")

	frame := csmsg.CSMsgHttpFrame{
		FrameType: csmsg.CSFrameType_CS_FRAME_REQ,
		ReqHead: &csmsg.CSMsgReqHead {
			Time: time.Now().Unix(),
			Context: &s.BaseCtx,
		},
	}
	frame.Payload,_ = proto.Marshal(req)
	nb, ne := proto.Marshal(&frame)
	if ne != nil {
		fmt.Fprintf(os.Stderr, "call path:%s request mashal2 error:%v\n", path, ne)
		return ne
	}
	hreq, he := http.NewRequest("POST", url, bytes.NewReader(nb))
	if he != nil {
		fmt.Fprintf(os.Stderr, "http request path:%s req len:%d error:%v\n", url, len(b), he)
		return he
	}
	if s.JwtTicket != " "{
		hreq.Header.Set("Authorization", "Bearer " + s.JwtTicket)
	}
	resp, he2 := http.DefaultClient.Do(hreq)
	if he2 != nil {
		fmt.Fprintf(os.Stderr, "http client do request error:%v\n", he2)
		return he2
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	var allrsp csmsg.CSMsgHttpFrame
	e = proto.Unmarshal(respBody, &allrsp)
	if e != nil {
		fmt.Fprintf(os.Stderr, "proto resp unmarshal rsp len:%d error:%v\n", len(respBody), e)
		return e
	}
	fmt.Fprintf(os.Stdout, "Header:%v\nBase:%v\n", resp.Header, allrsp.RspHead)

	if allrsp.RspHead.Code == 0 {
		e = proto.Unmarshal(allrsp.Payload, rsp)
		if e != nil {
			fmt.Fprint(os.Stderr, "app resp unmashal error len:%d !", len(allrsp.Payload))
			return e
		}
	} else {
		return &cserr.CSErrCodeError{Code: allrsp.RspHead.Code, Message:allrsp.RspHead.Msg,}
	}
	return nil
}
