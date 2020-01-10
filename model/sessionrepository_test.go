package model

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func setupSessionRepository(t *testing.T) *SessionRepository {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Can't open sqlite3 db: %v", err)
	}
	db.AutoMigrate(Session{})

	return NewSessionRepository(db)
}

func TestSessionRepositoryCreate(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession

	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	assert.NoError(t, err, "Can't create session")
}

func TestSessionRepositoryGetByID(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession
	
	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	assert.NoError(t, err, "Can't create session")

	newSession, err := sessionRepository.GetByID(session.ID)
	assert.NoError(t, err, "Can't get session by ID")

	assert.NotNil(t, newSession)
	if newSession != nil {
		assert.NotEmpty(t, newSession.ID)
		assert.NotEmpty(t, newSession.ContentHash)

		assert.NotEqual(t, session.CreatedAt, newSession.CreatedAt)
		assert.NotEqual(t, session.UpdatedAt, newSession.UpdatedAt)

		// Make sure the rest is the same
		session.CreatedAt = newSession.CreatedAt
		session.UpdatedAt = newSession.UpdatedAt
		assert.Equal(t, &session, newSession)
	}
}

func TestSessionRepositoryGetAll(t *testing.T) {
	sessionRepository := setupSessionRepository(t)
	session := testSession
	
	session.CalculateID()
	session.CalculateContentHash()
	err := sessionRepository.Create(&session)
	assert.NoError(t, err, "Can't create session")

	session.Username = "aladin"
	session.CalculateID()
	session.CalculateContentHash()
	err = sessionRepository.Create(&session)
	assert.NoError(t, err, "Can't create session")

	sessions, err := sessionRepository.GetAll()
	assert.NoError(t, err, "Can't get session by ID")

	assert.NotNil(t, sessions)
	if sessions != nil {
		assert.NotEmpty(t, len(sessions), 2)
		assert.Less(t, sessions[0].Username, sessions[1].Username, "sessions are not ordered by username")
	}
}

// TODO add tests for Update, Touch and PurgeOld
