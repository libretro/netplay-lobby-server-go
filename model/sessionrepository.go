package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// SessionRepository abstracts the database operation for Sessions.
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository returns a new SessionRepository.
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db}
}

// GetAllValid returns all sessions currently beeing hosted that are not older than the provided deadline.
func (r *SessionRepository) GetAllValid(deadline time.Time) ([]Session, error) {
	var s []Session
	if err := r.db.Where("updated_at > ?", deadline).Order("username").Find(&s).Error; err != nil {
		return nil, fmt.Errorf("can't query for all sessions: %w", err)
	}
	return s, nil
}

// GetByID returns the session with the given ID. Returns nil if session can't be found.
func (r *SessionRepository) GetByID(id string) (*Session, error) {
	var s Session
	if err := r.db.Where("id = ?", id).First(&s).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't query session with ID %s: %w", id, err)
	}
	return &s, nil
}

// Create creates a new session.
func (r *SessionRepository) Create(s *Session) error {
	if err := r.db.Create(s).Error; err != nil {
		return fmt.Errorf("can't create session %v: %w", s, err)
	}
	return nil
}

// Update updates a session.
func (r *SessionRepository) Update(s *Session) error {
	if err := r.db.Model(&s).Updates(&s).Error; err != nil {
		return fmt.Errorf("can't update session %v: %w", s, err)
	}
	return nil
}

// Touch updates the UpdatedAt timestamp.
func (r *SessionRepository) Touch(id string) error {
	if err := r.db.Model(&Session{}).Update("id", id).Error; err != nil {
		return fmt.Errorf("can't touch session with ID %s: %w", id, err)
	}
	return nil
}

// PurgeOld purges all sessions older than the given timestamp.
func (r *SessionRepository) PurgeOld(deadline time.Time) error {
	if err := r.db.Where("updated_at < ?", deadline).Delete(Session{}).Error; err != nil {
		return fmt.Errorf("can't delete old sessions: %w", err)
	}
	return nil
}
