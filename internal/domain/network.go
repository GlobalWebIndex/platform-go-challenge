package domain

type NotificationMessage struct {
	Token   string `json:"token"`
	MsgType string `json:"msg_type"`
	Msg     string `json:"msg"` // json string.
}
