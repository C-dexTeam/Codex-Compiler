package private

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	serviceErrors "github.com/C-dexTeam/codex-compiler/internal/errors"
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/internal/http/response"
	"github.com/C-dexTeam/codex-compiler/internal/http/sessionStore"
	"github.com/C-dexTeam/codex-compiler/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// PrivateHandler struct
type PrivateHandler struct {
	sess_store *session.Store
	services   *services.Services
	dtoManager dto.IDTOManager
}

// NewPrivateHandler creates a new instance of PrivateHandler
func NewPrivateHandler(
	sessStore *session.Store,
	dtoManager dto.IDTOManager,
	services *services.Services,
) *PrivateHandler {
	return &PrivateHandler{
		sess_store: sessStore,
		dtoManager: dtoManager,
		services:   services,
	}
}

// Init initializes the routes
func (h *PrivateHandler) Init(router fiber.Router) {
	root := router.Group("/private")
	root.Use(h.authMiddleware)

	// Default welcome route
	root.Get("/", func(c *fiber.Ctx) error {
		data := sessionStore.GetSessionData(c)
		return response.Response(200, fmt.Sprintf("Dear %s %s Welcome to Codex-Compiler API (Private Zone)", data.Name, data.Surname), nil)
	})

	// Initialize additional routes
	h.initUserRoutes(root)
}

// authMiddleware checks the user's session and verifies it
func (h *PrivateHandler) authMiddleware(c *fiber.Ctx) error {
	secretHeader := c.Get("Codex-Compiler")
	if secretHeader != "b77759141fc85bf31e75b1d9c48bbe67" {
		return response.Response(403, "Forbidden!", nil)
	}

	// Get session ID from cookies
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return serviceErrors.NewServiceErrorWithMessage(401, "Unauthorized")
	}

	// Prepare the request for verification, passing session ID as a query parameter
	verifyURL := "http://api:8080/api/v1/private/user/profile"
	reqURL, err := url.Parse(verifyURL)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessage(500, "Error parsing verify URL")
	}

	// Create a new GET request to verify session
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessage(500, "Error creating GET request")
	}

	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	})

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessage(500, "Error making GET request to verify session")
	}
	defer resp.Body.Close()

	// Check if verification is successful
	if resp.StatusCode != http.StatusOK {
		return response.Response(401, "Unauthorized", nil)
	}

	// Log the response status and body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessage(500, "Error reading response body")
	}

	// Decode session data from the response body
	var data response.BaseResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error decoding session data:", err)
		return serviceErrors.NewServiceErrorWithMessage(500, "Error decoding session data")
	}

	userInfoJson, err := json.Marshal(data.Data)
	if err != nil {
		fmt.Println("Error marshalling data.Data:", err)
		return serviceErrors.NewServiceErrorWithMessage(500, "Error marshalling session data")
	}

	// Now unmarshal into the sessionStore.SessionData struct
	var userInfo sessionStore.SessionData
	if err := json.Unmarshal(userInfoJson, &userInfo); err != nil {
		fmt.Println("Error unmarshalling session data:", err)
		return serviceErrors.NewServiceErrorWithMessage(500, "Error unmarshalling session data")
	}

	// Check if the user is banned
	if userInfo.Role == "Banned" {
		return serviceErrors.NewServiceErrorWithMessage(403, "Banned")
	}

	// Attach user data to request context
	c.Locals("user", userInfo)

	// Proceed to the next middleware/handler
	return c.Next()
}
