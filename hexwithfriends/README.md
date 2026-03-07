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
- Single-file SVG frontend

## Status

This is an MVP and not production ready.

### TODO

- Add in-game chat (port from [chatroom](../chatroom/))
- Clean up the codebase
- Improve the UI — mobile layout, animations, polish
- Handle player disconnects and reconnects gracefully
- Add spectator mode
- Add move take-backs
- Add swap rule
