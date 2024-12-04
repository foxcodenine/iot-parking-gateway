1. Settings Table: 
    a. Add app model for whitelist or blacklist.
    b. Enable/disable debug mode.
    c. Manage logs and errors:
        Print to terminal.
        Save to file.
        Disable printing and saving.

2. User View: 
    a. Implement table sorting and search functionality.

3. Sigfox: 
        a. Create router for Sigfox integration.
        b. Develop handler for processing Sigfox data.
        c. Implement firmware handling.
        d. Forward messages to RabbitMQ.
        e. Design models and create database tables.
        f. Save data to Redis and PostgreSQL.

4. UDP:
    a. Implement whitelist and blacklist functionality.

5. Device Page:
    a. Add a table for devices.
    b. Include a form for device creation.
    c. Add edit functionality.
    d. Develop device location model.
    e. Implement hiding, whitelisting, and blacklisting.

6. Device API:
    a. Add authentication and check user access_level.

7. Socket.IO:
    a. Start implementing from both the backend and client side.

8. Map:
    a. Begin by loading devices onto the map.
    b. Update device data dynamically using Socket.IO.