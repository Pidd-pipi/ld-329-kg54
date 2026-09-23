package errors

import "net/http"

// 预约协商流程的业务错误，错误码与消息集中维护
var (
	ErrAppointmentMatchNotFound  = NewBusinessError(http.StatusNotFound, "MATCH_NOT_FOUND", "匹配不存在")
	ErrAppointmentNotFound       = NewBusinessError(http.StatusNotFound, "APPOINTMENT_NOT_FOUND", "预约不存在或已撤下")
	ErrAppointmentNotParticipant = NewBusinessError(http.StatusForbidden, "NOT_PARTICIPANT", "只有匹配双方才能操作该预约")
	ErrAppointmentOnlyResponder  = NewBusinessError(http.StatusForbidden, "ONLY_RESPONDER", "只有对方才能拒绝该预约")
	ErrAppointmentOnlyInitiator  = NewBusinessError(http.StatusForbidden, "ONLY_INITIATOR", "只有发起方才能修改预约")
	ErrAppointmentActiveExists   = NewBusinessError(http.StatusConflict, "ACTIVE_APPOINTMENT_EXISTS", "同一匹配已有进行中的预约")
	ErrAppointmentNotPending     = NewBusinessError(http.StatusConflict, "APPOINTMENT_NOT_PENDING", "预约已生效，无法再修改或拒绝")
	ErrAppointmentInvalidSlot    = NewBusinessError(http.StatusBadRequest, "INVALID_TIME_SLOT", "只能选择双方共同可约的时间")
	ErrAppointmentPlaceRequired  = NewBusinessError(http.StatusBadRequest, "PLACE_REQUIRED", "请填写交换地点")
	ErrAppointmentAgendaRequired = NewBusinessError(http.StatusBadRequest, "AGENDA_REQUIRED", "请填写协商议程")
	ErrAppointmentFieldTooLong   = NewBusinessError(http.StatusBadRequest, "FIELD_TOO_LONG", "地点或议程超出长度限制")
	ErrAppointmentInvalidRequest = NewBusinessError(http.StatusBadRequest, "INVALID_REQUEST", "请求参数不完整")
)
