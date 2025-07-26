package sessions

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/willtrojniak/TabAppBackend/cache"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services"
	"github.com/willtrojniak/TabAppBackend/util"
	"slices"
)

const (
	session_cookie = "session"
	csrf_header    = "X-CSRF-TOKEN"
	csrf_field     = "xcsrftoken"
)

var (
	safe_methods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}
)

type sessionData struct {
	UserId    string
	CSRFToken string
	Ip        string
}

type Session struct {
	data sessionData
	Id   string
	ttl  time.Duration
}

type AuthedSession struct {
	Id     string
	UserId string
}

type Handler struct {
	logger      *slog.Logger
	store       cache.Cache
	authTTL     time.Duration
	unauthTTL   time.Duration
	handleError services.HTTPErrorHandler
}

func New(store cache.Cache, authTTL time.Duration, unauthTTL time.Duration, h services.HTTPErrorHandler, logger *slog.Logger) *Handler {

	return &Handler{
		logger:      logger,
		store:       store,
		authTTL:     authTTL,
		unauthTTL:   unauthTTL,
		handleError: h,
	}
}

func (s *Handler) SetNewSession(w http.ResponseWriter, r *http.Request, user *models.User) (*Session, error) {
	ip := readUserIP(r)

	session, err := s.newSessionFromUser(ip, user)
	if err != nil {
		s.logger.Debug("Error while creating new session.")
		return nil, err
	}

	err = s.saveSessionToStore(r.Context(), session)
	if err != nil {
		return nil, err
	}

	oldSessionId, err := s.getSessionIdFromRequest(r)
	if err == nil { // i.e The client has a previous saved session
		err = s.store.Delete(r.Context(), oldSessionId)
		if err != nil {
			s.logger.Warn("Failed to delete old session from store.", "sessionId", oldSessionId, "error", err)
		}
	}

	s.saveSessionToResponse(w, session)

	s.logger.Debug("Session created", "sessionId", session.Id)

	return session, nil
}

func (s *Handler) GetSession(r *http.Request) (*Session, error) {
	sessionId, err := s.getSessionIdFromRequest(r)
	if err != nil {
		return nil, services.NewUnauthenticatedServiceError(err)
	}

	session, err := s.getSessionFromStore(r.Context(), sessionId)
	if err != nil {
		return nil, err
	}

	if session.data.Ip != readUserIP(r) {
		s.logger.Warn("Attempted to access session with different ip", "stored-ip", session.data.Ip, "request-ip", readUserIP(r))
		return nil, services.NewUnauthenticatedServiceError(err)
	}

	return session, nil
}

func (s *Handler) WithAuthedSession(next func(w http.ResponseWriter, r *http.Request, session *AuthedSession)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.GetSession(r)
		if err != nil {
			s.handleError(w, err)
			return
		}

		authed, err := session.Authed()
		if err != nil {
			s.handleError(w, err)
			return
		}

		next(w, r, authed)
	}
}

func (s *Handler) RequireAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.GetSession(r)
		if err != nil {
			s.handleError(w, err)
			return
		}
		if _, err := session.Authed(); err != nil {
			s.handleError(w, err)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (s *Handler) RequireCSRFToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestToken := getCSRFTokenFromRequest(r)
		safeMethod := slices.Contains(safe_methods, r.Method)

		session, sessionErr := s.GetSession(r)
		validSession := sessionErr == nil

		// Make sure header has token on future requests
		if validSession {
			s.saveSessionToResponse(w, session)
		} else {
			// Create a new anonymous session if there is no active session
			s.SetNewSession(w, r, nil)
		}

		if !safeMethod {
			if !validSession {
				s.handleError(w, services.NewServiceError(errors.New("No CSRF token to match"), http.StatusForbidden, nil))
				return
			}

			if requestToken != session.data.CSRFToken {
				s.logger.Warn("CSRF Tokens did not match", "incoming-token", requestToken, "stored-token", session.data.CSRFToken)
				s.handleError(w, services.NewServiceError(errors.New("CSRF Tokens did not match"), http.StatusForbidden, nil))
				return
			}

		}
		next.ServeHTTP(w, r)

	}
}

func (s *Handler) newSessionFromUser(ip string, user *models.User) (*Session, error) {
	id, err := util.RandString(32)
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	csrf, err := util.RandString(32)
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	ttl := s.unauthTTL
	if user != nil {
		ttl = s.authTTL
	}

	userId := ""
	if user != nil {
		userId = user.Id
	}

	return &Session{
		Id:  id,
		ttl: ttl,
		data: sessionData{
			UserId:    userId,
			CSRFToken: csrf,
			Ip:        ip,
		},
	}, nil
}

func (s *Handler) saveSessionToResponse(w http.ResponseWriter, session *Session) {
	c := &http.Cookie{
		Name:     session_cookie,
		Value:    session.Id,
		MaxAge:   int(session.ttl.Seconds()),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: 4,
	}
	http.SetCookie(w, c)
	w.Header().Set(csrf_header, session.data.CSRFToken)
}

func (s *Handler) saveSessionToStore(ctx context.Context, session *Session) error {
	data, err := json.Marshal(session.data)
	if err != nil {
		s.logger.Warn("Failed to marshal session data.", "sessionId", session.Id)
		return services.NewInternalServiceError(err)
	}

	err = s.store.Set(ctx, session.Id, data, session.ttl)
	if err != nil {
		s.logger.Warn("Failed to save session to store.", "sessionId", session.Id)
		return services.NewInternalServiceError(err)
	}

	return nil
}

func (s *Handler) getSessionIdFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(session_cookie)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (s *Handler) getSessionFromStore(ctx context.Context, id string) (*Session, error) {
	data, err := s.store.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, cache.ErrNotFound):
			return nil, services.NewUnauthenticatedServiceError(err)
		default:
			return nil, services.NewInternalServiceError(err)
		}
	}

	session := &Session{
		Id:  id,
		ttl: s.unauthTTL,
	}
	err = json.Unmarshal(data, &session.data)
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	if session.data.UserId != "" {
		session.ttl = s.authTTL
	}

	return session, nil
}

func (s *Session) IsAuthed() bool {
	_, err := s.Authed()
	return err == nil
}
func (s *Session) Authed() (*AuthedSession, error) {
	if s.data.UserId == "" {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}
	return &AuthedSession{Id: s.Id, UserId: s.data.UserId}, nil
}

func readUserIP(r *http.Request) string {
	addr := r.RemoteAddr
	ip := strings.Split(addr, ":")[0]
	return ip
}

func getCSRFTokenFromRequest(r *http.Request) string {
	token := r.Header.Get(csrf_header)
	if token != "" {
		return token
	}

	token = r.PostFormValue(csrf_field)
	if token != "" {
		return token
	}

	return ""
}
