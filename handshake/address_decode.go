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
	"github.com/blocktree/openwallet/openwallet"
)

func init() {

}

type AddressDecoder interface {
	openwallet.AddressDecoder
	ScriptPubKeyToBech32Address(scriptPubKey []byte) (string, error)
}

type addressDecoder struct {
	wm *WalletManager //钱包管理者
}

//NewAddressDecoder 地址解析器
func NewAddressDecoder(wm *WalletManager) *addressDecoder {
	decoder := addressDecoder{}
	decoder.wm = wm
	return &decoder
}

//PrivateKeyToWIF 私钥转WIF
func (decoder *addressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {
	return "", nil
}

//PublicKeyToAddress 公钥转地址
func (decoder *addressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {

	pkHash := owcrypt.Hash(pub, 20, owcrypt.HASH_ALG_BLAKE2B)

	address := handshakeTransaction.AddressEncode(pkHash)

	return address, nil
}

//RedeemScriptToAddress 多重签名赎回脚本转地址
func (decoder *addressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {
	return "", nil
}

//WIFToPrivateKey WIF转私钥
func (decoder *addressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {
	return nil, nil
}

//ScriptPubKeyToBech32Address scriptPubKey转Bech32地址
func (decoder *addressDecoder) ScriptPubKeyToBech32Address(scriptPubKey []byte) (string, error) {
	return "", nil
}

//ScriptPubKeyToBech32Address scriptPubKey转Bech32地址
func scriptPubKeyToBech32Address(scriptPubKey []byte, isTestNet bool) (string, error) {
	return "", nil
}
