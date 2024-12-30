import socket
import time
import csv
import random
import base64

# CSV file path
csv_file_path = './udp_raw_data2.csv'

# Function to read hex strings from a CSV file
def read_hex_strings_from_csv(file_path):
    hex_strings = []
    try:
        with open(file_path, mode='r', newline='') as file:
            reader = csv.reader(file)
            for row in reader:
                hex_strings.extend(row)  # Assuming each row is a list of hex strings
    except Exception as e:
        print(f"Error reading from CSV file: {e}")
    return hex_strings

# Function to convert hex string to Base64
def hex_to_base64(hex_string):
    try:
        # Convert hex string to bytes
        byte_data = bytes.fromhex(hex_string.strip())
        # Encode bytes to Base64
        base64_string = base64.b64encode(byte_data).decode('utf-8')
        return base64_string
    except ValueError as e:
        print(f"Error converting hex to Base64: {e}")
        return None

# List of hex strings to send
hex_strings = read_hex_strings_from_csv(csv_file_path)

# Check if the list is empty
if not hex_strings:
    print("No hex strings found in the file.")
    exit(1)

# UDP target IP and port
UDP_IP = "localhost"
UDP_PORT = 1234

# Create a socket
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

# Use a loop to iterate over the hex strings
for hex_string in hex_strings:
    base64_message = hex_to_base64(hex_string)
    if base64_message:
        # Print the Base64 message for debugging
        print(f"Sending (Base64): {base64_message}")
        # Send the Base64 message via UDP
        sock.sendto(base64_message.encode('utf-8'), (UDP_IP, UDP_PORT))
        # Wait for a specific time between 1 and 2 seconds
        time.sleep(random.uniform(1, 2))

# Optionally, loop continuously through the list
while True:
    for hex_string in hex_strings:
        base64_message = hex_to_base64(hex_string)
        if base64_message:
            print(f"Sending (Base64): {base64_message}")
            sock.sendto(base64_message.encode('utf-8'), (UDP_IP, UDP_PORT))
            time.sleep(random.uniform(5, 10))
