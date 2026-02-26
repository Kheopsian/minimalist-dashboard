<div align="center">
  <img src="frontend/assets/images/favicon.png" alt="Minimalist Dashboard" width="120" height="120">
  <h1>Minimalist Dashboard</h1>
  
  [![Lighthouse Scores](https://img.shields.io/badge/Lighthouse_(Desktop)-100%25-brightgreen?logo=lighthouse)](https://googlechrome.github.io/lighthouse/viewer/)
  [![Docker Image CI](https://github.com/kheopsian/minimalist-dashboard/actions/workflows/docker-publish.yml/badge.svg?branch=main)](https://github.com/kheopsian/minimalist-dashboard/actions/workflows/docker-publish.yml)
  [![Docker Image](https://img.shields.io/badge/docker-ghcr.io-blue?logo=docker)](https://ghcr.io/Kheopsian/minimalist-dashboard)
  [![Go Version](https://img.shields.io/badge/go-1.25+-blue?logo=go)](https://golang.org/)
  [![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
  [![Live Demo](https://img.shields.io/badge/demo-dashboard.kheopsian.com-brightgreen)](https://dashboard.kheopsian.com)
  
  <p>A lightweight, real-time server monitoring dashboard built with Go and vanilla JavaScript. This project provides a clean web interface to display live system metrics with real-time updates via WebSockets.</p>
</div>

## ✨ Live Demo

🔗 **[dashboard.kheopsian.com](https://dashboard.kheopsian.com)**

## 📋 Features

### Real-time System Monitoring
- **CPU Usage & Temperature**: Live CPU utilization with thermal monitoring
- **Memory (RAM)**: Used, total, and percentage with visual progress bars
- **Storage**: Disk usage monitoring with mount point detection
- **Network Traffic**: Real-time network I/O with interactive charts
- **System Information**: OS, kernel, CPU model, and uptime details

### Advanced Storage Features
- **ZFS Pool Management**: Complete ZFS pool configuration display
  - Data VDEVs and Cache VDEV visualization
  - Pool status and device health monitoring
- **ARC & L2ARC Cache**: ZFS cache performance metrics
  - Hit rates, cache sizes, and efficiency monitoring

### Container & Media Services
- **Docker Integration**: Container, image, and volume counts
- **Media Library Monitoring**: Track movies, TV shows, and anime collections
- **Streaming Activity**: Real-time Plex integration for active playback and transcoding

### Technical Highlights
- **Real-time Updates**: Metrics refreshed every 2 seconds via WebSockets
- **Lightweight Architecture**: Pure Go backend with vanilla JavaScript frontend
- **Responsive Design**: Modern, clean interface with dark/light theme support
- **Containerized**: Production-ready Docker deployment

## 🛠 Tech Stack

### Backend
- **Go 1.25+** - High-performance server
- **Gorilla WebSocket** - Real-time communication
- **Native system calls** - Direct OS integration for metrics collection

### Frontend
- **HTML5/CSS3** - Modern semantic markup
- **Vanilla JavaScript (ES6+)** - No framework dependencies
- **Chart.js** - Interactive network traffic visualization
- **Font Awesome** - Professional iconography

### Infrastructure
- **Docker** - Containerized deployment
- **Multi-stage builds** - Optimized production images
- **Alpine Linux** - Minimal security footprint

## 🚀 Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (recommended)
- [Go 1.25+](https://golang.org/dl/) (for local development)

### Docker Deployment (Recommended)

#### Using Pre-built Image (GitHub Container Registry)

```bash
# Pull and run the latest image (Don't forget to mount your config.json)
docker run -p 9254:9254 -v /path/to/your/config.json:/app/config.json ghcr.io/Kheopsian/minimalist-dashboard:main
```

#### Building from Source

```bash
# Clone the repository
git clone https://github.com/Kheopsian/minimalist-dashboard.git
cd minimalist-dashboard

# Build and run with Docker
docker build -t minimalist-dashboard .
docker run -p 9254:9254 -v $(pwd)/config.json:/app/config.json minimalist-dashboard
```

Access the dashboard at [http://localhost:9254](http://localhost:9254)

### Local Development

```bash
# Clone and enter directory
git clone https://github.com/Kheopsian/minimalist-dashboard.git
cd minimalist-dashboard

# Install dependencies
go mod download

# Run the application
go run main.go
```

## ⚙️ Configuration

Configure the application using a `config.json` file. You can find a template in `config.example.json`.

```json
{
  "WebUIPort": "9254",
  "PathFilms": "/path/to/films",
  "PathSeries": "/path/to/series",
  "PathAnimes": "/path/to/animes",
  "NetInterface": "eth0",
  "PlexURL": "http://localhost:32400",
  "PlexToken": "your-plex-token-here"
}
```

### Example Docker Configuration

Make sure to map the path to your `config.json` when running the container.

```bash
docker run -p 9254:9254 \
  -v /path/to/your/config.json:/app/config.json:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /mnt/media:/mnt/media:ro \
  ghcr.io/Kheopsian/minimalist-dashboard:main
```

### Docker Compose

```yaml
version: '3.8'
services:
  dashboard:
    image: ghcr.io/Kheopsian/minimalist-dashboard:main
    ports:
      - "9254:9254"
    volumes:
      - /path/to/your/config.json:/app/config.json:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /mnt/media:/mnt/media:ro
```

### Unraid Installation

For Unraid users, an XML template is provided for a seamless setup:
1. Copy the `unraid-template.xml` file into your Unraid flash drive under `config/plugins/dockerMan/templates-user/`.
2. Go to your Unraid **Docker** tab.
3. Click **Add Container** and select `minimalist-dashboard` from the dropdown list of *User Templates*.
4. Ensure you place your `config.json` (based on `config.example.json`) in the `/mnt/user/appdata/minimalist-dashboard/` folder before starting the container.

#### ZFS Monitoring on Unraid (Important)

Due to container limitations, the dashboard needs to read ZFS status from a file generated by the Unraid host.
1. Install the **User Scripts** plugin on Unraid if you haven't already.
2. Create a new script, set it to run on a **Custom** schedule (e.g., `* * * * *` for every minute).
3. Paste the following script:

```bash
#!/bin/bash

# Configuration du dossier de destination
OUTPUT_DIR="/mnt/user/appdata/minimalist-dashboard"
OUTPUT_FILE="$OUTPUT_DIR/zpool_status.txt"

# Création du dossier si inexistant
mkdir -p "$OUTPUT_DIR"

# Génération du statut ZFS (Ajustez 'datapool' avec votre nom de pool ZFS)
if /sbin/zpool status datapool > "$OUTPUT_FILE"; then
    echo "Statut ZFS mis à jour : $(date)"
else
    echo "Erreur lors de la mise à jour du statut ZFS"
fi

# Ajustement des permissions
chmod 644 "$OUTPUT_FILE"
```

## 🏗 Architecture

### Backend Components
- **MetricsService**: Orchestrates all system metric collection
- **WebSocketHandler**: Manages real-time client connections
- **Service Layer**: Specialized services for different system components
  - `SystemService` - OS and hardware information
  - `CPUService` - Processor metrics and thermal data
  - `StorageService` - Disk usage and network statistics
  - `ZFSService` - ZFS pool and cache management
  - `DockerService` - Container platform integration
  - `StreamingService` - Media library and Plex monitoring

### Data Flow
1. **Metric Collection**: Backend services gather system data every 2 seconds
2. **WebSocket Broadcasting**: Real-time metrics pushed to connected clients
3. **Frontend Updates**: JavaScript dynamically updates UI components
4. **Chart Rendering**: Network traffic visualized with Chart.js

## 📊 Monitoring Capabilities

### System Metrics
- CPU usage percentage and temperature monitoring
- Memory utilization with detailed breakdown
- Disk space monitoring across mount points
- Network interface traffic analysis

### ZFS Integration
- Pool health and configuration display
- RAIDZ/Mirror VDEV visualization
- ARC and L2ARC cache performance
- Device-level status monitoring

### Container Platform
- Docker container lifecycle tracking
- Image and volume inventory
- Resource usage correlation

### Media Services
- Library size tracking (movies, shows, anime)
- Active streaming session monitoring
- Transcoding activity detection

## 🐳 Production Deployment

The included `Dockerfile` uses multi-stage builds for optimal production images:

```dockerfile
# Build stage with full Go toolchain
FROM golang:1.25-alpine AS builder
# ... build process

# Runtime stage with minimal dependencies
FROM alpine:latest
RUN apk --no-cache add ca-certificates docker-cli zfs
# ... final image
```

### Security Considerations
- Minimal Alpine Linux base image
- Non-root execution context
- Read-only filesystem mounts where possible
- Secure Docker socket access patterns

## 🔧 Development

### Project Structure
```
minimalist-dashboard/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP and WebSocket handlers
│   ├── models/            # Data structures
│   ├── services/          # Business logic layer
│   └── utils/             # Helper functions
├── frontend/              # Static web assets
│   ├── assets/
│   │   ├── css/          # Stylesheets and fonts
│   │   └── js/           # JavaScript modules
│   └── *.html            # HTML templates
└── Dockerfile            # Container build configuration
```

### Building from Source
```bash
# Development build
go build -o dashboard-api

# Production build (static binary)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dashboard-api
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 🔗 Related Projects

- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket implementation
- [Chart.js](https://www.chartjs.org/) - Interactive charts
- [Font Awesome](https://fontawesome.com/) - Icon library
