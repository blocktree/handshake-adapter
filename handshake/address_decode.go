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
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/handshake-adapter/handshakeTransaction"
	"github.com/blocktree/openwallet/v2/openwallet"
)

func init() {

}

type AddressDecoderV2 struct {

	openwallet.AddressDecoderV2Base
	//ScriptPubKeyToBech32Address(scriptPubKey []byte) (string, error)
}

//NewAddressDecoder 地址解析器
func NewAddressDecoderV2(wm *WalletManager) *AddressDecoderV2 {
	decoder := AddressDecoderV2{}
	return &decoder
}


//AddressDecode 地址解析
func (dec *AddressDecoderV2) AddressDecode(addr string, opts ...interface{}) ([]byte, error) {
	decodeHash, err := handshakeTransaction.AddressDecode(addr)
	if err != nil {
		return nil, err
	}
	return decodeHash, nil
}

//AddressEncode 地址编码
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {

	//公钥hash处理
	hash = owcrypt.Hash(hash, 20, owcrypt.HASH_ALG_BLAKE2B)

	address := handshakeTransaction.AddressEncode(hash)

	return address, nil
}

// AddressVerify 地址校验
func (dec *AddressDecoderV2) AddressVerify(address string, opts ...interface{}) bool {
	_, err := handshakeTransaction.AddressDecode(address)
	if err != nil {
		return false
	}
	return true
}




//
////PrivateKeyToWIF 私钥转WIF
//func (decoder *addressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {
//	return "", nil
//}
//
////PublicKeyToAddress 公钥转地址
//func (decoder *addressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {
//
//	pkHash := owcrypt.Hash(pub, 20, owcrypt.HASH_ALG_BLAKE2B)
//
//	address := handshakeTransaction.AddressEncode(pkHash)
//
//	return address, nil
//}
//
////RedeemScriptToAddress 多重签名赎回脚本转地址
//func (decoder *addressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {
//	return "", nil
//}
//
////WIFToPrivateKey WIF转私钥
//func (decoder *addressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {
//	return nil, nil
//}
//
//
