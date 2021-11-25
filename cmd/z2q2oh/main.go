package main

import (
	"flag"
	"zigbee2mqtt2openhab/service"
)

func main() {
	flag.Parse()
	svc, err := service.New()
	if err != nil {
		panic(err)
	}
	if err := svc.Run(); err != nil {
		panic(err)
	}
}
