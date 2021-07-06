/*
Package fcrlotusmgr - lotus manager handles the interaction with filecoin via lotus.
*/
package fcrlotusmgr

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

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/filecoin-project/go-address"
	lotusbig "github.com/filecoin-project/go-state-types/big"
	crypto2 "github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/chain/actors/builtin/paych"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

// FCRLotusMgr represents the manager that interacts with the lotus.
type FCRLotusMgr interface {
	// CreatePaymentChannel creates a payment channel using the given private key, recipient address and a given amount.
	CreatePaymentChannel(privKey string, recipientAddr string, amt *big.Int) (string, error)

	// TopupPaymentChannel topups a payment channel using the given private key, channel address and a given amount.
	TopupPaymentChannel(privKey string, chAddr string, amt *big.Int) error

	// SettlePaymentChannel settles a payment channel using the given private key, channel address and a final voucher.
	SettlePaymentChannel(privKey string, chAddr string, vouchers []string) error

	// CollectPaymentChannel collects a payment channel using the given private key, channel address.
	CollectPaymentChannel(privKey string, chAddr string) error

	// CheckPaymentChannel checks the state of a channel.
	// It returns a boolean indicating if the channel is settling/settled, the channel balance, the recipient address and error.
	CheckPaymentChannel(chAddr string) (bool, *big.Int, string, error)

	// GetCostToCreate gets the current cost to create a payment channel.
	GetCostToCreate(privKey string, recipientAddr string, amt *big.Int) (*big.Int, error)

	// GetCostToSettle gets the current cost to settle a payment channel + updating voucher.
	GetCostToSettle(privKey string, chAddr string, vouchers []string) (*big.Int, error)

	// GetPaymentChannelCreationBlock gets the block number at which given payment channel is created.
	GetPaymentChannelCreationBlock(chAddr string) (*big.Int, error)

	// GetPaymentChannelSettlementBlock gets the block number at which given payment channel is called to settle.
	GetPaymentChannelSettlementBlock(chAddr string) (*big.Int, error)
}

// LotusAPI is the minimum interface interacting with the Lotus to achieve payment function.
type LotusAPI interface {
	ChainReadObj(ctx context.Context, obj cid.Cid) ([]byte, error)

	GasEstimateFeeCap(ctx context.Context, msg *types.Message, maxqueueblks int64, tsk types.TipSetKey) (types.BigInt, error)

	GasEstimateGasLimit(ctx context.Context, msg *types.Message, tsk types.TipSetKey) (int64, error)

	GasEstimateGasPremium(ctx context.Context, nblocksincl uint64, sender address.Address, gaslimit int64, tsk types.TipSetKey) (types.BigInt, error)

	MpoolGetNonce(ctx context.Context, addr address.Address) (uint64, error)

	MpoolPush(ctx context.Context, smsg *types.SignedMessage) (cid.Cid, error)

	StateAccountKey(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error)

	StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error)

	StateGetReceipt(ctx context.Context, msg cid.Cid, tsk types.TipSetKey) (*types.MessageReceipt, error)
}

// GenerateVoucher generates a voucher by given private key, channel address, lane number and amount.
func GenerateVoucher(privKey string, chAddr string, lane uint64, nonce uint64, newRedeemed *big.Int) (string, error) {
	addr, err := address.NewFromString(chAddr)
	if err != nil {
		return "", err
	}
	sv := &paych.SignedVoucher{
		ChannelAddr: addr,
		Lane:        lane,
		Nonce:       nonce,
		Amount:      lotusbig.NewFromGo(newRedeemed),
	}
	vb, err := sv.SigningBytes()
	if err != nil {
		return "", err
	}
	pk, err := hex.DecodeString(privKey)
	if err != nil {
		return "", err
	}
	sig, err := Sign(pk, vb)
	sv.Signature = &crypto2.Signature{
		Type: crypto2.SigTypeSecp256k1,
		Data: sig,
	}
	voucher, err := encodedVoucher(sv)
	if err != nil {
		return "", err
	}
	return voucher, nil
}

// VerifyVoucher verifies a voucher by given voucher.
// It returns the sender's address, channel address, lane, nonce, new redeemed, and error)
func VerifyVoucher(voucher string) (string, string, uint64, uint64, *big.Int, error) {
	sv, err := paych.DecodeSignedVoucher(voucher)
	if err != nil {
		return "", "", 0, 0, nil, err
	}
	vb, err := sv.SigningBytes()
	if err != nil {
		return "", "", 0, 0, nil, err
	}
	if sv.Signature.Type != crypto2.SigTypeSecp256k1 {
		return "", "", 0, 0, nil, errors.New("Wrong signature type")
	}
	sender, err := Verify(sv.Signature.Data, vb)
	if err != nil {
		return "", "", 0, 0, nil, err
	}
	return sender, sv.ChannelAddr.String(), sv.Lane, sv.Nonce, sv.Amount.Int, nil
}
