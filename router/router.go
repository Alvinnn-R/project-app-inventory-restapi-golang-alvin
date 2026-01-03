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

	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	return r
}

func Apiv1(handler handler.Handler, mw mCostume.MiddlewareCostume) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mw.Logging)

	// Public routes - no authentication required
	r.Post("/login", handler.HandlerAuth.Login)

	// Protected routes - authentication required
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)

		// Logout endpoint
		r.Post("/logout", handler.HandlerAuth.Logout)

		// Items routes - CRUD for inventory items
		r.Route("/items", func(r chi.Router) {
			r.Get("/", handler.ItemHandler.List)
			r.Post("/", handler.ItemHandler.Create)
			r.Route("/{item_id}", func(r chi.Router) {
				r.Get("/", handler.ItemHandler.GetByID)
				r.Put("/", handler.ItemHandler.Update)
				r.Delete("/", handler.ItemHandler.Delete)
			})
		})

		// Categories routes - CRUD for item categories
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", handler.CategoryHandler.List)
			r.Post("/", handler.CategoryHandler.Create)
			r.Route("/{category_id}", func(r chi.Router) {
				r.Get("/", handler.CategoryHandler.GetByID)
				r.Put("/", handler.CategoryHandler.Update)
				r.Delete("/", handler.CategoryHandler.Delete)
			})
		})

		// Racks routes - CRUD for storage racks
		r.Route("/racks", func(r chi.Router) {
			r.Get("/", handler.RackHandler.List) // supports ?warehouse_id=X filter
			r.Post("/", handler.RackHandler.Create)
			r.Route("/{rack_id}", func(r chi.Router) {
				r.Get("/", handler.RackHandler.GetByID)
				r.Put("/", handler.RackHandler.Update)
				r.Delete("/", handler.RackHandler.Delete)
			})
		})

		// Warehouses routes - CRUD for warehouses
		r.Route("/warehouses", func(r chi.Router) {
			r.Get("/", handler.WarehouseHandler.List)
			r.Post("/", handler.WarehouseHandler.Create)
			r.Route("/{warehouse_id}", func(r chi.Router) {
				r.Get("/", handler.WarehouseHandler.GetByID)
				r.Put("/", handler.WarehouseHandler.Update)
				r.Delete("/", handler.WarehouseHandler.Delete)
			})
		})

		// Sales routes - CRUD for sales transactions
		r.Route("/sales", func(r chi.Router) {
			r.Get("/", handler.SaleHandler.List)
			r.Post("/", handler.SaleHandler.Create)
			r.Route("/{sale_id}", func(r chi.Router) {
				r.Get("/", handler.SaleHandler.GetByID)
				r.Put("/", handler.SaleHandler.Update)
				r.Delete("/", handler.SaleHandler.Delete)
			})
		})

		// Users routes - CRUD for user management
		r.Route("/users", func(r chi.Router) {
			r.Get("/", handler.UserHandler.List)
			r.Post("/", handler.UserHandler.Create)
			r.Route("/{user_id}", func(r chi.Router) {
				r.Get("/", handler.UserHandler.GetByID)
				r.Put("/", handler.UserHandler.Update)
				r.Delete("/", handler.UserHandler.Delete)
			})
		})

		// Reports routes - Summary reports for inventory system
		r.Route("/reports", func(r chi.Router) {
			r.Get("/summary", handler.ReportHandler.GetSummary)
		})

		// Assignment routes (example - will be replaced with inventory routes later)
		r.Route("/assignment", func(r chi.Router) {
			r.Get("/", handler.AssignmentHandler.List)
			r.Post("/", handler.AssignmentHandler.Create)
			r.Route("/{assignment_id}", func(r chi.Router) {
				r.Get("/", handler.AssignmentHandler.GetByID)
				r.Put("/", handler.AssignmentHandler.Update)
				r.Delete("/", handler.AssignmentHandler.Delete)
			})
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
