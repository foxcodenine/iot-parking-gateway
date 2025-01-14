package sigfox

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

// SigfoxData defines the structure for incoming Sigfox data.
type SigfoxData struct {
	Timestamp int64  `json:"timestamp"`
	Device    string `json:"device"`
	SeqNumber string `json:"seq_number"`
	Data      string `json:"data"`
}

// ConvertHexToBase64 converts the data from hex to Base64 format.
func (sd *SigfoxData) ConvertHexToBase64() error {
	decoded, err := hex.DecodeString(sd.Data)
	if err != nil {
		return fmt.Errorf("[SIGFOX] Error decoding hex: %v", err)
	}
	sd.Data = base64.StdEncoding.EncodeToString(decoded)
	return nil
}

// Run processes data from a file and sends it to an API.
func Run() {
	file, err := os.Open("sigfox_raw_data.txt")
	if err != nil {
		log.Fatalf("[SIGFOX] Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineCount int

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), `"`)
		var data SigfoxData

		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Fatalf("[SIGFOX] Line %d: Error decoding JSON: %v", lineCount, err)
		}

		err = data.ConvertHexToBase64()
		if err != nil {
			log.Fatalf("[SIGFOX] Line %d: Error converting hex to Base64: %v", lineCount, err)
		}

		err = sendDataToAPI(data)
		if err != nil {
			log.Printf("[SIGFOX] Line %d: Error sending data: %v", lineCount, err)
		} else {
			log.Printf("[SIGFOX] Line %d: Data sent successfully: %+v", lineCount, data)
		}

		fmt.Println("")

		lineCount++
		time.Sleep(time.Second) // Simulate throttling
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("[SIGFOX] Error reading file: %v", err)
	}
}

// sendDataToAPI sends the Sigfox data to an API via HTTP POST.
func sendDataToAPI(data SigfoxData) error {
	url := "http://localhost:8080/api/sigfox"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("[SIGFOX] Error marshalling data to JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("[SIGFOX] Error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
