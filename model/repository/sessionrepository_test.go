package repository

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/libretro/netplay-lobby-server-go/model"
	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

var testSession = entity.Session{
	ID:                  "",
	RoomID:              0,
	Username:            "zelda",
	Country:             "EN",
	GameName:            "supergame",
	GameCRC:             "FFFFFFFF",
	CoreName:            "unes",
	CoreVersion:         "0.2.1",
	SubsystemName:       "subsub",
	RetroArchVersion:    "1.1.1",
	Frontend:            "retro",
	IP:                  net.ParseIP("127.0.0.1"),
	Port:                55355,
	MitmAddress:         "hostname.com",
	MitmPort:            0,
	HostMethod:          entity.HostMethodUPNP,
	HasPassword:         false,
	HasSpectatePassword: false,
	CreatedAt:           time.Now(),
	UpdatedAt:           time.Now(),
	ContentHash:         "",
}

func setupSessionRepository(t *testing.T) *SessionRepository {
	db, err := model.GetSqliteDB(":memory:")
	//db = db.LogMode(true)
	if err != nil {
		t.Fatalf("Can't open sqlite3 db: %v", err)
	}
	db.AutoMigrate(entity.Session{})

	return NewSessionRepository(db)
}

func TestSessionRepositoryCreate(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")
}

func TestSessionRepositoryCreateIPNotNull(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.IP = nil
	session.MitmAddress = "newhost.com"
	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	require.Error(t, err, "Should not be able to create nil value for IP")
}

func TestSessionRepositoryCreateMITMNull(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.MitmAddress = ""
	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	require.NoError(t, err, "Should be able to create nil value for MITM IP")
}

func TestSessionRepositoryGetByID(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	newSession, err := sessionRepository.GetByID(session.ID)
	require.NoError(t, err, "Can't get session by ID")

	require.NotNil(t, newSession)

	assert.NotEmpty(t, newSession.ID)
	assert.NotEmpty(t, newSession.ContentHash)

	assert.NotEqual(t, session.CreatedAt, newSession.CreatedAt)
	assert.NotEqual(t, session.UpdatedAt, newSession.UpdatedAt)

	// Make sure the rest is the same
	session.CreatedAt = newSession.CreatedAt
	session.UpdatedAt = newSession.UpdatedAt
	assert.Equal(t, &session, newSession)

	noSession, err := sessionRepository.GetByID("should_not_exists")
	require.NoError(t, err, "Can't get nil session")
	assert.Nil(t, noSession)
}

func TestSessionRepositoryGetByRoomID(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.RoomID = 400
	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	newSession, err := sessionRepository.GetByRoomID(session.RoomID)
	require.NoError(t, err, "Can't get session by RoomID")

	require.NotNil(t, newSession)

	assert.NotEmpty(t, newSession.ID)
	assert.NotEmpty(t, newSession.ContentHash)

	assert.NotEqual(t, session.CreatedAt, newSession.CreatedAt)
	assert.NotEqual(t, session.UpdatedAt, newSession.UpdatedAt)

	// Make sure the rest is the same
	session.CreatedAt = newSession.CreatedAt
	session.UpdatedAt = newSession.UpdatedAt
	assert.Equal(t, &session, newSession)

	noSession, err := sessionRepository.GetByID("should_not_exists")
	require.NoError(t, err, "Can't get nil session")
	assert.Nil(t, noSession)
}

func TestSessionRepositoryGetAll(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.CalculateID()
	session.CalculateContentHash()
	session.RoomID = 0
	err := sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	session.Username = "aladin"
	session.CalculateID()
	session.CalculateContentHash()
	session.RoomID = 0
	err = sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	session.Username = "invalid"
	session.UpdatedAt = time.Now().Add(-2 * time.Minute)
	session.CalculateID()
	session.CalculateContentHash()
	session.RoomID = 0
	err = sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	deadline := time.Now().Add(-1 * time.Minute)
	sessions, err := sessionRepository.GetAll(deadline)
	require.NoError(t, err, "Can't get all sessions with deadline")

	require.NotNil(t, sessions)
	require.Equal(t, 2, len(sessions), "Query seems to include non valid entries.")
	assert.Less(t, sessions[0].Username, sessions[1].Username, "Sessions are not ordered by username")

	sessions, err = sessionRepository.GetAll(time.Time{})
	require.NoError(t, err, "Can't get all sessions without deadline")

	require.NotNil(t, sessions)
	require.Equal(t, 3, len(sessions), "Query seems to not include invalid entries.")
	assert.Less(t, sessions[0].Username, sessions[1].Username, "Sessions are not ordered by username")
	assert.Less(t, sessions[1].Username, sessions[2].Username, "Sessions are not ordered by username")
}

func TestSessionRepositoryUpdate(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	newIP := "83.12.41.222"
	session.MitmAddress = newIP

	session.CalculateContentHash()
	err = sessionRepository.Update(&session)
	require.NoError(t, err, "Can't update session")

	newSession, err := sessionRepository.GetByID(session.ID)
	require.NoError(t, err, "Can't get session by ID")

	require.NotNil(t, newSession)
	assert.Equal(t, newSession.MitmAddress, newIP)
}

func TestSessionRepositoryTouch(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	oldSession, err := sessionRepository.GetByID(session.ID)
	require.NoError(t, err, "Can't get session by ID")

	oldTimestamp := oldSession.UpdatedAt

	err = sessionRepository.Touch(oldSession.ID)
	require.NoError(t, err, "Can't touch session")

	newSession, err := sessionRepository.GetByID(session.ID)
	require.NoError(t, err, "Can't get session by ID")

	require.NotNil(t, newSession)
	assert.False(t, newSession.UpdatedAt.Equal(oldTimestamp), "New timestamp did not change after touch")
	assert.True(t, newSession.UpdatedAt.After(oldTimestamp), "New timestamp is not older after touch")

	assert.Equal(t, oldSession.ContentHash, newSession.ContentHash)
}

func TestSessionRepositoryPurgeOld(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.CalculateID()
	session.CalculateContentHash()
	session.RoomID = 0
	err := sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	session.Username = "aladin"
	session.CalculateID()
	session.CalculateContentHash()
	session.RoomID = 0
	err = sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	session.Username = "invalid"
	session.UpdatedAt = time.Now().Add(-2 * time.Minute)
	session.CalculateID()
	session.CalculateContentHash()
	session.RoomID = 0
	err = sessionRepository.Create(&session)
	require.NoError(t, err, "Can't create session")

	deadline := time.Now().Add(-1 * time.Minute)
	err = sessionRepository.PurgeOld(deadline)
	require.NoError(t, err, "Can't purge old sessions")

	sessions, err := sessionRepository.GetAll(time.Time{})
	require.NoError(t, err, "Can't get all sessions")

	require.NotNil(t, sessions)
	require.Equal(t, len(sessions), 2, "Query seems to include non valid entries.")
}
