package lora

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// DeviceInfo defines the structure for device-specific information.
type DeviceInfo struct {
	TenantID          string                 `json:"tenantId"`          // Unique identifier for the tenant
	TenantName        string                 `json:"tenantName"`        // Name of the tenant
	ApplicationID     string                 `json:"applicationId"`     // Identifier for the application
	ApplicationName   string                 `json:"applicationName"`   // Name of the application
	DeviceProfileID   string                 `json:"deviceProfileId"`   // ID of the device profile
	DeviceProfileName string                 `json:"deviceProfileName"` // Name of the device profile
	DeviceName        string                 `json:"deviceName"`        // Name of the device
	DevEui            string                 `json:"devEui"`            // Device EUI (Extended Unique Identifier)
	Tags              map[string]interface{} `json:"tags"`              // Arbitrary tags as a map for additional metadata
}

// ConvertHexToBase64 converts the Data field from a hex string to a base64 string.
func (ld *LoraData) ConvertHexToBase64() error {
	bytes, err := hex.DecodeString(ld.Data)
	if err != nil {
		return fmt.Errorf("error decoding hex: %v", err)
	}
	ld.Data = base64.StdEncoding.EncodeToString(bytes)
	return nil
}

// LoraData represents the structure for the incoming LoRa data.
type LoraData struct {
	DeviceInfo DeviceInfo `json:"deviceInfo"` // Embedded DeviceInfo structure
	Data       string     `json:"data"`       // Encoded data received from the device
}

// Run initiates the process of reading, processing, and sending data to an API.
func Run() {
	file, err := os.Open("lora_raw_data.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineCount int

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, `"`) // Clean the line from surrounding double quotes.

		var data LoraData
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Fatalf("> %d\nError decoding JSON: %v\n", lineCount, err)
		}

		// Convert hex data to base64 format.
		err = data.ConvertHexToBase64()
		if err != nil {
			log.Fatalf("> %d\nError converting hex to base64: %v\n", lineCount, err)
		}

		// Send the data to the remote API.
		err = sendDataToAPI(data)
		if err != nil {
			log.Printf("> %d\nError sending data to API: %v\n", lineCount, err)
			// Optionally handle the error, e.g., retry or log
		}

		fmt.Printf("(%d) ----------------------\nDecoded and sent data: %+v\n", lineCount, data)
		lineCount++

		time.Sleep(time.Second * 1) // Throttle the rate of data processing.
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
}

// sendDataToAPI sends the LoraData to the specified endpoint via HTTP POST.
func sendDataToAPI(data LoraData) error {
	url := "http://localhost:8080/api/lora/chirpstack"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data to JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	// Optionally, read and log the response body
	// responseBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	//     return fmt.Errorf("error reading response body: %v", err)
	// }
	// fmt.Printf("Received response: %s\n", responseBody)

	return nil
}
