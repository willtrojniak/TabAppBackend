package sessions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
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
  Ip string

}

const (session_cookie = "session")

func New(store *redis.Client, expiryTime time.Duration, logger *slog.Logger) *SessionManager {

  return &SessionManager{
    logger: logger,
    store: store,
    expirationTime: expiryTime,
  };
}


func (s *SessionManager) CreateSession(c context.Context, w http.ResponseWriter, r *http.Request, userId string) error {
  sessionID, err := randString(32);
  if err != nil {
    return services.NewInternalServiceError(err);
  }
  s.logger.Debug("Creating session", "sessionId", sessionID);

  sessionData := SessionData{UserId: userId, Ip: readUserIP(r)};
  jsonString, err := json.Marshal(sessionData);
  if err != nil {
    return services.NewInternalServiceError(err);
  }

  if err := s.store.Set(c, sessionID, jsonString, s.expirationTime).Err(); err != nil {
    s.logger.Error("Session Manager could not save session to redis");
    return services.NewInternalServiceError(err);
  }

  s.createSessionCookie(w, r, sessionID)
  s.logger.Debug("Session created", "sessionId", sessionID);

  return nil;
}

func (s *SessionManager) GetSession(c context.Context, r *http.Request) (*SessionData, error) {
  // TODO: Maybe just use the request's context?
  sessionCookie, err := r.Cookie(session_cookie);
  if err != nil {
    return nil, services.NewUnauthorizedServiceError(err);
  }

  sessionId := sessionCookie.Value;
  jsonString, err := s.store.Get(c, sessionId).Bytes();
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

func (s *SessionManager) ClearSession(c context.Context, w http.ResponseWriter, r *http.Request) error {
  // TODO: Maybe just use the request's context?
  sessionCookie, err := r.Cookie(session_cookie);
  if err != nil {
    // There is no session to clear
    return nil;
  }
  sessionId := sessionCookie.Value;
  err = s.store.Del(c, sessionId).Err();
  if err != nil {
    s.logger.Warn("Session manager could not delete session.", "sessionid", sessionId, "err", err);
    return services.NewInternalServiceError(err);
  }

  cookie := &http.Cookie{
    Name: session_cookie,
    Value: "",
    MaxAge: -1,
  };
  http.SetCookie(w, cookie);

  return nil;

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
  ipAddr := r.Header.Get("X-Real-Ip")
  if ipAddr == "" {
    ipAddr = r.Header.Get("X-Forwarded-For")
  }
  if ipAddr == "" {
    ipAddr = r.RemoteAddr;
  }
  return ipAddr;

}

func randString(nByte int) (string, error) {
  b := make([]byte, nByte);
  if _, err := io.ReadFull(rand.Reader, b); err != nil {
    return "", err;
  }
  return base64.RawURLEncoding.EncodeToString(b), nil;
}
