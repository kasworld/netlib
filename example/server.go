package main

import (
	"flag"
	"os"
	"runtime/pprof"
	"time"

	"github.com/kasworld/actionstat"
	"github.com/kasworld/log"
	"github.com/kasworld/netlib/gogueconn"
	"github.com/kasworld/netlib/gogueserver"
)

type DataPacket struct {
	Cmd string
	Arg int
}

var stat *actionstat.ActionStat

func main() {
	stat = actionstat.NewActionStat()
	var listenFrom = flag.String("listenFrom", ":6666", "server ip/port")
	var connCount = flag.Int("count", 1000, "connection count")
	var connThrottle = flag.Int("throttle", 10, "connection throttle")
	var rundur = flag.Int("rundur", 3600, "run sec")
	var profilefilename = flag.String("pfilename", "", "profile filename")
	flag.Parse()

	if *profilefilename != "" {
		f, err := os.Create(*profilefilename)
		if err != nil {
			log.Fatalf("profile %v", err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	go gogueserver.TCPServer(*listenFrom, *connCount, *connThrottle, servergo)

	go func() {
		timerInfoCh := time.Tick(time.Duration(1000) * time.Millisecond)
		for {
			select {
			case <-timerInfoCh:
				log.Info("%v", stat)
				stat.UpdateLap()
			}
		}
	}()
	time.Sleep(time.Duration(*rundur) * time.Second)
}

func servergo(gconn *gogueconn.GogueConn, clientQueue <-chan bool) {
	defer gconn.Close()
	// log.Info("client connected")
	for {
		var rdata DataPacket
		err := gconn.Recv(&rdata)
		if err != nil {
			if err.Error() != "EOF" {
				log.Error("recv %v", err)
			}
			break
		}
		err = gconn.Send(&rdata)
		if err != nil {
			if err.Error() != "EOF" {
				log.Error("send %v", err)
			}
			break
		}
		stat.Inc()
	}
}
