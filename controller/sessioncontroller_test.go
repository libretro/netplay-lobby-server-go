package controller

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/libretro/netplay-lobby-server-go/domain"
	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

var testSession = entity.Session{
	ID:                  "",
	RoomID:              0,
	Username:            "zelda",
	Country:             "en",
	GameName:            "supergame",
	GameCRC:             "FFFFFFFF",
	CoreName:            "unes",
	CoreVersion:         "0.2.1",
	SubsystemName:       "subsub",
	RetroArchVersion:    "1.1.1",
	Frontend:            "retro",
	IP:                  net.ParseIP("127.0.0.1"),
	Port:                55355,
	MitmAddress:         "0.0.0.0",
	MitmPort:            0,
	HostMethod:          entity.HostMethodUPNP,
	HasPassword:         false,
	HasSpectatePassword: false,
	CreatedAt:           time.Date(2010, 9, 12, 11, 33, 05, 0, time.UTC),
	UpdatedAt:           time.Date(2010, 9, 12, 11, 33, 05, 0, time.UTC),
	ContentHash:         "",
}

type SessionDomainMock struct {
	mock.Mock
}

func (m *SessionDomainMock) Add(request *domain.AddSessionRequest, ip net.IP) (*entity.Session, error) {
	args := m.Called(request, ip)
	session, _ := args.Get(0).(*entity.Session)
	return session, args.Error(1)
}

func (m *SessionDomainMock) Get(roomID int32) (*entity.Session, error) {
	args := m.Called(roomID)
	session, _ := args.Get(0).(*entity.Session)
	return session, args.Error(1)
}

func (m *SessionDomainMock) List() ([]entity.Session, error) {
	args := m.Called()
	sessions, _ := args.Get(0).([]entity.Session)
	return sessions, args.Error(1)
}

func (m *SessionDomainMock) PurgeOld() error {
	args := m.Called()
	return args.Error(0)
}

func TestSessionControllerIndex(t *testing.T) {
	domainMock := &SessionDomainMock{}

	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	handler := NewSessionController(domainMock)
	err := handler.PrerenderTemplates(server, "../web/templates/*.html")
	require.NoError(t, err)

	session1 := testSession
	session1.Username = "Player 1"
	session2 := testSession
	session2.Username = "Player 2"
	sessions := []entity.Session{session1, session2}
	domainMock.On("List").Return(sessions, nil)

	handler.Index(ctx)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Player 1")
	assert.Contains(t, rec.Body.String(), "Player 2")
}

func TestSessionControllerGet(t *testing.T) {
	domainMock := &SessionDomainMock{}

	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	ctx.SetPath("/:roomID")
	ctx.SetParamNames("roomID")
	ctx.SetParamValues("100")
	handler := NewSessionController(domainMock)

	session := testSession
	expectedResultBody := `{
  "fields": {
    "id": 0,
    "username": "zelda",
    "country": "en",
    "game_name": "supergame",
    "game_crc": "FFFFFFFF",
    "core_name": "unes",
    "core_version": "0.2.1",
    "subsystem_name": "subsub",
    "retroarch_version": "1.1.1",
    "frontend": "retro",
    "ip": "127.0.0.1",
    "port": 55355,
    "mitm_ip": "0.0.0.0",
    "mitm_port": 0,
    "host_method": 2,
    "has_password": false,
    "has_spectate_password": false,
    "created": "2010-09-12T11:33:05Z",
    "updated": "2010-09-12T11:33:05Z"
  }
}
`
	domainMock.On("Get", int32(100)).Return(&session, nil)

	handler.Get(ctx)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResultBody, rec.Body.String())
}

func TestSessionControllerList(t *testing.T) {
	domainMock := &SessionDomainMock{}

	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	handler := NewSessionController(domainMock)

	sessions := make([]entity.Session, 1)
	expectedResultBody := `[
  {
    "fields": {
      "id": 0,
      "username": "",
      "country": "",
      "game_name": "",
      "game_crc": "",
      "core_name": "",
      "core_version": "",
      "subsystem_name": "",
      "retroarch_version": "",
      "frontend": "",
      "ip": "",
      "port": 0,
      "mitm_ip": "",
      "mitm_port": 0,
      "host_method": 0,
      "has_password": false,
      "has_spectate_password": false,
      "created": "0001-01-01T00:00:00Z",
      "updated": "0001-01-01T00:00:00Z"
    }
  }
]
`
	domainMock.On("List").Return(sessions, nil)

	handler.List(ctx)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResultBody, rec.Body.String())
}

func TestSessionControllerListError(t *testing.T) {
	domainMock := &SessionDomainMock{}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	handler := NewSessionController(domainMock)

	domainMock.On("List").Return(nil, errors.New("test error"))

	handler.List(ctx)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "", rec.Body.String())
}
