package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	hs1xx "github.com/ppacher/tplink-hs1xx"
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
		err = plug.TurnOn(ctx)
	case "off":
		err = plug.TurnOff(ctx)
	case "sysinfo":
		output, err = plug.SysInfo(ctx)
	case "meter":
		var sys interface{}
		var meter interface{}
		sys, meter, err = plug.MeterInfo(ctx)
		output = map[string]interface{}{
			"sysinfo": sys,
			"meter":   meter,
		}
	default:
		log.Fatalf("Unknown command. Valid commands are: on, off, sysinfo, meter")
	}

	if err != nil {
		log.Fatal(err)
	}

	if output != nil {
		blob, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(blob))
	}
}
