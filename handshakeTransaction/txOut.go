package handshakeTransaction

import "errors"

type TxOut struct {
	amount     []byte
	lockScript []byte
}

func newTxOutForEmptyTrans(vout []Vout) ([]TxOut, error) {
	if vout == nil || len(vout) == 0 {
		return nil, errors.New("No address to send when create an empty transaction!")
	}

	var ret []TxOut

	for _, v := range vout {
		if v.Amount == 0 {
			return nil, errors.New("Invalid amount to send!")
		}
		amount := uint64ToLittleEndianBytes(v.Amount)

		hash, err := AddressDecode(v.Address)
		if err != nil {
			return nil, err
		}
		hash = append([]byte{byte(len(hash))}, hash...)
		hash = append([]byte{AddressVersion}, hash...)
		hash = append(hash, TypeSend, ActionNone)
		
		ret = append(ret, TxOut{
			amount:     amount,
			lockScript: hash,
		})
	}

	return ret, nil
}

func (out TxOut) toBytes() ([]byte, error) {
	if out.amount == nil || len(out.amount) != 8 {
		return nil, errors.New("Invalid amount for a transaction output!")
	}
	if out.lockScript == nil || len(out.lockScript) == 0 {
		return nil, errors.New("Invalid lock script for a transaction output!")
	}

	ret := []byte{}
	ret = append(ret, out.amount...)
	ret = append(ret, out.lockScript...)

	return ret, nil
}