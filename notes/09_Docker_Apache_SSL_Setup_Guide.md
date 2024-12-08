
# Setup Guide: Dockerized App with Apache and SSL on iotrack.pro

## Step 1: Install Apache
On your cloud server, ensure Apache is installed and running.
```bash
sudo apt update
sudo apt install apache2 -y
sudo systemctl start apache2
sudo systemctl enable apache2
```

## Step 2: Enable Apache Modules
Enable the necessary Apache modules for proxying, SSL, and WebSocket support.
```bash
sudo a2enmod proxy proxy_http ssl headers proxy_wstunnel rewrite
sudo systemctl restart apache2
```

## Step 3: Configure a Virtual Host
Create an Apache virtual host file for your domain.

1. Open a new configuration file for your domain using `vi`:
   ```bash
   sudo vi /etc/apache2/sites-available/iotrack.pro.conf
   ```
2. Add the following configuration for handling both HTTP and WebSocket traffic:
   ```apache
   <VirtualHost *:80>
       ServerName iotrack.pro
       ServerAlias www.iotrack.pro

       ProxyPreserveHost On
       ProxyRequests Off
       ProxyPass / http://127.0.0.1:8080/
       ProxyPassReverse / http://127.0.0.1:8080/

       ErrorLog ${APACHE_LOG_DIR}/iotrack.pro-error.log
       CustomLog ${APACHE_LOG_DIR}/iotrack.pro-access.log combined

       RewriteEngine on
       RewriteCond %{SERVER_NAME} =iotrack.pro [OR]
       RewriteCond %{SERVER_NAME} =www.iotrack.pro
       RewriteRule ^ https://%{SERVER_NAME}%{REQUEST_URI} [END,NE,R=permanent]
   </VirtualHost>
   ```
3. Enable the site and reload Apache:
   ```bash
   sudo a2ensite iotrack.pro.conf
   sudo systemctl reload apache2
   ```

## Step 4: Test the HTTP Setup
Ensure your domain (`iotrack.pro`) points to the server's IP in your DNS settings.
Visit `http://iotrack.pro` to confirm the site is accessible over HTTP.

## Step 5: Install Certbot
Install Certbot to manage SSL certificates.

1. Install Certbot:
   ```bash
   sudo apt install certbot python3-certbot-apache -y
   ```
2. Obtain and install the SSL certificate:
   ```bash
   sudo certbot --apache -d iotrack.pro -d www.iotrack.pro
   ```
   Certbot will automatically configure Apache for HTTPS and redirect HTTP to HTTPS.

3. Test the SSL setup by visiting `https://iotrack.pro`.

## Step 6: Verify Certificate Renewal
Check if Certbot’s renewal is working:
```bash
sudo certbot renew --dry-run
```
This command simulates the renewal process to ensure it's working.

## Step 7: Add WebSocket SSL Configuration
Adjust your SSL virtual host to include WebSocket handling.

1. Edit the SSL virtual host file using `vi`:
   ```bash
   sudo vi /etc/apache2/sites-available/iotrack.pro-le-ssl.conf
   ```
2. Add the following lines to handle WebSocket traffic:
   ```apache
   RewriteEngine On
   RewriteCond %{HTTP:Upgrade} =websocket [NC]
   RewriteRule /(.*)           ws://localhost:8080/$1 [P,L]

   ProxyPass /socket.io/ ws://127.0.0.1:8080/socket.io/
   ProxyPassReverse /socket.io/ ws://127.0.0.1:8080/socket.io/
   ```
3. Save the changes and restart Apache:
   ```bash
   sudo systemctl restart apache2
   ```
4. Test WebSocket functionality to ensure that WebSocket connections are properly secured and functional.

## Complete SSL Virtual Host Configuration

Here is the full SSL Virtual Host configuration incorporating WebSocket support:

```apache
<VirtualHost *:443>
    ServerName iotrack.pro
    ServerAlias www.iotrack.pro

    SSLEngine on
    SSLCertificateFile /etc/letsencrypt/live/iotrack.pro-0001/fullchain.pem
    SSLCertificateKeyFile /etc/letsencrypt/live/iotrack.pro-0001/privkey.pem

    ErrorLog ${APACHE_LOG_DIR}/iotrack.pro-error.log
    CustomLog ${APACHE_LOG_DIR}/iotrack.pro-access.log combined

    ProxyRequests Off
    ProxyPreserveHost On

    # WebSocket configuration
    RewriteEngine On
    RewriteCond %{HTTP:Upgrade} =websocket [NC]
    RewriteRule /(.*)           ws://localhost:8080/$1 [P,L]

    ProxyPass /socket.io/ ws://127.0.0.1:8080/socket.io/
    ProxyPassReverse /socket.io/ ws://127.0.0.1:8080/socket.io/

    # HTTP fallback for non-WebSocket connections
    ProxyPass / http://127.0.0.1:8080/
    ProxyPassReverse / http://127.0.0.1:8080/

    Include /etc/letsencrypt/options-ssl-apache.conf
</VirtualHost>
```
