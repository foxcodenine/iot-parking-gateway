import requests
import time
import json
import random
import base64

# File path
txt_file_path = './lora_raw_data.txt'

# Function to read data entries from a TXT file and parse JSON
def read_data_from_txt(file_path):
    data_entries = []
    try:
        with open(file_path, mode='r', newline='') as file:
            for line in file:
                try:
                    # Parse JSON from each line
                    json_data = json.loads(line.strip().strip('"'))
                    data_entries.append(json_data)
                except json.JSONDecodeError as e:
                    print(f"Error decoding JSON: {e}")
    except Exception as e:
        print(f"Error reading from TXT file: {e}")
    return data_entries

# Function to convert hex string to Base64
def hex_to_base64(hex_string):
    try:
        # Convert hex string to bytes
        byte_data = bytes.fromhex(hex_string)
        # Encode bytes to Base64
        base64_string = base64.b64encode(byte_data).decode('utf-8')
        return base64_string
    except ValueError as e:
        print(f"Error converting hex to Base64 for '{hex_string}': {e}")
        return None

# List of data entries to send
data_entries = read_data_from_txt(txt_file_path)

# Check if the list is empty
if not data_entries:
    print("No data found in the file.")
    exit(1)

# REST API URL
API_URL = "http://localhost:8080/api/sigfox"

# Use a loop to iterate over the data entries
for entry in data_entries:
    if 'data' in entry:
        base64_message = hex_to_base64(entry['data'])
        if base64_message:
            entry['data'] = base64_message
            try:
                # Send the POST request
                response = requests.post(API_URL, json=entry)
                # Print the response for debugging
                print(f"Sent: {entry}, Response: {response.status_code} {response.text}")
            except requests.RequestException as e:
                print(f"Error sending POST request: {e}")
            # Wait for a specific time between 1 and 2 seconds
            time.sleep(random.uniform(1, 2))
        else:
            print(f"Skipping entry due to Base64 conversion failure: {entry}")
    else:
        print(f"Skipping entry without 'data' field: {entry}")
