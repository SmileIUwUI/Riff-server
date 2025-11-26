package tcpserver

type TCPServer struct {
	config ServerConfig
}

func NewTCPServer(config ServerConfig) *TCPServer {
	return &TCPServer{
		config: config,
	}
}
