package main

import (
	"Goose47/storage/config"
	"Goose47/storage/server"
)

func main() {
	config.Init()
	server.Init()
}
