package sqlstore

import (
	"database/sql"

	"github.com/http-rest-API/internal/app/store"
	_ "github.com/lib/pq"
)

// Store is a storage that includes the following fields:
// - db: the database that uses for storing information about users.
// - userRepository: the interface for calling function.
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New returns new store with specified database.
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User uses for calling UserRepository.
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
