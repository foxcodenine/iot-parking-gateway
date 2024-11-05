## Goroutines for UDP Server in iot-parking-gateway Project

This section focuses on the goroutines used specifically in the UDP server, which handles communication with NB-IoT devices in this project. The UDP server relies on goroutines to keep the main application unblocked and to handle each incoming data packet from IoT devices concurrently. In the future, additional protocols (e.g., Sigfox, LoRa) may also use goroutines, though they will have distinct handling methods.

### Main Goroutine (in `main.go`)

- **Purpose**: Start the UDP server in a way that doesn’t block other parts of the program.
    
- **How It Works**: In `main.go`, we start the UDP server with:

```go
go func() {
    if err := udpServer.Start(); err != nil {
        app.ErrorLog.Fatalf("Failed to start UDP server: %v", err)
    }
}()
```
    
- **Explanation**:
    
    - By wrapping `udpServer.Start()` in a goroutine, we let `main.go` continue running without waiting for the UDP server to stop.
    - This allows the main application to initialize other components (like the HTTP server) without being blocked by the UDP server’s listening process.
    - The `defer udpServer.Stop()` line ensures the server stops gracefully when the program exits.

### Concurrent Goroutine for Each Packet (in `server.go`)

- **Purpose**: Process each incoming UDP packet independently.
    
- **How It Works**: Inside `server.go`, in the `listen` function, we start a new goroutine for each message:
    
```go
go handleUDPMessage(s.Connection, buffer[:n], addr)

```
    
- **Explanation**:
    
    - Every time a UDP packet arrives, `listen` reads it and immediately spawns a new goroutine to handle it.
    - This ensures that while one packet is being processed, other packets can still be received and handled simultaneously.
    - The function `handleUDPMessage` does the actual work for each packet: parsing the data, logging information, and potentially saving data to Redis or another database.

### Why Use These Goroutines?

- **Non-blocking Start**: The main UDP server starts in a goroutine so that it won’t hold up other tasks in `main.go`.
- **Parallel Packet Handling**: Each packet is handled in its own goroutine, ensuring the server can process multiple packets at the same time. This prevents bottlenecks, especially when multiple IoT devices are sending data at once.

### Summary

1. **Main Goroutine in `main.go`**: Keeps the main thread unblocked by running `udpServer.Start()` in a separate goroutine.
2. **Per-Packet Goroutine in `server.go`**: Allows each incoming packet to be processed concurrently, providing a scalable solution for handling many IoT devices.

This setup is efficient and keeps your server responsive, even when handling a high volume of incoming data.