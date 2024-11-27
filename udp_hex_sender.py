import socket
import time
import csv
import random

# CSV file path
csv_file_path = './udp_raw_data.csv'

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

# List of hex strings to send
hex_strings = read_hex_strings_from_csv(csv_file_path)

# Check if the list is empty
if not hex_strings:
    print("No hex strings found in the file.")
    exit(1)

# UDP target IP and port
UDP_IP = "localhost"
UDP_PORT = 1236

# Create a socket
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

# Use a loop to iterate over the hex strings
for hex_string in hex_strings:
    # Encode the string as UTF-8 (plain text)
    message = hex_string.encode('utf-8')
    # Print the message for debugging
    print(f"Sending: {message}")
    # Send the plain-text message via UDP
    sock.sendto(message, (UDP_IP, UDP_PORT))
    # Wait for a specific time between 5 and 10 seconds
    time.sleep(random.uniform(1, 2))

# Optionally, loop continuously through the list
while True:
    for hex_string in hex_strings:
        message = hex_string.encode('utf-8')
        print(f"Sending: {message}")
        sock.sendto(message, (UDP_IP, UDP_PORT))
        time.sleep(random.uniform(5, 10))
