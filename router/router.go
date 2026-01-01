package router

import (
	"net/http"
	"project-app-inventory/handler"
	mCostume "project-app-inventory/middleware"
	"project-app-inventory/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(handler handler.Handler, service service.Service, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	// middleware
	mw := mCostume.NewMiddlewareCustome(service, log)

	r.Mount("/api/v1", Apiv1(handler, mw))
	r.Mount("/api/v2", Apiv2(handler))

	// //menu
	// r.Route("/user", func(r chi.Router) {
	// 	r.Use(middleware.AuthMiddleware)
	// 	r.Get("/assignments", handler.AssignmentHandler.List)
	// 	r.Get("/success-submit", handler.AssignmentHandler.SuccessSubmit)
	// 	r.Post("/submit-assignment", handler.AssignmentHandler.SubmitAssignment)
	// 	r.Get("/grade", handler.HandlerMenu.GradeView)
	// 	r.Get("/logout", handler.HandlerAuth.LogoutView)
	// })
	// r.Get("/page401", handler.HandlerMenu.PageUnauthorized)

	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	return r
}

func Apiv1(handler handler.Handler, mw mCostume.MiddlewareCostume) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mw.Logging)
	//authentication
	r.Post("/login", handler.HandlerAuth.Login)
	// r.Post("/logout", handler.HandlerAuth.Logout)

	r.Route("/assignment", func(r chi.Router) {
		r.Get("/", handler.AssignmentHandler.List)
		r.Post("/", handler.AssignmentHandler.Create)
		r.Route("/{assignment_id}", func(r chi.Router) {
			r.Get("/", handler.AssignmentHandler.GetByID)
			r.Put("/", handler.AssignmentHandler.Update)
			r.Delete("/", handler.AssignmentHandler.Delete)
		})
	})

	return r
}

func Apiv2(handler handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/assignment", func(r chi.Router) {
		r.Post("/", handler.AssignmentHandler.Create)
	})

	return r
}
