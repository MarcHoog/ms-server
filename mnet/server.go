package mnet

type MapleServer interface {
	MapleConn
}

type serverConnection struct {
	baseConnection
}
