package udp

import (
	"bytes"
	"encoding/json"
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

// nbMessageHandler processes incoming UDP messages and logs data to Redis.
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
	firmwareVersionTmp, nextOffset, err := helpers.ParseHexSubstring(hexStr, 0, 1)
	if err != nil {
		handleErrorSendResponse(err, "Failed to parse firmware version", conn, addr, reply)
		return
	}
	// Divide by 10 to convert to float64 and shift decimal place
	firmwareVersion := float64(firmwareVersionTmp) / 10.0

	// Parse device ID
	deviceID, _, err := helpers.ParseHexSubstring(hexStr, nextOffset, 7)
	if err != nil {
		handleErrorSendResponse(err, "Failed to parse device ID", conn, addr, reply)
		return
	}

	// TODO:  Implement checks against a blacklist or whitelist for the device ID.

	// Check if the device ID is already in the Bloom Filter
	networkTypeAndDeviceID := fmt.Sprintf("NB-IoT %d", deviceID)
	isDeviceRegistered, err := s.cache.CheckItemInBloomFilter("device-id", networkTypeAndDeviceID)
	if err != nil {
		helpers.LogError(err, "Failed to check Bloom Filter for device ID")
	}

	// If the device ID is not registered, add it to the set and Bloom Filter
	if !isDeviceRegistered {
		// Add the device ID to a Redis set for later processing
		if err := s.cache.SAdd("device-to-create", networkTypeAndDeviceID); err != nil {
			helpers.LogError(err, "Failed to add device ID to the set")
		}

		// Add the device ID to the Bloom Filter to register it
		if _, err := s.cache.AddItemToBloomFilter("device-id", networkTypeAndDeviceID); err != nil {
			helpers.LogError(err, "Failed to add device ID to the Bloom Filter")
		}
	}

	// Generate a new UUID for the RawDataLog entry
	rawUUID, err := uuid.NewV7()
	if err != nil {
		handleErrorSendResponse(err, "Failed to generate UUID for RawDataLog entry", conn, addr, reply)
		return
	}

	// Create a new RawDataLog object to store in Redis.
	rawDataLog := models.RawDataLog{
		ID:              rawUUID,
		DeviceID:        strconv.Itoa(deviceID),
		FirmwareVersion: firmwareVersion,
		NetworkType:     "NB-IoT",
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
	helpers.LogInfo("Firmware Version: %.2f Device ID: %d", firmwareVersion, deviceID)

	// Process firmware-specific data parsing based on the firmware version.
	var parsedData map[string]any
	switch firmwareVersion {
	case 5.3:
		parsedData, err = firmware.NB_53(hexStr)
	case 5.8, 5.9:
		parsedData, err = firmware.NB_58(hexStr)

	default:
		// Send a default response if the firmware version is not handled.
		sendResponse(conn, addr, reply)
		return
	}

	if err != nil {
		handleErrorSendResponse(err, "Failed to parse data from NB_53 firmware", conn, addr, reply)
		return
	}

	// Push parsed parking data packages to Redis.

	for _, i := range parsedData["parking_packages"].([]map[string]any) {
		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = fmt.Sprintf("%d", parsedData["device_id"])
		i["raw_id"] = rawUUID
		i["event_id"] = 26
		i["network_type"] = "NB-IoT"

		// socketserver.IOService.SocketServer.BroadcastToNamespace("/", "update", i)

		err := s.cache.RPush("nb-activity-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push parking package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}
		s.mqProducer.SendMessage("nb_iot_event_logs_exchange", "nb_iot_event_logs_queue", string(messageData))
	}

	// Push parsed keepalive data to Redis.
	for _, i := range parsedData["keep_alive_packages"].([]map[string]any) {
		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = fmt.Sprintf("%d", parsedData["device_id"])
		i["raw_id"] = rawUUID
		i["event_id"] = 6
		i["network_type"] = "NB-IoT"

		err := s.cache.RPush("nb-keepalive-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push keepalive package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}
		s.mqProducer.SendMessage("nb_iot_event_logs_exchange", "nb_iot_event_logs_queue", string(messageData))
	}

	// Push parsed settings data to Redis.
	for _, i := range parsedData["settings_packages"].([]map[string]any) {
		// Add common fields to each individual package
		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = fmt.Sprintf("%d", parsedData["device_id"])
		i["raw_id"] = rawUUID
		i["event_id"] = 25 // Assuming 25 is the event ID for setting logs
		i["network_type"] = "NB-IoT"

		// Push the package to Redis
		err := s.cache.RPush("nb-setting-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push setting package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}
		s.mqProducer.SendMessage("nb_iot_event_logs_exchange", "nb_iot_event_logs_queue", string(messageData))
	}

	// time.Sleep(1 * time.Second)
	// s.services.RegisterNewDevices()
	// s.services.SyncActivityLogsAndDevices()
	// s.services.SyncNBIoTKeepaliveLogs()
	// s.services.SyncNBIoTSettingLogs()

	// Send the final response back to the UDP client to confirm processing.
	sendResponse(conn, addr, reply)
}

// sendResponse sends a structured reply back to the UDP client.
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, reply []string) {
	response := []byte(strings.Join(reply, "") + "\n")
	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		helpers.LogError(err, "Failed to send response to client")
	}
}

// handleErrorSendResponse logs an error, sends a response to the client, and exits the handler.
func handleErrorSendResponse(err error, message string, conn *net.UDPConn, addr *net.UDPAddr, reply []string) {
	helpers.LogError(err, message, 3)
	sendResponse(conn, addr, reply)
}
