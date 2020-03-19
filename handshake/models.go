/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package handshake

import (
	"encoding/hex"
	"errors"
	"github.com/blocktree/handshake-adapter/handshakeTransaction"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/btcsuite/btcd/txscript"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"strings"
)

//BlockchainInfo 本地节点区块链信息
type BlockchainInfo struct {
	Chain                string `json:"chain"`
	Blocks               uint64 `json:"blocks"`
	Headers              uint64 `json:"headers"`
	Bestblockhash        string `json:"bestblockhash"`
	Difficulty           string `json:"difficulty"`
	Mediantime           uint64 `json:"mediantime"`
	Verificationprogress string `json:"verificationprogress"`
	Chainwork            string `json:"chainwork"`
	Pruned               bool   `json:"pruned"`
}

func NewBlockchainInfo(json *gjson.Result) *BlockchainInfo {
	b := &BlockchainInfo{}
	//解析json
	b.Chain = gjson.Get(json.Raw, "chain").String()
	b.Blocks = gjson.Get(json.Raw, "blocks").Uint()
	b.Headers = gjson.Get(json.Raw, "headers").Uint()
	b.Bestblockhash = gjson.Get(json.Raw, "bestblockhash").String()
	b.Difficulty = gjson.Get(json.Raw, "difficulty").String()
	b.Mediantime = gjson.Get(json.Raw, "mediantime").Uint()
	b.Verificationprogress = gjson.Get(json.Raw, "verificationprogress").String()
	b.Chainwork = gjson.Get(json.Raw, "chainwork").String()
	b.Pruned = gjson.Get(json.Raw, "pruned").Bool()
	return b
}

//Unspent 未花记录
type Unspent struct {

	/*
			{
		        "txid" : "d54994ece1d11b19785c7248868696250ab195605b469632b7bd68130e880c9a",
		        "vout" : 1,
		        "address" : "mgnucj8nYqdrPFh2JfZSB1NmUThUGnmsqe",
		        "account" : "test label",
		        "scriptPubKey" : "76a9140dfc8bafc8419853b34d5e072ad37d1a5159f58488ac",
		        "amount" : 0.00010000,
		        "confirmations" : 6210,
		        "spendable" : true,
		        "solvable" : true
		    }
	*/
	Key           string `storm:"id"`
	TxID          string `json:"txid"`
	Vout          uint64 `json:"vout"`
	Address       string `json:"address"`
	AccountID     string `json:"account" storm:"index"`
	ScriptPubKey  string `json:"scriptPubKey"`
	Amount        string `json:"amount"`
	Confirmations uint64 `json:"confirmations"`
	Type          uint64 `json:"type"`
	Action        string `json:"action"`
	Spendable     bool   `json:"spendable"`
	Solvable      bool   `json:"solvable"`
	HDAddress     openwallet.Address
}

func NewUnspent(json *gjson.Result) *Unspent {
	/*

	  {
	    "version": 0,
	    "height": 3975,
	    "value": 1039500000,
	    "address": "hs1qtsevyrskarucasazwgs7rk8stc36lky7wqrh5k",
	    "covenant": {
	      "type": 0,
	      "action": "NONE",
	      "items": []
	    },
	    "coinbase": false,
	    "hash": "00218895b13d384bba2c0bfda09d8154ddf938d1ec3e5619ed7c887703ec750f",
	    "index": 1
	  }

	*/
	obj := &Unspent{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "hash").String()
	obj.Vout = gjson.Get(json.Raw, "index").Uint()
	obj.Address = gjson.Get(json.Raw, "address").String()
	obj.AccountID = gjson.Get(json.Raw, "address").String()
	hash, _ := handshakeTransaction.AddressDecode(obj.Address)
	obj.ScriptPubKey = hex.EncodeToString(hash)

	amountDecimal, _ := decimal.NewFromString(gjson.Get(json.Raw, "value").String())
	div, _ := decimal.NewFromString("1000000")

	obj.Amount = amountDecimal.Div(div).String()
	//obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	//obj.Spendable = gjson.Get(json.Raw, "spendable").Bool()
	obj.Type = gjson.Get(json.Raw, "covenant").Get("type").Uint()
	obj.Action = gjson.Get(json.Raw, "covenant").Get("action").String()
	obj.Spendable = true
	obj.Solvable = gjson.Get(json.Raw, "solvable").Bool()

	return obj
}

type UnspentSort struct {
	Values     []*Unspent
	Comparator func(a, b *Unspent) int
}

func (s UnspentSort) Len() int {
	return len(s.Values)
}
func (s UnspentSort) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
}
func (s UnspentSort) Less(i, j int) bool {
	return s.Comparator(s.Values[i], s.Values[j]) < 0
}

