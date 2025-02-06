package main

import (
	// "iot-data-streamer/lora"
	// "iot-data-streamer/sigfox"
	"iot-data-streamer/nbiot"
)

func main() {
	// go lora.Run()
	// go sigfox.Run()
	go nbiot.Run()

	for {
	}
}
