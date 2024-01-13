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
	reader    func()
	closed    bool

	cryptSend *crypt.Maple
	cryptRecv *crypt.Maple

	interServer bool

	latency int
	pSend   chan func()
}

func (bc *baseConn) Writer() {

	for {
		p, ok := <-bc.eventSend
		if !ok {
			return
		}

		tmp := make(mpacket.Packet, len(p))
		copy(tmp, p)

		// TODO(Marc) I think
		// Encrypt the packet if we have a crypter
		if bc.cryptSend != nil {
			bc.cryptSend.Encrypt(tmp, true, false)
		}

		// TODO(Marc) No clue what this means :D
		if bc.interServer {
			tmp[0] = byte(len(tmp) - 1)

			bc.Conn.Write(tmp)

			// TODO(Marc) Removed all the jitter latency stuff :D cause who cares)

			/*

				if bc.latency > 0 {
					now := time.Now().UnixNano()
					sendTime := now + int64(rand.Intn(bc.jitter)+bc.latency)*1000000
					bc.pSend <- func() {
						now := time.Now().UnixNano()
						delta := sendTime - now

						if delta > 0 {
							time.Sleep(time.Duration(delta))
						}

						bc.Conn.Write(tmp)
					}
				} else {
					bc.Conn.Write(tmp)
				}
			*/

		}

	}

}

func (bc *baseConn) Send(p mpacket.Packet) {
	// If the connection has been closed we don't wanne send anything
	if bc.closed {
		return
	}

	// put a package into the event send channel
	bc.eventSend <- p
}

func (bc *baseConn) Reader() {
	bc.reader()
}

// Gets the address on the otherside as string probarly the IP
func (bc *baseConn) String() string {
	return bc.Conn.RemoteAddr().String()
}

// Todo(Marc)closes the event send connection but not the event recv connection for
// some reason lol?
func (bc *baseConn) Cleanup() {
	bc.closed = true
	close(bc.eventSend)
	// close(bc.eventRecv)
}

// Gonna be honest boo not really sure what these readers do
// Todo(Marc)
// Somathing but then for a client but why we have a client??
// we are  a server? idk?
// I think that the client is for testing?
func clientReader(conn net.Conn, eRecv chan *Event, mapleVersion int16, headerSize int, cryptRecv *crypt.Maple) {
	eRecv <- &Event{Type: MapleEventClientConnected, Conn: conn}

	header := true
	readSize := headerSize

	for {
		buffer := make([]byte, readSize)

		if _, err := conn.Read(buffer); err != nil {
			eRecv <- &Event{Type: MapleEventClientDisconnect, Conn: conn}
			break
		}

		if header {
			readSize = crypt.GetPacketLength(buffer)
		} else {
			readSize = headerSize

			if cryptRecv != nil {
				cryptRecv.Decrypt(buffer, true, false)
			}

			eRecv <- &Event{Type: MapleEventClientPacket, Conn: conn, Packet: buffer}
		}

		header = !header
	}
}

func serverReader(conn net.Conn, eRecv chan *Event, headerSize int) {
	// So I think this fucks something into the EventRecv channel at start up
	eRecv <- &Event{Type: MapleEventServerConnected, Conn: conn}

	// SO this is clearly a stream of data that's why it does the header first
	// and then later disables the header
	header := true
	readSize := headerSize

	// And then it keeps looping over all stuff that ithas gotten
	for {
		buffer := make([]byte, readSize)

		if _, err := conn.Read(buffer); err != nil {
			eRecv <- &Event{Type: MapleEventServerDisconnect, Conn: conn}
			break
		}

		if header {
			readSize = int(buffer[0])
		} else {
			readSize = headerSize
			eRecv <- &Event{Type: MapleEventServerPacket, Conn: conn, Packet: buffer}
		}

		header = !header
	}
}
