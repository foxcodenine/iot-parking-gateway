1. Settings Table: 

    b. Enable/disable debug mode.
    c. Manage logs and errors:
        Print to terminal.
        Save to file.
        Disable printing and saving.
    d. Clear cashed users and devices


3. Sigfox: 
        a. Create router for Sigfox integration.
        b. Develop handler for processing Sigfox data.
        c. Implement firmware handling.
        d. Forward messages to RabbitMQ.
        e. Design models and create database tables.
        f. Save data to Redis and PostgreSQL.

4. UDP:
    a. Implement hidden functionality, in map view and device view.

6. Keepalive:
    a. Save last keeapalive happened_at in device.

7. Truncate db:
    a. Add auto db truncate in postgress.

8. Map:
    a. Add side bar
    

9. 
    a. Save authUser data during login.
    b. Implement favorites functionality in the map's info window and backend.
    c. Add functionality for users to change their passwords.
    d. Implement "Forgot Password" functionality.
    e. Add the ability to hide devices from the map and device views/page.
    f. Include a sidebar in the map view for better navigation and control.


10. NbiotDeviceSettings.BulkUpdate
    a. add settings_at in device
    b. update device setting only if there are more recent then device settings_at