package sessionStore

import (
	"encoding/gob"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type SessionData struct {
	UserID        string `json:"userAuthID"`
	UserProfileID string
	PublicKey     string
	RoleID        string
	Role          string
	Username      string
	Email         string
	Name          string
	Surname       string
	Level         int
	Experience    int
	NextLevelExp  int
}

func GetSessionData(c *fiber.Ctx) *SessionData {
	user := c.Locals("user")
	if user == nil {
		return nil
	}
	sessionData, ok := user.(SessionData)
	if !ok {
		return nil
	}
	return &sessionData
}

func NewSessionStore(storage ...fiber.Storage) *session.Store {
	if len(storage) <= 0 {
		storage = append(storage, session.ConfigDefault.Storage)
	}
	gob.Register(SessionData{})
	return session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: false,
		Storage:        storage[0],
	})
}

// UnmarshalSessionData: Redis'ten gelen JSON verisini çözümler
func (s *SessionData) Unmarshal(sessionData string) error {
	if s == nil {
		return errors.New("destination struct cannot be nil")
	}
	err := json.Unmarshal([]byte(sessionData), s)
	if err != nil {
		return err
	}
	return nil
}
