# File Server Upload Service

A Go-based HTTP server using Echo framework for handling file uploads with metadata tracking and deduplication support.

## Prerequisites

- Go 1.23.3 or higher
- Make (for running make commands)
- SQLite3

## Features

- File upload and download capabilities
- File deduplication using MD5 hash checking
- Metadata storage in SQLite database
- User-specific storage directories
- Graceful shutdown handling
- Middleware support:
  - Logging
  - Panic recovery
  - CORS enabled
- Configurable via environment variables
- API routing with grouped endpoints
- Structured logging using zerolog

## API Endpoints

### Health Check
- **GET** `/api/ping`
- Returns: "pong" with 200 OK status
- Used to verify server is running

### File Upload
- **POST** `/api/upload/:username/:filename`
- Uploads a file for a specific user
- Performs MD5 hash checking for deduplication
- Returns:
  - `fileId` if file is newly uploaded
  - `exists: true` if file already exists

### File Download
- **GET** `/api/download/:username/:filename`
- Downloads a specific file for a user

## Project Structure

## Shutdown

The server can be gracefully shutdown by sending an interrupt signal (Ctrl+C). It will:
1. Stop accepting new connections
2. Complete any in-flight requests
3. Perform cleanup operations
4. Exit cleanly