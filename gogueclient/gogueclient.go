package gogueclient

import (
	"net"
	"time"

	"github.com/kasworld/log"
	"github.com/kasworld/netlib/gogueconn"
	"github.com/kasworld/rand"
)

func NewClientGogueConn(connectTo string) *gogueconn.GogueConn {
	conn, err := net.Dial("tcp", connectTo)
	if err != nil {
		log.Error("client %v", err)
		return nil
	}
	return gogueconn.New(conn)
}

type ClientGoFn func(connectTo string, num int, dur int, endch chan bool)

func MultiClient(connectTo string, count int, rundur int, fn ClientGoFn) {
	rnd := rand.New()
	endch := make(chan bool, count)
	go func() {
		for i := 0; i < count; i++ {
			endch <- true
			go fn(connectTo, i, rnd.Intn(rundur), endch)
			time.Sleep(1 * time.Millisecond)
		}
		for i := count; ; i++ {
			endch <- true
			go fn(connectTo, i, rnd.Intn(rundur), endch)
		}
	}()
	time.Sleep(time.Duration(rundur) * time.Second)
}
