package mnet

import (
	"net"

	"github.com/MarcHoog/ms-server/mpacket"
)

type Server interface {
	MapleConn
}

// TODO(Marc) Why rename it if we even just use it as
// a basic CONNECTION it MAKES NO SENCE
type server struct {
	baseConn
}

func NewServer(conn net.Conn, eRecv chan *Event, queueSize int) *server {
	s := &server{}
	s.Conn = conn

	s.eventSend = make(chan mpacket.Packet, queueSize)
	s.eventRecv = eRecv

	s.reader = func() {
		serverReader(s, s.eventRecv, 1)
	}

	s.interServer = true

	return s
}
