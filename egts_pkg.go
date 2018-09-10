package main

import (
	"bytes"
	"encoding/binary"
)

// BinaryData интерфейс для работы с бинарными секциями
type BinaryData interface {
	Encode() ([]byte, error)
	Length() uint16
}

// EgtsPackage стуркура для описания пакета ЕГТС
type EgtsPackage struct {
	EgtsHeader
	ServicesFrameData         BinaryData
	ServicesFrameDataCheckSum uint16
}

// Encode кодирует струткуру в байтовую строку
func (p *EgtsPackage) Encode() ([]byte, error) {
	var result []byte
	buf := new(bytes.Buffer)

	hdr, err := p.EgtsHeader.Encode()
	if err != nil {
		return result, err
	}
	buf.Write(hdr)

	sfrd, err := p.ServicesFrameData.Encode()
	if err != nil {
		return result, err
	}

	if len(sfrd) > 0 {
		buf.Write(sfrd)

		if err := binary.Write(buf, binary.LittleEndian, crc16(sfrd)); err != nil {
			return result, err
		}
	}

	result = buf.Bytes()
	return result, nil
}
