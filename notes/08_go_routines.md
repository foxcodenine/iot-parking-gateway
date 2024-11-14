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




i have this go func 

func (b *BeaconLog) BulkInsert(beacons []BeaconLog) error {

	for _, i := range beacons {
		fmt.Println(i.ActivityID, i.BeaconNumber)
	}

	// If there are no records to insert, exit early
	if len(beacons) == 0 {
		return nil
	}

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(beacons))      // Holds the placeholder for each row
	args := make([]interface{}, 0, len(beacons)*6) // Holds the actual values for each column

	for i, beacon := range beacons {
		// For each beacon, create a placeholder with indexed arguments, e.g., ($1, $2, $3, $4, $5, $6)
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))

		// Append the actual values for each placeholder in the same order as the columns
		args = append(args, beacon.ActivityID, beacon.HappenedAt, beacon.BeaconNumber, beacon.Major, beacon.Minor, beacon.RSSI)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf("INSERT INTO %s (activity_id, happened_at, beacon_number, major, minor, rssi) VALUES %s",
		b.TableName(), strings.Join(values, ", "))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert for beacons: %w", err)
	}

	return nil
}


and i am passing this slice to it:

[
  {
    "activity_id": "01932bf8-b5ae-7891-a7fc-ed6ad5d755bd",
    "happened_at": "2023-09-01T15:22:00+02:00",
    "beacon_number": 1,
    "major": 11233,
    "minor": 45666,
    "rssi": 204
  },
  {
    "activity_id": "01932bf8-b5ae-7891-a7fc-ed6ad5d755bd",
    "happened_at": "2023-09-01T15:22:00+02:00",
    "beacon_number": 2,
    "major": 23456,
    "minor": 65432,
    "rssi": 202
  },
  {
    "activity_id": "01932bf8-b5ae-7891-a7fc-ed6ad5d755bd",
    "happened_at": "2023-09-01T15:22:00+02:00",
    "beacon_number": 3,
    "major": 44551,
    "minor": 22334,
    "rssi": 176
  },
  {
    "activity_id": "01932bf8-b5ae-78ea-a036-a660e437f531",
    "happened_at": "2023-09-01T15:20:18+02:00",
    "beacon_number": 1,
    "major": 23456,
    "minor": 65432,
    "rssi": 205
  },
  {
    "activity_id": "01932bf8-b5ae-78ea-a036-a660e437f531",
    "happened_at": "2023-09-01T15:20:18+02:00",
    "beacon_number": 2,
    "major": 11233,
    "minor": 45666,
    "rssi": 198
  },
  {
    "activity_id": "01932bf8-b5ae-78ea-a036-a660e437f531",
    "happened_at": "2023-09-01T15:20:18+02:00",
    "beacon_number": 3,
    "major": 44551,
    "minor": 22334,
    "rssi": 167
  }
]

this is the table i am trying to populate:
CREATE TABLE IF NOT EXISTS parking.beacon_logs (
    activity_id UUID NOT NULL,
    happened_at TIMESTAMP NOT NULL,            -- Added to align with the composite key in activity_logs
    beacon_number INTEGER NOT NULL,            -- Sequence number for the beacon within the activity
    major INTEGER NOT NULL,
    minor INTEGER NOT NULL,
    rssi INTEGER,
    PRIMARY KEY (activity_id, happened_at, beacon_number),  -- Updated composite primary key
    FOREIGN KEY (activity_id, happened_at) REFERENCES parking.activity_logs(id, happened_at) ON DELETE CASCADE
);

