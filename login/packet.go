package login

import (
	"github.com/MarcHoog/elesia/common/opcode"
	"github.com/MarcHoog/elesia/mpacket"
)

func NewLoginResponse(result byte, userID int32, gender byte, isAdmin bool, username string, isBanned int) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendLoginResponse)
	p.WriteByte(result)
	p.WriteByte(0x00)
	p.WriteInt32(0)

	if result <= 0x01 {
		p.WriteInt32(userID)
		p.WriteByte(gender)
		// pac.WriteByte(isAdmin)
		p.WriteBool(isAdmin)
		p.WriteByte(0x01)
		p.WriteString(username)
	} else if result == 0x02 {
		p.WriteByte(byte(isBanned))
		p.WriteInt64(0) // Expire time, for now let set this to epoch
	}

	p.WriteInt64(0)
	p.WriteInt64(0)
	p.WriteInt64(0)

	return p
}

func NewLoginTOS() mpacket.Packet {
	p := mpacket.CreateWithOpcode(0x00)
	p.WriteByte(23)
	p.WriteByte(0x00)
	p.WriteInt32(0)

	return p
}

func NewLoginWorldListing() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendLoginWorldList)
	// world stuff
	p.WriteByte(1)           // world id
	p.WriteString("Testpia") // World name -
	p.WriteByte(3)           // Ribbon on world - 0 = normal, 1 = event, 2 = new, 3 = hot
	p.WriteString("Testpia") // World event message
	p.WriteByte(0)           // ? exp event notification?
	p.WriteByte(byte(1))     // number of channels

	// Channel Stuff
	p.WriteString("testpia" + "-" + "1") // channel name
	p.WriteInt32(32)                     //  Channel Population population
	p.WriteByte(1)
	p.WriteByte(1) // channel id
	p.WriteByte(0) // ?

	return p
}

func NewLoginEndWorldList() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendLoginWorldList)
	p.WriteByte(0xFF)

	return p
}

func NewLoginWorldInfo(warning, population byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendLoginWorldMeta)
	p.WriteByte(warning)    // Warning - 0 = no warning, 1 - high amount of concurent users, 2 = max users in world
	p.WriteByte(population) // Population marker - 0 = No maker, 1 = Highly populated, 2 = over populated

	return p
}
