1. Settings Table: 
    a. Add app model for whitelist or blacklist.
    b. Enable/disable debug mode.
    c. Manage logs and errors:
        Print to terminal.
        Save to file.
        Disable printing and saving.
    d. Clear cashed users and devices


    Create a settings table that holds a key-value pair together with a description column.
    Populate the settings table from .env settings when you initialize the app for the first time.
    Each time the app is started, create a Redis key-value pair for each record in the settings table.
    Retrieve settings from Redis when a user logs in; if not available, fallback to PostgreSQL.
    When app settings are changed, push updates to the user frontend.
    Require all other users to sign out and sign back in when settings are updated to ensure they receive the new configurations.

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

    a. Add access level in form and table
    c  Implement table sorting and search functionality.
  


6. Device API:
    a. Add authentication and check user access_level.

7. Socket.IO:
    a. Start implementing from both the backend and client side.

8. Map:
    a. Begin by loading devices onto the map.
    b. Update device data dynamically using Socket.IO.