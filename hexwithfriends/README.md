# hexwithfriends

Two-player [Hex](https://en.wikipedia.org/wiki/Hex_(board_game)) board game. Play online at [hexwithfriends.com](https://hexwithfriends.com).

Inspired by [setwithfriends.com](https://setwithfriends.com).

## How it works

- Create or join a room by visiting `hexwithfriends.com/<room-name>`
- Share the link with a friend
- Red connects top to bottom, blue connects left to right
- First to connect their two sides wins
- Colors swap between games

## Tech

- Go backend with WebSocket multiplayer
- Union-Find for win detection
- Room-based matchmaking
- SVG frontend

## Project structure

```
main.go                     # entrypoint — routing and server startup
internal/
  game/                     # game logic (no networking)
    hex.go                  # board, nodes, neighbors, win detection, moves
    union_find.go           # union-find for connected component tracking
    README.md               # explanation of the win-detection algorithm
  hub/
    hub.go                  # Hub struct, main game loop (register/broadcast/unregister)
    room.go                 # room management — create/lookup by name
    player.go               # WebSocket handler, client read loop
    gameplay.go             # color assignment
static/
  landing.html              # home page
  index.html                # game board UI
  names.json                # random player names
```

## Status

This is an MVP and not production ready.

### TODO

- Add in-game chat (port from [chatroom](../chatroom/))
-Improve the UI — mobile layout, animations, polish
- Handle player disconnects and reconnects gracefully
- Add spectator mode
- Add move take-backs
- Add swap rule
