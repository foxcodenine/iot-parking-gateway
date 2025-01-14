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

// DeviceInfo contains device metadata for LoRa devices.
type DeviceInfo struct {
	TenantID          string                 `json:"tenantId"`
	TenantName        string                 `json:"tenantName"`
	ApplicationID     string                 `json:"applicationId"`
	ApplicationName   string                 `json:"applicationName"`
	DeviceProfileID   string                 `json:"deviceProfileId"`
	DeviceProfileName string                 `json:"deviceProfileName"`
	DeviceName        string                 `json:"deviceName"`
	DevEui            string                 `json:"devEui"`
	Tags              map[string]interface{} `json:"tags"`
}

// LoraData represents the structure for incoming LoRa data.
type LoraData struct {
	DeviceInfo DeviceInfo `json:"deviceInfo"`
	Data       string     `json:"data"`
}

// ConvertHexToBase64 converts the data from hex to Base64 format.
func (ld *LoraData) ConvertHexToBase64() error {
	decoded, err := hex.DecodeString(ld.Data)
	if err != nil {
		return fmt.Errorf("[LORA] Error decoding hex: %v", err)
	}
	ld.Data = base64.StdEncoding.EncodeToString(decoded)
	return nil
}

// Run processes data from a file and sends it to an API.
func Run() {
	file, err := os.Open("lora_raw_data.txt")
	if err != nil {
		log.Fatalf("[LORA] Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineCount int

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), `"`)
		var data LoraData

		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Fatalf("[LORA] Line %d: Error decoding JSON: %v", lineCount, err)
		}

		err = data.ConvertHexToBase64()
		if err != nil {
			log.Fatalf("[LORA] Line %d: Error converting hex to Base64: %v", lineCount, err)
		}

		err = sendDataToAPI(data)
		if err != nil {
			log.Printf("[LORA] Line %d: Error sending data: %v", lineCount, err)
		} else {
			log.Printf("[LORA] Line %d: Data sent successfully: %+v", lineCount, data)
		}

		fmt.Println("")

		lineCount++
		time.Sleep(time.Second) // Simulate throttling
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("[LORA] Error reading file: %v", err)
	}
}

// sendDataToAPI sends the LoRa data to an API via HTTP POST.
func sendDataToAPI(data LoraData) error {
	url := "http://localhost:8080/api/lora/chirpstack"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("[LORA] Error marshalling data to JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("[LORA] Error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
