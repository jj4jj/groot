package constk

type BrokerTopicType string

const (
	BRK_TOPIC_SERVICE_EVENT  BrokerTopicType = "ServiceEvent"
	BRK_TOPIC_CLIENT_EVENT   BrokerTopicType = "ClientEvent"
	BRK_TOPIC_TASK_SERVICE   BrokerTopicType = "TaskService"
	BRK_TOPIC_NOTIFY_SERVICE BrokerTopicType = "NotifyService"
	BRK_TOPIC_CALL_SERVICE   BrokerTopicType = "CallService"
)
