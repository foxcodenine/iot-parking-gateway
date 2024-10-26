iot-parking-gateway

## Project Overview

**iot-parking-gateway** is a Go-based server built to manage and communicate with IoT-enabled parking sensors. These sensors send and receive data over UDP, and the server acts as a central gateway for handling data flows, interfacing with a database, and ensuring real-time updates are accessible to connected services and clients. This project is Dockerized for ease of deployment and integrates various components such as PostgreSQL, Redis, RabbitMQ, and Socket.io to create a robust microservice architecture.

## Project Aim

The primary goal of **iot-parking-gateway** is to serve as a lightweight, scalable server to manage IoT parking sensors, offering:

1. **Two-way UDP Communication**: Direct communication with parking sensors to receive and respond to real-time data.
2. **Database Storage**: Efficient storage of structured data using PostgreSQL.
3. **Real-time Data Stream**: Enable instant data flow to front-end applications using Socket.io for a seamless user experience.
4. **Message Queue Integration**: Leverage RabbitMQ to connect and sync with other microservices.
5. **REST API**: Additional data access and control through a RESTful interface.

This project serves as a personal side project, aiming to deepen my knowledge of Go, PostgreSQL, Docker, and microservice communication, as well as to improve my skills in building high-performance server architectures suitable for IoT applications.

## Key Technologies

- **Go**: The main language used to develop a highly efficient, concurrent server.
- **Docker**: For containerizing the application and simplifying deployment.
- **PostgreSQL**: A reliable, SQL-based database for structured data storage.
- **Redis**: For quick-access caching and optimized data retrieval.
- **RabbitMQ**: To handle messaging between microservices.
- **Socket.io**: For real-time communication with front-end clients.
- **REST API**: Standardized API endpoints for extended application functionality.

## Project Structure

The project is organized into a modular folder structure:

    iot-parking-gateway/
    ├── cmd/                     # Main entry point for the server
    ├── config/                  # Configuration files
    ├── internal/                # Core application modules
    │   ├── db/                  # Database connections
    │   ├── models/              # Data models representing PostgreSQL tables
    │   ├── services/            # Business logic for UDP, DB interactions, etc.
    │   ├── iot/                 # IoT-specific logic for UDP communication
    │   ├── api/                 # REST and WebSocket routes and handlers
    │   └── mq/                  # RabbitMQ integration
    ├── pkg/                     # Reusable utility functions
    ├── scripts/                 # Database migration and setup scripts
    ├── Dockerfile               # Docker configuration for the Go app
    ├── docker-compose.yml       # Docker Compose setup for multi-service deployment
    └── go.mod                   # Go module file


## Getting Started

### Prerequisites

Ensure you have the following installed:

- **Docker** and **Docker Compose**
- **Go** (if running locally)
- **PostgreSQL** and **Redis** (if not using Docker)

### Running the Project


1. Clone the repository:

```bash
git clone https://github.com/yourusername/iot-parking-gateway.git
cd iot-parking-gateway
```
    
2. Build and start the services using Docker Compose:
    
```bash
docker-compose up --build
```
    
3. The server will be accessible at `http://localhost:4080` by default.

### Configuration

All configurations (e.g., database credentials, UDP ports, RabbitMQ settings) are managed through `config/config.yaml`.

## Future Goals

1. Expand REST API functionality to support a broader range of data access.
2. Add monitoring and logging features for better insight into system performance.
3. Integrate authentication and authorization mechanisms for secure data handling.
