package entity

import (
	"time"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSession = Session{
	ID:                  "",
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
	MitmIP:              net.ParseIP("0.0.0.0"),
	MitmPort:            0,
	HostMethod:          HostMethodUPNP,
	HasPassword:         false,
	HasSpectatePassword: false,
	CreatedAt:           time.Now(),
	UpdatedAt:           time.Now(),
	ContentHash:         "",
}

func TestIDCreationFoEmptySession(t *testing.T) {
	session := &Session{}
	session.CalculateID()

	assert.Equal(t, "7e8b1406d903bc9137fb69e769742c8d3e36f1c4fed51a608809b08de9f3e4a0", session.ID)
}

func TestIDCreationForTestSession(t *testing.T) {
	session := testSession
	session.CalculateID()

	assert.NotEqual(t, "7e8b1406d903bc9137fb69e769742c8d3e36f1c4fed51a608809b08de9f3e4a0", session.ID)
	assert.Equal(t, "b78a35ff8be6cc104cce6ef1c3ab631621456a475d11d4df9612274285a48843", session.ID)
}

func TestIDDoesNotChangeID(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ID

	session.ID = "CHANGED ID"
	session.CalculateID()
	newHash := session.ID

	assert.Equal(t, oldHash, newHash)
}

func TestCreatedAtDoesNotChangeID(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ID

	session.CreatedAt = time.Now()
	session.CalculateID()
	newHash := session.ID

	assert.Equal(t, oldHash, newHash)
}

func TestUpdatedAtDoesNotChangeID(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ID

	session.UpdatedAt = time.Now()
	session.CalculateID()
	newHash := session.ID

	assert.Equal(t, oldHash, newHash)
}

func TestContentHashCreationFoEmptySession(t *testing.T) {
	session := &Session{}
	session.CalculateContentHash()

	assert.Equal(t, "d89f176c5afab7c6604184c30dbb51b9791940f9a0e9bfd21e0c9f86520fd958", session.ContentHash)
}

func TestContentHashCreationForTestSession(t *testing.T) {
	session := testSession
	session.CalculateContentHash()

	assert.NotEqual(t, "d89f176c5afab7c6604184c30dbb51b9791940f9a0e9bfd21e0c9f86520fd958", session.ContentHash)
	assert.Equal(t, "179ca1a118b42f910f782a001c39824338bc7506bb11c6a956057921e1f84605", session.ContentHash)
}

func TestIDDoesNotChangeContentHash(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ContentHash

	session.ID = "CHANGED ID"
	session.CalculateID()
	newHash := session.ContentHash

	assert.Equal(t, oldHash, newHash)
}

func TestCreatedAtDoesNotChangeContentHash(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ContentHash

	session.CreatedAt = time.Now()
	session.CalculateID()
	newHash := session.ContentHash

	assert.Equal(t, oldHash, newHash)
}

func TestUpdatedAtDoesNotChangeContentHash(t *testing.T) {
	session := testSession
	session.CalculateID()
	oldHash := session.ContentHash

	session.UpdatedAt = time.Now()
	session.CalculateID()
	newHash := session.ContentHash

	assert.Equal(t, oldHash, newHash)
}
