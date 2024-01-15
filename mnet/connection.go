package mnet

import (
	"net"

	"github.com/MarcHoog/ms-server/mnet/crypt"
	"github.com/MarcHoog/ms-server/mpacket"
)

// An interface for a maplestory connection I think???
type MapleConn interface {
	String() string
	Send(mpacket.Packet)
	Cleanup()
}

// What the difference is within a mapleconnection and a base con connect that is something the gods may know
type baseConn struct {
	net.Conn
	eventSend chan mpacket.Packet
	eventRecv chan *Event
	Reader    func()
	closed    bool

	cryptSend *crypt.Crypt
	cryptRecv *crypt.Crypt
}

func (bc *baseConn) Send(p mpacket.Packet) {
	// put a package into the event send channel
	bc.eventSend <- p
}

// Gets the address on the otherside as string probarly the IP
func (bc *baseConn) String() string {
	return bc.Conn.RemoteAddr().String()
}

// Todo(Marc)closes the event send connection but not the event recv connection for
// some reason lol?
func (bc *baseConn) Close() {
	bc.closed = true
	close(bc.eventSend)
	// close(bc.eventRecv)
}

// Gonna be honest boo not really sure what these readers do
// Todo(Marc)
// Somathing but then for a client but why we have a client??
