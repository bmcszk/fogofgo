# Fog of Go (fogofgo)

An experimental multiplayer Real-Time Strategy (RTS) game built with Go, featuring fog of war mechanics and real-time gameplay.

## About

Fog of Go is a tile-based RTS game that uses:
- **Go** for backend and game logic
- **Ebiten** for 2D graphics and rendering
- **WebSocket** for real-time multiplayer communication
- **Fog of War** system for strategic gameplay

## Quick Start

1. Build and run the server:
   ```bash
   go build -o bin/server ./server
   ./bin/server
   ```

2. Build and run the client (requires player name):
   ```bash
   go build -o bin/client ./client
   ./bin/client YourPlayerName
   ```

## Features

- Real-time multiplayer gameplay
- Tile-based world with fog of war/visibility system
- Unit selection and movement via mouse controls
- Camera controls with arrow keys
- Dynamic map loading from external world service
- Action-based game architecture for networked play

For detailed development information, see [CLAUDE.md](./CLAUDE.md).
