package login

import (
	"github.com/MarcHoog/elesia/common/opcode"
	"github.com/MarcHoog/elesia/mpacket"
)

func NewLoginResponse(result byte, userID int32, gender byte, isAdmin bool, username string, isBanned int) mpacket.Packet {
	pac := mpacket.CreateWithOpcode(opcode.SendLoginResponse)
	pac.WriteByte(result)
	pac.WriteByte(0x00)
	pac.WriteInt32(0)

	if result <= 0x01 {
		pac.WriteInt32(userID)
		pac.WriteByte(gender)
		// pac.WriteByte(isAdmin)
		pac.WriteBool(isAdmin)
		pac.WriteByte(0x01)
		pac.WriteString(username)
	} else if result == 0x02 {
		pac.WriteByte(byte(isBanned))
		pac.WriteInt64(0) // Expire time, for now let set this to epoch
	}

	pac.WriteInt64(0)
	pac.WriteInt64(0)
	pac.WriteInt64(0)

	return pac
}
