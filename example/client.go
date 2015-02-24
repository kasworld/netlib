package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/kasworld/log"
	"github.com/kasworld/netlib/gogueclient"
)

type DataPacket struct {
	Cmd string
	Arg int
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var connectTo = flag.String("connectTo", "localhost:6666", "server ip/port")
	var count = flag.Int("count", 1000, "client count")
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

	gogueclient.MultiClient(*connectTo, *count, *rundur, clientMain)
}

func clientMain(connectTo string, clientnum int, dur int, endch chan bool) {
	defer func() {
		<-endch
	}()
	clientname := fmt.Sprintf("test%v", clientnum)
	go func() {
		gconn := gogueclient.NewClientGogueConn(connectTo)
		if gconn == nil {
			return
		}
		defer gconn.Close()
		for i := 0; ; i++ {
			tosend := DataPacket{
				clientname, i,
			}
			err := gconn.Send(tosend)
			if err != nil {
				if err.Error() != "EOF" {
					log.Error("send %v", err)
				}
				break
			}

			var torecv DataPacket
			err = gconn.Recv(&torecv)
			if err != nil {
				if err.Error() != "EOF" {
					log.Error("recv %v", err)
				}
				break
			}
		}
	}()
	time.Sleep(time.Duration(dur) * time.Second)
}
