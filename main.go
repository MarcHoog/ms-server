package main

import (
	"crypto/rand"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/MarcHoog/elesia/login"
	"github.com/MarcHoog/elesia/mnet"
	"github.com/MarcHoog/elesia/mpacket"
)

type loginServer struct {
	// config    loginConfig
	// dbConfig  dbConfig
	toMainThread chan *mnet.Event
	wg           *sync.WaitGroup
	mt           login.MainThread
}

func newLoginServer() *loginServer {
	//	config, dbConfig := loginConfigFromFile(configFile)

	return &loginServer{
		toMainThread: make(chan *mnet.Event),
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
	go ls.processEvent()

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

	for {

		Conn, err := listener.Accept()
		log.Println("Client attempting to to connect")

		if err != nil {
			log.Println("Error in accepting client", err)
			close(ls.toMainThread)
			return
		}

		log.Println("Client Connected Succesfully")

		// ls.gameState.ClientConnected(conn, ls.eRecv, ls.config.PacketQueueSize)

		log.Println("Creating 2, 4 Byte long encryption keys")
		fromClientKey := [4]byte{}
		rand.Read(fromClientKey[:])
		toClientKey := [4]byte{}
		rand.Read(toClientKey[:])

		client := mnet.NewClientConn(Conn, ls.toMainThread, 512, fromClientKey, toClientKey)

		go client.Reader()
		go client.Writer()

		log.Println("Creating handshake packet")
		handshakePacket := mnet.NewClientHandshakePacket(28, fromClientKey[:], toClientKey[:])

		log.Println("Sending handshake packet naar de client over de connection")
		Conn.Write(handshakePacket)
	}
}

func (ls *loginServer) processEvent() {
	defer ls.wg.Done()

	for {
		e, ok := <-ls.toMainThread

		if !ok {
			log.Println("Something went wrong with reading from toMainThread channel")
			return
		}

		reader := mpacket.NewReader(&e.Packet, time.Now().Unix())

		switch conn := e.Conn.(type) {
		case mnet.MapleClient:
			switch e.Type {
			case mnet.MapleEventClientConnected:
				log.Println("Client connected")
			case mnet.MapleEventClientDisconnect:
				log.Println("Client disconnected")
			case mnet.MapleEventClientPacket:
				ls.mt.HandleClientEvents(reader, conn)
			}
		}
	}
}

func main() {
	s := newLoginServer()
	s.run()
}
