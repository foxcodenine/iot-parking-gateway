package udp

import (
	"fmt"
	"net"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// handleUDPMessage processes incoming UDP messages.
func handleUDPMessage(conn *net.UDPConn, data []byte, addr *net.UDPAddr) {

	rawDataString := string(data)

	// helpers.LogInfo("Received message from %s: %s", addr, rawDataString)

	rawDataArray := helpers.SplitIntoPairs(rawDataString)

	firmwareVersionHex, rawDataArray := helpers.Splice(rawDataArray, 0, 1, []string{})
	deviceIDHex, rawDataArray := helpers.Splice(rawDataArray, 0, 7, []string{})
	timestampHex, rawDataArray := helpers.Splice(rawDataArray, 0, 4, []string{})
	eventIDHex, rawDataArray := helpers.Splice(rawDataArray, 0, 1, []string{})

	fmt.Println(firmwareVersionHex, deviceIDHex, timestampHex, eventIDHex, rawDataArray)

	firmwareVersion, _ := helpers.HexSliceToBase10(firmwareVersionHex)
	deviceID, _ := helpers.HexSliceToBase10(deviceIDHex)
	timestamp, _ := helpers.HexSliceToBase10(timestampHex)
	eventID, _ := helpers.HexSliceToBase10(eventIDHex)

	fmt.Println(firmwareVersion, deviceID, timestamp, eventID)

	response := []byte("Acknowledged\n")

	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		helpers.LogError(err, "Error sending response")
	}
}
