package main

import (
	"gitea.com/ha666/sync-mysql/service"
)

func main() {
	go service.StartWrite()
	go service.StartRead()
	select {}
}
