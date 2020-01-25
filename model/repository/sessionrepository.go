package repository

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// SessionRepository abstracts the database operation for Sessions.
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository returns a new SessionRepository.
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db}
}

// GetAll returns all sessions currently beeing hosted. Deadline is used to filter our old sessions. Deadline of zero value deactivates this filter.
func (r *SessionRepository) GetAll(deadline time.Time) ([]entity.Session, error) {
	var s []entity.Session
	if deadline.IsZero() {
		if err := r.db.Order("username").Find(&s).Error; err != nil {
			return nil, fmt.Errorf("can't query for all sessions: %w", err)
		}
	} else if err := r.db.Where("updated_at > ?", deadline).Order("username").Find(&s).Error; err != nil {
		return nil, fmt.Errorf("can't query for all sessions with deadline %s: %w", deadline, err)
	}

	return s, nil
}

// GetByID returns the session with the given ID. Returns nil if session can't be found.
func (r *SessionRepository) GetByID(id string) (*entity.Session, error) {
	var s entity.Session
	if err := r.db.Where("id = ?", id).First(&s).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't query session with ID %s: %w", id, err)
	}

	return &s, nil
}

// Create creates a new session.
func (r *SessionRepository) Create(s *entity.Session) error {
	if err := r.db.Create(s).Error; err != nil {
		return fmt.Errorf("can't create session %v: %w", s, err)
	}

	return nil
}

// Update updates a session.
func (r *SessionRepository) Update(s *entity.Session) error {
	if err := r.db.Model(&s).Save(&s).Error; err != nil {
		return fmt.Errorf("can't update session %v: %w", s, err)
	}

	return nil
}

// Touch updates the UpdatedAt timestamp.
func (r *SessionRepository) Touch(id string) error {
	if err := r.db.Model(&entity.Session{}).Where("id = ?", id).Update("updated_at", time.Now()).Error; err != nil {
		return fmt.Errorf("can't touch session with ID %s: %w", id, err)
	}

	return nil
}

// PurgeOld purges all sessions older than the given timestamp.
func (r *SessionRepository) PurgeOld(deadline time.Time) error {
	if err := r.db.Where("updated_at < ?", deadline).Delete(entity.Session{}).Error; err != nil {
		return fmt.Errorf("can't delete old sessions: %w", err)
	}

	return nil
}
