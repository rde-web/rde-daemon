package config

var Config config

func init() {
	Config = config{
		ProjectPath: "/project",
	}
}

type config struct {
	ProjectPath string
}
