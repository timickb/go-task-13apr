package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
	"refactoring/internal/config"
	"refactoring/internal/models"
	"time"
)

type UserService interface {
	CreateUser(name, email string) (*models.User, error)
	DeleteUser(id string) error
	SearchUsers() []*models.User
	GetUser(id string) (*models.User, error)

	// UpdateUser If newName or newEmail equaled to empty string,
	// it wouldn't be updated.
	UpdateUser(id, newName, newEmail string) error
}

type Server struct {
	router *chi.Mux
	logger *logrus.Logger
	cfg    *config.Config
	user   UserService
}

func New(cfg *config.Config, log *logrus.Logger, user UserService) *Server {
	srv := &Server{
		logger: log,
		cfg:    cfg,
		user:   user,
	}
	srv.configureRouter()
	return srv
}

func (s *Server) Run() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.AppPort), s.router)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) configureRouter() {
	s.router = chi.NewRouter()

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)

	s.router.Use(middleware.Timeout(60 * time.Second))

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.Data(w, r, []byte(time.Now().String()))
	})

	s.router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.searchUsers)
				r.Post("/", s.createUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", s.getUser)
					r.Patch("/", s.updateUser)
					r.Delete("/", s.deleteUser)
				})
			})
		})
	})
}
