# Minimalist Dashboard

Minimalist Dashboard is a lightweight, real-time server monitoring tool. It provides a clean web interface to display live system metrics from your server.

 

## Features

- **Real-time Monitoring**: Metrics are updated every 2 seconds via WebSockets.
- **Comprehensive Metrics**:
    - **System**: OS, Kernel, CPU Model, Uptime
    - **CPU**: Usage percentage, Temperature
    - **RAM**: Used, Total, and Percentage
    - **Disk**: Total, Used, Free, and Percentage
    - **Network**: Inbound and Outbound traffic
    - **ZFS**: Pool configuration and status (Data VDEVs, Cache VDEV)
    - **ARC Cache**: Size, Target Size, Hit Rate
    - **L2ARC Cache**: Size, Hit Rate
    - **Docker**: Container, Image, and Volume counts
    - **Streaming Activity**: Media library counts (movies, shows, animes) and playback/transcoding status.
- **Lightweight & Efficient**: Built with Go on the backend and vanilla HTML/CSS/JS on the frontend.
- **Dockerized**: Easy to deploy with Docker. The multi-stage `Dockerfile` ensures a small and secure final image.

## Tech Stack

- **Backend**: Go
  - `net/http` for the web server
  - `github.com/gorilla/websocket` for WebSocket communication
- **Frontend**: HTML5, CSS3, JavaScript (ES6)
  - No frameworks, just vanilla JS.
- **Deployment**: Docker

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started) installed on your machine.
- [Go](https://golang.org/dl/) (v1.22+) installed for local development.

### Installation & Running

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/username/minimalist-dashboard.git
    cd minimalist-dashboard
    ```

2.  **Using Docker (Recommended):**
    Build and run the Docker container:
    ```bash
    docker build -t minimalist-dashboard .
    docker run -p 8080:8080 minimalist-dashboard
    ```
    The dashboard will be available at [http://localhost:8080](http://localhost:8080).

3.  **Running Locally:**
    Run the Go backend server:
    ```bash
    go run main.go
    ```
    The dashboard will be available at [http://localhost:8080](http://localhost:8080).

## Configuration

The application can be configured via environment variables. This is particularly useful for the "Streaming Activity" panel, which counts files in specified directories.

- `PATH_FILMS`: Absolute path to your movies directory.
- `PATH_SERIES`: Absolute path to your TV shows directory.
- `PATH_ANIMES`: Absolute path to your animes directory.

**Example with Docker:**

```bash
docker run -p 8080:8080 \
  -e PATH_FILMS="/mnt/media/movies" \
  -e PATH_SERIES="/mnt/media/tv" \
  -e PATH_ANIMES="/mnt/media/animes" \
  minimalist-dashboard
```

## How It Works

1.  The **Go backend** starts a web server on port 8080.
2.  It serves the static `frontend` directory (`index.html`, `style.css`).
3.  A WebSocket endpoint is available at `/ws`.
4.  When a client connects to `/ws`, the backend starts sending a JSON object with all system metrics every 2 seconds.
5.  The **JavaScript frontend** receives this JSON object and dynamically updates the HTML elements on the page to reflect the new data.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
