package tcpserver

type ServerConfig struct {
	host string
	port uint16
}

func NewServerConfig(host string, port uint16) (*ServerConfig, error) {
	if !IsValidIP(host) {
		return nil, ErrInvalidHost
	}

	return &ServerConfig{
		host: host,
		port: port,
	}, nil
}
