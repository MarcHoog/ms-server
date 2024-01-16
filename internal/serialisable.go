package internal

import (
	"github.com/MarcHoog/elesia/mnet"
)

type Rates struct {
	Exp   float32
	Drop  float32
	Mesos float32
}

type World struct {
	Conn          mnet.MapleServer
	Icon          byte
	Name, Message string
	Ribbon        byte
	Channels      []Channel
	Rates         Rates
	DefaultRates  Rates
}

type Channel struct {
	Conn        mnet.MapleServer
	IP          []byte
	Port        int16
	MaxPop, Pop int16
}
