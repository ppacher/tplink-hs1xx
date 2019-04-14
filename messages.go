package hs1xx

import "encoding/json"

// RelayState describes the state of a HS1xx plug
type RelayState bool

const (
	ON  = RelayState(true)
	OFF = RelayState(false)
)

func (state RelayState) MarshalJSON() ([]byte, error) {
	msg := SetStateRequest{}

	if bool(state) {
		msg.System.SetRelayState.State = 1
	}

	return json.Marshal(msg)
}

type GetDailyStatRequest struct {
	Emeter struct {
		GetDaystat struct {
			Month int `json:"month"`
			Year  int `json:"year"`
		} `json:"get_daystat"`
	} `json:"emeter"`
}

type GetSysInfoRequest struct {
	System struct {
		GetSysinfo struct{} `json:"get_sysinfo"`
	} `json:"system"`
}

type SysInfo struct {
	SWVer      string `json:"sw_ver"`
	HWVer      string `json:"hw_ver"`
	Type       string `json:"type"`
	Model      string `json:"model"`
	MAC        string `json:"mac"`
	DeviceName string `json:"dev_name"`
	Alias      string `json:"alias"`
	RelayState int    `json:"relay_state"`
	OnTime     int    `json:"on_time"`
	ActiveMode string `json:"active_mode"`
	Feature    string `json:"feature"`
	Updating   int    `json:"updating"`
	IconHash   string `json:"icon_hash"`
	Rssi       int    `json:"rssi"`
	LedOff     int    `json:"led_off"`
	Longitude  int    `json:"longitude_i"`
	Latitude   int    `json:"latitude_i"`
	HwID       string `json:"hwId"`
	FwID       string `json:"fwId"`
	DeviceID   string `json:"deviceId"`
	OEMID      string `json:"oemID"`
}

type SysInfoResponse struct {
	System struct {
		SysInfo SysInfo `json:"get_sysinfo"`
	} `json:"system"`
}

type GetMeterInfoRequest struct {
	Emeter struct {
		GetRealtime   struct{} `json:"get_realtime"`
		GetVgainIgain struct{} `json:"get_vgain_igain"`
	} `json:"emeter"`
	System struct {
		GetSysinfo struct{} `json:"get_sysinfo"`
	} `json:"system"`
}

type MeterInfoResponse struct {
	SysInfoResponse
}

type SetStateRequest struct {
	System struct {
		SetRelayState struct {
			State int `json:"state"`
		} `json:"set_relay_state"`
	} `json:"system"`
}
