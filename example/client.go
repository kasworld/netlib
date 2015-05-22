// Copyright 2015 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	var clientdur = flag.Int("clientdur", 3600, "run sec")
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

	c := Clients{
		*clientdur,
	}
	gogueclient.MultiClient(*connectTo, *count, *rundur, c.clientMain)
}

type Clients struct {
	RunDur int
}

func (c Clients) clientMain(connectTo string, clientnum int, endch chan bool) {
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
	time.Sleep(time.Duration(c.RunDur) * time.Second)
}
