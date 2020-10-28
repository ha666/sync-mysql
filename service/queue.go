package service

import (
	"time"

	"github.com/ha666/logs"
)

var (
	receiveQueue = make(chan *sqlAndArgs, 100000)
)

//统计队列
func StatisticQueues() {
	for {
		time.Sleep(5 * time.Second)
		logs.Info("统计队列:len=%d", len(receiveQueue))
	}
}
