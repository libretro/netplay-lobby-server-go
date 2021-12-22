package controller

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/libretro/netplay-lobby-server-go/domain"
	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// Template abspracts the template rendering.
type Template struct {
	templates *template.Template
}

// Render implements the echo template rendering interface
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// SessionDomain interface to decouple the controller logic from the domain code.
type SessionDomain interface {
	Add(request *domain.AddSessionRequest, ip net.IP) (*entity.Session, error)
	Get(roomID int32) (*entity.Session, error)
	List() ([]entity.Session, error)
	GetTunnel(tunnelName string) *domain.MitmInfo
	PurgeOld() error
}

// ListSessionsResponse is a custom DTO for backward compatability.
type SessionsResponse struct {
	Fields entity.Session `json:"fields"`
}

// SessionController handles all session related request
type SessionController struct {
	sessionDomain SessionDomain
}

// NewSessionController returns a new session controller
func NewSessionController(sessionDomain SessionDomain) *SessionController {
	return &SessionController{sessionDomain}
}

// RegisterRoutes registers all controller routes at an echo framework instance.
func (c *SessionController) RegisterRoutes(server *echo.Echo) {
	server.POST("/add", c.Add)
	server.POST("/add/", c.Add) // Legacy path
	server.GET("/list", c.List)
	server.GET("/list/", c.List) // Legacy path
	server.GET("/tunnel", c.Tunnel)
	server.GET("/tunnel/", c.Tunnel) // Legacy path
	server.GET("/", c.Index)
	server.GET("/:roomID", c.Get)
	server.GET("/:roomID/", c.Get) // Legacy path
}

// PrerenderTemplates prerenders all templates
func (c *SessionController) PrerenderTemplates(server *echo.Echo, filePattern string) error {
	templates, err := template.New("").Funcs(
		template.FuncMap{
			"prettyBool": func(b bool) string {
				if b {
					return "Yes"
				}
				return "No"
			},
			"prettyDate": func(d time.Time) string {
				utc, _ := time.LoadLocation("UTC")
				return d.In(utc).Format(time.RFC822)
			},
		},
	).ParseGlob(filePattern)

	if err != nil {
		return fmt.Errorf("Can't parse template: %w", err)
	}

	t := &Template{
		templates: templates,
	}
	server.Renderer = t
	return nil
}

// Index handler
// GET /
func (c *SessionController) Index(ctx echo.Context) error {
	logger := ctx.Logger()

	sessions, err := c.sessionDomain.List()
	if err != nil {
		logger.Errorf("Can't render session list: %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.Render(http.StatusOK, "index.html", sessions)
}

// Get handler
// GET /:roomID
func (c *SessionController) Get(ctx echo.Context) error {
	logger := ctx.Logger()

	roomIDString := ctx.Param("roomID")
	roomID, err := strconv.ParseInt(roomIDString, 10, 32)
	if err != nil {
		logger.Errorf("Can't get session by roomID. The RoomID wasn't a valid int32: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	session, err := c.sessionDomain.Get(int32(roomID))
	if err != nil || session == nil {
		logger.Errorf("Can't get session: %v", err)
		return ctx.NoContent(http.StatusNotFound)
	}

	// For legacy reasons, we need to put the sessions inside a wrapper object
	// that has the session accessible under the key "fields". The old implementation
	// also returned a list for one entry...
	response := make([]SessionsResponse, 1)
	response[0] = SessionsResponse{*session}
	return ctx.JSONPretty(http.StatusOK, response, "  ")
}

// List handler
// GET /list
func (c *SessionController) List(ctx echo.Context) error {
	logger := ctx.Logger()

	sessions, err := c.sessionDomain.List()
	if err != nil {
		logger.Errorf("Can't render session list: %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	// For legacy reasons, we need to put the sessions inside a wrapper object
	// that has the session accessible under the key "fields"
	response := make([]SessionsResponse, len(sessions))
	for i, session := range sessions {
		response[i].Fields = session
	}

	return ctx.JSONPretty(http.StatusOK, response, "  ")
}

// Add handler
// POST /add
func (c *SessionController) Add(ctx echo.Context) error {
	logger := ctx.Logger()
	var err error
	var session *entity.Session

	var req domain.AddSessionRequest
	if err := ctx.Bind(&req); err != nil {
		logger.Errorf("Can't parse incomming session: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	ip := net.ParseIP(ctx.RealIP())

	if session, err = c.sessionDomain.Add(&req, ip); err != nil {
		logger.Errorf("Won't add session: %v", err)

		if errors.Is(err, domain.ErrSessionRejected) {
			logger.Errorf("Rejected session: %v", session)
			return ctx.NoContent(http.StatusBadRequest)
		} else if errors.Is(err, domain.ErrRateLimited) {
			return ctx.NoContent(http.StatusTooManyRequests)
		}
		return ctx.NoContent(http.StatusBadRequest)
	}

	result := "status=OK\n"
	result += session.PrintForRetroarch()
	return ctx.String(http.StatusOK, result)
}

// Tunnel handler
// GET /tunnel
func (c *SessionController) Tunnel(ctx echo.Context) error {
	logger := ctx.Logger()

	tunnelName := ctx.QueryParam("name")
	if tunnelName == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	tunnel := c.sessionDomain.GetTunnel(tunnelName)
	if tunnel == nil {
		logger.Errorf("Can't find tunnel server: '%s'", tunnelName)
		return ctx.NoContent(http.StatusNotFound)
	}

	result := "status=OK\n"
	result += tunnel.PrintForRetroarch()
	return ctx.String(http.StatusOK, result)
}