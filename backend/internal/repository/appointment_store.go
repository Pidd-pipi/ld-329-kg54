package repository

import (
	"sync"

	"cyskillswap/internal/constants"
	apperrors "cyskillswap/internal/errors"
	"cyskillswap/internal/model"
)

var (
	appointmentMu    sync.Mutex
	appointmentSeq   = 3
	appointmentStore = seedAppointments()
)

func seedAppointments() []model.Appointment {
	return []model.Appointment{
		{ID: 1, MatchID: 1, Pair: "林澈 ↔ 孟野", Initiator: "林澈", Responder: "孟野", Time: "周六上午", Place: "东校区湖边", Agenda: "先拍宣传照，再约 2 次吉他课", Status: constants.AppointmentStatusConfirmed, ConfirmedBy: []string{"林澈", "孟野"}, Revision: 1},
		{ID: 2, MatchID: 2, Pair: "周芮 ↔ 许安", Initiator: "周芮", Responder: "许安", Time: "周二晚", Place: "线上会议室", Agenda: "导入问卷 CSV 并完成基础可视化", Status: constants.AppointmentStatusPending, ConfirmedBy: []string{"周芮"}, Revision: 1},
	}
}

// ListAppointments 返回有效预约（待确认 + 已确认），已拒绝撤下的不再返回。
func ListAppointments() []model.Appointment {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()
	result := make([]model.Appointment, 0, len(appointmentStore))
	for _, item := range appointmentStore {
		if item.Status == constants.AppointmentStatusRejected {
			continue
		}
		result = append(result, cloneAppointment(item))
	}
	return result
}

func GetAppointment(id int) (model.Appointment, bool) {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()
	for _, item := range appointmentStore {
		if item.ID == id {
			return cloneAppointment(item), true
		}
	}
	return model.Appointment{}, false
}

// AddAppointmentIfNoActive 在同一匹配没有有效预约时原子创建，保证只留一条有效预约。
func AddAppointmentIfNoActive(item model.Appointment) (model.Appointment, bool) {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()
	for _, existing := range appointmentStore {
		if existing.MatchID == item.MatchID && existing.Status != constants.AppointmentStatusRejected {
			return model.Appointment{}, false
		}
	}
	item.ID = appointmentSeq
	appointmentSeq++
	appointmentStore = append(appointmentStore, item)
	return cloneAppointment(item), true
}

// UpdateAppointment 在锁内对预约执行原子修改，业务校验放在 mutate 回调中。
func UpdateAppointment(id int, mutate func(*model.Appointment) error) (model.Appointment, error) {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()
	for i := range appointmentStore {
		if appointmentStore[i].ID == id {
			if err := mutate(&appointmentStore[i]); err != nil {
				return model.Appointment{}, err
			}
			return cloneAppointment(appointmentStore[i]), nil
		}
	}
	return model.Appointment{}, apperrors.ErrAppointmentNotFound
}

func cloneAppointment(item model.Appointment) model.Appointment {
	item.ConfirmedBy = append([]string(nil), item.ConfirmedBy...)
	return item
}
