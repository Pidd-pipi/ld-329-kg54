package service

import (
	"strings"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/model"
)

// validateAppointmentFields checks time/place/agenda against the match.
func validateAppointmentFields(match model.Match, request model.AppointmentRequest) errors.BusinessError {
	if !contains(match.CommonSlots, strings.TrimSpace(request.Time)) {
		return errors.New(constants.ErrInvalidField, constants.MsgInvalidSlot)
	}
	if strings.TrimSpace(request.Place) == "" {
		return errors.New(constants.ErrInvalidField, constants.MsgInvalidPlace)
	}
	if strings.TrimSpace(request.Agenda) == "" {
		return errors.New(constants.ErrInvalidField, constants.MsgInvalidAgenda)
	}
	return errors.BusinessError{}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
