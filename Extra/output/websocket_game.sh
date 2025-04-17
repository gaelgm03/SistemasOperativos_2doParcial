#!/bin/sh
echo -ne '\033c\033]0;WebSocket Multiplayer Demo\a'
base_path="$(dirname "$(realpath "$0")")"
"$base_path/websocket_game.x86_64" "$@"
