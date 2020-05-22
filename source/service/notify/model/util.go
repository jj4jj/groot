package model

import (
	"fmt"
	"groot/proto/comm"
)

const DbUserNotifyEventIdReserveBase = 0x1000
func GetEventIdFromDbEventRoute(dbId uint64,dbidx,tbidx uint32) string {
	return fmt.Sprintf("e.%d.%d.%x",dbidx, tbidx, DbUserNotifyEventIdReserveBase +dbId)
}
func GetDbEventRouteFromEventId(eventId string) (dbId uint64,dbidx,tbidx uint32) {
	n,e := fmt.Sscanf(eventId,"e.%d.%d.%x", &dbidx, &tbidx, &dbId)
	if e != nil || n != 3 {
		return 0,0,0
	}
	dbId -= DbUserNotifyEventIdReserveBase
	return
}

func DbNotifyUserEventToCsUserEvent(event *DbNotifyUserEvent) *comm.CSUserEvent {
	return &comm.CSUserEvent{
		Uin :      event.Uin,
		EventId: event.UserEventId,
		EvtType:   event.EventType,
		IntParam:  event.IntParam,
		StrParam:  event.StrParam,
		EvtParam:  event.EvtParam,
		EvtState:  event.EvtState,
		EvtStateParam: event.EvtStateParam,
	}
}