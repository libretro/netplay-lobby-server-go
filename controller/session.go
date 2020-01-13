package controller


import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/libretro/netplay-lobby-server-go/domain"
	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// SessionController handles all session related request
type SessionController struct {
	sessionDomain domain.SessionDomain
}

// AddSessionRequest defines the request for the AddSession request.
type AddSessionRequest struct {
	Username string `form:"username"`
	CoreName string `form:"core_name"`
	CoreVersion string `form:"core_version"`
	GameName string `form:"game_name"`
	GameCRC  string `form:"game_crc"`
	Port uint16 `form:"port"`
	MITMServer string `form:"mitm_server"`
	HasPassword bool  `form:"has_password"` // 1/0
	HasSpectatePassword bool `form:"has_spectate_password"`
	ForceMITM bool `form:"force_mitm"`
	RetroarchVersion string `form:"retroarch_version"`
	Frontend string `form:"frontend"`
	SubsystemName string  `form:"subsystem_name"`
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

	var req AddSessionRequest
	if err := ctx.Bind(&req); err != nil {
		logger.Errorf("Can't parse incomming session: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest);
	}

	// TODO MITM server logic
	session := &entity.Session{
		Username: req.Username,
		GameName: req.GameName,
		GameCRC: strings.ToUpper(req.GameCRC),
		CoreName: req.CoreName,
		CoreVersion: req.CoreVersion,
		SubsystemName: req.SubsystemName,
		RetroArchVersion: req.RetroarchVersion,
		Frontend: req.Frontend,
		IP: net.ParseIP(ctx.RealIP()),
		Port: req.Port,
		MitmAddress: "", // TODO
		MitmPort: 0, // TODO
		HostMethod: entity.HostMethodUnknown,
		HasPassword: req.HasPassword,
		HasSpectatePassword: req.HasSpectatePassword,
	}

	if session, err = c.sessionDomain.Add(session); err != nil {
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
