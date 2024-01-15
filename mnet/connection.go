package mnet

import (
	"log"
	"net"

	"github.com/MarcHoog/elesia/mnet/crypt"
	"github.com/MarcHoog/elesia/mpacket"
)

type baseConnection struct {
	net.Conn

	toMainThread chan *Event
	toClient     chan mpacket.Packet

	fromClientCrypt *crypt.Crypt
	toClientCrypt   *crypt.Crypt

	active bool
}

func (bc *baseConnection) String() string {
	return bc.RemoteAddr().String()
}

func (bc *baseConnection) Cleanup() {
	bc.active = false
	close(bc.toClient)
}

func (bc *baseConnection) Writer() {

	for {
		p, ok := <-bc.toClient

		if !ok {
			log.Println("Something went wrong with reading from toClient channel")
			return
		}

		if bc.toClientCrypt != nil {
			const maple = true
			const aes = false
			bc.toClientCrypt.Encrypt(p, maple, aes)
		}

		if _, err := bc.Conn.Write(p); err != nil {
			log.Println("Something went wrong with writing to over Connection:", err)
			return
		}
	}
}
