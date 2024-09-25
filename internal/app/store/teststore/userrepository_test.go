package teststore_test

import (
	"database/sql"
	"testing"

	"github.com/http-rest-API/internal/app/model"
	"github.com/http-rest-API/internal/app/store"
	"github.com/http-rest-API/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_Find(t *testing.T) {
	s := teststore.New()
	u1 := model.TestUser(t)
	s.User().Create(u1)
	u2, err := s.User().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	email := "user@example.org"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email.String = email
	u.Email.Valid = true
	s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByIDTelegram(t *testing.T) {
	s := teststore.New()
	idTelegram := 12345678
	_, err := s.User().FindByIDTelegram(idTelegram)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUserWithTelegram(t)
	u.IDTelegram = sql.NullInt64{Int64: int64(idTelegram), Valid: true}
	s.User().Create(u)

	u, err = s.User().FindByIDTelegram(idTelegram)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
