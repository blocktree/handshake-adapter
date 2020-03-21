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

package openwtester

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/v2/common/file"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openw"
	"github.com/blocktree/openwallet/v2/openwallet"
	"path/filepath"
	"testing"
)

////////////////////////// 测试单个扫描器 //////////////////////////

type subscriberSingle struct {
}

//BlockScanNotify 新区块扫描完成通知
func (sub *subscriberSingle) BlockScanNotify(header *openwallet.BlockHeader) error {
	log.Notice("header:", header)
	return nil
}

//BlockTxExtractDataNotify 区块提取结果通知
func (sub *subscriberSingle) BlockExtractDataNotify(sourceKey string, data *openwallet.TxExtractData) error {
	log.Notice("account:", sourceKey)

	for i, input := range data.TxInputs {
		log.Std.Notice("data.TxInputs[%d]: %+v", i, input)
	}

	for i, output := range data.TxOutputs {
		log.Std.Notice("data.TxOutputs[%d]: %+v", i, output)
	}

	log.Std.Notice("data.Transaction: %+v", data.Transaction)

	return nil
}


func (sub *subscriberSingle) BlockExtractSmartContractDataNotify(sourceKey string, data *openwallet.SmartContractReceipt) error {
	return nil
}

func TestSubscribeAddress_HNS(t *testing.T) {

	var (
		endRunning = make(chan bool, 1)
		symbol     = "HNS"
		addrs      = map[string]string{
			"hs1qqseet76e8v0met8lgcvznwurlycdukqtdkx6fg": "sender",
			//"hs1qkvpfvra3vvj4uw4ls4d6h4ra58vpnwu9hpr9lt": "reciver",
			//"hs1qqt5y73p5ahvfne4e428lkn8gtuppvfpa6974t6": "reciver",
		}

	)

	var scanAddressFunc openwallet.BlockScanTargetFunc
	scanAddressFunc = func (target openwallet.ScanTarget) (string, bool) {
			key, ok := addrs[target.Address]
			if !ok {
				return "", false
			}
			return key, true
	}
	//GetSourceKeyByAddress 获取地址对应的数据源标识
	//scanAddressFunc := func(address string) (string, bool) {
	//	key, ok := addrs[address]
	//	if !ok {
	//		return "", false
	//	}
	//	return key, true
	//}

	assetsMgr, err := openw.GetAssetsAdapter(symbol)
	if err != nil {
		log.Error(symbol, "is not support")
		return
	}

	//读取配置
	absFile := filepath.Join(configFilePath, symbol+".ini")

	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return
	}
	assetsMgr.LoadAssetsConfig(c)

	assetsLogger := assetsMgr.GetAssetsLogger()
	if assetsLogger != nil {
		assetsLogger.SetLogFuncCall(true)
	}

	//log.Debug("already got scanner:", assetsMgr)
	scanner := assetsMgr.GetBlockScanner()

	if scanner.SupportBlockchainDAI() {
		file.MkdirAll(dbFilePath)
		dai, err := openwallet.NewBlockchainLocal(filepath.Join(dbFilePath, dbFileName), false)
		if err != nil {
			log.Error("NewBlockchainLocal err: %v", err)
			return
		}

		scanner.SetBlockchainDAI(dai)
	}

	scanner.SetRescanBlockHeight(7091)
	//scanner.SetBlockScanAddressFunc(scanAddressFunc)
	scanner.SetBlockScanTargetFunc(scanAddressFunc)

	sub := subscriberSingle{}
	scanner.AddObserver(&sub)

	scanner.Run()

	<-endRunning
}
