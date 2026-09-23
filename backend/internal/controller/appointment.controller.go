package controller

import (
	"net/http"
	"strconv"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/logger"
	"cyskillswap/internal/model"
	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

func writeBusinessError(c *gin.Context, err errors.BusinessError) {
	logger.Warn("appointment request rejected", err.Code, err.Message)
	c.JSON(http.StatusUnprocessableEntity, err)
}

func bindAppointmentRequest(c *gin.Context) (model.AppointmentRequest, bool) {
	var request model.AppointmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.New(constants.ErrInvalidJSON, constants.MsgInvalidJSON))
		return request, false
	}
	return request, true
}

func bindAppointmentAction(c *gin.Context) (model.AppointmentAction, bool) {
	var action model.AppointmentAction
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, errors.New(constants.ErrInvalidJSON, constants.MsgInvalidJSON))
		return action, false
	}
	return action, true
}

func appointmentID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errors.New(constants.ErrAppointmentGone, constants.MsgAppointmentGone))
		return 0, false
	}
	return id, true
}

// CreateAppointment handles POST /api/appointments.
func CreateAppointment(c *gin.Context) {
	request, ok := bindAppointmentRequest(c)
	if !ok {
		return
	}
	appointment, err := service.CreateAppointment(request)
	if err.Code != "" {
		writeBusinessError(c, err)
		return
	}
	c.JSON(http.StatusCreated, appointment)
}

// ReviseAppointment handles PUT /api/appointments/:id.
func ReviseAppointment(c *gin.Context) {
	id, ok := appointmentID(c)
	if !ok {
		return
	}
	request, ok := bindAppointmentRequest(c)
	if !ok {
		return
	}
	appointment, err := service.ReviseAppointment(id, request)
	if err.Code != "" {
		writeBusinessError(c, err)
		return
	}
	c.JSON(http.StatusOK, appointment)
}

// ConfirmAppointment handles POST /api/appointments/:id/confirm.
func ConfirmAppointment(c *gin.Context) {
	id, ok := appointmentID(c)
	if !ok {
		return
	}
	action, ok := bindAppointmentAction(c)
	if !ok {
		return
	}
	appointment, err := service.ConfirmAppointment(id, action)
	if err.Code != "" {
		writeBusinessError(c, err)
		return
	}
	c.JSON(http.StatusOK, appointment)
}

// RejectAppointment handles POST /api/appointments/:id/reject.
func RejectAppointment(c *gin.Context) {
	id, ok := appointmentID(c)
	if !ok {
		return
	}
	action, ok := bindAppointmentAction(c)
	if !ok {
		return
	}
	if err := service.RejectAppointment(id, action); err.Code != "" {
		writeBusinessError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "removed": true})
}
