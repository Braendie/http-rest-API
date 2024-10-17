package teststore

import (
	"github.com/http-rest-API/internal/app/model"
	"github.com/http-rest-API/internal/app/store"
)

// UserRepository uses for manipulating with test store.
// It including:
// - store: it is test store.
// - users: it is map that uses how database for testing.
type UserRepository struct {
	store *Store
	users map[int]*model.User
}

// Create adds a new user into map (it validates before adding).
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u
	u.ID = len(r.users)

	return nil
}

// Find finds the user in map by using his id.
func (r *UserRepository) Find(id int) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

// Find finds the user in map by using his email.
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email.String == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// Find finds the user in database by using his telegram id.
func (r *UserRepository) FindByIDTelegram(idTelegram int) (*model.User, error) {
	for _, u := range r.users {
		if u.IDTelegram.Int64 == int64(idTelegram) {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
