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

// import necessary standard and third-party packages
type SigfoxData struct {
	Timestamp string `json:"timestamp"`
	Device    string `json:"device"`
	SeqNumber string `json:"seq_number"`
	Data      string `json:"data"`
}

// ConvertHexToBase64 converts the Data field from a hex string to a base64 string.
func (sd *SigfoxData) ConvertHexToBase64() error {
	bytes, err := hex.DecodeString(sd.Data)
	if err != nil {
		return fmt.Errorf("error decoding hex: %v", err)
	}
	sd.Data = base64.StdEncoding.EncodeToString(bytes)
	return nil
}

// Run initiates the process of reading, processing, and sending data to an API.
func Run() {
	file, err := os.Open("sigfox_raw_data.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineCount int

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, `"`) // Clean the line from surrounding double quotes.

		var data SigfoxData
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Fatalf("> %d\nError decoding JSON: %v\n", lineCount, err)
		}

		// Convert hex data to base64 format.
		err = data.ConvertHexToBase64()
		if err != nil {
			log.Fatalf("> %d\nError converting hex to base64: %v\n", lineCount, err)
		}

		fmt.Printf("----------------------\nDecoded and sent data: %+v\n", data)
		lineCount++

		time.Sleep(time.Second * 1) // Throttle the rate of data processing.
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
}

// {"timestamp":1735373204,"device":"034D35AE","seq_number":"186","data":"PAoBABcgTgA="}

// sendDataToAPI sends the SigfoxData to the specified endpoint via HTTP POST.
func sendDataToAPI(data SigfoxData) error {
	url := "http://localhost:8080/api/sigfox"
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