//type Address struct {
//	Address   string `json:"address" storm:"id"`
//	Account   string `json:"account" storm:"index"`
//	HDPath    string `json:"hdpath"`
//	CreatedAt time.Time
//}

type User struct {
	UserKey string `storm:"id"`     // primary key
	Group   string `storm:"index"`  // this field will be indexed
	Email   string `storm:"unique"` // this field will be indexed with a unique constraint
	Name    string // this field will not be indexed
	Age     int    `storm:"index"`
}

type Block struct {

	/*

		"hash": "000000000000000127454a8c91e74cf93ad76752cceb7eb3bcff0c398ba84b1f",
		"confirmations": 2,
		"strippedsize": 191875,
		"size": 199561,
		"weight": 775186,
		"height": 1354760,
		"version": 536870912,
		"versionHex": "20000000",
		"merkleroot": "48239e76f8b37d9c8824fef93d42ac3d7c433029c1e9fa23b6416dd0356f3e57",
		"tx": ["c1e12febeb58aefb0b01c04360262138f4ee0faeb207276e79ea3866608ed84f"]
		"time": 1532143012,
		"mediantime": 1532140298,
		"nonce": 3410287696,
		"bits": "19499855",
		"difficulty": 58358570.79038175,
		"chainwork": "00000000000000000000000000000000000000000000006f68c43926cd6c2d1f",
		"previousblockhash": "00000000000000292d142fcc1ddbd9dafd4518310009f152bdca2a66cc589f97",
		"nextblockhash": "0000000000004a50ef5733ab333f718e6ef5c1995e2cfd5a7caa0875f118cd30"

	*/

	Hash              string
	Confirmations     uint64
	Merkleroot        string
	tx                []string
	Previousblockhash string
	Height            uint64 `storm:"id"`
	Version           uint64
	Time              uint64
	Fork              bool
	txDetails         []*Transaction
	isVerbose         bool
}

func (wm *WalletManager) NewBlock(json *gjson.Result) *Block {
	obj := &Block{}
	//解析json
	obj.Height = gjson.Get(json.Raw, "height").Uint()
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Merkleroot = gjson.Get(json.Raw, "merkleroot").String()
	obj.Previousblockhash = gjson.Get(json.Raw, "previousblockhash").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.Time = gjson.Get(json.Raw, "time").Uint()

	txs := make([]string, 0)
	txDetails := make([]*Transaction, 0)
	for _, tx := range gjson.Get(json.Raw, "tx").Array() {
		if tx.IsObject() {
			obj.isVerbose = true
			txObj := wm.newTxByCore(&tx)
			txObj.BlockHeight = obj.Height
			txObj.BlockHash = obj.Hash
			txObj.Blocktime = int64(obj.Time)
			txDetails = append(txDetails, txObj)
		} else {
			obj.isVerbose = false
			txs = append(txs, tx.String())
		}

	}

	obj.tx = txs
	obj.txDetails = txDetails

	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader(symbol string) *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	obj.Confirmations = b.Confirmations
	obj.Merkleroot = b.Merkleroot
	obj.Previousblockhash = b.Previousblockhash
	obj.Height = b.Height
	obj.Version = b.Version
	obj.Time = b.Time
	obj.Symbol = symbol

	return &obj
}

type Transaction struct {
	TxID          string
	Size          uint64
	Version       uint64
	LockTime      int64
	Hex           string
	BlockHash     string
	BlockHeight   uint64
	Confirmations uint64
	Blocktime     int64
	IsCoinBase    bool
	Fees          string
	Decimals      int32

	Vins  []*Vin
	Vouts []*Vout
}

type Vin struct {
	Coinbase string
	TxID     string
	Vout     uint64
	N        uint64
	Addr     string
	Value    string
}

type Vout struct {
	N            uint64
	Addr         string
	Value        string
	ScriptPubKey string
	Type         string
	Action       string
}

