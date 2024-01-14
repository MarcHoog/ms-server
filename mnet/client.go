package mnet

import (
	"net"

	"github.com/MarcHoog/ms-server/common/constant"
	"github.com/MarcHoog/ms-server/mnet/crypt"
)

type clientConn struct {
	baseConn
}

func NewClientConn(Conn net.Conn, eRecv chan *Event, queueSize int, keySend, keyRecv [4]byte) *clientConn {
	c := &clientConn{}
	c.Conn = Conn

	c.eventRecv = eRecv

	c.cryptSend = crypt.New(keySend, constant.MapleVersion)
	c.cryptRecv = crypt.New(keyRecv, constant.MapleVersion)

	c.Reader = func() {
		Reader(Conn, c.eventRecv, constant.MapleVersion, constant.ClientHeaderSize, c.cryptRecv)
	}

	return c
}

func Reader(Conn net.Conn, eventRecv chan *Event, mapleVersion int16, headerSize int, cryptRecv *crypt.Crypt) {

	// When the reader is started it sends an event that a client has connected Succesfully
	eventRecv <- &Event{Type: MapleEventClientConnected, Conn: Conn}

	header := true
	readSize := headerSize

	for {
		buffer := make([]byte, readSize)

		// Fill the buffer until it's full little bad bo
		if _, err := Conn.Read(buffer); err != nil {
			eventRecv <- &Event{Type: MapleEventClientDisconnect, Conn: Conn}
			break
		}

		if header {
			readSize = crypt.GetPacketLength(buffer)
		} else {
			readSize = headerSize

			if cryptRecv != nil {
				cryptRecv.Decrypt(buffer, true, false)
			}

			eventRecv <- &Event{Type: MapleEventClientPacket, Conn: Conn, Packet: buffer}
		}

		header = !header
	}
}
