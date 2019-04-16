package hs1xx

import (
	"context"

	tpsmartapi "github.com/ppacher/tplink-hs1xx"
	tpshp "github.com/ppacher/tplink-smart-home-protocol"
)

// HS1xx allows to interact with HS100, HS105 and HS110 Smart-WiFi-Plugs from
// TP-Link using it's Smart-Home Protocol
type HS1xx interface {
	Device() tpsmartapi.Device

	// EMeter returns a client for the built-in energy meter of some HS1xx plugs
	EMeter() tpsmartapi.EMeter

	// Client returns the tpshp client used
	Client() tpshp.Client

	// TurnOn activates the relay of a HS1xx smart plug
	TurnOn(context.Context) <-chan tpsmartapi.RPCError

	// TurnOff de-activates the relay of a HS1xx smart plug
	TurnOff(context.Context) <-chan tpsmartapi.RPCError

	// SetRelayState sets the state of the HS1xx relay
	SetRelayState(context.Context, RelayState) <-chan tpsmartapi.RPCError

	// GetRelayState queries the relay state of the HS1xx plug
	GetRelayState(context.Context) (RelayState, error)

	// SysInfo returns the system information of the device
	SysInfo(ctx context.Context) <-chan *SysInfo

	// SetLedState enables or disables the LED on the HS1xx smart plug
	SetLedState(ctx context.Context, on bool) <-chan tpsmartapi.RPCError
}

var hs1xxNamespace = map[string]string{
	"system": "system",
	"netif":  "netif",
	"emeter": "emeter",
}

type hs1xx struct {
	client tpshp.Client
}

func (h *hs1xx) SetLedState(ctx context.Context, on bool) <-chan tpsmartapi.RPCError {
	res := make(chan tpsmartapi.RPCError, 1)
	var result tpsmartapi.ErrorHandler

	v := 0
	if !on {
		v = 1
	}

	req := tpshp.NewRequest().
		AddCommand("system", "set_led_off", map[string]interface{}{
			"off": v,
		}, nil)

	h.Device().Call(ctx, req, &result.NetErr, func() { res <- &result })

	return res
}

// New creates a new HS1xx client
func New(IP string) HS1xx {
	return &hs1xx{
		client: tpshp.New(IP),
	}
}

func (h *hs1xx) Device() tpsmartapi.Device {
	return tpsmartapi.NewDevice(h.client, hs1xxNamespace)
}

func (h *hs1xx) EMeter() tpsmartapi.EMeter {
	return tpsmartapi.NewEmeter(h.client, hs1xxNamespace)
}

func (h *hs1xx) Client() tpshp.Client {
	return h.client
}

func (h *hs1xx) TurnOn(ctx context.Context) <-chan tpsmartapi.RPCError {
	return h.SetRelayState(ctx, ON)
}

func (h *hs1xx) TurnOff(ctx context.Context) <-chan tpsmartapi.RPCError {
	return h.SetRelayState(ctx, OFF)
}

func (h *hs1xx) SetRelayState(ctx context.Context, state RelayState) <-chan tpsmartapi.RPCError {
	res := make(chan tpsmartapi.RPCError, 1)
	var result tpsmartapi.ErrorHandler

	req := tpshp.NewRequest().AddCommand("system", "set_relay_state", map[string]RelayState{
		"state": state,
	}, &result)

	tpsmartapi.Call(ctx, h.client, req, &result.NetErr, func() { res <- &result })

	return res
}

func (h *hs1xx) GetRelayState(ctx context.Context) (RelayState, error) {
	s := <-h.SysInfo(ctx)

	return s.RelayState, s.Err()
}

func (h *hs1xx) SysInfo(ctx context.Context) <-chan *SysInfo {
	res := make(chan *SysInfo, 1)
	var result SysInfo
	req := tpshp.NewRequest().AddCommand("system", "get_sysinfo", struct{}{}, &result)

	tpsmartapi.Call(ctx, h.client, req, &result.NetErr, func() { res <- &result })
	return res
}
