package udp

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/google/uuid"
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

	uuid1, err := uuid.NewUUID()

	if err != nil {
		helpers.LogError(err, "Failed to generate a new UUID for RawDataLog entry")
	}

	rawDataLog := models.RawDataLog{
		Uuid:            uuid1,
		DeviceID:        strconv.Itoa(deviceID),
		FirmwareVersion: firmwareVersion,
		NetworkType:     "nb",
		RawData:         rawDataString,
		CreatedAt:       time.Now(),
	}

	_ = s.cache.RPush("raw-data-logs", rawDataLog)

	// timestampHex, rawDataArray := helpers.Splice(rawDataArray, 0, 4, []string{})
	// eventIDHex, rawDataArray := helpers.Splice(rawDataArray, 0, 1, []string{})

	// timestamp, _ := helpers.HexSliceToBase10(timestampHex)
	// eventID, _ := helpers.HexSliceToBase10(eventIDHex)

	response := []byte("Acknowledged\n")

	_, err = conn.WriteToUDP(response, addr)
	if err != nil {
		helpers.LogError(err, "Error sending response")
	}
}
