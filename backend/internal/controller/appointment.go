package controller

import (
	"net/http"
	"strconv"

	apperrors "cyskillswap/internal/errors"
	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

type appointmentCreateRequest struct {
	MatchID int    `json:"matchId" binding:"required"`
	Actor   string `json:"actor" binding:"required"`
	Time    string `json:"time" binding:"required"`
	Place   string `json:"place"`
	Agenda  string `json:"agenda"`
}

type appointmentActorRequest struct {
	Actor string `json:"actor" binding:"required"`
}

type appointmentReviseRequest struct {
	Actor  string `json:"actor" binding:"required"`
	Time   string `json:"time" binding:"required"`
	Place  string `json:"place"`
	Agenda string `json:"agenda"`
}

func Appointments(c *gin.Context) { c.JSON(http.StatusOK, service.Appointments()) }

func CreateAppointment(c *gin.Context) {
	var req appointmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAppointmentError(c, apperrors.ErrAppointmentInvalidRequest)
		return
	}
	item, err := service.CreateAppointment(req.Actor, req.MatchID, req.Time, req.Place, req.Agenda)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func ConfirmAppointment(c *gin.Context) {
	id, ok := parseAppointmentID(c)
	if !ok {
		return
	}
	var req appointmentActorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAppointmentError(c, apperrors.ErrAppointmentInvalidRequest)
		return
	}
	item, err := service.ConfirmAppointment(id, req.Actor)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func RejectAppointment(c *gin.Context) {
	id, ok := parseAppointmentID(c)
	if !ok {
		return
	}
	var req appointmentActorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAppointmentError(c, apperrors.ErrAppointmentInvalidRequest)
		return
	}
	item, err := service.RejectAppointment(id, req.Actor)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func ReviseAppointment(c *gin.Context) {
	id, ok := parseAppointmentID(c)
	if !ok {
		return
	}
	var req appointmentReviseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAppointmentError(c, apperrors.ErrAppointmentInvalidRequest)
		return
	}
	item, err := service.ReviseAppointment(id, req.Actor, req.Time, req.Place, req.Agenda)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func parseAppointmentID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		respondAppointmentError(c, apperrors.ErrAppointmentNotFound)
		return 0, false
	}
	return id, true
}

func respondAppointmentError(c *gin.Context, err error) {
	if business, ok := err.(apperrors.BusinessError); ok {
		c.JSON(business.Status, business)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "服务器内部错误"})
}
