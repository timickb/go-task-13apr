package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"refactoring/internal/models"
	"strconv"
)

var (
	ErrEntityNotFound   = errors.New("err entity not found")
	ErrInvalidIncrement = errors.New("err invalid increment")
)

type fileData struct {
	List      map[string]*models.User `json:"list,omitempty"`
	Increment int                     `json:"increment,omitempty"`
}

type Store struct {
	data     fileData
	fileName string
}

func New(fileName string) (*Store, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("err create file store: %w", err)
	}

	st := &Store{
		fileName: fileName,
	}

	if err := json.Unmarshal(file, &st.data); err != nil {
		return nil, fmt.Errorf("err create file store: %w", err)
	}

	return st, nil
}

func (s *Store) CreateUser(user models.User) (string, error) {
	userId := strconv.Itoa(s.data.Increment + 1)

	if _, ok := s.data.List[userId]; ok {
		return "", ErrInvalidIncrement
	}

	user.Id = userId

	s.data.List[userId] = &user
	s.data.Increment++

	if err := s.commit(); err != nil {
		delete(s.data.List, userId) // Rollback changes.
		s.data.Increment--
		return "", fmt.Errorf("err write store file: %w", err)
	}
	return userId, nil
}

func (s *Store) DeleteUser(id string) error {
	user, ok := s.data.List[id]
	if !ok {
		return errors.New("err user not found")
	}

	saved := &models.User{
		Id:          user.Id,
		CreatedAt:   user.CreatedAt,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}

	delete(s.data.List, id)

	if err := s.commit(); err != nil {
		s.data.List[id] = saved // Rollback changes.
		return fmt.Errorf("err write store file: %w", err)
	}
	return nil
}

func (s *Store) GetUserById(id string) (*models.User, error) {
	user, ok := s.data.List[id]
	if !ok {
		return nil, ErrEntityNotFound
	}

	return user, nil
}

func (s *Store) ListUsers() []*models.User {
	length := len(s.data.List)
	users := make([]*models.User, length)

	i := 0
	for id, user := range s.data.List {
		user.Id = id
		users[i] = user
		i++
		if i >= length {
			break
		}
	}

	return users
}

func (s *Store) UpdateUser(updated models.User) error {
	user, ok := s.data.List[updated.Id]
	if !ok {
		return ErrEntityNotFound
	}

	saved := &models.User{
		Id:          user.Id,
		CreatedAt:   user.CreatedAt,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}

	if updated.Email != "" {
		s.data.List[updated.Id].Email = updated.Email
	}

	if updated.DisplayName != "" {
		s.data.List[updated.Id].DisplayName = updated.DisplayName
	}

	if err := s.commit(); err != nil {
		s.data.List[updated.Id] = saved // Rollback changes.
		return fmt.Errorf("err write store file: %w", err)
	}

	return nil
}

func (s *Store) commit() error {
	b, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.fileName, b, fs.ModePerm); err != nil {
		return err
	}
	return nil
}
