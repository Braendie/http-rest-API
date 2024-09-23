package sqlstore

import (
	"database/sql"

	"github.com/http-rest-API/internal/app/model"
	"github.com/http-rest-API/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	if u.IDTelegram.Valid {
		if u.PhoneNumber.Valid {
			return r.store.db.QueryRow(
				"INSERT INTO users (id_telegram, height, age, weight, gender, phone_number) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
				u.IDTelegram,
				u.Height,
				u.Age,
				u.Weight,
				u.Gender,
				u.PhoneNumber,
			).Scan(&u.ID)
		} else {
			return r.store.db.QueryRow(
				"INSERT INTO users (id_telegram, height, age, weight, gender) VALUES($1, $2, $3, $4, $5) RETURNING id",
				u.IDTelegram,
				u.Height,
				u.Age,
				u.Weight,
				u.Gender,
			).Scan(&u.ID)
		}
	} else {
		if u.PhoneNumber.Valid {
			return r.store.db.QueryRow(
				"INSERT INTO users (email, encrypted_password, height, age, weight, gender, phone_number) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
				u.Email,
				u.EncryptedPassword,
				u.Height,
				u.Age,
				u.Weight,
				u.Gender,
				u.PhoneNumber,
			).Scan(&u.ID)
		} else {
			return r.store.db.QueryRow(
				"INSERT INTO users (email, encrypted_password, height, age, weight, gender) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
				u.Email,
				u.EncryptedPassword,
				u.Height,
				u.Age,
				u.Weight,
				u.Gender,
			).Scan(&u.ID)
		}
	}
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, id_telegram, email, encrypted_password, height, age, weight, gender, phone_number FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.IDTelegram,
		&u.Email,
		&u.EncryptedPassword,
		&u.Height,
		&u.Age,
		&u.Weight,
		&u.Gender,
		&u.PhoneNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, id_telegram, email, encrypted_password, height, age, weight, gender, phone_number FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.IDTelegram,
		&u.Email,
		&u.EncryptedPassword,
		&u.Height,
		&u.Age,
		&u.Weight,
		&u.Gender,
		&u.PhoneNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByIDTelegram(idTelegram int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, id_telegram, email, encrypted_password, height, age, weight, gender, phone_number FROM users WHERE id_telegram = $1",
		idTelegram,
	).Scan(
		&u.ID,
		&u.IDTelegram,
		&u.Email,
		&u.EncryptedPassword,
		&u.Height,
		&u.Age,
		&u.Weight,
		&u.Gender,
		&u.PhoneNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
