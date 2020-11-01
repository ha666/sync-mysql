package main

import (
	"gitea.com/ha666/sync-mysql/service"
)

func main() {
	go service.StartDBWrite()
	go service.StartKafkaWrite()
	go service.StatisticQueues()
	service.StartRead()
}
