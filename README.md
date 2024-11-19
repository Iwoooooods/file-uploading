# File Server Upload Service

A Go-based HTTP server using Echo framework for handling file uploads.

## Prerequisites

- Go 1.23.3 or higher
- Make (for running make commands)

## API Endpoints

### Health Check
- **GET** `/api/ping`
- Returns: "pong" with 200 OK status
- Used to verify server is running

## Features

- Graceful shutdown handling
- Middleware support:
  - Logging
  - Panic recovery
- Configurable port via command line flags
- API routing with grouped endpoints
- Structured logging using zerolog

## Project Structure

## Shutdown

The server can be gracefully shutdown by sending an interrupt signal (Ctrl+C). It will:
1. Stop accepting new connections
2. Complete any in-flight requests
3. Perform cleanup operations
4. Exit cleanly