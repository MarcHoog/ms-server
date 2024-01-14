package main

import (
	"crypto/rand"
	"log"
	"net"
	"os"
	"sync"

	"github.com/MarcHoog/ms-server/common/constant"
	"github.com/MarcHoog/ms-server/mpacket"

	"github.com/MarcHoog/ms-server/mnet"
)

type loginServer struct {
	// config    loginConfig
	// dbConfig  dbConfig
	eventRecv chan *mnet.Event
	wg        *sync.WaitGroup
	// gameState login.Server
}

func NewClientHandshakePacket(mapleVersion int16, keyRecv, KeySend []byte) mpacket.Packet {
	p := mpacket.NewPacket()

	p.WriteInt16(13)
	p.WriteInt16(mapleVersion)
	p.WriteString("")
	p.Append(keyRecv)
	p.Append(KeySend)
	p.WriteByte(8)

	return p

}

func newLoginServer() *loginServer {
	//	config, dbConfig := loginConfigFromFile(configFile)

	return &loginServer{
		eventRecv: make(chan *mnet.Event),
		//		config:   config,
		//		dbConfig: dbConfig,
		wg: &sync.WaitGroup{},
	}
}

func (ls *loginServer) run() {
	log.Println("Login Server")

	//	start := time.Now()
	//	nx.LoadFile("Data.nx")
	//	elapsed := time.Since(start)

	//	log.Println("Loaded and parsed Wizet data (NX) in", elapsed)

	//	ls.gameState.Initialise(ls.dbConfig.User, ls.dbConfig.Password, ls.dbConfig.Address, ls.dbConfig.Port, ls.dbConfig.Database, ls.config.WithPin)

	ls.wg.Add(2)
	go ls.acceptNewClientConnections()
	//go ls.processEvent()

	ls.wg.Wait()
}

func (ls *loginServer) acceptNewClientConnections() {
	defer ls.wg.Done()

	listener, err := net.Listen("tcp", "0.0.0.0"+":"+"8484")

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Client listener ready:", "0.0.0.0"+":"+"8484")

	// GO lang for?  nog even uitzoeken
	for {

		// Wacht op een client die probeert te connecten en accepteert deze
		Conn, err := listener.Accept()
		log.Println("Client attempting to to connect")

		if err != nil {
			log.Println("Error in accepting client", err)
			close(ls.eventRecv)
			return
		}

		log.Println("Client Connected Succesfully")

		// ls.gameState.ClientConnected(conn, ls.eRecv, ls.config.PacketQueueSize)

		log.Println("Creating 2, 4 Byte long encryption keys")
		keyRecv := [4]byte{}
		rand.Read(keyRecv[:])
		keySend := [4]byte{}
		rand.Read(keySend[:])

		client := mnet.NewClientConn(Conn, ls.eventRecv, 512, keySend, keyRecv)

		go client.Reader()

		log.Println("Creating handshake packet")
		handshakePacket := NewClientHandshakePacket(constant.MapleVersion, keyRecv[:], keySend[:])

		log.Println("Sending handshake packet naar de client over de connection")
		Conn.Write(handshakePacket)
	}
}

func main() {
	s := newLoginServer()
	s.run()

}

/*
func (ls *loginServer) processEvent() {
	defer ls.wg.Done()

	for {
		select {
		case e, ok := <-ls.eventRecv:

			if !ok {
				log.Println("Stopping event handling due to channel read error")
				return
			}

			switch conn := e.Conn.(type) {
			case mnet.Client:
				switch e.Type {
				case mnet.MapleEventClientConnected:
					log.Println("New client from", conn)
				case mnet.MClientDisconnect:
					log.Println("Client at", conn, "disconnected")
					ls.gameState.ClientDisconnected(conn)
				case mnet.MEClientPacket:
					ls.gameState.HandleClientPacket(conn, mpacket.NewReader(&e.Packet, time.Now().Unix()))
				}
			}
		}
	}
}
*/
