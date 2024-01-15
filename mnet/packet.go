package mnet

import (
	"github.com/MarcHoog/elesia/mpacket"
)

func NewClientHandshakePacket(mapleVersion int16, fromClientKey, toClientKey []byte) mpacket.Packet {
	p := mpacket.NewPacket(0)

	p.WriteInt16(13)
	p.WriteInt16(mapleVersion)
	p.WriteString("")
	p.Append(fromClientKey)
	p.Append(toClientKey)
	p.WriteByte(8)

	return p

}
