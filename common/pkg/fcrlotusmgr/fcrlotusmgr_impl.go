/*
Package fcrlotusmgr - lotus manager handles the interaction with filecoin via lotus.
*/
package fcrlotusmgr

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
	"github.com/filecoin-project/go-jsonrpc"
	lotusbig "github.com/filecoin-project/go-state-types/big"
	crypto2 "github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/actors"
	"github.com/filecoin-project/lotus/chain/actors/builtin/paych"
	"github.com/filecoin-project/lotus/chain/types"
	init4 "github.com/filecoin-project/specs-actors/v4/actors/builtin/init"
	paych2 "github.com/filecoin-project/specs-actors/v4/actors/builtin/paych"
	"github.com/ipfs/go-cid"
	"github.com/minio/blake2b-simd"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
)

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

type FCRLotusMgrImpl struct {
	lotusAPIAddr string
	authToken    string
	getLotusAPI  func(authToken, lotusAPIAddr string) (LotusAPI, jsonrpc.ClientCloser, error)
}

func NewFCRLotusMgrImpl(lotusAPIAddr string, authToken string, getLotusAPI func(authToken, lotusAPIAddr string) (LotusAPI, jsonrpc.ClientCloser, error)) FCRLotusMgr {
	if getLotusAPI == nil {
		getLotusAPI = getRemoteLotusAPI
	}
	return &FCRLotusMgrImpl{lotusAPIAddr: lotusAPIAddr, authToken: authToken, getLotusAPI: getLotusAPI}
}

func (mgr *FCRLotusMgrImpl) CreatePaymentChannel(prvKey string, to string, amt *big.Int) (string, error) {
	_, from, err := fcrcrypto.GetPublicKey(prvKey)
	if err != nil {
		return "", err
	}
	fromAddr, err := address.NewFromString(from)
	if err != nil {
		return "", err
	}
	toAddr, err := address.NewFromString(to)
	if err != nil {
		return "", err
	}
	// Get API
	api, closer, err := mgr.getLotusAPI(mgr.authToken, mgr.lotusAPIAddr)
	if err != nil {
		return "", err
	}
	if closer != nil {
		defer closer()
	}
	// Message builder
	builder := paych.Message(actors.Version4, fromAddr)
	msg, err := builder.Create(toAddr, lotusbig.NewFromGo(amt))
	if err != nil {
		return "", err
	}
	// Get signed message
	signedMsg, err := fillMsg(prvKey, api, msg)
	if err != nil {
		return "", err
	}
	contentID, err := api.MpoolPush(context.Background(), signedMsg)
	if err != nil {
		return "", err
	}
	receipt, err := waitReceipt(&contentID, api)
	if err != nil {
		return "", err
	}
	if receipt.ExitCode != 0 {
		return "", fmt.Errorf("Transaction fails to execute: %s", receipt.ExitCode.Error())
	}
	var decodedReturn init4.ExecReturn
	err = decodedReturn.UnmarshalCBOR(bytes.NewReader(receipt.Return))
	if err != nil {
		return "", fmt.Errorf("Payment manager has error unmarshal receipt: %v", receipt)
	}
	return decodedReturn.RobustAddress.String(), nil
}

func (mgr *FCRLotusMgrImpl) TopupPaymentChannel(prvKey string, chAddr string, amt *big.Int) error {
	_, from, err := fcrcrypto.GetPublicKey(prvKey)
	if err != nil {
		return err
	}
	fromAddr, err := address.NewFromString(from)
	if err != nil {
		return err
	}
	toAddr, err := address.NewFromString(chAddr)
	if err != nil {
		return err
	}
	// Get API
	api, closer, err := mgr.getLotusAPI(mgr.authToken, mgr.lotusAPIAddr)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer()
	}
	msg := &types.Message{
		To:     toAddr,
		From:   fromAddr,
		Value:  lotusbig.NewFromGo(amt),
		Method: 0,
	}
	// Get signed message
	signedMsg, err := fillMsg(prvKey, api, msg)
	if err != nil {
		return err
	}
	contentID, err := api.MpoolPush(context.Background(), signedMsg)
	if err != nil {
		return err
	}
	receipt, err := waitReceipt(&contentID, api)
	if err != nil {
		return err
	}
	if receipt.ExitCode != 0 {
		return fmt.Errorf("Transaction fails to execute: %s", receipt.ExitCode.Error())
	}
	return nil
}

