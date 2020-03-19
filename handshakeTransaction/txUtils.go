package handshakeTransaction

import "encoding/binary"

//uint32ToLittleEndianBytes
func uint32ToLittleEndianBytes(data uint32) []byte {
	tmp := [4]byte{}
	binary.LittleEndian.PutUint32(tmp[:], data)
	return tmp[:]
}

//littleEndianBytesToUint32
func littleEndianBytesToUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

//uint64ToLittleEndianBytes
func uint64ToLittleEndianBytes(data uint64) []byte {
	tmp := [8]byte{}
	binary.LittleEndian.PutUint64(tmp[:], data)
	return tmp[:]
}

//littleEndianBytesToUint64
func littleEndianBytesToUint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}
