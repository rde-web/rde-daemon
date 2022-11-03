package daemon

import (
	"bytes"
	"fmt"
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

const (
	remoteSetupOK   uint8 = iota
	rmeoteSetupFail uint8 = iota
)

type info struct {
	ID string `msgpack:"daemon_id"`
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
	s.rmeoteSetup()

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

func (s *Streamer) rmeoteSetup() error {
	var setupMsg info
	setupMsg.ID = "daemon" //@todo
	data, errEncode := encode(setupMsg)
	if errEncode != nil {
		return errEncode
	}
	if _, errSetup := s.conn.Write(data); errSetup != nil {
		return errSetup
	}
	var buf []byte = make([]byte, 1)
	readN, errRead := s.conn.Read(buf)
	if errRead != nil {
		return errRead
	}
	log.Printf("reded %d bytes", readN)
	switch buf[0] {
	case remoteSetupOK:
		return nil
	default:
		return fmt.Errorf("comutator returned non-ok response (%d)", buf[0])
	}
}

func (s *Streamer) Shutdown() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
