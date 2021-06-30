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
	"math/big"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

// FCRLotusMgr represents the manager that interacts with the lotus.
type FCRLotusMgr interface {
	// CreatePaymentChannel creates a payment channel using the given private key, recipient address and a given amount.
	CreatePaymentChannel(prvKey string, to string, amt *big.Int) (string, error)

	// TopupPaymentChannel topups a payment channel using the given private key, channel address and a given amount.
	TopupPaymentChannel(prvKey string, chAddr string, amt *big.Int) error

	// SettlePaymentChannel settles a payment channel using the given private key, channel address and a final voucher.
	SettlePaymentChannel(prvKey string, chAddr string, voucher string) error

	// CollectPaymentChannel collects a payment channel using the given private key, channel address.
	CollectPaymentChannel(prvKey string, chAddr string) error

	// CheckPaymentChannel checks the state of a channel.
	// It returns a boolean indicating if the channel is settling/settled, the channel balance, the recipient address and error.
	CheckPaymentChannel(chAddr string) (bool, *big.Int, string, error)

	// GetCostToCreate gets the current cost to create a payment channel.
	GetCostToCreate(fromAddr string, to string, amt *big.Int) (*big.Int, error)

	// GetCostToSettle gets the current cost to settle a payment channel + updating voucher.
	GetCostToSettle(fromAddr string, chAddr string) (*big.Int, error)

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
