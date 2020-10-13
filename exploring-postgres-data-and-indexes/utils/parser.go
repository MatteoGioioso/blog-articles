package utils

import "encoding/binary"

const (
	uint8Len = 1
	uint16Len = 2
	uint32Len = 4
)

type Parser struct {
	Buff []byte
}

func (p *Parser) ShiftBy(len int) []byte {
	bytes := p.Buff[:len]
	p.Buff = p.Buff[len:]
	return bytes
}

func (p *Parser) GetInt16() uint16 {
	u := binary.LittleEndian.Uint16(p.Buff[:uint16Len])
	p.Buff = p.Buff[uint16Len:]
	return u
}

func (p *Parser) GetInt32() uint32 {
	u := binary.LittleEndian.Uint32(p.Buff[:uint32Len])
	p.Buff = p.Buff[uint32Len:]
	return u
}

func (p *Parser) GetInt8() uint8 {
	u := p.Buff[0]
	p.Buff = p.Buff[uint8Len:]
	return u
}
