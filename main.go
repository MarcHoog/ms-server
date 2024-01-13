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
	eRecv chan *mnet.Event
	wg    *sync.WaitGroup
	// gameState login.Server
}

func packetClientHandshake(mapleVersion int16, recv, send []byte) mpacket.Packet {
	p := mpacket.NewPacket()

	p.WriteInt16(13)
	p.WriteInt16(mapleVersion)
	p.WriteString("")
	p.Append(recv)
	p.Append(send)
	p.WriteByte(8)

	return p

}

func newLoginServer() *loginServer {
	//	config, dbConfig := loginConfigFromFile(configFile)

	return &loginServer{
		eRecv: make(chan *mnet.Event),
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

	ls.wg.Add(1)
	go ls.acceptNewClientConnections()

	ls.wg.Add(1)
	go ls.acceptNewServerConnections()

	ls.wg.Wait()
}

func (ls *loginServer) acceptNewServerConnections() {
	defer ls.wg.Done()

	listener, err := net.Listen("tcp", "0.0.0.0"+":"+"8485")

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Server listener ready:", "0.0.0.0"+":"+"8485")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Error in accepting client", err)
			close(ls.eRecv)
			return
		}

		serverConn := mnet.NewServer(conn, ls.eRecv, 512)

		go serverConn.Reader()
		go serverConn.Writer()
	}
}

func (ls *loginServer) acceptNewClientConnections() {
	defer ls.wg.Done()

	listener, err := net.Listen("tcp", "0.0.0.0"+":"+"8484")

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Client listener ready:", "0.0.0.0"+":"+"8484")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Error in accepting client", err)
			close(ls.eRecv)
			return
		}

		// ls.gameState.ClientConnected(conn, ls.eRecv, ls.config.PacketQueueSize)
		keySend := [4]byte{}
		rand.Read(keySend[:])
		keyRecv := [4]byte{}
		rand.Read(keyRecv[:])

		client := mnet.NewClient(conn, ls.eRecv, 512, keySend, keyRecv, 0, 0)

		go client.Reader()
		go client.Writer()

		conn.Write(packetClientHandshake(constant.MapleVersion, keyRecv[:], keySend[:]))
	}
}

func main() {
	s := newLoginServer()
	s.run()

}
