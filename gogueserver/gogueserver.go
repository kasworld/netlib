package gogueserver

import (
	"net"
	"time"

	"github.com/kasworld/log"
	"github.com/kasworld/netlib/gogueconn"
)

type ServerGoFn func(conn *gogueconn.GogueConn, clientQueue <-chan bool)

func TCPServer(listenString string, connNum int, connThrottle int, runConn ServerGoFn) {
	log.Info("Start server %v", listenString)
	// concurrent connection count control
	clientQueue := make(chan bool, connNum)

	listener, err := net.Listen("tcp", listenString)
	if err != nil {
		log.Error("%v", err)
		return
	}
	defer listener.Close()

	for {
		time.Sleep(time.Duration(connThrottle) * time.Millisecond)
		clientQueue <- true
		conn, err := listener.Accept()
		if err != nil {
			log.Error("%v", err)
		} else {
			go runConn(gogueconn.NewGogueConn(conn), clientQueue)
		}
	}
}
