package model

// Appointment 表示一次技能交换预约协商。
// 发起方创建后自动确认，对方确认前可修改，双方确认后生效；
// ConfirmedBy 为去重集合，重复确认不会累计。
type Appointment struct {
	ID          int      `json:"id"`
	MatchID     int      `json:"matchId"`
	Pair        string   `json:"pair"`
	Initiator   string   `json:"initiator"`
	Responder   string   `json:"responder"`
	Time        string   `json:"time"`
	Place       string   `json:"place"`
	Agenda      string   `json:"agenda"`
	Status      string   `json:"status"`
	ConfirmedBy []string `json:"confirmedBy"`
	Revision    int      `json:"revision"`
}
