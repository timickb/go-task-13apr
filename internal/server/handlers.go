package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	request := CreateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := s.user.CreateUser(request.DisplayName, request.Email)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": user.Id,
	})
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := s.user.GetUser(id)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, user)
}

func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	request := UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id := chi.URLParam(r, "id")
	if err := s.user.UpdateUser(id, request.DisplayName, ""); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := s.user.DeleteUser(id); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (s *Server) searchUsers(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, s.user.SearchUsers())
}
