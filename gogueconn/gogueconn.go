package gogueconn

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/kasworld/log"
)

var PK_type = PT_gob

const (
	PT_gob = iota
	PT_json
)

type IDecoder interface {
	Decode(v interface{}) error
}
type IEncoder interface {
	Encode(v interface{}) error
}

// func printjson(v interface{}) {
// 	data, err := json.MarshalIndent(v, "", "    ")
// 	if err == nil {
// 		log.Info("%v", string(data))
// 	} else {
// 		log.Error("Marshal err ")
// 	}
// }

type GogueConn struct {
	conn       net.Conn
	packettype int
	enc        IEncoder
	dec        IDecoder
}

func New(conn net.Conn) *GogueConn {
	c := GogueConn{
		conn:       conn,
		packettype: PK_type,
	}
	switch c.packettype {
	default:
		return nil
	case PT_gob:
		c.dec = gob.NewDecoder(c.conn)
		c.enc = gob.NewEncoder(c.conn)
	case PT_json:
		c.dec = json.NewDecoder(c.conn)
		c.enc = json.NewEncoder(c.conn)
	}
	return &c
}
func (c *GogueConn) Send(v interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("Connection error %v", e)
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()
	switch c.packettype {
	default:
		return errors.New("Unknown packet type")
	case PT_gob:
		err = c.enc.Encode(v)
	case PT_json:
		err = c.enc.Encode(v)
	}
	return
}
func (c *GogueConn) Recv(v interface{}) (err error) {
	switch c.packettype {
	default:
		return errors.New("Unknown packet type")
	case PT_gob:
		err = c.dec.Decode(v)
	case PT_json:
		err = c.dec.Decode(v)
	}
	return
}

func (c *GogueConn) Close() {
	c.conn.Close()
}

// func (c *GogueConn) RunRecv(v interface{}) <-chan interface{} {
// 	rtn := make(chan interface{})
// 	go func() {
// 		c.Recv(v)
// 		rtn <- v
// 	}()
// 	return rtn
// }

// func (c *GogueConn) RunRecv2(packet interface{}, rtn chan<- error) {
// 	go func() {
// 		rtn <- c.Recv(packet)
// 	}()
// }
