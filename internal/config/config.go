package config

type Config struct {
	Service  Service
	Upstream Upstream
}

type Service struct {
	Port             uint16
	MaxMessageLength uint
}
type Upstream struct {
	Hosts       []string
	Connections uint
}

func LoadFromFile(filePath string) (*Config, error) {
	panic("not implemented")
}
