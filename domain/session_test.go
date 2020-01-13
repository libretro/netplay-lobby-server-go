package domain


import (
	"net"
	"errors"
	"time"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

var testSession = entity.Session{
	ID:                  "",
	Username:            "zelda",
	Country:             "EN",
	GameName:            "supergame",
	GameCRC:             "FFFFFFFF",
	CoreName:            "bsnes",
	CoreVersion:         "0.2.1",
	SubsystemName:       "subsub",
	RetroArchVersion:    "1.1.1",
	Frontend:            "retro",
	IP:                  net.ParseIP("192.168.178.2"),
	Port:                55355,
	MitmAddress:         "hostname.com",
	MitmPort:            0,
	HostMethod:          entity.HostMethodUPNP,
	HasPassword:         false,
	HasSpectatePassword: false,
	CreatedAt:           time.Now().Add(-5 * time.Minute),
	UpdatedAt:           time.Now().Add(-5 * time.Minute),
	ContentHash:         "",
}

type SessionRepositoryMock struct{
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
func (m *SessionRepositoryMock) Touch(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *SessionRepositoryMock) GetByID(id string) (*entity.Session, error) {
	args := m.Called(id)
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

	validationDomain, err := NewValidationDomain(testCoreWhitelist, testStringBlacklist, testIPBlacklist)
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
			before := time.Now().Add(-(SessionDeadline-1) * time.Second)
			after := time.Now().Add(-(SessionDeadline+1) * time.Second)
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
			before := time.Now().Add(-(SessionDeadline-1) * time.Second)
			after := time.Now().Add(-(SessionDeadline+1) * time.Second)
        	return d.Before(before) && d.After(after)
		})).Return(make([]entity.Session, 3), nil)

	sessions, err := sessionDomain.List()
	require.NoError(t, err, "Can't list sessions")
	require.NotNil(t, sessions)
	assert.Equal(t, 3, len(sessions))
}

func TestSessionDomainValidateSessionAtCreate(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	session := testSession
	comp := session
	comp.CalculateID()
	comp.CalculateContentHash()
	session.GameCRC = "123456789"

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
        	return s == comp.ID
		})).Return(nil, nil)

	newSession, err := sessionDomain.Add(&session)
	require.Error(t, err)
	assert.Nil(t, newSession)
	assert.True(t, errors.Is(err, ErrSessionRejected))
}

func TestSessionDomainValidateSessionAtUpdate(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	session := testSession
	comp := session
	comp.CalculateID()
	comp.CalculateContentHash()
	session.RetroArchVersion = "0123456789ABCDEF0123456789ABCDEF_INVALID"

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
        	return s == comp.ID
		})).Return(&comp, nil)

	newSession, err := sessionDomain.Add(&session)
	require.Error(t, err)
	assert.Nil(t, newSession)
	assert.True(t, errors.Is(err, ErrSessionRejected))
}

func TestSessionDomainAddSessionTypeCreate(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	session := testSession
	comp := session
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

	newSession, err := sessionDomain.Add(&session)
	require.NoError(t, err)
	require.NotNil(t, newSession)
	assert.Equal(t, comp.ID, newSession.ID)
	assert.Equal(t, comp.ContentHash, newSession.ContentHash)
}

func TestSessionDomainAddSessionTypeUpdate(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	session := testSession
	comp := session
	comp.CalculateID()
	comp.CalculateContentHash()

	session.GameCRC = "88888888"

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
        	return s == comp.ID
		})).Return(&comp, nil)

	repoMock.On("Update", mock.MatchedBy(
		func(s *entity.Session) bool {
			return s.ID == comp.ID && s.ContentHash != comp.ContentHash
		})).Return(nil)

	newSession, err := sessionDomain.Add(&session)
	require.NoError(t, err)
	require.NotNil(t, newSession)
	assert.Equal(t, comp.ID, newSession.ID)
	assert.NotEqual(t, comp.ContentHash, newSession.ContentHash)
}

func TestSessionDomainAddSessionTypeTouch(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	session := testSession
	comp := session
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

	newSession, err := sessionDomain.Add(&session)
	require.NoError(t, err)
	require.NotNil(t, newSession)
	assert.Equal(t, comp.ID, newSession.ID)
	assert.Equal(t, comp.ContentHash, newSession.ContentHash)
}

func TestSessionDomainAddSessionTypeUpdateRateLimit(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	session := testSession
	comp := session
	comp.UpdatedAt = time.Now().Add(-5 * time.Second)
	comp.CalculateID()
	comp.CalculateContentHash()

	session.GameCRC = "88888888"

	repoMock.On("GetByID", mock.MatchedBy(
		func(s string) bool {
        	return s == comp.ID
		})).Return(&comp, nil)

	repoMock.On("Update", mock.MatchedBy(
		func(s *entity.Session) bool {
			return s.ID == comp.ID && s.ContentHash != comp.ContentHash
		})).Return(nil)

	newSession, err := sessionDomain.Add(&session)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrRateLimited))
	assert.Nil(t, newSession)
}

func TestSessionDomainAddSessionTypeTouchRateLimit(t *testing.T) {
	sessionDomain, repoMock := setupSessionDomain(t)

	session := testSession
	comp := session
	comp.UpdatedAt = time.Now().Add(-5 * time.Second)
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

	newSession, err := sessionDomain.Add(&session)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrRateLimited))
	assert.Nil(t, newSession)
}

// TODO test MITM Codepath
