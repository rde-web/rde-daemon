package config

var Config config

func init() {
	Config = config{
		ProjectPath: "tmp/",
		// ProjectPath: "/project",
	}
}

type config struct {
	ProjectPath string
}
