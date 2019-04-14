package hs1xx

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TurnOn(t *testing.T) {
	plug := New("10.8.1.103")

	assert.Nil(t, plug.TurnOn(context.Background()))
}

func Test_TurnOff(t *testing.T) {
	plug := New("10.8.1.103")

	assert.Nil(t, plug.TurnOff(context.Background()))
}

func Test_Output(t *testing.T) {
	plug := New("10.8.1.103")

	var p []byte
	err := plug.SendCommand(context.Background(), GetMeterInfoRequest{}, &p)
	assert.Nil(t, err)
	log.Println(string(p))
}

func Test_MeterInfo(t *testing.T) {
	plug := New("10.8.1.103")

	info, err := plug.MeterInfo(context.Background())
	assert.Nil(t, err)
	log.Println(info)
}

func Test_SysInfo(t *testing.T) {
	plug := New("10.8.1.103")

	info, err := plug.SysInfo(context.Background())
	assert.Nil(t, err)
	log.Println(info)
}