01932bf8-b5ae-7891-a7fc-ed6ad5d755bd 1
01932bf8-b5ae-7891-a7fc-ed6ad5d755bd 2
01932bf8-b5ae-7891-a7fc-ed6ad5d755bd 3
01932bf8-b5ae-78ea-a036-a660e437f531 1
01932bf8-b5ae-78ea-a036-a660e437f531 2
01932bf8-b5ae-78ea-a036-a660e437f531 3
2024/11/14 19:39:00     Session ID:     00001
        Query:          INSERT INTO parking.beacon_logs (activity_id, happened_at, beacon_number, major, minor, rssi) VALUES ($1, $2, $3, $4, $5, $6), ($7, $8, $9, $10, $11, $12), ($13, $14, $15, $16, $17, $18), ($19, $20, $21, $22, $23, $24), ($25, $26, $27, $28, $29, $30), ($31, $32, $33, $34, $35, $36)
        Arguments:      []interface {}{uuid.UUID{0x1, 0x93, 0x2b, 0xf8, 0xb5, 0xae, 0x78, 0x91, 0xa7, 0xfc, 0xed, 0x6a, 0xd5, 0xd7, 0x55, 0xbd}, time.Date(2023, time.September, 1, 15, 22, 0, 0, time.Local), 1, 11233, 45666, 204, uuid.UUID{0x1, 0x93, 0x2b, 0xf8, 0xb5, 0xae, 0x78, 0x91, 0xa7, 0xfc, 0xed, 0x6a, 0xd5, 0xd7, 0x55, 0xbd}, time.Date(2023, time.September, 1, 15, 22, 0, 0, time.Local), 2, 23456, 65432, 202, uuid.UUID{0x1, 0x93, 0x2b, 0xf8, 0xb5, 0xae, 0x78, 0x91, 0xa7, 0xfc, 0xed, 0x6a, 0xd5, 0xd7, 0x55, 0xbd}, time.Date(2023, time.September, 1, 15, 22, 0, 0, time.Local), 3, 44551, 22334, 176, uuid.UUID{0x1, 0x93, 0x2b, 0xf8, 0xb5, 0xae, 0x78, 0xea, 0xa0, 0x36, 0xa6, 0x60, 0xe4, 0x37, 0xf5, 0x31}, time.Date(2023, time.September, 1, 15, 20, 18, 0, time.Local), 1, 23456, 65432, 205, uuid.UUID{0x1, 0x93, 0x2b, 0xf8, 0xb5, 0xae, 0x78, 0xea, 0xa0, 0x36, 0xa6, 0x60, 0xe4, 0x37, 0xf5, 0x31}, time.Date(2023, time.September, 1, 15, 20, 18, 0, time.Local), 2, 11233, 45666, 198, uuid.UUID{0x1, 0x93, 0x2b, 0xf8, 0xb5, 0xae, 0x78, 0xea, 0xa0, 0x36, 0xa6, 0x60, 0xe4, 0x37, 0xf5, 0x31}, time.Date(2023, time.September, 1, 15, 20, 18, 0, time.Local), 3, 44551, 22334, 167}
        Stack:          
                fmt.(*pp).handleMethods@/usr/local/go/src/fmt/print.go:673
                fmt.(*pp).printArg@/usr/local/go/src/fmt/print.go:756
                fmt.(*pp).doPrint@/usr/local/go/src/fmt/print.go:1208
                fmt.Append@/usr/local/go/src/fmt/print.go:289
                log.(*Logger).Print.func1@/usr/local/go/src/log/log.go:261
                log.(*Logger).output@/usr/local/go/src/log/log.go:238
                log.(*Logger).Print@/usr/local/go/src/log/log.go:260
                github.com/foxcodenine/iot-parking-gateway/internal/models.(*BeaconLog).BulkInsert@/home/foxcodenine/foxfiles/git/repo/iot-parking-gateway/internal/models/beacons_log.go:57
                github.com/foxcodenine/iot-parking-gateway/internal/services.(*Service).RedisToPostgresActivityLags@/home/foxcodenine/foxfiles/git/repo/iot-parking-gateway/internal/services/services.go:158
                main.main.func1@/home/foxcodenine/foxfiles/git/repo/iot-parking-gateway/cmd/app/main.go:47
                github.com/robfig/cron/v3.FuncJob.Run@/home/foxcodenine/go/pkg/mod/github.com/robfig/cron/v3@v3.0.0/cron.go:131
                github.com/robfig/cron/v3.(*Cron).startJob.func1@/home/foxcodenine/go/pkg/mod/github.com/robfig/cron/v3@v3.0.0/cron.go:307
                runtime.goexit@/usr/local/go/src/runtime/asm_amd64.s:1700
        Error:          ERROR: insert or update on table "beacon_logs" violates foreign key constraint "beacon_logs_activity_id_happened_at_fkey" (SQLSTATE 23503)
        Time taken:     0.00547s
        Context:        context.Background

ERROR   2024/11/14 19:39:00 services.go:160: Failed to insert beacon logs to PostgreSQL: failed to execute bulk insert for beacons: ERROR: insert or update on table "beacon_logs" violates foreign key constraint "beacon_logs_activity_id_happened_at_fkey" (SQLSTATE 23503)

i do not know y becase the uuid and beacon number combinations are all unique as you can see:

01932bf8-b5ae-7891-a7fc-ed6ad5d755bd 1
01932bf8-b5ae-7891-a7fc-ed6ad5d755bd 2
01932bf8-b5ae-7891-a7fc-ed6ad5d755bd 3
01932bf8-b5ae-78ea-a036-a660e437f531 1
01932bf8-b5ae-78ea-a036-a660e437f531 2
01932bf8-b5ae-78ea-a036-a660e437f531 3
