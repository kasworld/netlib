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

package gogueclient

import (
	"net"
	"time"

	"github.com/kasworld/log"
	"github.com/kasworld/netlib/gogueconn"
)

func NewClientGogueConn(connectTo string) *gogueconn.GogueConn {
	conn, err := net.Dial("tcp", connectTo)
	if err != nil {
		log.Error("client %v", err)
		return nil
	}
	return gogueconn.New(conn)
}

type ClientGoFn func(connectTo string, num int, endch chan bool)

func MultiClient(connectTo string, count int, rundur int, fn ClientGoFn) {
	go func() {
		endch := make(chan bool, count)
		for i := 0; ; i++ {
			endch <- true
			go fn(connectTo, i, endch)
			time.Sleep(1 * time.Millisecond)
		}
	}()
	time.Sleep(time.Duration(rundur) * time.Second)
}
