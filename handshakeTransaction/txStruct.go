package handshakeTransaction

import (
	"encoding/hex"
	"errors"
	"strings"
)

type Transaction struct {
	Version  []byte
	Vins     []TxIn
	Vouts    []TxOut
	LockTime []byte
}

func newEmptyTransaction(vins []Vin, vouts []Vout, lockTime uint32) (*Transaction, error) {

	txIn, err := newTxInForEmptyTrans(vins)
	if err != nil {
		return nil, err
	}

	txOut, err := newTxOutForEmptyTrans(vouts)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Version:  uint32ToLittleEndianBytes(TxVersion),
		Vins:     txIn,
		Vouts:    txOut,
		LockTime: uint32ToLittleEndianBytes(lockTime),
	}, nil

}

func (t Transaction) toBytes() ([]byte, [][]byte, error)  {

	if t.Vins == nil || len(t.Vins) == 0 {
		return nil, nil, errors.New("No input found in the transaction struct!")
	}

	if t.Vouts == nil || len(t.Vouts) == 0 {
		return nil, nil, errors.New("No output found in the transaction struct!")
	}

	if t.Version == nil || len(t.Version) != 4 {
		return nil, nil, errors.New("Invalid transaction version data!")
	}

	if t.LockTime == nil || len(t.LockTime) != 4 {
		return nil, nil, errors.New("Invalid loack time data!")
	}

	ret := []byte{}

	ret = append(ret, t.Version...)

	ret = append(ret, byte(len(t.Vins)))
	for _, in := range t.Vins {
		inBytes, err := in.toBytes()
		if err != nil {
			return nil, nil, err
		}

		ret = append(ret, inBytes...)
	}

	ret = append(ret, byte(len(t.Vouts)))
	for _, out := range t.Vouts {
		outBytes, err := out.toBytes()
		if err != nil {
			return nil, nil, err
		}

		ret = append(ret, outBytes...)
	}

	ret = append(ret, t.LockTime...)

	scripts := [][]byte{}

	for _, in := range t.Vins {
		script, err := in.getScript()
		if err != nil {
			return nil, nil, err
		}

		scripts = append(scripts, script)
	}

	return ret, scripts, nil
}

func decodeScript(script string) ([]byte, []byte, error)  {
	bytes, err := hex.DecodeString(script)
	if err != nil {
		return nil, nil, errors.New("Invalid script!")
	}

	if len(bytes) != 34 {
		return nil, nil, errors.New("Invalid script!")
	}

	return bytes[:26], bytes[26:], nil
}

func DecodeEmptyTransaction(emptyTrans string) (*Transaction, error) {

	txs := strings.Split(emptyTrans, ":")

	tx, err := hex.DecodeString(txs[0])
	if err != nil {
		return nil, errors.New("Invalid transaction data!")
	}

	limit := len(tx)
	index := 0
	var t Transaction

	if index + 4 > limit {
		return nil, errors.New("Invalid transaction data!")
	}
	t.Version = tx[index: index + 4]
	index += 4

	if index + 1 > limit {
		return nil, errors.New("Invalid transaction data!")
	}
	inCount := tx[index]
	index ++
	if inCount == 0 {
		return nil, errors.New("Invalid transaction data!")
	}

	for i := 0; i < int(inCount); i ++ {
		if index + 40 > limit {
			return nil, errors.New("Invalid transaction data!")
		}

		script, amount, err := decodeScript(txs[i + 1])
		if err != nil {
			return nil, err
		}

		t.Vins = append(t.Vins, TxIn{
			TxID:     tx[index:index+32],
			Vout:     tx[index+32:index+36],
			Sequence: tx[index+36:index+40],
			Script:   script,
			Amount:   amount,
		})
		index += 40
	}

	if index + 1 > limit {
		return nil, errors.New("Invalid transaction data!")
	}
	outCount := tx[index]
	index ++
	if outCount == 0 {
		return nil, errors.New("Invalid transaction data!")
	}

	for i := 0; i < int(outCount); i ++ {
		if index + 32 > limit {
			return nil, errors.New("Invalid transaction data!")
		}

		t.Vouts = append(t.Vouts, TxOut{amount:tx[index:index+8], lockScript:tx[index+8:index+32]})
		index += 32
	}

	if index + 4 != limit {
		return nil, errors.New("Invalid transaction data!")
	}
	t.LockTime = tx[index:]

	return &t, nil
}