func (wm *WalletManager) newTxByCore(json *gjson.Result) *Transaction {

	/*
		{
			"txid": "6595e0d9f21800849360837b85a7933aeec344a89f5c54cf5db97b79c803c462",
			"hash": "f758cb5181d51f8bee1512b4a862faad5b51c7c85a1a11cd6092ffc1c6649bc5",
			"version": 2,
			"size": 249,
			"vsize": 168,
			"locktime": 1414190,
			"vin": [],
			"vout": [],
			"hex": "02000000000101cc8a3077023c08040e677647ad0e528564764f456b01d8519828df165ab3c4550100000017160014aa59f94152351c79b57b14a53e538a923e332468feffffff02a716167c6f00000017a914a0fe07f130a36d9c7581ccd2886895c049b0cc8287ece29c00000000001976a9148c0bceb59d452b3e077f73a420b8bfe09e0550a788ac0247304402205e667171c1798cde426282bb8bff45901866ad6bf0d209e856c1765eda65ba4802203aaa319ea3de00eccef0006e6ee2089aed4b91ada7953f420a47c9c258d424ca0121033cfda2f93d13b01d46ecc406b03ebaba3e1bd526d2148a0a5d579d52f8c7cf022e941500",
			"blockhash": "0000000040730ea7935cce346ce68bf4c07c10b137ba31960bf8a47c4f7da4ec",
			"confirmations": 20076,
			"time": 1537841342,
			"blocktime": 1537841342
		}
	*/

	obj := Transaction{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.LockTime = gjson.Get(json.Raw, "locktime").Int()
	obj.BlockHash = gjson.Get(json.Raw, "blockhash").String()
	//obj.BlockHeight = gjson.Get(json.Raw, "blockheight").Uint()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Blocktime = gjson.Get(json.Raw, "blocktime").Int()
	obj.Size = gjson.Get(json.Raw, "size").Uint()
	//obj.Fees = gjson.Get(json.Raw, "fees").String()
	obj.Decimals = wm.Decimal()
	obj.Vins = make([]*Vin, 0)
	if vins := gjson.Get(json.Raw, "vin"); vins.IsArray() {
		for i, vin := range vins.Array() {
			input := newTxVinByCore(&vin)
			input.N = uint64(i)
			obj.Vins = append(obj.Vins, input)
		}
	}

	obj.Vouts = make([]*Vout, 0)
	if vouts := gjson.Get(json.Raw, "vout"); vouts.IsArray() {
		for _, vout := range vouts.Array() {
			output := newTxVoutByCore(&vout)
			obj.Vouts = append(obj.Vouts, output)
		}
	}

	return &obj
}

func newTxVinByCore(json *gjson.Result) *Vin {

	/*
		{
			"txid": "55c4b35a16df289851d8016b454f766485520ead4776670e04083c0277308acc",
			"vout": 1,
			"scriptSig": {
				"asm": "0014aa59f94152351c79b57b14a53e538a923e332468",
				"hex": "160014aa59f94152351c79b57b14a53e538a923e332468"
			},
			"txinwitness": ["304402205e667171c1798cde426282bb8bff45901866ad6bf0d209e856c1765eda65ba4802203aaa319ea3de00eccef0006e6ee2089aed4b91ada7953f420a47c9c258d424ca01", "033cfda2f93d13b01d46ecc406b03ebaba3e1bd526d2148a0a5d579d52f8c7cf02"],
			"sequence": 4294967294
		}
	*/
	obj := Vin{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Vout = gjson.Get(json.Raw, "vout").Uint()
	obj.Coinbase = gjson.Get(json.Raw, "coinbase").String()
	//obj.Addr = gjson.Get(json.Raw, "addr").String()
	//obj.Value = gjson.Get(json.Raw, "value").String()

	return &obj
}

func newTxVoutByCore(json *gjson.Result) *Vout {

	/*
		{
			"value": 4788.23192231,
			"n": 0,
			"scriptPubKey": {
				"asm": "OP_HASH160 a0fe07f130a36d9c7581ccd2886895c049b0cc82 OP_EQUAL",
				"hex": "a914a0fe07f130a36d9c7581ccd2886895c049b0cc8287",
				"reqSigs": 1,
				"type": "scripthash",
				"addresses": ["2N7vURMwMDjqgijLNFsErFLAWtAg58S6qNv"]
			}
		}
	*/
	obj := Vout{}
	//解析json
	obj.Value = gjson.Get(json.Raw, "value").String()
	obj.N = gjson.Get(json.Raw, "n").Uint()
	obj.ScriptPubKey = gjson.Get(json.Raw, "scriptPubKey.hex").String()

	//提取地址
	if addresses := gjson.Get(json.Raw, "scriptPubKey.addresses"); addresses.IsArray() {
		obj.Addr = addresses.Array()[0].String()
	}

	obj.Type = gjson.Get(json.Raw, "scriptPubKey.type").String()

	//if len(obj.Addr) == 0 {
	//	scriptBytes, _ := hex.DecodeString(obj.ScriptPubKey)
	//	obj.Addr, _ = wm.Decoder.ScriptPubKeyToBech32Address(scriptBytes)
	//}

	return &obj
}

func DecodeScript(script string) ([]byte, error) {
	opcodes := strings.Split(script, " ")
	scriptBuilder := txscript.NewScriptBuilder()
	for _, codeName := range opcodes {
		code, ok := txscript.OpcodeByName[codeName]
		if ok {
			scriptBuilder.AddOp(code)
		} else {
			if len(codeName)%2 != 0 {
				codeName = "0" + codeName
			}
			data, err := hex.DecodeString(codeName)
			if err != nil {
				return nil, err
			}
			scriptBuilder.AddData(data)
		}
	}
	return scriptBuilder.Script()
}

func newBlock(json *gjson.Result) *Block {

	/*
		{
				"hash": "0000000000000141e0f71fe490cded22def37cd08accf19f25dc98b681512d5a",
				"confirmations": 51,
				"strippedsize": 3380,
				"size": 6840,
				"weight": 16980,
				"height": 6000,
				"version": 0,
				"versionHex": "00000000",
				"merkleroot": "cfd2d710bb6743e975f15ffb48e08b683e977007f4dcaa3e8080206f20543782",
				"witnessroot": "f087dd353c9b02828d8d4f427a50970dab85ef84796a580412596e6c8fd3c401",
				"treeroot": "6527db76db0ac3b352d8a89e62c95ab1dd8eab5f3720746847f065674b1adc93",
				"reservedroot": "0000000000000000000000000000000000000000000000000000000000000000",
				"mask": "0000000000000000000000000000000000000000000000000000000000000000",
				"coinbase": ["6632706f6f6c", "5b1c242744be67e7", "0000000000000000"],
				"tx": ["032be866fa778e291bfc50de9d3e35a932e9084f6ca89744947095f8be83d276", "cbeb59df3a503606c85674402fc2298e2d8e2bd0a94b5dc3fa27d855d1c19c76", "c5fd97bffa30a9dc5aadb96f05ee2cc9c13ed18e6a7e2033bad33a04c5e36c1c", "19e73a61ebe25f7e87c631e5d279686a6c7de01b89bcf046598ecd04978b12c3", "aee36844e97c81b01c1ac74e86689f4fa0c0954fa5716427565470f2006a1bdc", "5609f14928a343cdfd9ba9dcedc57662fe2fc8337efb060a35c31b6f30320dfa", "59d3de7cdf620c34f227a9d346a192b5c596c4f544c759f576a3132f4edcd3bf", "1c7af566d5993f2fa460ad3fb1306e46879f824105b99c4fff15a4ef3db18d70", "de7cc9bc616ecc052e93edadcec853ae2e980427c6b318f7bd0ac9ab319c9bb3", "0f137d0a32c0176179fe795b83834ac2e807e6eeed9eee62f83585f3edd5f164", "81fd681b68bd83476df7a816b989628ed65c8f2eeb673e0dfddeacabe2e569c8", "9675e8a59b0e8351e0f7f441522cd073906f0bf08af629ecb262a3256b7cfed3", "5359060f9f7ba724ea14956c9d779636795ef880d74c063d87ea213ce346511f", "d8e1169eaff4163a2fb35390e98c603a242b99db5b4ac0ab69c184c1925abbfe", "e162c56579d019f40c749ba9f491d1e0f41e47ffb24b180e00a6118a55e5d14f", "55f215e8aa16012ae545087772d66d56c9d4f4e8a61ac10bace1112082727056"],
				"time": 1584030329,
				"mediantime": 1584024525,
				"bits": 436342722,
				"difficulty": 8138016.450490726,
				"chainwork": "000000000000000000000000000000000000000000000003c55a2782ba15f9bd",
				"previousblockhash": "00000000000000fdd5ab59e35682e99fb5a58cc5b09057445f9984b6040914c0",
				"nextblockhash": "0000000000000094548ff6a51169e08136136adf69e3f00ef3ed1af6c2550ebe"
			}
	*/
	obj := &Block{}
	//解析json
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Merkleroot = gjson.Get(json.Raw, "merkleroot").String()

	txs := make([]string, 0)
	for _, tx := range gjson.Get(json.Raw, "tx").Array() {
		txs = append(txs, tx.String())
	}

	obj.tx = txs
	obj.Previousblockhash = gjson.Get(json.Raw, "previousblockhash").String()
	obj.Height = gjson.Get(json.Raw, "height").Uint()
	//obj.Version = gjson.Get(json.Raw, "version").String()
	obj.Time = gjson.Get(json.Raw, "time").Uint()

	return obj
}

func (c Client) newTx(json *gjson.Result) (*Transaction, error) {
	/*

		{
					"txid": "ec823cbfcd7e6e49491e5d3c2ad09d0b76f770bfa24d3cd877e2ab323674d522",
					"hash": "8920c00875ff95e06d8542d250c7b41793a871d4ba926f4ddce25fb6e831154e",
					"size": 215,
					"vsize": 140,
					"version": 0,
					"locktime": 0,
					"vin": [],
					"vout": [],
					"blockhash": "00000000000001575695356d50a4bde8269c9aa6748e505748e4a034f981430e",
					"confirmations": 12,
					"time": 1583989399,
					"blocktime": 1583989365,
					"hex": "0000000001fa1e8de6158618ec5e278db5d647ff0175be004cb913a7a0565875194fe19fc401000000ffffffff0240420f00000000000014b302960fb163255e3abf855babd47da1d819bb85000020444f0100000000001402e84f4434edd899e6b9aa8ffb4ce85f0216243d00000000000002413bd1043b98792587d96572edb04ec367f92bde48767ee3138e8081a1fd9704326da384381958a2128f29872932d64fe9cc7af62b3b9cd8ec83d6dda28686441101210264d953039023df3424f2471b8269d67d1c714e94c707908cb8efffbb7cb9cc23"
				},

	*/
	obj := Transaction{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.LockTime = gjson.Get(json.Raw, "locktime").Int()
	obj.BlockHash = gjson.Get(json.Raw, "blockhash").String()
	if obj.BlockHash != "" {
		block, err := c.getBlock(obj.BlockHash)
		if err != nil {
			return &Transaction{}, err
		}
		obj.BlockHeight = block.Height
	}

	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Blocktime = gjson.Get(json.Raw, "blocktime").Int()
	obj.Size = gjson.Get(json.Raw, "size").Uint()

	obj.Vins = make([]*Vin, 0)
	if vins := gjson.Get(json.Raw, "vin"); vins.IsArray() {
		for _, vin := range vins.Array() {
			if vin.Get("coinbase").String() == "true" {
				break
			}
			input, err := c.newTxVin(&vin)
			if err != nil {
				return &Transaction{}, err
			}
			if input != nil {
				obj.Vins = append(obj.Vins, input)
			}
		}
	}

	obj.Vouts = make([]*Vout, 0)
	if vouts := gjson.Get(json.Raw, "vout"); vouts.IsArray() {
		for _, vout := range vouts.Array() {
			output, err := newTxVout(&vout)
			if err != nil {
				return &Transaction{}, err
			}
			if output != nil {
				obj.Vouts = append(obj.Vouts, output)
			}
		}
	}

	return &obj, nil
}

func (c Client) newTxVin(json *gjson.Result) (*Vin, error) {
	/*

		{
			"coinbase": false,
			"txid": "fa1e8de6158618ec5e278db5d647ff0175be004cb913a7a0565875194fe19fc4",
			"vout": 1,
			"txinwitness": ["3bd1043b98792587d96572edb04ec367f92bde48767ee3138e8081a1fd9704326da384381958a2128f29872932d64fe9cc7af62b3b9cd8ec83d6dda28686441101", "0264d953039023df3424f2471b8269d67d1c714e94c707908cb8efffbb7cb9cc23"],
			"sequence": 4294967295
		}
	*/

	obj := Vin{}
	var err error
	//解析json
	obj.Coinbase = gjson.Get(json.Raw, "coinbase").String()
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Vout = gjson.Get(json.Raw, "vout").Uint()
	obj.N = obj.Vout
	obj.Addr, obj.Value, err = c.getAddressAndAmoutOfInput(obj.TxID, obj.Vout)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func newTxVout(json *gjson.Result) (*Vout, error) {
	/*

		{
			"value": 1,
			"n": 0,
			"address": {
				"version": 0,
				"hash": "b302960fb163255e3abf855babd47da1d819bb85"
			},
			"covenant": {
				"type": 0,
				"action": "NONE",
				"items": []
			}
		}
	*/

	obj := Vout{}
	//解析json
	obj.Value = gjson.Get(json.Raw, "value").String()
	obj.N = gjson.Get(json.Raw, "n").Uint()
	obj.ScriptPubKey = gjson.Get(json.Raw, "address").Get("hash").String()
	if gjson.Get(json.Raw, "address").Get("version").Uint() != 0 {
		return nil, errors.New("unknown address version")
	}

	hash, _ := hex.DecodeString(obj.ScriptPubKey)
	obj.Addr = handshakeTransaction.AddressEncode(hash)
	obj.Type = gjson.Get(json.Raw, "covenant").Get("type").String()
	obj.Action = gjson.Get(json.Raw, "covenant").Get("action").String()

	return &obj, nil
}
