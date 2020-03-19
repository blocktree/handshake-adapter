package handshakeTransaction

const (
	TxVersion = uint32(0)
	AddressPrefix = "hs"
	AddressVersion = byte(0)
	TypeSend = byte(0)
	ActionNone = byte(0)
	SigHashAll = byte(1)
	OP_DUP = byte(0x76)
	OP_BLAKE160 = byte(0xc0)
	OP_EQUALVERIFY = byte(0x88)
	OP_CHECKSIG = byte(0xac)

)
