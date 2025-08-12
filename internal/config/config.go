package config

type Config struct {
	ListenAddress string
	TLSCertPath   string
	TLSKeyPath    string
	Docker        Docker
	YouTube       YouTube
	Kemono        Kemono
}

type Docker struct {
	Enabled   bool
	GHCR      GHCR
	DockerHub DockerHub
}

type DockerAuth struct {
	Username, Password string
}

type DockerHub struct {
	DockerAuth
}

type GHCR struct {
	DockerAuth
	AppID          int64
	InstallationID int64
	PrivateKey     string
	PrivateKeyPath string
}

type YouTube struct {
	Enabled bool
	APIKey  string
}

type Kemono struct {
	Enabled bool
	Hosts   map[string]string
}

func New() *Config {
	return &Config{
		ListenAddress: ":3000",
		Docker: Docker{
			Enabled: true,
		},
		YouTube: YouTube{
			Enabled: true,
		},
		Kemono: Kemono{
			Enabled: true,
			Hosts:   map[string]string{"kemono": "kemono.su"},
		},
	}
}
