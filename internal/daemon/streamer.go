package daemon

import (
	"bytes"
	"log"
	"net"
	"rde-daemon/internal/config"
	"time"
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
	s.remeoteSetup()

	for {
		buff, errRead := s.read()
		if errRead != nil {
			*errChan <- errRead
			return
		}
		var msg message
		if errDecode := decode(buff, &msg); errDecode != nil {
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

func (s *Streamer) remeoteSetup() error {
	var setupMsg info
	setupMsg.ID = "daemon" //@todo
	data, errEncode := encode(setupMsg)
	if errEncode != nil {
		return errEncode
	}
	return s.write(data)
}

func (s *Streamer) write(data []byte) error {
	var batchSize int = config.Instance.BufferSize
	if len(data) <= batchSize {
		_, err := s.conn.Write(data)
		return err
	}
	for {
		if errSetDeadline := s.conn.SetWriteDeadline(
			time.Now().Add(config.Instance.StreamerWriteTimeout * time.Second),
		); errSetDeadline != nil {
			return errSetDeadline
		}
		var size int = batchSize
		var delta int = len(data) - size
		if delta < 0 {
			size += delta
		}
		var payload []byte = data[:size]
		writedN, err := s.conn.Write(payload)
		if err != nil {
			return err
		}
		if delta <= 0 {
			break
		}
		data = data[writedN:]
	}
	_, errSendEOF := s.conn.Write([]byte{0}) // EOF
	return errSendEOF
}

func (s *Streamer) read() ([]byte, error) {
	var result []byte = make([]byte, 0)
	for {
		var buff []byte = make([]byte, config.Instance.BufferSize)

		if errSetDeadline := s.conn.SetReadDeadline(
			time.Now().Add(config.Instance.StreamerReadTimeout * time.Second),
		); errSetDeadline != nil {
			return nil, errSetDeadline
		}

		readedN, err := s.conn.Read(buff)
		if err != nil {
			return nil, err
		}
		if readedN == 1 && buff[0] == 0 {
			break
		}
		result = append(result, buff[:readedN]...)
	}
	return result, nil
}

func (s *Streamer) Shutdown() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
