# Web Crawler

A distributed web crawler built in Go that crawls websites and stores the data in MongoDB. The crawler is designed to be efficient and scalable, with features like rate limiting and concurrent processing.

## Features

- Concurrent web crawling
- MongoDB integration for data storage
- Docker support for easy deployment
- Rate limiting and politeness policies
- Real-time crawling statistics
- Configurable crawling depth and limits

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- MongoDB (automatically set up with Docker)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/mauFade/web-crawler.git
cd web-crawler
```

2. Create a `.env` file in the root directory (optional):

```bash
MONGODB_URI=mongodb://localhost:27017
```

## Running the Application

### Using Docker (Recommended)

1. Build and start the containers:

```bash
docker-compose up --build
```

This will start both the web crawler application and MongoDB in separate containers.

### Running Locally

1. Install dependencies:

```bash
go mod download
```

2. Run the application:

```bash
go run cmd/main.go
```

## Configuration

The crawler can be configured through environment variables:

- `MONGODB_URI`: MongoDB connection string (default: mongodb://localhost:27017)

## Project Structure

```
.
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── db/              # Database connection and operations
│   ├── models/          # Data models and structures
│   └── utils/           # Utility functions
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── go.mod              # Go module file
└── go.sum              # Go module checksum
```

## How It Works

1. The crawler starts with a seed URL (default: https://www.cc.gatech.edu/)
2. It fetches the page content and parses all links
3. New URLs are added to a queue for processing
4. The crawler continues until it reaches the maximum number of pages (5000 by default)
5. Crawling statistics are displayed every minute
6. All crawled data is stored in MongoDB

## Monitoring

The crawler provides real-time statistics including:

- Total queued URLs
- Current queue size
- Number of crawled pages
- Crawling rate

## License

This project is open source and available under the MIT License.
