package mnet

import (
	"github.com/MarcHoog/elesia/mpacket"
)

const (
	MapleEventClientConnected = iota
	MapleEventClientDisconnect
	MapleEventClientPacket
	MapleEventServerConnected
	MapleEventServerDisconnect
	MapleEventServerPacket
)

type Event struct {
	Type   int
	Packet mpacket.Packet
	Conn   MapleConn
}
