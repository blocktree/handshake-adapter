package handshakeTransaction

import (
	"encoding/hex"
	"errors"
)

type TxIn struct {
	TxID []byte
	Vout []byte
	Sequence []byte
	Script []byte
	Amount []byte
}

func (in TxIn) GetTxID() string {
	return hex.EncodeToString(in.TxID)
}

func (in TxIn) GetVout() uint32 {
	return littleEndianBytesToUint32(in.Vout)
}

func newTxInForEmptyTrans(vin []Vin) ([]TxIn, error) {
	if vin == nil || len(vin) == 0 {
		return nil, errors.New("No input found when create an empty transaction!")
	}

	var ret []TxIn

	for _, v := range vin {
		if v.TxID == "" || len(v.TxID) != 64 {
			return nil, errors.New("Invalid previous txID!")
		}
		txid, err := hex.DecodeString(v.TxID)
		if err != nil || len(txid) != 32 {
			return nil, errors.New("Invalid previous transaction id!")
		}

		vout := uint32ToLittleEndianBytes(v.Vout)

		if v.LockScript == "" {
			return nil, errors.New("Invalid lock script!")
		}

		script, err := hex.DecodeString(v.LockScript)
		if err != nil {
			return nil, errors.New("Invalid lock script!")
		}
		script = append([]byte{OP_DUP,OP_BLAKE160, byte(len(script))}, script...)
		script = append(script, OP_EQUALVERIFY, OP_CHECKSIG)
		script = append([]byte{byte(len(script))}, script...)
		if v.Amount == 0 {
			return nil, errors.New("Invalid amount of previous out!")
		}

		ret = append(ret, TxIn{
			TxID:txid,
			Vout:vout,
			Sequence:[]byte{0xFF, 0xFF, 0xFF, 0xFF},
			Script:script,
			Amount:uint64ToLittleEndianBytes(v.Amount),
		})
	}
	return ret, nil
}

func (in TxIn) toBytes() ([]byte, error) {
	var ret []byte

	if in.TxID == nil || len(in.TxID) != 32 {
		return nil, errors.New("Invalid previous transaction id!")
	}
	if in.Vout == nil || len(in.Vout) != 4 {
		return nil, errors.New("Invalid previous transaction vout!")
	}

	ret = append(ret, in.TxID...)
	ret = append(ret, in.Vout...)
	ret= append(ret, in.Sequence...)
	return ret, nil
}

func (in TxIn) getScript() ([]byte, error) {
	if in.Script == nil || len(in.Script) == 0 {
		return nil, errors.New("Invalid script data!")
	}

	if in.Amount == nil || len(in.Amount) == 0 {
		return nil, errors.New("Invalid amount data!")
	}

	ret := []byte{}
	ret = append(ret, in.Script...)
	ret = append(ret, in.Amount...)

	return ret, nil
}