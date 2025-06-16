# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an experimental multiplayer Real-Time Strategy (RTS) game built with Go and ChatGPT assistance. The game uses Ebiten for 2D graphics and WebSocket for real-time multiplayer communication.

## Architecture

### Client-Server Structure
- **Client** (`/client/`): Ebiten-based game client handling rendering, input, and WebSocket communication
- **Server** (`/server/`): WebSocket server managing game state and player connections  
- **Shared packages** (`/pkg/`): Common game logic, communication layer, and utilities

### Key Components
- **Game Logic** (`/pkg/game/`): Action-based game system with serializable game events
- **Communication** (`/pkg/comm/`): WebSocket client wrapper for bidirectional messaging
- **World Service** (`/pkg/world/`): External map data integration (port 8080)
- **Actions System**: All game events (movement, spawning, etc.) are represented as Actions

### Game Features
- Tile-based world with fog of war/visibility system
- Real-time unit selection and movement via mouse
- Camera controls with arrow keys
- Player identification via MD5 hash of player name
- External world service integration for dynamic map loading

## Development Commands

### Build and Run
```bash
# Build client
go build -o bin/client ./client

# Build server  
go build -o bin/server ./server

# Run client (requires server to be running)
./bin/client

# Run server (default port 8081)
./bin/server
```

### Development
```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Tidy dependencies
go mod tidy

# Vendor dependencies (if needed)
go mod vendor
```

### Testing
Currently no unit tests exist in the codebase. When adding tests, use Go's standard testing framework with `*_test.go` files.

## External Dependencies

The game expects an external world service running on port 8080 that provides map tile data via HTTP API.

## File Structure Notes

- Game assets (PNG files) are stored in `/client/` directory
- All dependencies are managed via `go.mod` 
- The project uses Go 1.23 with vendored dependencies
- WebSocket communication happens on port 8081 by default