package udp

import (
	"bytes"
	"errors"
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

	// Prepare initial reply with message type and timestamp
	reply := []string{"0106"}
	hexTimestamp := helpers.GetCurrentTimestampHex()
	reply = append(reply, hexTimestamp)

	// Convert data to a string and trim any newlines
	hexStr := string(bytes.TrimSpace(data))

	// Validate minimum hex string length
	if len(hexStr) < 14 {
		handleErrorSendResponse(errors.New("incoming data too short for parsing"), "Invalid message length", conn, addr, reply)
		return
	}

	// Parse firmware version
	firmwareVersion, nextOffset, err := helpers.ParseHexSubstring(hexStr, 0, 1)
	if err != nil {
		handleErrorSendResponse(err, "Failed to parse firmware version", conn, addr, reply)
		return
	}

	// Parse device ID
	deviceID, _, err := helpers.ParseHexSubstring(hexStr, nextOffset, 7)
	if err != nil {
		handleErrorSendResponse(err, "Failed to parse device ID", conn, addr, reply)
		return
	}

	// Generate a new UUID for the RawDataLog entry
	uuid1, err := uuid.NewUUID()
	if err != nil {
		handleErrorSendResponse(err, "Failed to generate UUID for RawDataLog entry", conn, addr, reply)
		return
	}

	// Construct a new raw data log entry
	rawDataLog := models.RawDataLog{
		Uuid:            uuid1,
		DeviceID:        strconv.Itoa(deviceID),
		FirmwareVersion: firmwareVersion,
		NetworkType:     "nb",
		RawData:         hexStr,
		CreatedAt:       time.Now(),
	}

	// Push the raw data log entry to Redis
	err = s.cache.RPush("raw-data-logs", rawDataLog)
	if err != nil {
		handleErrorSendResponse(err, "Failed to push raw data log to Redis", conn, addr, reply)
		return
	}

	// Debug output for parsed values
	fmt.Println("Firmware Version:", firmwareVersion, "Device ID:", deviceID)

	parsedData, err := firmware.NB_53(hexStr)
	if err != nil {
		handleErrorSendResponse(err, "Failed to parse data from NB_53 firmware", conn, addr, reply)
		return
	} else {
		helpers.PrettyPrintJSON(parsedData)
	}

	// Send response back to the client
	sendResponse(conn, addr, reply)
}

// sendResponse sends a reply to the client over UDP.
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, reply []string) {
	response := []byte(strings.Join(reply, "") + "\n")
	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		helpers.LogError(err, "Failed to send response to client")
	}
}

// handleErrorSendResponse logs the error, sends a response, and returns to exit the function.
func handleErrorSendResponse(err error, message string, conn *net.UDPConn, addr *net.UDPAddr, reply []string) {
	helpers.LogError(err, message, 3)
	sendResponse(conn, addr, reply)
}
