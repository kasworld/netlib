package gogueclient

import (
	"fmt"
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
	return gogueconn.NewGogueConn(conn)
}

type ClientGoFn func(clientname string, connectTo string) chan<- bool

func MultiClient(connectTo string, count int, rundur int, fn ClientGoFn) {
	rnd := rand.New()
	endch := make(chan bool, count)
	go func() {
		for i := 0; i < count; i++ {
			endch <- true
			go Client(connectTo, fmt.Sprintf("test%v", i), rnd.Intn(rundur), endch, fn)
			time.Sleep(1 * time.Millisecond)
		}
		for i := count; ; i++ {
			endch <- true
			go Client(connectTo, fmt.Sprintf("test%v", i), rnd.Intn(rundur), endch, fn)
		}
	}()
	time.Sleep(time.Duration(rundur) * time.Second)
}

func Client(connectTo string, name string, dur int, endch chan bool, fn ClientGoFn) {
	defer func() {
		<-endch
	}()
	quit := fn(name, connectTo)
	time.Sleep(time.Duration(dur) * time.Second)
	quit <- true
}
