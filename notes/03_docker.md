
## Development Workflow

For efficient development with continuous updates, try these options to rebuild or restart containers:

1. **Auto-Restart with 
    Rebuilds images as needed and restarts containers, keeping everything fresh without a manual teardown.

    ```bash
    docker-compose up --build
    ```

    
2. **Detached Mode with 
    Runs containers in detached mode. 

    ```bash
    docker-compose up --build -d
    ```
    
    To restart only a specific service, use:

    ```bash
    docker-compose restart app
    ```

3. **Bind-Mount Source Code for Real-Time Changes**  
    In development, bind-mount the source code to see changes instantly in the container. For production, comment out the bind mounts to rely on Docker volumes or images instead.

### Development Setup

Uncomment the following lines to enable bind mounts for development:

    ``` yaml
    services:
    app:
        volumes:
        - .:/app  # (development only) - Mounts the current directory into /app inside the container for live code updates 

    postgres:
        volumes:
        - /srv/docker/bind-mounts/iot-parking-gateway/postgres-data/:/var/lib/postgresql/data/  # (development only) - Bind mount for dev, enables direct file access on host

    redis:
        volumes:
        - /srv/docker/bind-mounts/iot-parking-gateway/redis-data/:/data  # (development only) - Bind mount for dev, enables direct file access on host

    ```


### Production Setup

Comment out the development bind mounts and uncomment the Docker-managed volumes for production:

    ``` yaml
    # For production, use Docker-managed volumes:
    services:
    app:
        # volumes:
        #   - .:/app  # Comment out this line in production

    postgres:
        # volumes:
        #   - /srv/docker/bind-mounts/iot-parking-gateway/postgres-data/:/var/lib/postgresql/data/  # Comment out in production
        - postgres-data:/var/lib/postgresql/data  # Uncomment this line in production to persist data

    redis:
        # volumes:
        #   - /srv/docker/bind-mounts/iot-parking-gateway/redis-data/:/data  # Comment out in production
        - redis-data:/data  # Uncomment this line in production to persist data

    volumes: # Uncomment this line in production to persist data
    postgres-data:
    redis-data:
    ```

This approach:

- **Enables bind mounts** in development for easy, real-time updates.
- **Switches to Docker-managed volumes** in production for secure and isolated data persistence.


<!-- --------------------------------------------------------------- -->

### Upload app to Docker Hub

1. Build and Tag Your Docker Image:

    $ docker build -t foxcodenine/iot-parking-gateway_app:latest -f config/app/Dockerfile .

2. Log in to Docker Hub:

    $ docker login

3. Push the Image to Docker Hub:

    docker push foxcodenine/iot-parking-gateway_app:latest

4. Update the docker-compose.yml:

    app:
        image: foxcodenine/iot-parking-gateway_app:latest

<!-- --------------------------------------------------------------- -->


