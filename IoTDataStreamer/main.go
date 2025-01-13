package main

import (
	"iot-data-streamer/lora"
	"iot-data-streamer/sigfox"
)

func main() {
	go lora.Run()
	go sigfox.Run()

	for {

	}
}
