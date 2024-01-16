package login

import (
	"log"

	"github.com/MarcHoog/elesia/common/opcode"
	"github.com/MarcHoog/elesia/mnet"
	"github.com/MarcHoog/elesia/mpacket"
)

type MainThread struct {

	// migrating map[mnet.Client]bool
	// worlds []internal.World
	withPin bool
}

func (mt *MainThread) HandleClientEvents(reader mpacket.Reader, Conn mnet.MapleClient) {

	switch reader.ReadByte() {
	case opcode.RecvLoginRequest:
		mt.handleLoginRequest(reader, Conn)
	case opcode.RecvLoginCheckLogin:
		mt.handleGoodLogin(reader, Conn)
		mt.handleTos(reader, Conn)
	case opcode.RecvLoginWorldSelect:
		mt.handleWorldSelect(reader, Conn)
	default:
		log.Println("No handler could be found for the package: ", reader)
	}
}

func (mt *MainThread) handleLoginRequest(reader mpacket.Reader, Conn mnet.MapleClient) {
	log.Println("handleLoginRequest")

	Conn.Send(NewLoginResponse(0x00, 1, 0x00, false, "bobdylan", 0))

}

func (mt *MainThread) handleGoodLogin(reader mpacket.Reader, Conn mnet.MapleClient) {
	Conn.SetLoggedIn(true)

	Conn.Send(NewLoginWorldListing())
	Conn.Send(NewLoginEndWorldList())

}

func (mt *MainThread) handleTos(reader mpacket.Reader, Conn mnet.MapleClient) {
	Conn.Send(NewLoginTOS())
}

func (mt *MainThread) handleWorldSelect(reader mpacket.Reader, Conn mnet.MapleClient) {
	Conn.Send(NewLoginWorldInfo(1, 1))
}
