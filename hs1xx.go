// Package hs1xx provides methods to interact with TPLink HS100 and HS110 smart WiFi plugs.
// It is based on the reverse engineering done in
// https://www.softscheck.com/en/reverse-engineering-tp-link-hs110/ and
// https://georgovassilis.blogspot.sg/2016/05/controlling-tp-link-hs100-wi-fi-smart.html
// as well as the hs1xx library https://github.com/sausheong/hs1xxplug
package hs1xx

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"net"
)

// Plug allows to control and query a TP-Link HS100 or HS110 smart-wifi plug
type Plug interface {
	// TurnOn sets the plugs relay state to on
	TurnOn(context.Context) error

	// TurnOff sets the plugs relay state to off
	TurnOff(context.Context) error

	// SetRelayState sets the state of the hs1xx plug relay
	SetRelayState(context.Context, RelayState) error

	// MeterInfo returns data about the engery meter included in HS110 plugs
	MeterInfo(context.Context) (*MeterInfoResponse, error)

	// SendCommand sends an abritrary command to the hs1xx plug
	SendCommand(ctx context.Context, request interface{}, response interface{}) error
}

// New creates a new HS1xx plug
func New(ip string) Plug {
	return &plug{ip}
}

// plug implements the Plug interface
type plug struct {
	// IPAddress holds the IP address assigned to the smart plug
	ip string
}

// dial connects the the plug on port :9999
func (p *plug) dial(ctx context.Context) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "tcp", p.ip+":9999")
}

func (p *plug) MeterInfo(ctx context.Context) (*MeterInfoResponse, error) {
	var payload MeterInfoResponse

	if err := p.SendCommand(ctx, GetMeterInfoRequest{}, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (p *plug) SetRelayState(ctx context.Context, state RelayState) error {
	return p.SendCommand(ctx, state, nil)
}

func (p *plug) TurnOn(ctx context.Context) error {
	return p.SetRelayState(ctx, ON)
}

func (p *plug) TurnOff(ctx context.Context) error {
	return p.SetRelayState(ctx, OFF)
}

func (p *plug) SendCommand(ctx context.Context, request interface{}, response interface{}) error {
	conn, err := p.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	blob, err := json.Marshal(request)
	if err != nil {
		return err
	}
	cipher := Encrypt(blob)

	n, err := conn.Write(cipher)
	if n != len(cipher) {
		return err
	}

	if response != nil {
		var size uint32
		if err := binary.Read(conn, binary.BigEndian, &size); err != nil {
			return err
		}

		var payload = make([]byte, size)

		n, err := conn.Read(payload)
		if n == int(size) {
			plain := Decrypt(payload)
			if v, ok := response.(*[]byte); ok {
				*v = []byte(plain)
				return nil
			}
			return json.Unmarshal([]byte(plain), response)
		}

		return err
	}

	return nil
}

// Encrypt a plaintext byte slice using the encryption algorithm used by HS1xx plugs
func Encrypt(plaintext []byte) []byte {
	n := len(plaintext)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(n))
	ciphertext := []byte(buf.Bytes())

	key := byte(0xAB)
	payload := make([]byte, n)
	for i := 0; i < n; i++ {
		payload[i] = plaintext[i] ^ key
		key = payload[i]
	}

	for i := 0; i < len(payload); i++ {
		ciphertext = append(ciphertext, payload[i])
	}

	return ciphertext
}

// Decrypt a ciphertext byte slice using the decryption algorithm used by HS1xx plugs
func Decrypt(ciphertext []byte) string {
	n := len(ciphertext)
	key := byte(0xAB)
	var nextKey byte
	for i := 0; i < n; i++ {
		nextKey = ciphertext[i]
		ciphertext[i] = ciphertext[i] ^ key
		key = nextKey
	}
	return string(ciphertext)
}

// compile time check
var _ Plug = &plug{}
