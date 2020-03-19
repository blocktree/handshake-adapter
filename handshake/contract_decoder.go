package handshake

import "github.com/blocktree/openwallet/v2/openwallet"

type ContractDecoder struct {
	*openwallet.SmartContractDecoderBase
	wm *WalletManager
}

func NewContractDecoder(wm *WalletManager) *ContractDecoder {
	decoder := ContractDecoder{}
	decoder.wm = wm
	return &decoder
}

func (decoder *ContractDecoder) GetTokenBalanceByAddress(contract openwallet.SmartContract, address ...string) ([]*openwallet.TokenBalance, error) {
	return nil, nil
}
