package hs1xx

import (
	"encoding/json"
)

// RelayState describes the state of a HS1xx plug
type RelayState bool

const (
	// ON indicates that the HS1xx plug is ON
	ON = RelayState(true)

	// OFF indicates that the HS1xx plug in OFF
	OFF = RelayState(false)
)

func (state RelayState) String() string {
	if bool(state) {
		return "ON"
	}
	return "OFF"
}

// MarshalJSON implements the json.Marshaler interface
func (state RelayState) MarshalJSON() ([]byte, error) {
	msg := SetStateRequest{}

	if bool(state) {
		msg.System.SetRelayState.State = 1
	}

	return json.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (state *RelayState) UnmarshalJSON(b []byte) error {
	var i int
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	if i == 1 {
		*state = true
	} else {
		*state = false
	}
	return nil
}

type GetDailyStatRequest struct {
	Emeter struct {
		GetDaystat struct {
			Month int `json:"month"`
			Year  int `json:"year"`
		} `json:"get_daystat"`
	} `json:"emeter"`
}

type GetRealtimeResponse struct {
	VoltageMilli float64 `json:"voltage_mv"`
	CurrentMilli float64 `json:"current_ma"`
	PowerMilli   float64 `json:"power_mw"`
	TotalWH      float64 `json:"total_wh"`
	ErrorCode    int     `json:"err_code"`
}

// Voltage returns the current voltage in V
func (res *GetRealtimeResponse) Voltage() float64 {
	return float64(res.VoltageMilli) / 1000
}

// Current returns the current in A
func (res *GetRealtimeResponse) Current() float64 {
	return float64(res.CurrentMilli) / 1000
}

// Power returns the current power consumption in W
func (res *GetRealtimeResponse) Power() float64 {
	return float64(res.PowerMilli) / 1000
}

// Total returns the total power consumption in Wh
func (res *GetRealtimeResponse) Total() float64 {
	return float64(res.TotalWH)
}

type GetSysInfoRequest struct {
	System struct {
		GetSysinfo struct{} `json:"get_sysinfo"`
	} `json:"system"`
}

// SysInfo describes the JSON object returned for the get_sysinfo command
type SysInfo struct {
	SWVer      string     `json:"sw_ver"`
	HWVer      string     `json:"hw_ver"`
	Type       string     `json:"type"`
	Model      string     `json:"model"`
	MAC        string     `json:"mac"`
	DeviceName string     `json:"dev_name"`
	Alias      string     `json:"alias"`
	RelayState RelayState `json:"relay_state"`
	OnTime     int        `json:"on_time"`
	ActiveMode string     `json:"active_mode"`
	Feature    string     `json:"feature"`
	Updating   int        `json:"updating"`
	IconHash   string     `json:"icon_hash"`
	Rssi       int        `json:"rssi"`
	LedOff     int        `json:"led_off"`
	Longitude  int        `json:"longitude_i"`
	Latitude   int        `json:"latitude_i"`
	HwID       string     `json:"hwId"`
	FwID       string     `json:"fwId"`
	DeviceID   string     `json:"deviceId"`
	OEMID      string     `json:"oemID"`
}

type SysInfoResponse struct {
	System struct {
		SysInfo SysInfo `json:"get_sysinfo"`
	} `json:"system"`
}

type GetMeterInfoRequest struct {
	Emeter struct {
		Realtime      GetRealtimeResponse `json:"get_realtime"`
		GetVgainIgain struct{}            `json:"get_vgain_igain"`
	} `json:"emeter"`
	System struct {
		GetSysinfo struct{} `json:"get_sysinfo"`
	} `json:"system"`
}

type MeterInfoResponse struct {
	SysInfoResponse
	Emeter struct {
		Realtime      GetRealtimeResponse `json:"get_realtime"`
		GetVgainIgain struct{}            `json:"get_vgain_igain"`
	} `json:"emeter"`
}

type SetStateRequest struct {
	System struct {
		SetRelayState struct {
			State int `json:"state"`
		} `json:"set_relay_state"`
	} `json:"system"`
}
