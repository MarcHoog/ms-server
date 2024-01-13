package mnet

import (
	"net"

	"github.com/MarcHoog/ms-server/common/constant"
	"github.com/MarcHoog/ms-server/mnet/crypt"
	"github.com/MarcHoog/ms-server/mpacket"
)

type Client interface {
	MapleConn

	GetLogedIn() bool
	SetLogedIn(bool)
	GetAccountID() int32
	SetAccountID(int32)
	GetGender() byte
	SetGender(byte)
	GetWorldID() byte
	SetWorldID(byte)
	GetChannelID() byte
	SetChannelID(byte)
	GetAdminLevel() int
	SetAdminLevel(int)
}

type client struct {
	baseConn

	logedIn    bool
	accountID  int32
	gender     byte
	worldID    byte
	channelID  byte
	adminLevel int
}

func NewClient(conn net.Conn, eRecv chan *Event, queueSize int, keySend, keyRecv [4]byte, latency, jitter int) *client {
	c := &client{}
	c.Conn = conn

	c.eventSend = make(chan mpacket.Packet, queueSize)
	c.eventRecv = eRecv

	c.cryptSend = crypt.New(keySend, constant.MapleVersion)
	c.cryptRecv = crypt.New(keyRecv, constant.MapleVersion)

	c.reader = func() {
		clientReader(c, c.eventRecv, constant.MapleVersion, constant.ClientHeaderSize, c.cryptRecv)
	}

	c.interServer = false
	c.latency = latency
	c.pSend = make(chan func(), queueSize*10) // Used only when simulating latency
	return c
}

func (c *client) GetLogedIn() bool {
	return c.logedIn
}

func (c *client) SetLogedIn(logedIn bool) {
	c.logedIn = logedIn
}

func (c *client) GetAccountID() int32 {
	return c.accountID
}

func (c *client) SetAccountID(accountID int32) {
	c.accountID = accountID
}

func (c *client) GetGender() byte {
	return c.gender
}

func (c *client) SetGender(gender byte) {
	c.gender = gender
}

func (c *client) GetWorldID() byte {
	return c.worldID
}

func (c *client) SetWorldID(id byte) {
	c.worldID = id
}

func (c *client) GetChannelID() byte {
	return c.channelID
}

func (c *client) SetChannelID(id byte) {
	c.channelID = id
}

func (c *client) GetAdminLevel() int {
	return c.adminLevel
}

func (c *client) SetAdminLevel(level int) {
	c.adminLevel = level
}
