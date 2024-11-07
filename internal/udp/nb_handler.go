package udp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/firmware"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/google/uuid"
)

// nbMessageHandler processes incoming UDP messages.
func (s *UDPServer) nbMessageHandler(conn *net.UDPConn, data []byte, addr *net.UDPAddr) {

	rawDataString := string(data)

	rawDataArray := helpers.SplitIntoPairs(rawDataString)

	firmwareVersionHex, rawDataArray := helpers.Splice(rawDataArray, 0, 1, []string{})
	deviceIDHex, rawDataArray := helpers.Splice(rawDataArray, 0, 7, []string{})

	firmwareVersion, _ := helpers.HexSliceToBase10(firmwareVersionHex)
	deviceID, _ := helpers.HexSliceToBase10(deviceIDHex)

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

	reply := []string{"0106"}
	hexTimestamp := helpers.GetCurrentTimestampHex()
	reply = append(reply, hexTimestamp)

	// Attempt to push the rawDataLog entry to Redis
	err = s.cache.RPush("raw-data-logs", rawDataLog)
	if err != nil {
		helpers.LogError(err, "Failed to push raw data log to Redis")
		sendResponse(conn, addr, reply)
		return
	}

	logDataMap, _ := firmware.NB_53(rawDataString)
	fmt.Println(logDataMap)

	// timestampHex, rawDataArray := helpers.Splice(rawDataArray, 0, 4, []string{})
	// eventIDHex, rawDataArray := helpers.Splice(rawDataArray, 0, 1, []string{})

	// timestamp, _ := helpers.HexSliceToBase10(timestampHex)
	// eventID, _ := helpers.HexSliceToBase10(eventIDHex)

	sendResponse(conn, addr, reply)
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, reply []string) {
	response := []byte(strings.Join(reply, "") + "\n")
	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		helpers.LogError(err, "Failed to send error response to client")
	}
}
