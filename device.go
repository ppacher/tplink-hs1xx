package tpsmartapi

import (
	"context"
	"math"
	"time"

	tpshp "github.com/ppacher/tplink-smart-home-protocol"
)

// RPCError is implemented by all RPC responses
// and can be used to check for RPC errors
type RPCError interface {
	// Err return an error if the RPC failed, nil on success
	Err() error
}

// Device provides access to methods common for TP-Link smart home devices
type Device interface {
	// Call, do not use
	Call(ctx context.Context, req *tpshp.Request, err *error, fn func())

	// GetSysInfo returns the device's system information
	GetSysInfo(context.Context) <-chan *SysInfo

	// SetAlias sets the device's alias
	SetAlias(context.Context, string) <-chan RPCError

	// SetLocation sets the device's location (latitude/longitude)
	SetLocation(context.Context, float64, float64) <-chan RPCError

	// Reboot reboots the device
	Reboot(context.Context, time.Duration) <-chan RPCError

	// Reset resets the device to it's factory defaults
	Reset(context.Context, time.Duration) <-chan RPCError

	// GetWiFiScanInfo returns a list of wireless networks found by the device
	// If refresh is set to true, the device will start re-scanning for networks
	GetWiFiScanInfo(context.Context, bool, time.Duration) <-chan *WiFiScanResult
}

type Namespaces map[string]string

type device struct {
	ns     Namespaces
	client tpshp.Client
}

type emptyResponse struct {
	err error
}

func (e *emptyResponse) Err() error { return e.err }

// NewDevice creates a new device for the tpshp Client
func NewDevice(client tpshp.Client, ns Namespaces) Device {
	if ns["system"] == "" || ns["netif"] == "" {
		panic("no all namespaces configured")
	}

	return &device{
		ns:     ns,
		client: client,
	}
}

func Call(ctx context.Context, client tpshp.Client, req *tpshp.Request, err *error, fn func()) {
	go func() {
		*err = client.Call(ctx, req)
		fn()
	}()
}

func (dev *device) Call(ctx context.Context, req *tpshp.Request, err *error, fn func()) {
	Call(ctx, dev.client, req, err, fn)
}

func (dev *device) GetSysInfo(ctx context.Context) <-chan *SysInfo {
	var result SysInfo
	res := make(chan *SysInfo, 1)

	req := tpshp.NewRequest()
	req.AddCommand(dev.ns["system"], "get_sysinfo", struct{}{}, &result)

	dev.Call(ctx, req, &result.NetErr, func() {
		res <- &result
	})

	return res
}

func (dev *device) SetAlias(ctx context.Context, alias string) <-chan RPCError {
	res := make(chan RPCError, 1)

	req := tpshp.NewRequest().
		AddCommand(dev.ns["system"], "set_dev_alias", map[string]string{
			"alias": alias,
		}, nil)

	var result ErrorHandler
	dev.Call(ctx, req, &result.NetErr, func() {
		res <- &result
	})

	return res
}

func (dev *device) SetLocation(ctx context.Context, lat, lng float64) <-chan RPCError {
	res := make(chan RPCError, 1)

	latInt := math.Round(lat * 10000)
	lngInt := math.Round(lng * 10000)

	req := tpshp.NewRequest().AddCommand(dev.ns["system"], "set_dev_location", []interface{}{
		lat,
		lng,
		latInt,
		lngInt,
	}, nil)

	var result ErrorHandler
	dev.Call(ctx, req, &result.NetErr, func() {
		res <- &result
	})

	return res
}

func (dev *device) Reboot(ctx context.Context, d time.Duration) <-chan RPCError {
	res := make(chan RPCError, 1)

	req := tpshp.NewRequest().AddCommand(dev.ns["system"], "reboot", float64(d/time.Second), nil)
	var result ErrorHandler
	dev.Call(ctx, req, &result.NetErr, func() { res <- &result })

	return res
}

func (dev *device) Reset(ctx context.Context, d time.Duration) <-chan RPCError {
	res := make(chan RPCError, 1)

	req := tpshp.NewRequest().AddCommand(dev.ns["system"], "reset", float64(d/time.Second), nil)
	var result ErrorHandler
	dev.Call(ctx, req, &result.NetErr, func() { res <- &result })

	return res
}

func (dev *device) GetWiFiScanInfo(ctx context.Context, refresh bool, timeout time.Duration) <-chan *WiFiScanResult {
	var result WiFiScanResult
	res := make(chan *WiFiScanResult, 1)

	r := 0
	if refresh {
		r = 1
	}

	req := tpshp.NewRequest().AddCommand(dev.ns["netif"], "get_scaninfo", map[string]interface{}{
		"refresh": r,
		"timeout": float64(timeout / time.Second),
	}, &result)

	dev.Call(ctx, req, &result.NetErr, func() {
		res <- &result
	})

	return res
}
