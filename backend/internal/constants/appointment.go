package constants

// 预约协商状态：待确认 -> 双方确认生效 / 被拒绝撤下
const (
	AppointmentStatusPending   = "pending"
	AppointmentStatusConfirmed = "confirmed"
	AppointmentStatusRejected  = "rejected"
)

// 预约表单字段长度限制
const (
	AppointmentPlaceMaxLen  = 60
	AppointmentAgendaMaxLen = 200
)
