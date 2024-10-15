package utils

import (
	"slices"
	"time"

	"github.com/doptime/logger"
	"github.com/go-ping/ping"
)

var pingTaskServers = []string{}

func PingServer(domain string, skipRepeat bool) {
	var (
		pinger *ping.Pinger
		err    error
	)
	if skipRepeat && slices.Index(pingTaskServers, domain) != -1 {
		return
	}
	pingTaskServers = append(pingTaskServers, domain)

	if pinger, err = ping.NewPinger(domain); err != nil {
		logger.Info().AnErr("ERROR NewPinger", err).Send()
	}
	pinger.Count = 4
	pinger.Timeout = time.Second * 10
	pinger.OnRecv = func(pkt *ping.Packet) {}

	pinger.OnFinish = func(stats *ping.Statistics) {
		logger.Info().Str("Ping", domain).Str("Addr", stats.Addr).Any("Sent", stats.PacketsSent).Any("Recv", stats.PacketsRecv).Any("Loss", stats.PacketLoss).Any("Min", stats.MinRtt).Any("Avg", stats.MaxRtt).Any("Std", stats.StdDevRtt).Send()
	}
	go func() {
		if err := pinger.Run(); err != nil {
			logger.Info().AnErr("ERROR Ping", err).Send()
		}
	}()
}
