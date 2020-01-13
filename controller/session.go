package controller


import (
	"errors"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/libretro/netplay-lobby-server-go/domain"
	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// SessionController handles all session related request
type SessionController struct {
	sessionDomain domain.SessionDomain
}

// ListSessionsResponse is a custom DTO for backward compatability.
type ListSessionsResponse struct {
	Fields entity.Session `json:"fields"`
}

// RegisterRoutes registers all controller routes at an echo framework instance.
func (c *SessionController) RegisterRoutes(e *echo.Echo) {
	e.POST("/add", c.Add)
	e.GET("/list", c.List)
	e.GET("/", c.Index)
}

// Index handler
// GET /
// TODO testme
func (c *SessionController) Index(ctx echo.Context) error {
	logger := ctx.Logger()

	_, err := c.sessionDomain.List()
	if err != nil {
		logger.Errorf("Can't render session list: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError);
	}
	// TODO template rendering

	return ctx.String(http.StatusOK, "TODO")
}

// List handler
// GET /list
// TODO testme
func (c *SessionController) List(ctx echo.Context) error {
	logger := ctx.Logger()

	sessions, err := c.sessionDomain.List()
	if err != nil {
		logger.Errorf("Can't render session list: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError);
	}

	// For legacy reasons, we need to put the sessions inside a wrapper object
	// that has the session accessible under the key "fields"
	response := make([]ListSessionsResponse, len(sessions))
	for i, session := range sessions {
		response[i].Fields = session
	}

	return ctx.JSONPretty(http.StatusOK, response, "  ")
}

// Add handler
// GET /add
// TODO testme
func (c *SessionController) Add(ctx echo.Context) error {
	logger := ctx.Logger()
	var err error
	var session *entity.Session

	var req domain.AddSessionRequest
	if err := ctx.Bind(&req); err != nil {
		logger.Errorf("Can't parse incomming session: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest);
	}

	ip := net.ParseIP(ctx.RealIP())

	if session, err = c.sessionDomain.Add(&req, ip); err != nil {
		logger.Errorf("Won't add session: %w", err)

		if errors.Is(err, domain.ErrSessionRejected) {
			logger.Errorf("Rejected session: %v", session)
			return echo.NewHTTPError(http.StatusBadRequest);
		} else if errors.Is(err, domain.ErrRateLimited) {
			return echo.NewHTTPError(http.StatusTooManyRequests);
		}
		return echo.NewHTTPError(http.StatusBadRequest);
	}

	result := "status=OK\n"
	result += session.PrintForRetroarch()
	return ctx.String(http.StatusOK, result)
}
