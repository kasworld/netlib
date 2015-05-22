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
	PT_gob = uint8(iota)
	PT_json
)

type IDecoder interface {
	Decode(v interface{}) error
}
type IEncoder interface {
	Encode(v interface{}) error
}

type GogueConn struct {
	conn       net.Conn
	packettype uint8
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