func (mgr *FCRLotusMgrImpl) SettlePaymentChannel(prvKey string, chAddr string, voucher string) error {
	return errors.New("No implementation")
}

func (mgr *FCRLotusMgrImpl) CollectPaymentChannel(prvKey string, chAddr string) error {
	return errors.New("No implementation")
}

func (mgr *FCRLotusMgrImpl) CheckPaymentChannel(chAddr string) (bool, *big.Int, string, error) {
	to, err := address.NewFromString(chAddr)
	if err != nil {
		return false, nil, "", err
	}
	// Get API
	api, closer, err := mgr.getLotusAPI(mgr.authToken, mgr.lotusAPIAddr)
	if err != nil {
		return false, nil, "", err
	}
	if closer != nil {
		defer closer()
	}
	// Get actor state
	actor, err := api.StateGetActor(context.Background(), to, types.EmptyTSK)
	if err != nil {
		return false, nil, "", err
	}
	data, err := api.ChainReadObj(context.Background(), actor.Head)
	if err != nil {
		return false, nil, "", err
	}
	state := paych2.State{}
	err = state.UnmarshalCBOR(bytes.NewReader(data))
	if err != nil {
		return false, nil, "", err
	}
	recipient, err := api.StateAccountKey(context.Background(), state.To, types.EmptyTSK)
	if err != nil {
		return false, nil, "", err
	}

	return state.SettlingAt != 0, actor.Balance.Int, recipient.String(), nil
}

func (mgr *FCRLotusMgrImpl) GetCostToCreate(fromAddr string, to string, amt *big.Int) (*big.Int, error) {
	return nil, errors.New("No implementation")
}

func (mgr *FCRLotusMgrImpl) GetCostToSettle(fromAddr string, chAddr string) (*big.Int, error) {
	return nil, errors.New("No implementation")
}

func (mgr *FCRLotusMgrImpl) GetPaymentChannelCreationBlock(chAddr string) (*big.Int, error) {
	return nil, errors.New("No implementation")
}

func (mgr *FCRLotusMgrImpl) GetPaymentChannelSettlementBlock(chAddr string) (*big.Int, error) {
	return nil, errors.New("No implementation")
}

// fillMsg will fill the gas and sign a given message
func fillMsg(prvKeyStr string, api LotusAPI, msg *types.Message) (*types.SignedMessage, error) {
	prvKey, err := hex.DecodeString(prvKeyStr)
	if err != nil {
		return nil, err
	}
	// Get nonce
	nonce, err := api.MpoolGetNonce(context.Background(), msg.From)
	if err != nil {
		return nil, err
	}
	msg.Nonce = nonce

	// Calculate gas
	limit, err := api.GasEstimateGasLimit(context.Background(), msg, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasLimit = int64(float64(limit) * 1.25)

	premium, err := api.GasEstimateGasPremium(context.Background(), 10, msg.From, msg.GasLimit, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasPremium = premium

	feeCap, err := api.GasEstimateFeeCap(context.Background(), msg, 20, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasFeeCap = feeCap

	// Sign message
	sig, err := Sign(prvKey, msg.Cid().Bytes())
	if err != nil {
		return nil, err
	}
	return &types.SignedMessage{
		Message: *msg,
		Signature: crypto2.Signature{
			Type: crypto2.SigTypeSecp256k1,
			Data: sig,
		},
	}, nil
}

// wait receipt will wait until receipt is received for a given cid
func waitReceipt(cid *cid.Cid, api LotusAPI) (*types.MessageReceipt, error) {
	// Return until recipient is returned (transaction is processed)
	var receipt *types.MessageReceipt
	var err error
	for {
		receipt, err = api.StateGetReceipt(context.Background(), *cid, types.EmptyTSK)
		if err != nil {
			return nil, err
		}
		if receipt != nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return receipt, nil
}

// get lotus api returns the api that interacts with lotus for a given lotus api addr and access token
func getRemoteLotusAPI(authToken, lotusAPIAddr string) (LotusAPI, jsonrpc.ClientCloser, error) {
	var api apistruct.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	closer, err := jsonrpc.NewMergeClient(context.Background(), lotusAPIAddr, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		return nil, nil, err
	}
	return &api, closer, nil
}

// Sign signs the given msg using given private key
func Sign(pk []byte, msg []byte) ([]byte, error) {
	b2sum := blake2b.Sum256(msg)
	sig, err := crypto.Sign(pk, b2sum[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}
