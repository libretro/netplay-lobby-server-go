package controller

import (
	"net"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/libretro/netplay-lobby-server-go/domain"
	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

type SessionDomainMock struct {
	mock.Mock
}

func (m *SessionDomainMock) Add(request *domain.AddSessionRequest, ip net.IP) (*entity.Session, error) {
	args := m.Called(request, ip)
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

// TODO test Index

func TestSessionControllerList(t *testing.T) {
	// TODO move setup code in setup method
	domainMock := &SessionDomainMock{}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	handler := NewSessionController(domainMock)

	sessions := make([]entity.Session, 1)
	expectedResultBody := `[
  {
    "fields": {
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

	err := handler.List(ctx)
	require.NoError(t, err, "Can't make request to list")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResultBody, rec.Body.String())
}

func TestSessionControllerListError(t *testing.T) {
	domainMock := &SessionDomainMock{}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	handler := NewSessionController(domainMock)

	domainMock.On("List").Return(nil, errors.New("just a simple test"))

	err := handler.List(ctx)
	require.Error(t, err, "Expected error is not occuring")
	assert.Equal(t, http.StatusInternalServerError, rec.Code) // TODO
	assert.Equal(t, "", rec.Body.String())
}

// TODO test Add
