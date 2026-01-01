package handler

import (
	"project-app-inventory/service"
	"project-app-inventory/utils"
)

type Handler struct {
	HandlerAuth       AuthHandler
	HandlerMenu       MenuHandler
	AssignmentHandler AssignmentHandler
	ItemHandler       ItemHandler
	CategoryHandler   CategoryHandler
	RackHandler       RackHandler
}

func NewHandler(service service.Service, config utils.Configuration) Handler {
	return Handler{
		HandlerAuth: NewAuthHandler(service),
		// HandlerMenu:       NewMenuHandler(),
		AssignmentHandler: NewAssignmentHandler(service.AssignmentService, config),
		ItemHandler:       NewItemHandler(service.ItemService, config),
		CategoryHandler:   NewCategoryHandler(service.CategoryService, config),
		RackHandler:       NewRackHandler(service.RackService, config),
	}
}
