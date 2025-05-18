package server

import (
	"net/http"
	"todo/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
	customMiddleware "todo/internal/middleware"
	"todo/internal/repository"
)

// Server represents the HTTP server
type Server struct {
	router chi.Router
	db     *gorm.DB
}

// NewServer creates a new HTTP server
func NewServer(db *gorm.DB) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     db,
	}
	s.setupRoutes()
	return s
}

// setupRoutes initializes the routes
func (s *Server) setupRoutes() {
	// Middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(customMiddleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(customMiddleware.CORS)

	// Create repositories
	todoRepo := repository.NewTodoRepository(s.db)

	// Create handlers
	todoHandler := handlers.NewTodoHandler(todoRepo)

	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Todo routes
		r.Route("/todos", func(r chi.Router) {
			r.Get("/", todoHandler.GetAll)
			r.Post("/", todoHandler.Create)
			r.Get("/{id}", todoHandler.GetByID)
			r.Put("/{id}", todoHandler.Update)
			r.Delete("/{id}", todoHandler.Delete)
		})
	})

	// Swagger documentation
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
	))

	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
