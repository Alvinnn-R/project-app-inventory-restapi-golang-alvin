package handler

import (
	"project-app-inventory/service"
	"project-app-inventory/utils"
)

type Handler struct {
	HandlerAuth       AuthHandler
	HandlerMenu       MenuHandler
	AssignmentHandler AssignmentHandler
}

func NewHandler(service service.Service, config utils.Configuration) Handler {
	return Handler{
		HandlerAuth: NewAuthHandler(service),
		// HandlerMenu:       NewMenuHandler(),
		AssignmentHandler: NewAssignmentHandler(service.AssignmentService, config),
	}
}
