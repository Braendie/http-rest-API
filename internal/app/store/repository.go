package store

import "github.com/http-rest-API/internal/app/model"

//UserRepository is a interface that allows you to use functions for working with database or map.
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	FindByIDTelegram(int) (*model.User, error)
}
