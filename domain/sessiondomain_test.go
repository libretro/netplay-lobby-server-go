package domain

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

var testIP = net.ParseIP("192.168.178.2")

// testRequest and testSession should have the same values for the test below.
var playercnt int16 = 2
var spectacnt int16 = 1

var testRequest = AddSessionRequest{
	Username:            "zelda",
	CoreName:            "bsnes",
	CoreVersion:         "0.2.1",
	GameName:            "supergame",
	GameCRC:             "FFFFFFFF",
	Port:                55355,
	MITMServer:          "",
	HasPassword:         false,
	HasSpectatePassword: false,
	ForceMITM:           false,
	RetroArchVersion:    "1.1.1",
	Frontend:            "retro",
	SubsystemName:       "subsub",
	MITMSession:         "",
	MITMCustomServer:    "",
	MITMCustomPort:      0,
	PlayerCount:         &playercnt,
	SpectatorCount:      &spectacnt,
}

var testSession = entity.Session{
	ID:                  "",
	RoomID:              100,
	Username:            "zelda",
	Country:             "en",
	GameName:            "supergame",
	GameCRC:             "FFFFFFFF",
	CoreName:            "bsnes",
	CoreVersion:         "0.2.1",
	SubsystemName:       "subsub",
	RetroArchVersion:    "1.1.1",
	Frontend:            "retro",
	IP:                  net.ParseIP("192.168.178.2"),
	Port:                55355,
	MitmHandle:          "",
	MitmAddress:         "",
	MitmPort:            0,
	MitmSession:         "",
	HostMethod:          entity.HostMethodUnknown,
	HasPassword:         false,
	HasSpectatePassword: false,
	Connectable:         true,
	IsRetroArch:         true,
	PlayerCount:         2,
	SpectatorCount:      1,
	CreatedAt:           time.Now().Add(-5 * time.Minute),
	UpdatedAt:           time.Now().Add(-5 * time.Minute),
	ContentHash:         "",
}

type SessionRepositoryMock struct {
	mock.Mock
}

func (m *SessionRepositoryMock) Create(s *entity.Session) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *SessionRepositoryMock) Update(s *entity.Session) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *SessionRepositoryMock) Touch(s *entity.Session) error {
	args := m.Called(s.ID)
	return args.Error(0)
}

func (m *SessionRepositoryMock) GetByID(id string) (*entity.Session, error) {
	args := m.Called(id)
	session, _ := args.Get(0).(*entity.Session)
	return session, args.Error(1)
}

func (m *SessionRepositoryMock) GetByRoomID(roomID int32) (*entity.Session, error) {
	args := m.Called(roomID)
	session, _ := args.Get(0).(*entity.Session)
	return session, args.Error(1)
}

func (m *SessionRepositoryMock) GetAll(deadline time.Time) ([]entity.Session, error) {
	args := m.Called(deadline)
	sessions, _ := args.Get(0).([]entity.Session)
	return sessions, args.Error(1)
}

func (m *SessionRepositoryMock) PurgeOld(deadline time.Time) error {
	args := m.Called(deadline)
	return args.Error(0)
}

func setupSessionDomain(t *testing.T) (*SessionDomain, *SessionRepositoryMock) {
	repoMock := SessionRepositoryMock{}

	validationDomain, err := NewValidationDomain(testStringBlacklist, testIPBlacklist)
	require.NoError(t, err)

	geoip2Domain := setupGeoip2Domain(t)

	sessionDomain := NewSessionDomain(&repoMock, geoip2Domain, validationDomain, &MitmDomain{})
	require.NoError(t, err)

	return sessionDomain, &repoMock
}

func TestSessionDomainPurgeOld(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	// Test the deadline duration
	repoMock.On("PurgeOld", mock.MatchedBy(
		func(d time.Time) bool {
			before := time.Now().Add(-(SessionDeadline - 1) * time.Second)
			after := time.Now().Add(-(SessionDeadline + 1) * time.Second)
			return d.Before(before) && d.After(after)
		})).Return(nil)

	err := sessionDomain.PurgeOld()
	require.NoError(t, err, "Can't purge old sessions")
}

func TestSessionDomainList(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	// Test the deadline duration
	repoMock.On("GetAll", mock.MatchedBy(
		func(d time.Time) bool {
			before := time.Now().Add(-(SessionDeadline - 1) * time.Second)
			after := time.Now().Add(-(SessionDeadline + 1) * time.Second)
			return d.Before(before) && d.After(after)
		})).Return(make([]entity.Session, 3), nil)

	sessions, err := sessionDomain.List()
	require.NoError(t, err, "Can't list sessions")
	require.NotNil(t, sessions)
	assert.Equal(t, 3, len(sessions))
}

