package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
  logger *slog.Logger
  store *redis.Client
  expirationTime time.Duration
}

type SessionData struct {
  UserId string
  CSRFToken string
  Ip string

}

const (
  session_cookie = "session"
  csrf_header = "X-CSRF-TOKEN"
)

var (
  safe_methods = []string{"GET", "HEAD", "OPTIONS", "TRACE"};
)


func New(store *redis.Client, expiryTime time.Duration, logger *slog.Logger) *SessionManager {

  return &SessionManager{
    logger: logger,
    store: store,
    expirationTime: expiryTime,
  };
}


func (s *SessionManager) CreateSession(w http.ResponseWriter, r *http.Request, userId string) (*SessionData, error ){
  sessionID, err := randString(32);
  if err != nil {
    return nil, services.NewInternalServiceError(err);
  }
  csrfToken, err := randString(32);
  if err != nil {
    return nil, services.NewInternalServiceError(err);
  }
  s.logger.Debug("Creating session", "sessionId", sessionID);

  sessionData := SessionData{UserId: userId, Ip: readUserIP(r), CSRFToken: csrfToken};
  jsonString, err := json.Marshal(sessionData);
  if err != nil {
    return nil, services.NewInternalServiceError(err);
  }

  currentSession, err := r.Cookie(session_cookie);
  if err == nil {
    // i.e. The client has a previous session
    err := s.store.Del(r.Context(), currentSession.Value).Err();
    if err != nil {
      s.logger.Warn("Attempt to delete old session failed", "err", err);
    }
  }


  if err := s.store.Set(r.Context(), sessionID, jsonString, s.expirationTime).Err(); err != nil {
    s.logger.Error("Session Manager could not save session to redis");
    return nil, services.NewInternalServiceError(err);
  }

  s.createSessionCookie(w, r, sessionID)
  s.setCSRFHeader(w, &sessionData);
  s.logger.Debug("Session created", "sessionId", sessionID);

  return &sessionData, nil;
}

func (s *SessionManager) GetSession(r *http.Request) (*SessionData, error) {
  // TODO: Maybe just use the request's context?
  sessionCookie, err := r.Cookie(session_cookie);
  if err != nil {
    return nil, services.NewUnauthorizedServiceError(err);
  }

  sessionId := sessionCookie.Value;
  jsonString, err := s.store.Get(r.Context(), sessionId).Bytes();
  if err != nil {
    return nil, services.NewUnauthorizedServiceError(err);
  }

  sessionData := SessionData{};
  err = json.Unmarshal(jsonString, &sessionData);
  if err != nil {
    s.logger.Warn("Failed to parse json data from redis");
    return nil, services.NewUnauthorizedServiceError(err);
  }


  if sessionData.Ip != readUserIP(r) {
    s.logger.Debug("Attempted to access session with different ip", "stored-ip", sessionData.Ip, "request-ip", readUserIP(r));
    return nil, services.NewInternalServiceError(err);
  }

  return &sessionData, nil;
}

func (s *SessionManager) ClearSession(w http.ResponseWriter, r *http.Request) error {
  // TODO: Maybe just use the request's context?
  _, err := s.CreateSession(w, r, ""); // Create an anonymous session
  if err != nil {
    return err;
  }
  return nil;
}

func (s *SessionManager) RequireCSRFHeader(h services.HTTPErrorHandler) func(http.Handler) http.HandlerFunc {
  return func(next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      
      // Check for an active session
      session, err := s.GetSession(r);
      if err != nil {

        // If there is no session, create an unauthenticated session
        session, err = s.CreateSession(w, r, ""); // Empty string for no user id --> unauthenticated
        if err != nil {
          s.logger.Error("Failed to create anonymous session", "error", err);
          h(w, services.NewInternalServiceError(err));
          return;
        }
      }
      // Set the CSRF header in the response
      s.setCSRFHeader(w, session);

      requestToken := r.Header.Get(csrf_header);
      safeMethod := false;
      for _, val := range safe_methods {
        if val == r.Method {
          safeMethod = true;
          break;
        }
      }

      if !safeMethod && requestToken != session.CSRFToken {
        s.logger.Warn("CSRF Tokens did not match", "incoming-token", requestToken, "stored-token", session.CSRFToken);
        h(w, services.NewServiceError(errors.New("CSRF Tokens did not match"), http.StatusForbidden, "CSRF Tokens did not match", nil))
        return;
      }

      next.ServeHTTP(w, r);
    }
  }
}

func (s *SessionManager) setCSRFHeader(w http.ResponseWriter, data *SessionData) {
  w.Header().Set(csrf_header, data.CSRFToken);
}

func (s *SessionManager) createSessionCookie(w http.ResponseWriter, r *http.Request, sessionId string) {
  c := &http.Cookie{
    Name: session_cookie,
    Value: sessionId,
    MaxAge: int(s.expirationTime.Seconds()),
    Secure: r.TLS != nil,
    HttpOnly: true,
    Path: "/",
  }
  http.SetCookie(w, c);
}

func readUserIP(r *http.Request) string {
  addr := r.RemoteAddr;
  ip := strings.Split(addr, ":")[0];
  return ip;

}

func randString(nByte int) (string, error) {
  b := make([]byte, nByte);
  if _, err := io.ReadFull(rand.Reader, b); err != nil {
    return "", err;
  }
  return base64.RawURLEncoding.EncodeToString(b), nil;
}
