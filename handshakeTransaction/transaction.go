package handshakeTransaction

import (
	"encoding/hex"
	"errors"
	"github.com/blocktree/go-owcrypt"
	"strings"
)

type Vin struct {
	TxID string
	Vout uint32
	LockScript string
	Amount uint64
}

type Vout struct {
	Address string
	Amount  uint64
}

func CreateEmptyRawTransactionAndHash(vins []Vin, vouts []Vout, lockTime uint32) (string, []TxHash, error) {
	trans, err := newEmptyTransaction(vins, vouts, lockTime)
	if err != nil {
		return "", nil, err
	}

	emptyTrans, scripts, err := trans.toBytes()
	if err != nil {
		return "", nil, err
	}

	ret := hex.EncodeToString(emptyTrans)

	for _, script := range scripts {
		ret += ":"
		ret += hex.EncodeToString(script)
	}

	hashes := trans.getSigHashs()
	for i := 0; i < len(hashes); i ++ {
		ls, _ := hex.DecodeString(vins[i].LockScript)
		hashes[i].Address = AddressEncode(ls)
	}

	return ret, hashes, nil
}

func SignRawTransactionHash(txHash string, prikey []byte) ([]byte, error) {
	hash, err := hex.DecodeString(txHash)
	if err != nil || len(hash) != 32 {
		return nil, errors.New("Invalid transaction hash!")
	}

	if prikey == nil || len(prikey) != 32 {
		return nil, errors.New("Invalid private key data!")
	}

	sig, _, retCode := owcrypt.Signature(prikey, nil, hash,  owcrypt.ECC_CURVE_SECP256K1)
	if retCode != owcrypt.SUCCESS {
		return nil, errors.New("Sign failed!")
	}

	return sig, nil
}

func CombineAndVerifyRawTransaction(emptyTrans string, hashs []TxHash) (string, bool) {
	if hashs == nil || len(hashs) == 0 {
		return "", false
	}

	trans := strings.Split(emptyTrans, ":")
	if len(trans) - 1 != len(hashs) {
		return "", false
	}

	tx, err := DecodeEmptyTransaction(emptyTrans)
	if err != nil {
		return "", false
	}

	txHashs := tx.getSigHashs()

	for i, hash := range txHashs {
		if hash.Hash != hashs[i].Hash {
			return "", false
		}
		hashBytes, _ := hex.DecodeString(hash.Hash)
		pubkey := owcrypt.PointDecompress(hashs[i].PublicKey, owcrypt.ECC_CURVE_SECP256K1)[1:]
		if owcrypt.SUCCESS != owcrypt.Verify(pubkey, nil,  hashBytes,  hashs[i].Signature, owcrypt.ECC_CURVE_SECP256K1) {
			return "", false
		}
	}

	sigpubs := []byte{}

	for _, h := range hashs {
		sp, err := h.getSigScript()
		if err != nil {
			return "", false
		}
		sigpubs = append(sigpubs, sp...)
	}
	return trans[0] + hex.EncodeToString(sigpubs), true
}