package controller


import (
	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// SessionController handles all session related request
type SessionController struct {}

// AddSessionRequest defines the request for the AddSession request.
type AddSessionRequest struct {
	Username string `form:"username"`
	CoreName string `form:"core_name"`
	CoreVersion string `form:"core_version"`
	GameName string `form:"game_name"`
	GameCRC  string `form:"game_crc"`
	Port uint `form:"port"`
	MITMServer string `form:"mitm_server"`
	HasPassword bool  `form:"has_password"` // 1/0
	HasSpectatePassword bool `form:"has_spectate_password"`
	ForceMITM bool `form:"force_mitm"`
	RetroarchVersion string `form:"retroarch_version"`
	FrontendArchitecture string `form:"frontend_architecture"`
	SubsystemName string  `form:"subsystem_name"`
}

// ListSessionsResponse is a custom DTO for backward compatability.
type ListSessionsResponse struct {
	Fields entity.Session `json:"fields"`
}
