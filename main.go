package main

import (
	"gitea.com/ha666/sync-mysql/config"
	"gitea.com/ha666/sync-mysql/service"
)

func main() {
	go service.StartDBWrite()
	go service.StartKafkaWrite()
	go service.StatisticQueues()
	for i := 0; i < config.Conf.App.ThreadCount; i++ {
		go service.StartReadLoop(i)
	}
	select {}
}
