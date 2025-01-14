package main

import (
	"iot-data-streamer/lora"
	"iot-data-streamer/nbiot"
	"iot-data-streamer/sigfox"
)

func main() {
	go lora.Run()
	go sigfox.Run()
	go nbiot.Run()

	for {
	}
}
