package handshakeTransaction

import (
	"encoding/hex"
	"errors"
	"github.com/blocktree/go-owcrypt"
)

type TxHash struct {
	Hash string
	Address string
	Signature []byte
	PublicKey []byte
}

func (tx TxHash) GetTxHashHex() string {
	return tx.Hash
}

func (tx TxHash) GetTxAddress() string {
	return tx.Address
}

func (t Transaction) getSigHashs() []TxHash {
	var previous []byte
	var sequence []byte
	var outputs []byte

	for _, in := range t.Vins {
		previous = append(previous, in.TxID...)
		previous = append(previous, in.Vout...)

		sequence = append(sequence, in.Sequence...)
	}

	for _, out := range t.Vouts {
		outputs = append(outputs, out.amount...)
		outputs = append(outputs, out.lockScript...)
	}

	hashPrevious := owcrypt.Hash(previous, 32, owcrypt.HASH_ALG_BLAKE2B)
	hashSequence := owcrypt.Hash(sequence, 32, owcrypt.HASH_ALG_BLAKE2B)
	hashOutputs := owcrypt.Hash(outputs, 32, owcrypt.HASH_ALG_BLAKE2B)

	var hashs []TxHash
	for _, in := range t.Vins {
		txBytes := []byte{}
		txBytes = append(txBytes, t.Version...)
		txBytes = append(txBytes, hashPrevious...)
		txBytes = append(txBytes, hashSequence...)
		txBytes = append(txBytes, in.TxID...)
		txBytes = append(txBytes, in.Vout...)
		txBytes = append(txBytes, in.Script...)
		txBytes = append(txBytes, in.Amount...)
		txBytes = append(txBytes, in.Sequence...)
		txBytes = append(txBytes, hashOutputs...)
		txBytes = append(txBytes, t.LockTime...)
		txBytes = append(txBytes, uint32ToLittleEndianBytes(uint32(SigHashAll))...)

		hashs = append(hashs, TxHash{
			Hash:hex.EncodeToString(owcrypt.Hash(txBytes, 32, owcrypt.HASH_ALG_BLAKE2B)),
		})
	}

	return hashs
}

func (t TxHash) getSigScript() ([]byte, error) {
	if t.Signature == nil || len(t.Signature) != 64 {
		return nil, errors.New("check signature!")
	}

	if t.PublicKey == nil || len(t.PublicKey) != 33 {
		return nil, errors.New("check publick key!")
	}

	ret := []byte{}

	ret = append(ret, 0x02, 0x41)
	ret = append(ret, t.Signature...)
	ret = append(ret, SigHashAll, 0x21)
	ret = append(ret, t.PublicKey...)

	return ret, nil
}