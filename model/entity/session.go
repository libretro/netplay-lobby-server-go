package entity

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
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

// Session is the database presentation of a netplay session.
type Session struct {
	ID                  string     `json:"-" gorm:"primary_key;size:64"`
	ContentHash         string     `json:"-" gorm:"size:64"`
	RoomID              int32      `json:"id" gorm:"AUTO_INCREMENT;unique_index"`
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
	MitmHandle          string     `json:"-"`
	MitmAddress         string     `json:"mitm_ip"`
	MitmPort            uint16     `json:"mitm_port"`
	MitmSession         string     `json:"mitm_session"`
	HostMethod          HostMethod `json:"host_method"`
	HasPassword         bool       `json:"has_password"`
	HasSpectatePassword bool       `json:"has_spectate_password"`
	Connectable         bool       `json:"connectable"`
	IsRetroArch         bool       `json:"is_retroarch"`
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
	shake.Write([]byte(strconv.FormatUint(uint64(s.HostMethod), 10)))
	shake.Write([]byte(s.MitmHandle))
	shake.Write([]byte(s.MitmSession))
	shake.Write([]byte(strconv.FormatBool(s.HasPassword)))
	shake.Write([]byte(strconv.FormatBool(s.HasSpectatePassword)))

	shake.Read(hash)

	s.ContentHash = hex.EncodeToString(hash)
}

// PrintForRetroarch prints out the session information in a format that retroarch is expecting.
func (s *Session) PrintForRetroarch() string {
	var str string
	var hasPassword = 0
	var hasSpectatePassword = 0
	var connectable = 0

	if s.HasPassword {
		hasPassword = 1
	}

	if s.HasSpectatePassword {
		hasSpectatePassword = 1
	}

	if s.Connectable {
		connectable = 1
	}

	str += fmt.Sprintf("id=%d\nusername=%s\ncore_name=%s\ngame_name=%s\ngame_crc=%s\ncore_version=%s\nip=%s\nport=%d\nhost_method=%d\nhas_password=%d\nhas_spectate_password=%d\nretroarch_version=%s\nfrontend=%s\nsubsystem_name=%s\ncountry=%s\nconnectable=%d\n",
		s.RoomID,
		s.Username,
		s.CoreName,
		s.GameName,
		strings.ToUpper(s.GameCRC),
		s.CoreVersion,
		s.IP,
		s.Port,
		s.HostMethod,
		hasPassword,
		hasSpectatePassword,
		s.RetroArchVersion,
		s.Frontend,
		s.SubsystemName,
		strings.ToUpper(s.Country),
		connectable,
	)

	return str
}
