Here’s a step-by-step guide to host your Dockerized app on your domain `iotrack.pro` using Apache as a reverse proxy and securing it with SSL using Let's Encrypt and Certbot.

---

### **Step 1: Install Apache**

On your cloud server, ensure Apache is installed and running.

bash

Copy code

`sudo apt update sudo apt install apache2 -y sudo systemctl start apache2 sudo systemctl enable apache2`

---

### **Step 2: Enable Apache Modules**

Enable the necessary Apache modules for proxying and SSL.

bash

Copy code

`sudo a2enmod proxy proxy_http ssl headers sudo systemctl restart apache2`

---

### **Step 3: Configure a Virtual Host**

Create an Apache virtual host file for your domain.

1. Open a new configuration file for your domain:
    
    bash
    
    Copy code
    
    `sudo nano /etc/apache2/sites-available/iotrack.pro.conf`
    
2. Add the following configuration:
    
    apache
    
    Copy code
    
    `<VirtualHost *:80>     ServerName iotrack.pro     ServerAlias www.iotrack.pro      ProxyPreserveHost On     ProxyRequests Off     ProxyPass / http://127.0.0.1:8080/     ProxyPassReverse / http://127.0.0.1:8080/      ErrorLog ${APACHE_LOG_DIR}/iotrack.pro-error.log     CustomLog ${APACHE_LOG_DIR}/iotrack.pro-access.log combined </VirtualHost>`
    
    - Replace `127.0.0.1:8080` with the internal IP and port of your Docker container if it differs.
    - Save the file and exit (`CTRL+O`, `CTRL+X`).
3. Enable the site and reload Apache:
    
    bash
    
    Copy code
    
    `sudo a2ensite iotrack.pro.conf sudo systemctl reload apache2`
    

---

### **Step 4: Test the HTTP Setup**

- Ensure your domain (`iotrack.pro`) points to the server's IP in your DNS settings.
- Visit `http://iotrack.pro` in your browser to confirm the site is accessible over HTTP.

---

### **Step 5: Install Certbot**

Install Certbot to manage SSL certificates.

1. Install Certbot:
    
    bash
    
    Copy code
    
    `sudo apt install certbot python3-certbot-apache -y`
    
2. Obtain and install the SSL certificate:
    
    bash
    
    Copy code
    
    `sudo certbot --apache -d iotrack.pro -d www.iotrack.pro`
    
    - Certbot will automatically configure Apache for HTTPS and redirect HTTP to HTTPS.
3. Test the SSL setup:
    
    - Visit `https://iotrack.pro` in your browser.
    - Certbot will set up an automatic renewal cron job for your certificate.

---

### **Step 6: Verify Certificate Renewal**

Check if Certbot’s renewal is working:

bash

Copy code

`sudo certbot renew --dry-run`

This command simulates the renewal process to ensure it's working.

