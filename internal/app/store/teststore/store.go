package teststore

import (
	"github.com/http-rest-API/internal/app/model"
	"github.com/http-rest-API/internal/app/store"
)

// Store is a test storage that includes the following fields:
// - userRepository: the interface for calling function.
type Store struct {
	userRepository *UserRepository
}

// New returns a new Store.
func New() *Store {
	return &Store{}
}

// User uses for calling UserRepository.
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.userRepository
}
