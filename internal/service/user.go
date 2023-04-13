package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"refactoring/internal/models"
	"time"
)

type Store interface {
	CreateUser(user models.User) (string, error)
	DeleteUser(id string) error
	GetUserById(id string) (*models.User, error)
	ListUsers() []*models.User
	UpdateUser(user models.User) error
}

type User struct {
	store  Store
	logger *logrus.Logger
}

func NewUser(store Store, logger *logrus.Logger) *User {
	return &User{store: store, logger: logger}
}

func (s *User) CreateUser(name, email string) (*models.User, error) {
	s.logger.Infof("creating user: name=%s, email=%s", name, email)

	user := models.User{
		CreatedAt:   time.Now(),
		DisplayName: name,
		Email:       email,
	}

	userId, err := s.store.CreateUser(user)
	if err != nil {
		s.logger.Errorf("err create user: %s", err.Error())
		return nil, fmt.Errorf("err create user: %w", err)
	}

	user.Id = userId
	return &user, nil
}

func (s *User) DeleteUser(id string) error {
	s.logger.Infof("delete user: id=%s", id)

	if err := s.store.DeleteUser(id); err != nil {
		s.logger.Errorf("err delete user: %s", err.Error())
		return fmt.Errorf("err delete user: %w", err)
	}
	return nil
}

func (s *User) SearchUsers() []*models.User {
	return s.store.ListUsers()
}

func (s *User) GetUser(id string) (*models.User, error) {
	s.logger.Infof("get user: id=%s", id)

	user, err := s.store.GetUserById(id)
	if err != nil {
		s.logger.Errorf("err get user: %s", err.Error())
		return nil, fmt.Errorf("err get user: %w", err)
	}

	return user, nil
}

func (s *User) UpdateUser(id, newName, newEmail string) error {
	s.logger.Infof("updating user: id=%s, newName=%s, newEmail=%s", id, newName, newEmail)

	err := s.store.UpdateUser(models.User{
		Id:          id,
		DisplayName: newName,
		Email:       newEmail,
	})
	if err != nil {
		s.logger.Errorf("err update user: %s", err.Error())
		return err
	}
	return nil
}
