package tpsmartapi

import (
	"fmt"
	"strings"
	"time"
)

type ErrorHandler struct {
	ErrorCode int   `json:"err_code"`
	NetErr    error `json:"-"`
}

func (e ErrorHandler) Err() error {
	if e.ErrorCode != 0 {
		return fmt.Errorf("error-code: %d", e.ErrorCode)
	}

	return e.NetErr
}

// SysInfo describes the JSON object returned for the get_sysinfo command
type SysInfo struct {
	ErrorHandler

	SWVer      string `json:"sw_ver"`
	HWVer      string `json:"hw_ver"`
	Type       string `json:"type"`
	Model      string `json:"model"`
	MAC        string `json:"mac"`
	DeviceName string `json:"dev_name"`
	Alias      string `json:"alias"`
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

// EMeterSupported returns true if the device supports an engery meter
func (sys SysInfo) EMeterSupported() bool {
	return strings.Contains(sys.Feature, "EME")
}

// TimerSupported returns true if the device supports a timer
func (sys SysInfo) TimerSupported() bool {
	return strings.Contains(sys.Feature, "TIM")
}

// Features returns a list of feature supported by the device
func (sys SysInfo) Features() []string {
	return strings.Split(sys.Feature, ":")
}

type VgainIgain struct {
	Igain     int `json:"igain"`
	Vgain     int `json:"vgain"`
	ErrorCode int `json:"err_code"`
}

type DayStat struct {
	Day      int `json:"day"`
	EnergyWh int `json:"energy_wh"`
	Month    int `json:"month"`
	Year     int `json:"year"`
}

func (d DayStat) GetMonth() time.Month {
	return time.Month(d.Month)
}

type DailyStats struct {
	ErrorHandler

	Days []DayStat `json:"day_list"`
}

type MonthStat struct {
	EnergyWh int `json:"energy_wh"`
	Month    int `json:"month"`
	Year     int `json:"year"`
}

func (d MonthStat) GetMonth() time.Month {
	return time.Month(d.Month)
}

type MonthlyStats struct {
	ErrorHandler

	Months []MonthStat `json:"month_list"`
}

type WirelessNetwork struct {
	KeyType int    `json:"key_type"`
	SSID    string `json:"ssid"`
}

type WiFiScanResult struct {
	ErrorHandler

	APs []WirelessNetwork `json:"ap_list"`
}

type RealtimeInfo struct {
	ErrorHandler

	VoltageMilli float64 `json:"voltage_mv"`
	CurrentMilli float64 `json:"current_ma"`
	PowerMilli   float64 `json:"power_mw"`
	TotalWH      float64 `json:"total_wh"`
}

// Voltage returns the current voltage in V
func (res *RealtimeInfo) Voltage() float64 {
	return float64(res.VoltageMilli) / 1000
}

// Current returns the current in A
func (res *RealtimeInfo) Current() float64 {
	return float64(res.CurrentMilli) / 1000
}

// Power returns the current power consumption in W
func (res *RealtimeInfo) Power() float64 {
	return float64(res.PowerMilli) / 1000
}

// Total returns the total power consumption in Wh
func (res *RealtimeInfo) Total() float64 {
	return float64(res.TotalWH)
}
