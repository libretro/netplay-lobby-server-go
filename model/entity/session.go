package entity

import (
	"encoding/hex"
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

// HostMethod is the enum for the hosting method.
type HostMethod int64

// The enum values for HostMethod.
const (
	HostMethodUnknown = 0
	HostMethodManual  = 1
	HostMethodUPNP    = 2
	HostMethodMITM    = 3
)

// TODO remove "fixed" field from retroarch frontend code

// Session is the database presentation of a netplay session.
type Session struct {
	ID                  string     `json:"-" gorm:"primary_key;size:64"`
	ContentHash         string     `json:"-" gorm:"size:64"`
	Username            string     `json:"username"`
	Country             string     `json:"country" gorm:"size:2"`
	GameName            string     `json:"game_name"`
	GameCRC             string     `json:"game_crc"`
	CoreName            string     `json:"core_name"`
	CoreVersion         string     `json:"core_version"`
	SubsystemName       string     `json:"subsystem_name"`
	RetroArchVersion    string     `json:"retroarch_version"`
	Frontend            string     `json:"frontend"`
	IP                  net.IP     `json:"ip" gorm:"not null"`
	Port                uint16     `json:"port"`
	MitmIP              net.IP     `json:"mitm_ip"`
	MitmPort            uint16     `json:"mitm_port"`
	HostMethod          HostMethod `json:"host_method"`
	HasPassword         bool       `json:"has_password"`
	HasSpectatePassword bool       `json:"has_spectate_password"`
	CreatedAt           time.Time  `json:"created"`
	UpdatedAt           time.Time  `json:"updated" gorm:"index"`
}

// CalculateID creates a 32 byte SHAKE256 (SHA3) hash of the session for the db to use as PK.
func (s *Session) CalculateID() {
	hash := make([]byte, 32)
	shake := sha3.NewShake256()

	shake.Write([]byte(s.Username))
	shake.Write([]byte(s.IP))
	shake.Write([]byte(strconv.FormatUint(uint64(s.Port), 10)))

	shake.Read(hash)

	s.ID = hex.EncodeToString(hash)
}

// CalculateContentHash creates a 32 byte SHAKE256 (SHA3) hash of the session content.
func (s *Session) CalculateContentHash() {
	hash := make([]byte, 32)
	shake := sha3.NewShake256()

	shake.Write([]byte(s.Username))
	shake.Write([]byte(s.GameName))
	shake.Write([]byte(s.GameCRC))
	shake.Write([]byte(s.CoreName))
	shake.Write([]byte(s.CoreVersion))
	shake.Write([]byte(s.SubsystemName))
	shake.Write([]byte(s.RetroArchVersion))
	shake.Write([]byte(s.Frontend))
	shake.Write([]byte(s.IP))
	shake.Write([]byte(strconv.FormatUint(uint64(s.Port), 10)))
	shake.Write([]byte(s.MitmIP))
	shake.Write([]byte(strconv.FormatUint(uint64(s.MitmPort), 10)))
	shake.Write([]byte(strconv.FormatBool(s.HasPassword)))
	shake.Write([]byte(strconv.FormatBool(s.HasSpectatePassword)))

	shake.Read(hash)

	s.ContentHash = hex.EncodeToString(hash)
}
