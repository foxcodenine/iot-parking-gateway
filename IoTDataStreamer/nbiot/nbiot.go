package nbiot

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	// Configuration constants for flexibility
	serverAddr = "127.0.0.1:1234"   // UDP server address
	filePath   = "udp_raw_data.csv" // Path to the CSV file
)

// Run reads data from a CSV file and simulates an IoT NB device by sending the data line by line to a UDP server.
func Run() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("[NBIOT] Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineCount int

	// Iterate over each line in the file
	for scanner.Scan() {
		time.Sleep(330 * time.Millisecond) // Simulate device throttling

		line := scanner.Text()
		err = sendDataToUDP(line)
		if err != nil {
			log.Printf("[NBIOT] Line %d: Error sending data: %v", lineCount, err)
		} else {
			log.Printf("[NBIOT] Line %d: Data sent successfully", lineCount)
		}

		lineCount++
	}

	// Check for file scanning errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("[NBIOT] Error reading file: %v", err)
	}
}

// sendDataToUDP encodes the data in Base64 and sends it to the specified UDP server.
func sendDataToUDP(data string) error {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return fmt.Errorf("[NBIOT] Error resolving UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("[NBIOT] Error creating UDP connection: %v", err)
	}
	defer conn.Close()

	decodedData, err := hex.DecodeString(data)
	if err != nil {
		return fmt.Errorf("[NBIOT] Error decoding hex: %v", err)
	}

	// encodedData := base64.StdEncoding.EncodeToString(decodedData)
	// _, err = conn.Write([]byte(encodedData))

	_, err = conn.Write([]byte(decodedData))
	if err != nil {
		return fmt.Errorf("[NBIOT] Error sending data: %v\n", err)
	}

	return nil
}
