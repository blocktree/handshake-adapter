package handshakeTransaction

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_9a9e2b836d2b9640103633857febf62f488e78b4bc86806603382fcadc61a4f8(t*testing.T)  {
	in := Vin{
		TxID:       "ec823cbfcd7e6e49491e5d3c2ad09d0b76f770bfa24d3cd877e2ab323674d522",
		Vout:       0,
		LockScript: "b302960fb163255e3abf855babd47da1d819bb85",
		Amount:     1000000,
	}
	
	out1 := Vout{
		Address: "hs1q2vnxeuq4ueqh36hln642kmkjjpx26083upy9d5",
		Amount:  1000,
	}
	out2 := Vout{
		Address: "hs1qmhylkn9eg3fr0tushpkna0k9y9lqxzx4dzrpc7",
		Amount:  699000,
	}
	lockTime := uint32(0)

	emptyTrans, hashes, err := CreateEmptyRawTransactionAndHash([]Vin{in}, []Vout{out1, out2}, lockTime)
	if err != nil {
		t.Error("ceate tx failed! - ", err)
	} else {
		fmt.Println("empty tx : \n", emptyTrans)
		for i, h := range hashes {
			fmt.Println("the ", i + 1, " hash is : ")
			fmt.Println(h.Hash)
		}
	}

	prikey, _ := hex.DecodeString("370b3b5c6f74d0052b39982cd351d2d0901d821429e311a3df75515c40cceb68")

	sig, err := SignRawTransactionHash(hashes[0].Hash, prikey)
	if err != nil {
		t.Error("sign tx failed")
	} else {
		sig, _ = hex.DecodeString("55cfd748f2cb768ebad5601d164adf36e754f45adf69ea4cde6b885d12e9a95b05b023a4939d2e79b469a02f8a4e4dca1771945bb998cbe2dfc6384c12790af6")
		fmt.Println("sig : ", hex.EncodeToString(sig))
	}
	hashes[0].Signature = sig
	hashes[0].PublicKey, _ = hex.DecodeString("03ac2c33b23097cc8b442015f824fa90c1e2cd64b9a681add03aa1e82e7014edc1")
	signedTrans, pass := CombineAndVerifyRawTransaction(emptyTrans, hashes)
	if pass{
		fmt.Println(signedTrans)
	} else {
		t.Error("verify tx failed")
	}
}

func Test_4fadb5ceb93c405cb068f78d891edb1a8d32468a13eb0d228d7f0a380a885f03(t *testing.T)  {
	in1 := Vin{
		TxID:       "c9dbf90d3c5984883979e31a624b17678b0ec811dd2090db69ef4a28eda9d29a",
		Vout:       0,
		LockScript: "9cf4a7d906a8dd99638c4bdfa5984070d953cb92",
		Amount:     69821315720,
	}
	in2 := Vin{
		TxID:       "e9900720f2e2864cdbd0e3f66b0383bfafe43ca36ec45a1b1909e4e651561bff",
		Vout:       1,
		LockScript: "4a6b4a067f20ba8a4d3103f18f0e06103b3e7b03",
		Amount:     191650653400,
	}

	out1 := Vout{
		Address: "hs1qvyp48f8zkya3ldje8qmz2gf76ry77qpxjs4sw3",
		Amount:  61471948520,
	}
	out2 := Vout{
		Address: "hs1qvz2m5hte2pk377cn6ersh94cgvek8jt9374x4d",
		Amount:  200000000000,
	}
	lockTime := uint32(0)

	emptyTrans, hashes, err := CreateEmptyRawTransactionAndHash([]Vin{in1, in2}, []Vout{out1, out2}, lockTime)
	if err != nil {
		t.Error("ceate tx failed! - ", err)
	} else {
		fmt.Println("empty tx : \n", emptyTrans)
		for i, h := range hashes {
			fmt.Println("the ", i + 1, " hash is : ")
			fmt.Println(h.Hash)
		}
	}


	hashes[0].Signature, _ =  hex.DecodeString("cb83c91676308c6118817327823867e7365aca57734abf8275db015e958ce7b4611d53463f12cba76e495d11f4c237142f7db0fd5c0ad2ed158782f780c49608")
	hashes[0].PublicKey, _ = hex.DecodeString("034c34bd470c67bacdc075d3dce80a9d7eafe21f48f902ade6d570306faf30d2b0")


	hashes[1].Signature, _  = hex.DecodeString("4445bf956bc926e49ccd1819c90d957a4ff6dbe1e53ee27fcaa5b738d6431ade5b2bdeb1a76ef996b6240e30a7cf8a214e0bc08f7f2f8962e76dd4e9ab42dcf7")
	hashes[1].PublicKey, _ = hex.DecodeString("02920e84996ebbf09f3a67198e5db03e26dad989e78c25d11ee5d635b26169248a")

	signedTrans, pass := CombineAndVerifyRawTransaction(emptyTrans, hashes)
	if pass{
		fmt.Println(signedTrans)
	} else {
		t.Error("verify tx failed")
	}
}