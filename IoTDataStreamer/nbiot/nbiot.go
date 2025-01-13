package nbiot

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

// import necessary standard and third-party packages

type NbData string

func Run() {
	file, err := os.Open("udp_raw_data.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineCount int

	for scanner.Scan() {
		line := scanner.Text()

		var data NbData
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Fatalf("> %d\nError decoding JSON: %v\n", lineCount, err)
		}
	}
}
