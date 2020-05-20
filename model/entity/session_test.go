package entity

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testSession = Session{
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
	MitmHandle:          "",
	MitmAddress:         "0.0.0.0",
	MitmPort:            0,
	HostMethod:          HostMethodUPNP,
	HasPassword:         false,
	HasSpectatePassword: false,
	CreatedAt:           time.Now(),
	UpdatedAt:           time.Now(),
	ContentHash:         "",
}

func TestSessionIDCreationFoEmptySession(t *testing.T) {
	session := &Session{}
	session.CalculateID()

	assert.Equal(t, "7e8b1406d903bc9137fb69e769742c8d3e36f1c4fed51a608809b08de9f3e4a0", session.ID)
}

func TestSessionIDCreationForTestSession(t *testing.T) {
	session := testSession
	session.CalculateID()

	assert.NotEqual(t, "7e8b1406d903bc9137fb69e769742c8d3e36f1c4fed51a608809b08de9f3e4a0", session.ID)
	assert.Equal(t, "b78a35ff8be6cc104cce6ef1c3ab631621456a475d11d4df9612274285a48843", session.ID)
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

func TestSessionContentHashCreationFoEmptySession(t *testing.T) {
	session := &Session{}
	session.CalculateContentHash()

	assert.Equal(t, "d89f176c5afab7c6604184c30dbb51b9791940f9a0e9bfd21e0c9f86520fd958", session.ContentHash)
}

func TestSessionContentHashCreationForTestSession(t *testing.T) {
	session := testSession
	session.CalculateContentHash()

	assert.NotEqual(t, "d89f176c5afab7c6604184c30dbb51b9791940f9a0e9bfd21e0c9f86520fd958", session.ContentHash)
	assert.Equal(t, "2163d1d6642d8a0b4c2500ed0cf6d64c288b6649c9770210343890d7a6baef38", session.ContentHash)
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

func TestSessionPrintForRetroarch(t *testing.T) {
	session := testSession

	s := session.PrintForRetroarch()

	lines := strings.Split(s, "\n")

	assert.Equal(t, 18, len(lines))

	// Interate through all lines (except last empty line)
	// to see if it has the format of "%s=%s"
	for _, line := range lines[:len(lines)-1] {
		entries := strings.Split(line, "=")
		assert.Equal(t, 2, len(entries))
	}
}
