package mnet

import (
	"net"

	"github.com/MarcHoog/elesia/mnet/crypt"
	"github.com/MarcHoog/elesia/mpacket"
)

type clientConn struct {
	baseConnection

	headerSize int32

	loggedIn  bool
	accountID int32
	worldID   byte
	channelID byte
}

func NewClientConn(conn net.Conn, toMainThread chan *Event, queueSize int, fromClientKey, toClientKey [4]byte) *clientConn {

	cc := &clientConn{}
	cc.Conn = conn

	cc.toClient = make(chan mpacket.Packet, queueSize)
	cc.toMainThread = toMainThread

	cc.fromClientCrypt = crypt.NewCrypt(fromClientKey)
	cc.toClientCrypt = crypt.NewCrypt(toClientKey)

	cc.headerSize = 4
	cc.active = true

	return cc

}

func (cc *clientConn) Reader() {
	cc.toMainThread <- &Event{MapleEventClientConnected, nil, cc.toClient}

	header := true
	readSize := cc.headerSize

	for {

		packet := mpacket.NewPacket(readSize)

		if _, err := cc.Conn.Read(packet); err != nil {
			cc.toMainThread <- &Event{MapleEventClientDisconnect, nil, cc.toClient}
			return
		}

		if header {
			readSize = int32(crypt.GetPacketLength(packet))
		} else {
			readSize = cc.headerSize

			if cc.fromClientCrypt != nil {
				// our client does have AES turned on
				cc.fromClientCrypt.Decrypt(packet, true, false)
			}

			cc.toMainThread <- &Event{MapleEventClientPacket, packet, cc.toClient}

		}

		header = !header

	}

}
