package entity

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testSession = Session{
	ID:                  "",
	RoomID:              0,
	Username:            "zelda",
	Country:             "en",
	GameName:            "supergame",
	GameCRC:             "FFFFFFFF",
	CoreName:            "unes",
	CoreVersion:         "0.2.1",
	SubsystemName:       "subsub",
	RetroArchVersion:    "1.1.1",
	Frontend:            "retro",
	IP:                  net.ParseIP("127.0.0.1"),
	Port:                55355,
	MitmHandle:          "",
	MitmAddress:         "",
	MitmPort:            0,
	MitmSession:         "",
	HostMethod:          HostMethodUPNP,
	HasPassword:         false,
	HasSpectatePassword: false,
	Connectable:         true,
	IsRetroArch:         true,
	CreatedAt:           time.Now(),
	UpdatedAt:           time.Now(),
	ContentHash:         "",
}

func TestSessionIDDoesNotChangeID(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ID

	session.ID = "CHANGED ID"
	session.CalculateID()
	newHash := session.ID

	assert.Equal(t, oldHash, newHash)
}

func TestSessionCreatedAtDoesNotChangeID(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ID

	session.CreatedAt = time.Now()
	session.CalculateID()
	newHash := session.ID

	assert.Equal(t, oldHash, newHash)
}

func TestSessionUpdatedAtDoesNotChangeID(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ID

	session.UpdatedAt = time.Now()
	session.CalculateID()
	newHash := session.ID

	assert.Equal(t, oldHash, newHash)
}

func TestSessionIDDoesNotChangeContentHash(t *testing.T) {
	session := testSession
	session.CalculateContentHash()
	oldHash := session.ContentHash

	session.ID = "CHANGED ID"
	session.CalculateContentHash()
	newHash := session.ContentHash

	assert.Equal(t, oldHash, newHash)
}

func TestSessionCreatedAtDoesNotChangeContentHash(t *testing.T) {
	session := testSession
	session.CalculateContentHash()
	oldHash := session.ContentHash

	session.CreatedAt = time.Now()
	session.CalculateContentHash()
	newHash := session.ContentHash

	assert.Equal(t, oldHash, newHash)
}

func TestSessionUpdatedAtDoesNotChangeContentHash(t *testing.T) {
	session := testSession
	session.CalculateContentHash()
	oldHash := session.ContentHash

	session.UpdatedAt = time.Now()
	session.CalculateContentHash()
	newHash := session.ContentHash

	assert.Equal(t, oldHash, newHash)
}

func TestSessionHostMethodChangesContentHash(t *testing.T) {
	session := testSession
	session.CalculateContentHash()
	oldHash := session.ContentHash

	session.HostMethod = HostMethodManual
	session.CalculateContentHash()
	newHash := session.ContentHash

	assert.NotEqual(t, oldHash, newHash)
}
