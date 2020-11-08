package jsonrpc2net

import (
	"bufio"
	"net"
	"strings"

	"github.com/c0mm4nd/go-jsonrpc2"
)

type Logger interface {
	Debug(...interface{})
	Error(...interface{})
}

type Server struct {
	net      string
	listener net.Listener
	logger   Logger

	handlerMap map[string]JsonRpcHandler
}

type ServerConfig struct {
	Network string
	Addr    string
	Logger  Logger
}

func NewServer(config ServerConfig) (*Server, error) {
	var logger = config.Logger
	if logger == nil {
		logger = new(SimpleLogger)
	}

	listener, err := net.Listen(config.Network, config.Addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		net:      config.Network,
		listener: listener,
		logger:   logger,
	}, nil
}

func (s *Server) Serve() {
	if strings.Contains(s.net, "udp") {
		s.serveUDP()
	} else {
		s.serveTCP()
	}
}

func (s *Server) serveUDP() {
	panic("unsupported yet")
	//for {
	//	pack, err := s.listener.Accept()
	//	if err != nil {
	//		log.Println(err)
	//		continue
	//	}
	//
	//	raw, err := ioutil.ReadAll(pack.(*net.UDPConn))
	//	if err != nil {
	//		log.Println(err)
	//		continue
	//	}
	//
	//	if jsonrpc2.IsBatchMarshal(raw) {
	//		s.onBatchMsg(net.D, raw)
	//	}
	//}
}

func (s *Server) serveTCP() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error(err)
			continue
		}

		go s.handleTCPConn(conn.(*net.TCPConn))
	}
}

func (s *Server) handleTCPConn(conn *net.TCPConn) {
	r := bufio.NewReader(conn)
	for {
		raw, err := r.ReadBytes('\n')
		if err != nil {
			s.logger.Error(err)
			continue
		}

		if jsonrpc2.IsBatchMarshal(raw) {
			s.onBatchMsg(conn, raw)
		} else {
			s.onSingleMsg(conn, raw)
		}
	}
}
