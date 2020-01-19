package domain

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// SessionDeadline is lifespan of a session that hasn't recieved any updated in seconds.
const SessionDeadline = 60

// requestType enum
type requestType int

// SessionAddType enum value
const (
	SessionCreate requestType = iota
	SessionUpdate
	SessionTouch
)

// AddSessionRequest defines the request for the SessionDomain.Add() request.
type AddSessionRequest struct {
	Username            string `form:"username"`
	CoreName            string `form:"core_name"`
	CoreVersion         string `form:"core_version"`
	GameName            string `form:"game_name"`
	GameCRC             string `form:"game_crc"`
	Port                uint16 `form:"port"`
	MITMServer          string `form:"mitm_server"`
	HasPassword         bool   `form:"has_password"` // 1/0 (Can it be bound to bool?)
	HasSpectatePassword bool   `form:"has_spectate_password"`
	ForceMITM           bool   `form:"force_mitm"`
	RetroArchVersion    string `form:"retroarch_version"`
	Frontend            string `form:"frontend"`
	SubsystemName       string `form:"subsystem_name"`
}

// ErrSessionRejected is thrown when a session got rejected by the domain logic.
var ErrSessionRejected = errors.New("Session rejected")

// ErrRateLimited is thrown when the rate limit is reached for a particular session.
var ErrRateLimited = errors.New("Rate limit reached")

// SessionRepository interface to decouple the domain logic from the repository code.
type SessionRepository interface {
	Create(s *entity.Session) error
	GetByID(id string) (*entity.Session, error)
	GetAll(deadline time.Time) ([]entity.Session, error)
	Update(s *entity.Session) error
	Touch(id string) error
	PurgeOld(deadline time.Time) error
}

// SessionDomain abstracts the domain logic for netplay session handling.
type SessionDomain struct {
	sessionRepo      SessionRepository
	geopip2Domain    *GeoIP2Domain
	validationDomain *ValidationDomain
	mitmDomain       *MitmDomain
}

// NewSessionDomain returns an initalized SessionDomain struct.
func NewSessionDomain(
	sessionRepo SessionRepository,
	geoIP2Domain *GeoIP2Domain,
	validationDomain *ValidationDomain,
	mitmDomain *MitmDomain) *SessionDomain {
	return &SessionDomain{sessionRepo, geoIP2Domain, validationDomain, mitmDomain}
}

// Add adds or updates a session, based on the incomming request from the given IP.
// Returns ErrSessionRejected if session got rejected.
// Returns ErrRateLimited if rate limit for a session got reached.
func (d *SessionDomain) Add(request *AddSessionRequest, ip net.IP) (*entity.Session, error) {
	var err error
	var savedSession *entity.Session
	var requestType requestType = SessionCreate

	session := d.parseSession(request, ip)

	if session.IP == nil || session.Port == 0 {
		return nil, errors.New("IP or port not set")
	}

	// Decide if this is an CREATE, UPDATE or TOUCH operation
	session.CalculateID()
	if savedSession, err = d.sessionRepo.GetByID(session.ID); err != nil {
		return nil, fmt.Errorf("Can't get saved session: %w", err)
	}
	session.CalculateContentHash()
	if savedSession != nil {
		requestType = SessionTouch
		if savedSession.ContentHash != session.ContentHash {
			requestType = SessionUpdate
		}
	}

	// Ratelimit on UPDATE or TOUCH
	if requestType == SessionUpdate || requestType == SessionTouch {
		treshhold := time.Now().Add(-10 * time.Second)
		if savedSession.UpdatedAt.After(treshhold) {
			return nil, ErrRateLimited
		}
	}

	// Validate session on CREATE and UPDATE
	if requestType == SessionCreate || requestType == SessionUpdate {
		if !d.validateSession(session) {
			return nil, ErrSessionRejected
		}
	}

	// Open a game session on the selected MITM server if requested
	if (requestType == SessionCreate && session.HostMethod == entity.HostMethodMITM) ||
		requestType == SessionUpdate && session.HostMethod == entity.HostMethodMITM && savedSession.HostMethod != entity.HostMethodMITM {
		mitm, err := d.mitmDomain.OpenSession(request.MITMServer)
		if err != nil {
			return nil, fmt.Errorf("Can't open")
		}
		session.MitmAddress = mitm.Path
		session.MitmPort = mitm.Port
	}

	// Persist session changes
	switch requestType {
	case SessionCreate:
		if session.Country, err = d.geopip2Domain.GetCountryCodeForIP(session.IP); err != nil {
			return nil, fmt.Errorf("Can't find country for given IP %s: %w", session.IP, err)
		}

		if err = d.sessionRepo.Create(session); err != nil {
			return nil, fmt.Errorf("Can't create new session: %w", err)
		}
	case SessionUpdate:
		session.Country = savedSession.Country

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

// parseSession turns a request into a session information that can be compared to a persisted session
func (d *SessionDomain) parseSession(req *AddSessionRequest, ip net.IP) *entity.Session {
	var hostMethod entity.HostMethod = entity.HostMethodUnknown

	if req.ForceMITM {
		hostMethod = entity.HostMethodMITM
	}

	return &entity.Session{
		Username:            req.Username,
		GameName:            req.GameName,
		GameCRC:             strings.ToUpper(req.GameCRC),
		CoreName:            req.CoreName,
		CoreVersion:         req.CoreVersion,
		SubsystemName:       req.SubsystemName,
		RetroArchVersion:    req.RetroArchVersion,
		Frontend:            req.Frontend,
		IP:                  ip,
		Port:                req.Port,
		MitmAddress:         "",
		MitmPort:            0,
		HostMethod:          hostMethod,
		HasPassword:         req.HasPassword,
		HasSpectatePassword: req.HasSpectatePassword,
	}
}

// validateSession validaes an incomming session
func (d *SessionDomain) validateSession(s *entity.Session) bool {
	if len(s.Username) > 32 ||
		len(s.CoreName) > 255 ||
		len(s.GameName) > 255 ||
		len(s.GameCRC) != 8 ||
		len(s.RetroArchVersion) > 32 ||
		len(s.SubsystemName) > 255 ||
		len(s.Frontend) > 255 {
		return false
	}

	if !d.validationDomain.ValidateString(s.Username) ||
		!d.validationDomain.ValidateString(s.CoreName) ||
		!d.validationDomain.ValidateString(s.CoreVersion) ||
		!d.validationDomain.ValidateString(s.Frontend) ||
		!d.validationDomain.ValidateString(s.SubsystemName) ||
		!d.validationDomain.ValidateString(s.RetroArchVersion) {
		return false
	}

	return true
}

func (d *SessionDomain) getDeadline() time.Time {
	return time.Now().Add(-SessionDeadline * time.Second)
}
