package service

import (
	"fmt"
	"strings"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

func appointmentPair(match model.Match) string {
	return fmt.Sprintf("%s ↔ %s", match.Provider, match.Learner)
}

func resolveParties(match model.Match, actor string) (string, errors.BusinessError) {
	switch strings.TrimSpace(actor) {
	case match.Provider:
		return match.Learner, errors.BusinessError{}
	case match.Learner:
		return match.Provider, errors.BusinessError{}
	default:
		return "", errors.New(constants.ErrInvalidActor, constants.MsgInvalidActor)
	}
}

// CreateAppointment initiates a negotiation for a match.
func CreateAppointment(request model.AppointmentRequest) (model.Appointment, errors.BusinessError) {
	match, ok := repository.FindMatch(request.MatchID)
	if !ok {
		return model.Appointment{}, errors.New(constants.ErrMatchNotFound, constants.MsgMatchNotFound)
	}
	if _, exists := repository.FindActiveAppointment(match.ID); exists {
		return model.Appointment{}, errors.New(constants.ErrAppointmentExists, constants.MsgAppointmentExists)
	}
	responder, actorErr := resolveParties(match, request.Actor)
	if actorErr.Code != "" {
		return model.Appointment{}, actorErr
	}
	if fieldErr := validateAppointmentFields(match, request); fieldErr.Code != "" {
		return model.Appointment{}, fieldErr
	}

	appointment := model.Appointment{
		MatchID:   match.ID,
		Pair:      appointmentPair(match),
		Initiator: strings.TrimSpace(request.Actor),
		Responder: responder,
		Time:      strings.TrimSpace(request.Time),
		Place:     strings.TrimSpace(request.Place),
		Agenda:    strings.TrimSpace(request.Agenda),
		Status:    constants.AppointmentPending,
		Slots:     []string{strings.TrimSpace(request.Time)},
		Version:   1,
	}
	return repository.CreateAppointment(appointment), errors.BusinessError{}
}

// ReviseAppointment lets the initiator change time/place/agenda before confirmation.
func ReviseAppointment(id int, request model.AppointmentRequest) (model.Appointment, errors.BusinessError) {
	appointment, ok := repository.GetAppointment(id)
	if !ok {
		return model.Appointment{}, errors.New(constants.ErrAppointmentGone, constants.MsgAppointmentGone)
	}
	match, ok := repository.FindMatch(appointment.MatchID)
	if !ok {
		return model.Appointment{}, errors.New(constants.ErrMatchNotFound, constants.MsgMatchNotFound)
	}
	if strings.TrimSpace(request.Actor) != appointment.Initiator {
		return model.Appointment{}, errors.New(constants.ErrInvalidActor, constants.MsgInvalidActor)
	}
	if appointment.Status != constants.AppointmentPending {
		return model.Appointment{}, errors.New(constants.ErrNotPending, constants.MsgNotPending)
	}
	if fieldErr := validateAppointmentFields(match, request); fieldErr.Code != "" {
		return model.Appointment{}, fieldErr
	}

	updated, ok := repository.ReviseAppointment(id, request.Time, request.Place, request.Agenda)
	if !ok {
		return model.Appointment{}, errors.New(constants.ErrAppointmentGone, constants.MsgAppointmentGone)
	}
	return updated, errors.BusinessError{}
}

// ConfirmAppointment is performed by the responding party; it is idempotent.
func ConfirmAppointment(id int, action model.AppointmentAction) (model.Appointment, errors.BusinessError) {
	appointment, ok := repository.GetAppointment(id)
	if !ok {
		return model.Appointment{}, errors.New(constants.ErrAppointmentGone, constants.MsgAppointmentGone)
	}
	if strings.TrimSpace(action.Actor) != appointment.Responder {
		return model.Appointment{}, errors.New(constants.ErrNotResponder, constants.MsgNotResponder)
	}
	// 已生效的预约重复确认直接返回，不产生新的确认记录。
	if appointment.Status == constants.AppointmentConfirmed {
		return appointment, errors.BusinessError{}
	}
	confirmed, ok := repository.ConfirmAppointment(id)
	if !ok {
		return model.Appointment{}, errors.New(constants.ErrAppointmentGone, constants.MsgAppointmentGone)
	}
	return confirmed, errors.BusinessError{}
}

// RejectAppointment is performed by the responding party and takes it down.
func RejectAppointment(id int, action model.AppointmentAction) errors.BusinessError {
	appointment, ok := repository.GetAppointment(id)
	if !ok {
		return errors.New(constants.ErrAppointmentGone, constants.MsgAppointmentGone)
	}
	if strings.TrimSpace(action.Actor) != appointment.Responder {
		return errors.New(constants.ErrNotResponder, constants.MsgNotResponder)
	}
	if appointment.Status != constants.AppointmentPending {
		return errors.New(constants.ErrNotPending, constants.MsgNotPending)
	}
	repository.DeleteAppointment(id)
	return errors.BusinessError{}
}
