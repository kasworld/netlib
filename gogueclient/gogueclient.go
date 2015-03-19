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
