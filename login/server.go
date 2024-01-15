package login

import (
	"github.com/MarcHoog/elesia/common/opcode"
	"github.com/MarcHoog/elesia/mpacket"
)

type MainThread struct {

	// migrating map[mnet.Client]bool
	// worlds []internal.World
	withPin bool
}

func (mt *MainThread) HandleClientEvents(reader mpacket.Reader, toClient chan<- *mpacket.Packet) {

	switch reader.ReadByte() {
	case opcode.RecvLoginRequest:
		mt.handleLoginRequest(reader, toClient)

	}
}

func (mt *MainThread) handleLoginRequest(reader mpacket.Reader, toClient chan<- *mpacket.Packet) {

}
