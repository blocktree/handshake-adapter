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
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/blocktree/handshake-adapter/handshakeTransaction"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/imroc/req"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"strings"
)

type ClientInterface interface {
	Call(path string, request []interface{}) (*gjson.Result, error)
}

// A Client is a Handshake RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	client      *req.Req
	//Client *req.Req
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Id      string      `json:"id,omitempty"`
}

func NewClient(url, token string, debug bool) *Client {
	c := Client{
		BaseURL:     url,
		AccessToken: token,
		Debug:       debug,
	}

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *Client) Call(path string, request []interface{}) (*gjson.Result, error) {

	var (
		body = make(map[string]interface{}, 0)
	)

	if c.client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + c.AccessToken,
	}

	//json-rpc
	body["jsonrpc"] = "2.0"
	body["id"] = "1"
	body["method"] = path
	body["params"] = request

	if c.Debug {
		log.Std.Info("Start Request API...")
	}

	if request == nil && path != "getblockcount" {
		r, err := c.client.Get(c.BaseURL+path, nil, authHeader)
		if c.Debug {
			log.Std.Info("Request API Completed")
		}

		if c.Debug {
			log.Std.Info("%+v", r)
		}

		if err != nil {
			return nil, err
		}

		resp := gjson.ParseBytes(r.Bytes())
		if strings.Index(resp.String(), "[") != 0 {
			return nil, errors.New("respond invalid!")
		}
		return &resp, nil
	}

	r, err := c.client.Post(c.BaseURL, req.BodyJSON(&body), authHeader)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())
	err = isError(&resp)
	if err != nil {
		return nil, err
	}

	result := resp.Get("result")

	return &result, nil
}

// See 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//isError 是否报错
func isError(result *gjson.Result) error {
	var (
		err error
	)

	/*
		//failed 返回错误
		{
			"result": null,
			"error": {
				"code": -8,
				"message": "Block height out of range"
			},
			"id": "foo"
		}
	*/

	if !result.Get("error").IsObject() {

		if !result.Get("result").Exists() {
			return errors.New("Response is empty! ")
		}

		return nil
	}

	errInfo := fmt.Sprintf("[%d]%s",
		result.Get("error.code").Int(),
		result.Get("error.message").String())
	err = errors.New(errInfo)

	return err
}

func (c Client) getBlockHeight() (uint64, error) {

	result, err := c.Call("getblockcount", nil)

	if err != nil {
		result, err = c.Call("getblockcount", nil)
		if err != nil {
			result, err = c.Call("getblockcount", nil)
			if err != nil {
				return 0, err
			}
		}
	}

	return result.Uint(), nil

}

func (c Client) getBlockHash(height uint64) (string, error) {
	request := []interface{}{
		height,
	}

	result, err := c.Call("getblockhash", request)
	if err != nil {
		result, err = c.Call("getblockhash", request)
		if err != nil {
			result, err = c.Call("getblockhash", request)
			if err != nil {
				return "", err
			}
		}
	}

	return result.String(), nil
}

func (c Client) getBlock(hash string) (*Block, error) {
	request := []interface{}{
		hash,
		1,
		0,
	}
	result, err := c.Call("getblock", request)
	if err != nil {
		result, err = c.Call("getblock", request)
		if err != nil {
			result, err = c.Call("getblock", request)
			if err != nil {
				return nil, err
			}
		}
	}

	return newBlock(result), nil
}
func (c Client) getTxsInMemPool() ([]string, error) {
	request := []interface{}{
		0,
	}
	result, err := c.Call("getrawmempool", request)
	if err != nil {
		result, err = c.Call("getrawmempool", request)
		if err != nil {
			result, err = c.Call("getrawmempool", request)
			if err != nil {
				return nil, err
			}
		}
	}
	txs := result.Array()
	ret := []string{}

	for _, tx := range txs {
		ret = append(ret, tx.String())
	}
	return ret, nil
}

