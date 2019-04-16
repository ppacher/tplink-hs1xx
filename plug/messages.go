package hs1xx

import (
	"encoding/json"

	tpsmartapi "github.com/ppacher/tplink-hs1xx"
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
	v := 0
	if bool(state) {
		v = 1
	}
	return json.Marshal(v)
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

// SysInfo describes the JSON object returned for the get_sysinfo command
type SysInfo struct {
	tpsmartapi.SysInfo

	RelayState RelayState `json:"relay_state"`
}
