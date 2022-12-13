package config

import "time"

var Instance config

func init() {
	Instance = config{
		ProjectPath: "/project",
		SocketsPath: ".",
		// SocketsPath: "/var/rde",
		Commutator:           "localhost:8085",
		BufferSize:           256,
		StreamerReadTimeout:  5,
		StreamerWriteTimeout: 5,
	}
}

type config struct {
	ProjectPath          string
	SocketsPath          string
	Commutator           string
	BufferSize           int
	StreamerReadTimeout  time.Duration
	StreamerWriteTimeout time.Duration
}
