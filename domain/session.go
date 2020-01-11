package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// SessionDeadline is lifespan of a session that hasn't recieved any updated in seconds.
const SessionDeadline = 60

// RequestType enum
type RequestType int
// SessionAddType enum value
const (
    SessionCreate RequestType = iota
    SessionUpdate
    SessionTouch
)

// ErrSessionRejected is thrown when a session got rejected by the domain logic. 
var ErrSessionRejected = errors.New("Session rejected")

// ErrRateLimited is thrown when the rate limit is reached for a particular session.
var ErrRateLimited = errors.New("Rate limit reached")

// SessionDomain abrsracts the domain logic for netplay session handling.
type SessionDomain struct {
	sessionRepo SessionRepository
}

// NewSessionDomain returns an initalized SessionDomain struct.
func NewSessionDomain(sessionRepo SessionRepository) *SessionDomain {
	return &SessionDomain{sessionRepo}
}

// Add adds or updates a session, based on the incomming session information.
// Returns ErrSessionRejected if session got rejected.
// Returns ErrRateLimited if rate limit for a session got reached.
// TODO controller has to set the IP
// TODO middleware for rate-limiting an IP on request basis (not more than 10 adds a minute)
// TODO write test
func (d *SessionDomain) Add(session *entity.Session) (*entity.Session, error) {
	var err error
	var savedSession *entity.Session
	var requestType RequestType = SessionCreate

	// Decide if this is an CREATE, UPDATE or TOUCH operation
	session.CalculateID()
	if savedSession, err = d.sessionRepo.GetByID(session.ID); err != nil {
		return nil, fmt.Errorf("Can't get saved session: %w", err)
	}
	if savedSession != nil {
		requestType = SessionTouch
		session.CalculateContentHash()
		if savedSession.ContentHash != session.ContentHash {
			requestType = SessionUpdate
		}
	}

	// Validate session on CREATE and UPDATE
	if (requestType == SessionCreate || requestType == SessionUpdate) {
		if !d.validateSession(session) {
			return nil, ErrSessionRejected
		}
	}

	// Persist session changes
	switch (requestType) {
	case SessionCreate:
		if err = d.sessionRepo.Create(session); err != nil {
			return nil, fmt.Errorf("Can't create new session: %w", err)
		}
	case SessionUpdate:
		if err = d.sessionRepo.Update(session); err != nil {
			return nil, fmt.Errorf("Can't update old session: %w", err)
		}
	case SessionTouch:
		if err = d.sessionRepo.Touch(session.ID); err != nil {
			return nil, fmt.Errorf("Can't touch old session: %w", err)
		}
	}

	return session, nil
}

// List returns a list of all sessions that are currently beeing hosted
func (d *SessionDomain) List() ([]entity.Session, error) {
	sessions, err := d.sessionRepo.GetAll(d.getDeadline())
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// PurgeOld removes all sessions that have not been updated for longer than 45 seconds.
func (d *SessionDomain) PurgeOld() error {
	if err := d.sessionRepo.PurgeOld(d.getDeadline()); err != nil {
		return err
	}
	return nil
}

// TODO write test
func (d *SessionDomain) validateSession(s *entity.Session) bool {
	return false
}

func (d *SessionDomain) getDeadline() time.Time {
	return time.Now().Add(-SessionDeadline * time.Second)
}

// SessionRepository interface to decouple the domain logic from the repository code.
type SessionRepository interface {
	Create(s *entity.Session) error
	GetByID(id string) (*entity.Session, error)
	GetAll(deadline time.Time) ([]entity.Session, error)
	Update(s *entity.Session) error
	Touch(id string) error
	PurgeOld(deadline time.Time) error
}
