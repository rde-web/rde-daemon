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
	Service string      `json:"svc"`
	Method  string      `json:"mtd"`
	Payload interface{} `json:"pld"`
}

type info struct {
	ID string `json:"did"`
}

type errResponse struct {
	Err string `json:"err"`
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
		var msg message
		if errDecode := decode(buff[:readN], &msg); errDecode != nil {
			*errChan <- errDecode
			return
		}
		b, errEncode := encode(msg.Payload)
		if errEncode != nil {
			*errChan <- errEncode
			return
		}
		var reader bytes.Buffer = *bytes.NewBuffer(b)
		defer reader.Reset()
		rsp, errCall := NewClient(msg.Service).Post(msg.Method, &reader)
		if errCall != nil {
			var err errResponse = errResponse{Err: errCall.Error()}
			serialized, errEncodeError := encode(err)
			if errEncodeError != nil {
				log.Printf("Streamer encode error fail: %v", errEncodeError)
			}
			rsp = serialized
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
	return nil
}

func (s *Streamer) Shutdown() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