func TestSessionDomainValidateSessionAtCreate(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	request := testRequest
	comp := testSession
	comp.CalculateID()
	comp.CalculateContentHash()

	request.GameCRC = "123456789"

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
			return s == comp.ID
		})).Return(nil, nil)

	newSession, err := sessionDomain.Add(&request, testIP)
	require.Error(t, err)
	assert.Nil(t, newSession)
	assert.True(t, errors.Is(err, ErrSessionRejected))
}

func TestSessionDomainValidateSessionAtUpdate(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	request := testRequest
	comp := testSession
	comp.CalculateID()
	comp.CalculateContentHash()

	request.RetroArchVersion = "0123456789ABCDEF0123456789ABCDEF_INVALID"

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
			return s == comp.ID
		})).Return(&comp, nil)

	newSession, err := sessionDomain.Add(&request, testIP)
	require.Error(t, err)
	assert.Nil(t, newSession)
	assert.True(t, errors.Is(err, ErrSessionRejected))
}

func TestSessionDomainAddSessionTypeCreate(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	request := testRequest
	comp := testSession
	comp.CalculateID()
	comp.CalculateContentHash()

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
			return s == comp.ID
		})).Return(nil, nil)

	repoMock.On("Create", mock.MatchedBy(
		func(s *entity.Session) bool {
			return s.ID == comp.ID && s.ContentHash == comp.ContentHash
		})).Return(nil)

	newSession, err := sessionDomain.Add(&request, testIP)
	require.NoError(t, err)
	require.NotNil(t, newSession)
	assert.Equal(t, comp.ID, newSession.ID)
	assert.Equal(t, comp.ContentHash, newSession.ContentHash)
}

func TestSessionDomainAddSessionTypeCreateShouldSetDefaultUsername(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	request := testRequest
	request.Username = ""

	comp := testSession
	comp.Username = "Anonymous"
	comp.CalculateID()
	comp.CalculateContentHash()

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
			return s == comp.ID
		})).Return(nil, nil)

	repoMock.On("Create", mock.MatchedBy(
		func(s *entity.Session) bool {
			return s.ID == comp.ID && s.ContentHash == comp.ContentHash
		})).Return(nil)

	newSession, err := sessionDomain.Add(&request, testIP)
	require.NoError(t, err)
	require.NotNil(t, newSession)
	assert.Equal(t, comp.ID, newSession.ID)
	assert.Equal(t, comp.ContentHash, newSession.ContentHash)
}

func TestSessionDomainAddSessionTypeUpdate(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	request := testRequest
	comp := testSession
	comp.CalculateID()
	comp.CalculateContentHash()

	request.GameCRC = "88888888"

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
			return s == comp.ID
		})).Return(&comp, nil)

	repoMock.On("Update", mock.MatchedBy(
		func(s *entity.Session) bool {
			return s.ID == comp.ID && s.ContentHash != comp.ContentHash
		})).Return(nil)

	newSession, err := sessionDomain.Add(&request, testIP)
	require.NoError(t, err)
	require.NotNil(t, newSession)
	assert.Equal(t, comp.ID, newSession.ID)
	assert.NotEqual(t, comp.ContentHash, newSession.ContentHash)
}

func TestSessionDomainAddSessionTypeTouch(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	request := testRequest
	comp := testSession
	comp.CalculateID()
	comp.CalculateContentHash()

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
			return s == comp.ID
		})).Return(&comp, nil)

	repoMock.On("Touch", mock.MatchedBy(
		func(id string) bool {
			return id == comp.ID
		})).Return(nil)

	newSession, err := sessionDomain.Add(&request, testIP)
	require.NoError(t, err)
	require.NotNil(t, newSession)
	assert.Equal(t, comp.ID, newSession.ID)
	assert.Equal(t, comp.ContentHash, newSession.ContentHash)
}

func TestSessionDomainAddSessionTypeUpdateRateLimit(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	request := testRequest
	comp := testSession
	comp.UpdatedAt = time.Now().Add(-4 * time.Second)
	comp.CalculateID()
	comp.CalculateContentHash()

	request.GameCRC = "88888888"

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
			return s == comp.ID
		})).Return(&comp, nil)

	repoMock.On("Update", mock.MatchedBy(
		func(s *entity.Session) bool {
			return s.ID == comp.ID && s.ContentHash != comp.ContentHash
		})).Return(nil)

	newSession, err := sessionDomain.Add(&request, testIP)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrRateLimited))
	assert.Nil(t, newSession)
}

func TestSessionDomainAddSessionTypeTouchRateLimit(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	request := testRequest
	comp := testSession
	comp.UpdatedAt = time.Now().Add(-4 * time.Second)
	comp.CalculateID()
	comp.CalculateContentHash()

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
			return s == comp.ID
		})).Return(&comp, nil)

	repoMock.On("Touch", mock.MatchedBy(
		func(id string) bool {
			return id == comp.ID
		})).Return(nil)

	newSession, err := sessionDomain.Add(&request, testIP)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrRateLimited))
	assert.Nil(t, newSession)
}
