package hs1xx

import (
	"context"

	tpshp "github.com/ppacher/tplink-smart-home-protocol"
)

// HS1xx allows to interact with HS100, HS105 and HS110 Smart-WiFi-Plugs from
// TP-Link using it's Smart-Home Protocol
type HS1xx interface {
	// Client returns the tpshp client used
	Client() tpshp.Client

	// TurnOn activates the relay of a HS1xx smart plug
	TurnOn(context.Context) error

	// TurnOff de-activates the relay of a HS1xx smart plug
	TurnOff(context.Context) error

	// SysInfo queries the system information of the HS1xx device
	SysInfo(ctx context.Context) (*SysInfo, error)

	// SetRelayState sets the state of the HS1xx relay
	SetRelayState(context.Context, RelayState) error

	MeterInfo(ctx context.Context) (*SysInfo, *GetRealtimeResponse, error)
}

type hs1xx struct {
	client tpshp.Client
}

// New creates a new HS1xx client
func New(IP string) HS1xx {
	return &hs1xx{
		client: tpshp.New(IP),
	}
}

func (h *hs1xx) Client() tpshp.Client {
	return h.client
}

func (h *hs1xx) TurnOn(ctx context.Context) error {
	return h.SetRelayState(ctx, ON)
}

func (h *hs1xx) TurnOff(ctx context.Context) error {
	return h.SetRelayState(ctx, OFF)
}

func (h *hs1xx) SetRelayState(ctx context.Context, state RelayState) error {
	return h.client.Call(ctx, tpshp.NewRequest().AddCommand("system", "set_relay_state", map[string]RelayState{
		"state": state,
	}, nil))
}

func (h *hs1xx) SysInfo(ctx context.Context) (*SysInfo, error) {
	var res SysInfo

	if err := h.client.Call(ctx, tpshp.NewRequest().AddCommand("system", "get_sysinfo", struct{}{}, &res)); err != nil {
		return nil, err
	}

	return &res, nil
}

func (h *hs1xx) MeterInfo(ctx context.Context) (*SysInfo, *GetRealtimeResponse, error) {

	req := tpshp.NewRequest()

	var sysinfo SysInfo
	req.AddCommand("system", "get_sysinfo", struct{}{}, &sysinfo)

	var realtime GetRealtimeResponse
	req.AddCommand("emeter", "get_realtime", struct{}{}, &realtime)

	if err := h.client.Call(ctx, req); err != nil {
		return nil, nil, err
	}

	return &sysinfo, &realtime, nil
}
