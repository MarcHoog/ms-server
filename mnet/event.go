package mnet

import (
	"net"

	"github.com/MarcHoog/ms-server/mpacket"
)

// TODO(Marc) Event types I think?
const (
	MapleEventClientConnected = iota
	MapleEventClientDisconnect
	MapleEventClientPacket
	MapleEventServerConnected
	MapleEventServerDisconnect
	MapleEventServerPacket
)

// TODO(Marc) WTF actually does this?
// I think that an Event is a single package that has been read from the connection and is ready to be processed.
type Event struct {
	Type   int
	Packet mpacket.Packet
	Conn   net.Conn
}
