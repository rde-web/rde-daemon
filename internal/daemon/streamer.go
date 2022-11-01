package daemon

import (
	"bytes"
	"log"
	"net"
	"rde-daemon/internal/config"
)

const (
	bufferSize uint16 = 256
)

type message struct {
	Service string `msgpack:"svc"`
	Method  string `msgpack:"mtd"`
	Payload []byte `msgpack:"pld"`
}

type Streamer struct {
	conn net.Conn
}

func (s *Streamer) Run(errChan *chan error) {
	conn, errDial := net.Dial("tcp", config.Instance.Commutator)
	if errDial != nil {
		*errChan <- errDial
		return
	}
	s.conn = conn
	for {
		var buff []byte = make([]byte, bufferSize)
		readN, errRead := s.conn.Read(buff)
		if errRead != nil {
			*errChan <- errRead
			return
		}
		log.Printf("readed %d bytes", readN)
		var msg message
		if errDecode := decode(buff, &msg); errDecode != nil {
			*errChan <- errDecode
			return
		}
		var reader bytes.Buffer = *bytes.NewBuffer(msg.Payload)
		defer reader.Reset()
		rsp, errCall := NewClient(msg.Service).Post(msg.Method, &reader)
		if errCall != nil {
			rsp = []byte(errCall.Error())
		}
		if _, errWrite := conn.Write(rsp); errWrite != nil {
			*errChan <- errWrite
			return
		}
	}
}

func (s *Streamer) Shutdown() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
