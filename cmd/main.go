package main

import (
	_ "flag"

	"github.com/VarunSuddala/redis/internal/server"
)

func main() {
	server.RunSyncTCPServer()
}
