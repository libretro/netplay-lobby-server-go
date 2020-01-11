package domain

import (
	"time"

	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// SessionDomain abrsracts the domain logic for netplay session handling.
type SessionDomain struct {
	sessionRepo SessionRepository
}

// NewSessionDomain returns an initalized SessionDomain struct.
func NewSessionDomain(sessionRepo SessionRepository) *SessionDomain {
	return &SessionDomain{sessionRepo}
}

// AddOrUpdate adds or updates a session, based on the incomming session information.
func (d *SessionDomain) AddOrUpdate(s *entity.Session) error {
	// Parse POST body

	// Look if the CalculateHashID() is inside the database to see, if this is an ADD or UPDATE.
	// Look if the CalculateHashContent changed. If so: update_content=true

	// Do NOT updated a Session that comes from another IP.

	// For ADD and UPDATE with update_content == true, validate user input
	// Add/Update session in DB

	// If Update with update_content == false, touch session in DB


	// Return result in BODY as simple KWARGS as TEXT
	return nil
}

// PurgeOld removes all sessions that have not been updated for longer than 45 seconds.
func (d *SessionDomain) PurgeOld() error {
	if err := d.sessionRepo.PurgeOld(getDeadline()); err != nil {
		return err
	}
	return nil
}

// List returns a list of all sessions that are currently beeing hosted
func (d *SessionDomain) List() ([]entity.Session, error) {
	sessions, err := d.sessionRepo.GetAll(getDeadline())
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func getDeadline() time.Time {
	return time.Now().Add(-60 * time.Second)
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