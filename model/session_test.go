package model

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
	MitmIP:              nil,
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

	assert.Equal(t, "7e8b1406d903bc9137fb69e769742c8d3e36f1c4fed51a608809b08de9f3e4a0490648545b3b196b0f224c2ad37c66f68fad191d7d2d783e341c3dc75eb6f421", session.ID)
}

func TestIDCreationForTestSession(t *testing.T) {
	session := testSession
	session.CalculateID()

	assert.NotEqual(t, "7e8b1406d903bc9137fb69e769742c8d3e36f1c4fed51a608809b08de9f3e4a0490648545b3b196b0f224c2ad37c66f68fad191d7d2d783e341c3dc75eb6f421", session.ID)
	assert.Equal(t, "1e8c4230c29e84239fa024d437bcd0d6e6915d35084294c24c99a47486f698310531e44fcb62204cc2c5a11e345b45265590ef01145eec5cd8d0b32f486d5786", session.ID)
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

	assert.Equal(t, "d89f176c5afab7c6604184c30dbb51b9791940f9a0e9bfd21e0c9f86520fd958794aa335c452a8dce6d89e07a4c1cfbebe041d61cfc6962a0d54f79e8b65c6d5", session.ContentHash)
}

func TestContentHashCreationForTestSession(t *testing.T) {
	session := testSession
	session.CalculateContentHash()

	assert.NotEqual(t, "d89f176c5afab7c6604184c30dbb51b9791940f9a0e9bfd21e0c9f86520fd958794aa335c452a8dce6d89e07a4c1cfbebe041d61cfc6962a0d54f79e8b65c6d5", session.ContentHash)
	assert.Equal(t, "5a429c6f4fa0d1c7a253d896eec18c868eb00713b4d50f9b247da01a900a13eaa875237df3b5f43b7da56706d9865b4437d17cd3a55be0eff2b377b645a93fda", session.ContentHash)
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
