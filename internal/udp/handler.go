package udp

import (
	"fmt"
	"net"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// handleUDPMessage processes incoming UDP messages.
func (s *UDPServer) handleUDPMessage(conn *net.UDPConn, data []byte, addr *net.UDPAddr) {

	rawDataString := string(data)

	// helpers.LogInfo("Received message from %s: %s", addr, rawDataString)

	rawDataArray := helpers.SplitIntoPairs(rawDataString)

	firmwareVersionHex, rawDataArray := helpers.Splice(rawDataArray, 0, 1, []string{})
	deviceIDHex, rawDataArray := helpers.Splice(rawDataArray, 0, 7, []string{})

	firmwareVersion, _ := helpers.HexSliceToBase10(firmwareVersionHex)
	deviceID, _ := helpers.HexSliceToBase10(deviceIDHex)

	fmt.Println(firmwareVersion, deviceID)
	//  TODO:  save to redis

	// timestampHex, rawDataArray := helpers.Splice(rawDataArray, 0, 4, []string{})
	// eventIDHex, rawDataArray := helpers.Splice(rawDataArray, 0, 1, []string{})

	// timestamp, _ := helpers.HexSliceToBase10(timestampHex)
	// eventID, _ := helpers.HexSliceToBase10(eventIDHex)

	response := []byte("Acknowledged\n")

	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		helpers.LogError(err, "Error sending response")
	}
}
