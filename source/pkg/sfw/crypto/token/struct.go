package token

// AuthContext is the context of the JSON web token

type AuthContext struct {
	UIN           string
	DeviceID      string
	ClientVersion uint32
	DeviceType    string
	TokenSecKey	  string 	//for check
}

type LongConnAuthContext struct {
	UIN           string
	DeviceID      string
	ClientVersion uint32
	DeviceType    string
}

type FriendContext struct {
	SendUserName    string
	ReceiveUserName string
	Timestamp       int64
	EventType       int
}

type FileContext struct {
	UIN           string
	DeviceID      string
	ClientVersion uint32
	DeviceType    string
}

type SMSCodeContext struct {
	UserName  string
	DeviceID  string
	TimeStamp int64
}
