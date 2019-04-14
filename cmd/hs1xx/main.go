package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	hs1xx "github.com/ppacher/tplink-hs1xx"
)

func main() {
	flag.Parse()

	plug := hs1xx.New("10.8.1.103")

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
		log.Fatalf("Unknown command. Valid commands are: on, off")
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
