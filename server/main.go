package main

import (
	"github.com/MuriloUnten/udp-file-transfer/protocol"
	"fmt"
	"net"
)

type Server struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	quitChannel chan struct{}
	// TODO maybe add Directory for files
}

func NewServer(address string, port int) *Server {
	addr := &net.UDPAddr{
		IP: net.ParseIP(address),
		Port: port,
	}

	return &Server{
		addr: addr,
		quitChannel: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}
	s.conn = conn
	defer s.conn.Close()
	
	s.Serve()

	return nil
}

func (s *Server) Serve() {
	fmt.Printf("Server Running at %s:%d\n", string(s.addr.IP), s.addr.Port)

	buf := make([]byte, 2048)

	for {
		n, remoteAddress, err := s.conn.ReadFromUDP(buf)
		fmt.Printf("Received %d bytes from %s:%d\n", n, string(remoteAddress.IP), remoteAddress.Port)

		if err != nil {
			fmt.Println("Error reading packet:", err)
			continue
		}

		if n == 0 {
			fmt.Println("Ignoring empty request")
			continue
		}

		response := protocol.Response{
			SegmentNumber: 1,
			Error: nil,
			Content: "Hello from server!",
		}

		fmt.Printf("Response to %s:%d :\n    %s\n", string(remoteAddress.IP), remoteAddress.Port, response.Content)
		_, err = s.conn.WriteToUDP(response.Encode(), remoteAddress)
		if err != nil {
			fmt.Println("Error sending reponse:", err)
			continue
		}
	}
}

func main() {
	s := NewServer("127.0.0.1", 1337)
	s.Start()
}
