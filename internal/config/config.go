package config

var Instance config

func init() {
	Instance = config{
		ProjectPath: "/project",
		SocketsPath: ".",
		// SocketsPath: "/var/rde",
		Commutator: "localhost:8080",
	}
}

type config struct {
	ProjectPath string
	SocketsPath string
	Commutator  string
}
