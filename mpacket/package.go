package mpacket

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Packet []byte
type Opcode byte

// TODO(Marc) Notes why we keep adding to poss if we are not reusing the value
// A bit weird

// Defining

func NewPacket() Packet {
	return make(Packet, 0)
}

func (p *Packet) Append(data []byte) {
	*p = append(*p, data...)
}

// Writing

func (p *Packet) WriteByte(data byte) {
	*p = append(*p, data)
}

func (p *Packet) WriteBytes(data []byte) {
	p.Append(data)
}

func (p *Packet) WriteBool(data bool) {
	if data {
		*p = append(*p, 0x1)
	} else {
		*p = append(*p, 0x0)
	}
}

func (p *Packet) WriteUint16(data uint16) {
	*p = append(*p, byte(data), byte(data>>8))
}

func (p *Packet) WriteUint32(data uint32) {
	*p = append(*p, byte(data), byte(data>>8), byte(data>>16), byte(data>>24))
}

func (p *Packet) WriteUint64(data uint64) {
	*p = append(*p, byte(data), byte(data>>8), byte(data>>16), byte(data>>24),
		byte(data>>32), byte(data>>40), byte(data>>48), byte(data>>56))
}

func (p *Packet) WriteFloat32(data float32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], math.Float32bits(data))
	*p = append(*p, b[:]...)
}

func (p *Packet) WriteInt8(data int8) {
	*p = append(*p, byte(data))
}

func (p *Packet) WriteInt16(data int16) { p.WriteUint16(uint16(data)) }

func (p *Packet) WriteInt32(data int32) { p.WriteUint32(uint32(data)) }

func (p *Packet) WriteInt64(data int64) { p.WriteUint64(uint64(data)) }

func (p *Packet) WriteString(str string) {
	p.WriteUint16(uint16(len(str)))
	p.WriteBytes([]byte(str))
}

func (p *Packet) WritePaddedString(str string, number int) {
	if len(str) > number {
		p.WriteBytes([]byte(str)[:number])
	} else {
		p.WriteBytes([]byte(str))
		p.WriteBytes(make([]byte, number-len(str)))
	}
}

// TODO(Marc) Why even have size if we use len everywhere?
func (p *Packet) Size() int {
	return int(len(*p))
}

func (p Packet) String() string {
	return fmt.Sprintf("[Packet] (%d) : % X", len(p), string(p))

}

// Reading

func (p *Packet) readByte(pos *int) byte {
	r := byte((*p)[*pos])
	*pos++
	return r
}

func (p *Packet) readBool(pos *int) bool {
	r := ((*p)[*pos])
	*pos++
	return r != 0
}

func (p *Packet) readInt8(pos *int) int8 {
	r := int8((*p)[*pos])
	*pos++
	return r
}

func (p *Packet) readBytes(pos *int, length int) []byte {
	r := []byte((*p)[*pos : *pos+length])
	*pos += length
	return r

}

func (p *Packet) readInt16(pos *int) int16 {
	return int16(p.readByte(pos)) | (int16(p.readByte(pos)) << 8)
}

func (p *Packet) readInt32(pos *int) int32 {
	return int32(p.readByte(pos)) |
		int32(p.readByte(pos))<<8 |
		int32(p.readByte(pos))<<16 |
		int32(p.readByte(pos))<<24
}

func (p *Packet) readInt64(pos *int) int64 {
	return int64(p.readByte(pos)) |
		int64(p.readByte(pos))<<8 |
		int64(p.readByte(pos))<<16 |
		int64(p.readByte(pos))<<24 |
		int64(p.readByte(pos))<<32 |
		int64(p.readByte(pos))<<40 |
		int64(p.readByte(pos))<<48 |
		int64(p.readByte(pos))<<56
}

func (p *Packet) readUint16(pos *int) uint16 {
	return uint16(p.readByte(pos)) | (uint16(p.readByte(pos)) << 8)
}

func (p *Packet) readUint32(pos *int) uint32 {
	return uint32(p.readByte(pos)) |
		uint32(p.readByte(pos))<<8 |
		uint32(p.readByte(pos))<<16 |
		uint32(p.readByte(pos))<<24
}

func (p *Packet) readUint64(pos *int) uint64 {
	return uint64(p.readByte(pos)) |
		uint64(p.readByte(pos))<<8 |
		uint64(p.readByte(pos))<<16 |
		uint64(p.readByte(pos))<<24 |
		uint64(p.readByte(pos))<<32 |
		uint64(p.readByte(pos))<<40 |
		uint64(p.readByte(pos))<<48 |
		uint64(p.readByte(pos))<<56
}

func (p *Packet) readFloat32(pos *int) float32 {
	bits := binary.LittleEndian.Uint32((*p)[*pos:])
	f := math.Float32frombits(bits)
	*pos = *pos + 4

	return f

}

func (p *Packet) readString(pos *int, length int) string {
	return string(p.readBytes(pos, length))
}