func (c Client) getTransaction(txid string) (*Transaction, error) {
	request := []interface{}{
		txid,
		1,
	}
	result, err := c.Call("getrawtransaction", request)
	if err != nil {
		result, err = c.Call("getrawtransaction", request)
		if err != nil {
			result, err = c.Call("getrawtransaction", request)
			if err != nil {
				return nil, err
			}
		}
	}

	return c.newTx(result)
}

func (c Client) getAddressAndAmoutOfInput(txid string, vout uint64) (string, string, error) {
	request := []interface{}{
		txid,
		1,
	}
	result, err := c.Call("getrawtransaction", request)
	if err != nil {
		result, err = c.Call("getrawtransaction", request)
		if err != nil {
			result, err = c.Call("getrawtransaction", request)
			if err != nil {
				return "", "", err
			}
		}
	}

	outs := result.Get("vout").Array()

	if int(vout) > len(outs) {
		return "", "", errors.New("vout is too big")
	}

	hashStr := outs[int(vout)].Get("address").Get("hash").String()
	hash, _ := hex.DecodeString(hashStr)
	if outs[int(vout)].Get("address").Get("version").Uint() != 0 || err != nil {
		return "", "", errors.New("unknown address version")
	}

	return handshakeTransaction.AddressEncode(hash), outs[int(vout)].Get("value").String(), nil
}

func (c Client) getVout(txid string, vout uint64) (*Vout, error) {
	request := []interface{}{
		txid,
		1,
	}
	result, err := c.Call("getrawtransaction", request)
	if err != nil {
		result, err = c.Call("getrawtransaction", request)
		if err != nil {
			result, err = c.Call("getrawtransaction", request)
			if err != nil {
				return nil, err
			}
		}
	}

	outs := result.Get("vout").Array()

	if int(vout) > len(outs) {
		return nil, errors.New("vout is too big")
	}

	hashStr := outs[int(vout)].Get("address").Get("hash").String()
	hash, _ := hex.DecodeString(hashStr)
	if outs[int(vout)].Get("address").Get("version").Uint() != 0 || err != nil {
		return nil, errors.New("unknown address version")
	}

	return &Vout{
		N:            vout,
		Addr:         handshakeTransaction.AddressEncode(hash),
		Value:        outs[int(vout)].Get("value").String(),
		ScriptPubKey: hashStr,
		Type:         outs[int(vout)].Get("covenant").Get("type").String(),
		Action:       outs[int(vout)].Get("covenant").Get("action").String(),
	}, nil
}

func (c Client) listUnspend(addresses ...string) ([]*Unspent, error) {

	var (
		utxos = make([]*Unspent, 0)
	)
	for _, address := range addresses {
		path := "/coin/address/" + address

		result, err := c.Call(path, nil)
		if err != nil {
			result, err = c.Call(path, nil)
			if err != nil {
				result, err = c.Call(path, nil)
				if err != nil {
					return nil, err
				}
			}
		}

		array := result.Array()
		for _, a := range array {
			u := NewUnspent(&a)
			if u.Type == 0 && u.Action == "NONE" {
				utxos = append(utxos, NewUnspent(&a))
			}
		}
	}

	return utxos, nil
}

func (c Client) getEstimateFeeRate() (decimal.Decimal, error) {
	request := []interface{}{
		10,
	}
	result, err := c.Call("estimatesmartfee", request)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return decimal.NewFromString(result.Get("fee").String())
}

func (c Client) sendTransaction(rawHex string) (string, error) {
	request := []interface{}{
		rawHex,
	}

	result, err := c.Call("sendrawtransaction", request)
	if err != nil {
		result, err = c.Call("sendrawtransaction", request)
		if err != nil {
			result, err = c.Call("sendrawtransaction", request)
			if err != nil {
				return "", nil
			}
		}
	}

	return result.String(), nil
}
