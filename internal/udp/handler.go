package udp

import (
	"net"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// handleUDPMessage processes incoming UDP messages.
func handleUDPMessage(conn *net.UDPConn, data []byte, addr *net.UDPAddr) {

	rawDataString := string(data)

	helpers.LogInfo("Received message from %s: %s\n", addr, rawDataString)

	response := []byte("Acknowledged")

	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		helpers.LogError(err, "Error sending response")
	}
}
