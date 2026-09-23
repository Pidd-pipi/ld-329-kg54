package service

import (
	"strings"

	"cyskillswap/internal/constants"
	apperrors "cyskillswap/internal/errors"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

func Appointments() []model.Appointment { return repository.ListAppointments() }

// CreateAppointment 由匹配任一方发起预约，发起方自动确认，等待对方回应。
func CreateAppointment(actor string, matchID int, timeSlot, place, agenda string) (model.Appointment, error) {
	match, ok := repository.FindMatch(matchID)
	if !ok {
		return model.Appointment{}, apperrors.ErrAppointmentMatchNotFound
	}
	if actor != match.Provider && actor != match.Learner {
		return model.Appointment{}, apperrors.ErrAppointmentNotParticipant
	}
	if err := validateAppointmentInput(match, timeSlot, place, agenda); err != nil {
		return model.Appointment{}, err
	}
	responder := match.Provider
	if actor == match.Provider {
		responder = match.Learner
	}
	item := model.Appointment{
		MatchID:     match.ID,
		Pair:        match.Provider + " ↔ " + match.Learner,
		Initiator:   actor,
		Responder:   responder,
		Time:        timeSlot,
		Place:       strings.TrimSpace(place),
		Agenda:      strings.TrimSpace(agenda),
		Status:      constants.AppointmentStatusPending,
		ConfirmedBy: []string{actor},
		Revision:    1,
	}
	created, ok := repository.AddAppointmentIfNoActive(item)
	if !ok {
		return model.Appointment{}, apperrors.ErrAppointmentActiveExists
	}
	return created, nil
}

// ConfirmAppointment 双方确认后预约生效；确认集合去重，重复确认不累计。
func ConfirmAppointment(id int, actor string) (model.Appointment, error) {
	return repository.UpdateAppointment(id, func(item *model.Appointment) error {
		if item.Status == constants.AppointmentStatusRejected {
			return apperrors.ErrAppointmentNotFound
		}
		if actor != item.Initiator && actor != item.Responder {
			return apperrors.ErrAppointmentNotParticipant
		}
		if item.Status == constants.AppointmentStatusConfirmed || containsString(item.ConfirmedBy, actor) {
			return nil
		}
		item.ConfirmedBy = append(item.ConfirmedBy, actor)
		if containsString(item.ConfirmedBy, item.Initiator) && containsString(item.ConfirmedBy, item.Responder) {
			item.Status = constants.AppointmentStatusConfirmed
		}
		return nil
	})
}

// RejectAppointment 对方拒绝后预约撤下，匹配可重新发起。
func RejectAppointment(id int, actor string) (model.Appointment, error) {
	return repository.UpdateAppointment(id, func(item *model.Appointment) error {
		if item.Status == constants.AppointmentStatusRejected {
			return apperrors.ErrAppointmentNotFound
		}
		if actor != item.Initiator && actor != item.Responder {
			return apperrors.ErrAppointmentNotParticipant
		}
		if actor != item.Responder {
			return apperrors.ErrAppointmentOnlyResponder
		}
		if item.Status != constants.AppointmentStatusPending {
			return apperrors.ErrAppointmentNotPending
		}
		item.Status = constants.AppointmentStatusRejected
		return nil
	})
}

// ReviseAppointment 发起方在对方确认前修改时间或议程，修改后对方需重新确认。
func ReviseAppointment(id int, actor, timeSlot, place, agenda string) (model.Appointment, error) {
	item, ok := repository.GetAppointment(id)
	if !ok || item.Status == constants.AppointmentStatusRejected {
		return model.Appointment{}, apperrors.ErrAppointmentNotFound
	}
	match, ok := repository.FindMatch(item.MatchID)
	if !ok {
		return model.Appointment{}, apperrors.ErrAppointmentMatchNotFound
	}
	if err := validateAppointmentInput(match, timeSlot, place, agenda); err != nil {
		return model.Appointment{}, err
	}
	return repository.UpdateAppointment(id, func(current *model.Appointment) error {
		if current.Status == constants.AppointmentStatusRejected {
			return apperrors.ErrAppointmentNotFound
		}
		if actor != current.Initiator && actor != current.Responder {
			return apperrors.ErrAppointmentNotParticipant
		}
		if actor != current.Initiator {
			return apperrors.ErrAppointmentOnlyInitiator
		}
		if current.Status != constants.AppointmentStatusPending {
			return apperrors.ErrAppointmentNotPending
		}
		current.Time = timeSlot
		current.Place = strings.TrimSpace(place)
		current.Agenda = strings.TrimSpace(agenda)
		current.ConfirmedBy = []string{actor}
		current.Revision++
		return nil
	})
}

func validateAppointmentInput(match model.Match, timeSlot, place, agenda string) error {
	if !containsString(match.CommonSlots, timeSlot) {
		return apperrors.ErrAppointmentInvalidSlot
	}
	if strings.TrimSpace(place) == "" {
		return apperrors.ErrAppointmentPlaceRequired
	}
	if strings.TrimSpace(agenda) == "" {
		return apperrors.ErrAppointmentAgendaRequired
	}
	if len([]rune(place)) > constants.AppointmentPlaceMaxLen || len([]rune(agenda)) > constants.AppointmentAgendaMaxLen {
		return apperrors.ErrAppointmentFieldTooLong
	}
	return nil
}

func containsString(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
