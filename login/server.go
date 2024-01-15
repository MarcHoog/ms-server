package login

import (
	"log"

	"github.com/MarcHoog/elesia/common/opcode"
	"github.com/MarcHoog/elesia/mpacket"
)

type MainThread struct {

	// migrating map[mnet.Client]bool
	// worlds []internal.World
	withPin bool
}

func (mt *MainThread) HandleClientEvents(reader mpacket.Reader, toClient chan<- mpacket.Packet) {

	switch reader.ReadByte() {
	case opcode.RecvLoginRequest:
		mt.handleLoginRequest(reader, toClient)

	default:
		log.Println("No handler could be found for the package: ", reader)
	}
}

func (mt *MainThread) handleLoginRequest(reader mpacket.Reader, toClient chan<- mpacket.Packet) {
	log.Println("handleLoginRequest")

	toClient <- NewLoginResponse(0x00, 1, 0x00, false, "bobdylan", 0)

}
