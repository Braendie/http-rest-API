package apiserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/http-rest-API/internal/app/model"
	"github.com/http-rest-API/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "braendie"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
	domainURL = "http://localhost:8080"
)

var (
	errIncorrectEmailOrPassword  = errors.New("incorrect email or password")
	errNotAuthenticated          = errors.New("not authenticated")
	errConfirmPasswordIsRequired = errors.New("confirm password is required")
	errEasyPassword              = errors.New("password is easy to hack")
)

type ctxKey int8

// server represents a server application that handles HTTP requests
// It contains the following fields:
// - router: a request router using the gorilla/mux library for routing HTTP requests.
// - logger: a logger for recording server logs, using the logrus library.
// - store: an interface for working with the data store, providing access to data models.
// - sessionStore: an interface for working with session storage, enabling the management of user sessions.
type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

// newServer initializes a new server instance with the given store and session store,
// sets up routing and logging middleware, and returns the server instance.
func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

// ServeHTTP handles HTTP requests by passing them to the router for further handling.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter sets up the routing for the server by associating routes with
// their corresponding handler functions and applying middlewares.
func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	// Define public routes.
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	// Define routes under /enter prefix.
	enter := s.router.PathPrefix("/enter").Subrouter()
	enter.HandleFunc("/register", s.handleRegister()).Methods("GET")
	enter.HandleFunc("/login", s.handleLogin()).Methods("GET")
	enter.HandleFunc("/images", s.handleImage()).Methods("GET")

	// Define routes for Telegram-related actions.
	telegram := s.router.PathPrefix("/telegram").Subrouter()
	telegram.HandleFunc("/check", s.handleTelegramCheck()).Methods("POST")

	// Define private routes that require authorization.
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
	private.HandleFunc("/main", s.handleMain()).Methods("GET")
}

// setRequestID adds a unique request ID to each incoming request for tracking purposes.
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

// logRequest logs details about each incoming request, including the remote address,
// request method, URI, and the time taken to process the request.
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
		)
	})
}

// authenticateUser checks if the user is authenticated by verifying the session.
// If the session is valid, the user information is added to the request context.
func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

// handleMain serves the main page (HTML).
func (s *server) handleMain() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.sendHtmlFile(w, r, "main")
	}
}

// handleRegister serves the registration page (HTML).
func (s *server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.sendHtmlFile(w, r, "register")
	}
}

// handleLogin serves the login page (HTML).
func (s *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.sendHtmlFile(w, r, "login")
	}
}

// handleWhoami responds with information about the currently authenticated user.
func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

// handleImageRegister serves a specific image used on the registration page.
func (s *server) handleImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()

		imageName := queryParams.Get("image_name")

		filePath := "D:/GitHubProjects/http-rest-API/internal/app/htmlfiles/images/image_" + imageName + ".webp"

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		w.Header().Set("Content-Type", "image/webp")
		w.Header().Set("Content-Disposition", "inline; filename=background_image_register.webp")

		http.ServeFile(w, r, filePath)
	}
}

// handleTelegramCheck checks if a user exists based on their Telegram ID,
// and either logs them in or creates a new user.
func (s *server) handleTelegramCheck() http.HandlerFunc {
	type request struct {
		IDTelegram int `json:"id_telegram"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByIDTelegram(req.IDTelegram)
		if err != nil {
			if err == store.ErrRecordNotFound {
				u := &model.User{
					IDTelegram: sql.NullInt64{Int64: int64(req.IDTelegram), Valid: true},
				}
				if err := s.store.User().Create(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				s.createSessions(w, r, u)
				http.Redirect(w, r, domainURL+"/private/main", http.StatusFound)
				return
			}

			s.error(w, r, http.StatusInternalServerError, err)

		}

		s.createSessions(w, r, u)
		s.respond(w, r, http.StatusFound, u)
		http.Redirect(w, r, domainURL+"/private/main", http.StatusFound)
	}
}

// handleUsersCreate creates a new user based on the provided email, password.
func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.ConfirmPassword != req.Password {
			s.error(w, r, http.StatusBadRequest, errConfirmPasswordIsRequired)
			return
		}

		if !model.CheckPassword(req.Password) {
			s.error(w, r, http.StatusBadRequest, errEasyPassword)
			return
		}

		u := &model.User{
			IDTelegram: sql.NullInt64{Valid: false},
			Email:      sql.NullString{String: req.Email, Valid: req.Email != ""},
			Password:   req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

// handleSessionsCreate creates a new session for a user based on email and password.
func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		s.createSessions(w, r, u)
		s.respond(w, r, http.StatusOK, nil)
	}
}

// error calls respond function with error.
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond encodes the data into json format.
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// sendHtmlFile gives HTML file which is stored in specified folder.
func (s *server) sendHtmlFile(w http.ResponseWriter, r *http.Request, htmlName string) {
	filePath := "D:/GitHubProjects/http-rest-API/internal/app/htmlfiles/" + htmlName + ".html"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	http.ServeFile(w, r, filePath)
}

// createSessions creates a session for the authenticated user.
func (s *server) createSessions(w http.ResponseWriter, r *http.Request, u *model.User) {
	session := sessions.NewSession(s.sessionStore, sessionName)

	session.Values["user_id"] = u.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, nil)
}
