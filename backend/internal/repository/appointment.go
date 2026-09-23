package repository

import (
	"strings"
	"sync"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

var (
	appointmentMu    sync.Mutex
	appointmentSeq   = 1
	appointmentStore = map[int]model.Appointment{}
)

func init() {
	// 仅保留一条已生效的预约作为示例；拒绝的预约会被直接撤下。
	seed := model.Appointment{
		ID:        appointmentSeq,
		MatchID:   1,
		Pair:      "林澈 ↔ 孟野",
		Initiator: "林澈",
		Responder: "孟野",
		Time:      "周六 10:00",
		Place:     "东校区湖边",
		Status:    constants.AppointmentConfirmed,
		Agenda:    "先拍宣传照，再约 2 次吉他课",
		Slots:     []string{"周六上午"},
		Version:   1,
	}
	appointmentStore[seed.ID] = seed
	appointmentSeq++
}

// ListAppointments returns all active appointments ordered by id.
func ListAppointments() []model.Appointment {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()

	result := make([]model.Appointment, 0, len(appointmentStore))
	for id := 1; id <= appointmentSeq; id++ {
		if appointment, ok := appointmentStore[id]; ok {
			result = append(result, appointment)
		}
	}
	return result
}

// GetAppointment returns a single appointment by id.
func GetAppointment(id int) (model.Appointment, bool) {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()

	appointment, ok := appointmentStore[id]
	return appointment, ok
}

// FindActiveAppointment returns the single active appointment of a match.
func FindActiveAppointment(matchID int) (model.Appointment, bool) {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()

	for _, appointment := range appointmentStore {
		if appointment.MatchID == matchID {
			return appointment, true
		}
	}
	return model.Appointment{}, false
}

// CreateAppointment stores a newly initiated appointment.
func CreateAppointment(appointment model.Appointment) model.Appointment {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()

	appointment.ID = appointmentSeq
	appointmentSeq++
	appointmentStore[appointment.ID] = appointment
	return appointment
}

// ReviseAppointment changes time/place/agenda and restarts confirmation.
func ReviseAppointment(id int, time string, place string, agenda string) (model.Appointment, bool) {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()

	appointment, ok := appointmentStore[id]
	if !ok {
		return model.Appointment{}, false
	}
	appointment.Time = strings.TrimSpace(time)
	appointment.Place = strings.TrimSpace(place)
	appointment.Agenda = strings.TrimSpace(agenda)
	appointment.Status = constants.AppointmentPending
	appointment.Version++
	appointmentStore[id] = appointment
	return appointment, true
}

// ConfirmAppointment marks an appointment as effective.
func ConfirmAppointment(id int) (model.Appointment, bool) {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()

	appointment, ok := appointmentStore[id]
	if !ok {
		return model.Appointment{}, false
	}
	appointment.Status = constants.AppointmentConfirmed
	appointmentStore[id] = appointment
	return appointment, true
}

// DeleteAppointment removes an appointment (e.g. after rejection).
func DeleteAppointment(id int) bool {
	appointmentMu.Lock()
	defer appointmentMu.Unlock()

	if _, ok := appointmentStore[id]; !ok {
		return false
	}
	delete(appointmentStore, id)
	return true
}
