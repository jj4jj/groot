package stream

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"groot/proto/cserr"
	"groot/proto/csmsg"
)

type (
	CSStreamRpcHandler func(ctx context.Context, req proto.Message, rsp proto.Message) cserr.CSErrCode
	CSStreamLogic struct {
		LogicBase
		handler 	map[csmsg.CSMsgCmd]CSStreamRpcHandler
	}
)

