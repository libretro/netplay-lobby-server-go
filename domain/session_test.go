package domain


import (
	"time"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

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

// TODO test SessionDomain.AddOrUpdate()

func TestSessionDomainPurgeOld(t *testing.T) {
	repoMock := SessionRepositoryMock{}
	sessionDomain := NewSessionDomain(&repoMock, &GeoIP2Domain{}, &ValidationDomain{})

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
	repoMock := SessionRepositoryMock{}
	sessionDomain := NewSessionDomain(&repoMock, &GeoIP2Domain{}, &ValidationDomain{})

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
