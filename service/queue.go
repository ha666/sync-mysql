package service

import (
	"time"

	"gitea.com/ha666/sync-mysql/config"
	"github.com/ha666/golibs"
	"github.com/ha666/logs"
)

var (
	sequence       uint64
	receiveQueue0  = make(chan *sqlAndArgs, 10000)
	receiveQueue1  = make(chan *sqlAndArgs, 10000)
	receiveQueue2  = make(chan *sqlAndArgs, 10000)
	receiveQueue3  = make(chan *sqlAndArgs, 10000)
	receiveQueue4  = make(chan *sqlAndArgs, 10000)
	receiveQueue5  = make(chan *sqlAndArgs, 10000)
	receiveQueue6  = make(chan *sqlAndArgs, 10000)
	receiveQueue7  = make(chan *sqlAndArgs, 10000)
	receiveQueue8  = make(chan *sqlAndArgs, 10000)
	receiveQueue9  = make(chan *sqlAndArgs, 10000)
	receiveQueue10 = make(chan *sqlAndArgs, 10000)
	receiveQueue11 = make(chan *sqlAndArgs, 10000)
	receiveQueue12 = make(chan *sqlAndArgs, 10000)
	receiveQueue13 = make(chan *sqlAndArgs, 10000)
	receiveQueue14 = make(chan *sqlAndArgs, 10000)
	receiveQueue15 = make(chan *sqlAndArgs, 10000)
)

//统计队列
func StatisticQueues() {
	for {
		time.Sleep(5 * time.Second)
		data := golibs.NewStringBuilder()
		data.Append("统计队列:")
		if config.Conf.App.ThreadCount > 0 {
			data.Append(0).Append(":").Append(len(receiveQueue0))
		}
		if config.Conf.App.ThreadCount > 1 {
			data.Append("，").Append(1).Append(":").Append(len(receiveQueue1))
		}
		if config.Conf.App.ThreadCount > 2 {
			data.Append("，").Append(2).Append(":").Append(len(receiveQueue2))
		}
		if config.Conf.App.ThreadCount > 3 {
			data.Append("，").Append(3).Append(":").Append(len(receiveQueue3))
		}
		if config.Conf.App.ThreadCount > 4 {
			data.Append("，").Append(4).Append(":").Append(len(receiveQueue4))
		}
		if config.Conf.App.ThreadCount > 5 {
			data.Append("，").Append(5).Append(":").Append(len(receiveQueue5))
		}
		if config.Conf.App.ThreadCount > 6 {
			data.Append("，").Append(6).Append(":").Append(len(receiveQueue6))
		}
		if config.Conf.App.ThreadCount > 7 {
			data.Append("，").Append(7).Append(":").Append(len(receiveQueue7))
		}
		if config.Conf.App.ThreadCount > 8 {
			data.Append("，").Append(8).Append(":").Append(len(receiveQueue8))
		}
		if config.Conf.App.ThreadCount > 9 {
			data.Append("，").Append(9).Append(":").Append(len(receiveQueue9))
		}
		if config.Conf.App.ThreadCount > 10 {
			data.Append("，").Append(10).Append(":").Append(len(receiveQueue10))
		}
		if config.Conf.App.ThreadCount > 11 {
			data.Append("，").Append(11).Append(":").Append(len(receiveQueue11))
		}
		if config.Conf.App.ThreadCount > 12 {
			data.Append("，").Append(12).Append(":").Append(len(receiveQueue12))
		}
		if config.Conf.App.ThreadCount > 13 {
			data.Append("，").Append(13).Append(":").Append(len(receiveQueue13))
		}
		if config.Conf.App.ThreadCount > 14 {
			data.Append("，").Append(14).Append(":").Append(len(receiveQueue14))
		}
		if config.Conf.App.ThreadCount > 15 {
			data.Append("，").Append(15).Append(":").Append(len(receiveQueue15))
		}
		logs.Info(data.ToString())
	}
}
