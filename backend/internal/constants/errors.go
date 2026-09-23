package constants

const (
	ErrInvalidJSON       = "INVALID_REQUEST"
	ErrMatchNotFound     = "MATCH_NOT_FOUND"
	ErrInvalidActor      = "INVALID_ACTOR"
	ErrInvalidField      = "INVALID_FIELD"
	ErrAppointmentExists = "APPOINTMENT_EXISTS"
	ErrAppointmentGone   = "APPOINTMENT_NOT_FOUND"
	ErrNotResponder      = "NOT_RESPONDER"
	ErrNotPending        = "APPOINTMENT_NOT_PENDING"
)

const (
	MsgInvalidJSON       = "请求参数不合法"
	MsgMatchNotFound     = "匹配不存在或已失效"
	MsgInvalidActor      = "当前视角不属于该匹配的双方"
	MsgInvalidSlot       = "请选择双方的共同可约时间"
	MsgInvalidPlace      = "请填写交换地点"
	MsgInvalidAgenda     = "请填写交换议程"
	MsgAppointmentExists = "该匹配已有一条进行中的预约"
	MsgAppointmentGone   = "预约不存在或已被撤下"
	MsgNotResponder      = "只有回应一方可以确认或拒绝"
	MsgNotPending        = "该预约已确认，无法继续协商"
)
