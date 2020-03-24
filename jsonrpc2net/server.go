package jsonrpc2net

import (
	"bufio"
	"github.com/maoxs2/go-jsonrpc2"
	"log"
	"net"
	"strings"
)

type Server struct {
	net      string
	listener net.Listener

	handlerMap map[string]JsonRpcHandler
}

func NewServer(network string, addr string) (*Server, error) {
	listener, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		net:      network,
		listener: listener,
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
			log.Println(err)
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
			log.Println(err)
			continue
		}

		if jsonrpc2.IsBatchMarshal(raw) {
			s.onBatchMsg(conn, raw)
		} else {
			s.onSingleMsg(conn, raw)
		}
	}
}
