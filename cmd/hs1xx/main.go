package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	tpsmartapi "github.com/ppacher/tplink-hs1xx"

	hs1xx "github.com/ppacher/tplink-hs1xx/plug"
)

var deviceIP = flag.String("d", "", "IP of a TP-Link HS1xx smart plug")

func main() {
	flag.Parse()

	if *deviceIP == "" {
		log.Fatal("Missing IP address")
	}

	plug := hs1xx.New(*deviceIP)

	cmd := flag.Arg(0)

	var err error
	var output interface{}

	ctx := context.Background()

	switch cmd {
	case "on":
		output = <-plug.TurnOn(ctx)
	case "off":
		output = <-plug.TurnOff(ctx)
	case "sysinfo":
		output = <-plug.SysInfo(ctx)
	case "relay":
		output, err = plug.GetRelayState(ctx)
	case "wifi":
		output = <-plug.Device().GetWiFiScanInfo(ctx, true, time.Second*20)
	case "meter":
		output = <-plug.EMeter().GetRealtime(ctx)
	default:
		log.Fatalf("Unknown command. Valid commands are: on, off, sysinfo, meter")
	}

	if err != nil {
		log.Fatal(err)
	}

	if err, ok := output.(tpsmartapi.RPCError); ok && err.Err() != nil {
		log.Fatal(err.Err())
	}

	if output != nil {
		blob, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(blob))
	}
}
