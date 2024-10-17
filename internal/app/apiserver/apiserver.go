package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/http-rest-API/internal/app/store/sqlstore"
)

// Start creates a new server with new store and sessionStore.
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, s)
}

// newDB creates a new database for store and starts it.
func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
