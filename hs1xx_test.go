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

func Test_MeterInfo(t *testing.T) {
	plug := New("10.8.1.103")

	info, err := plug.MeterInfo(context.Background())
	assert.Nil(t, err)
	log.Println(info)
}
