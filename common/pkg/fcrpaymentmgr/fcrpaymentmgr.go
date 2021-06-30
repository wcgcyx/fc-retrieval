/*
Package fcrpaymentmgr - payment manager manages all payment related functions.
*/
package fcrpaymentmgr

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

import "math/big"

type FCRPaymentMgr interface {
	// Start starts the manager's routine.
	Start() error

	// Shutdown ends the manager's routine safely.
	Shutdown()

	/* For outbound payment */
	// Create creates a payment channel with given initial balance to given recipient.
	Create(recipientAddr string, amt *big.Int) error

	// Topup topups an existing payment channel with given amount.
	Topup(recipientAddr string, amt *big.Int) error

	// Pay pays a given recipient in given land with given amount.
	// It returns the voucher,
	// a boolean indicates whether or not to create a payment channel,
	// a boolean indicates whether or not to topup a payment channel,
	// error if any.
	Pay(recipientAddr string, lane uint64, amt *big.Int) (string, bool, bool, error)

	// ReceiveRefund receives a refund from recipient.
	// It returns the amount received from the refund.
	ReceiveRefund(recipientAddr string, voucher string) (*big.Int, error)

	// GetOutboundChStatus gets the outbound payment channel status by a given recipient addr.
	// It returns the payment channel address, balance and redeemed amount.
	GetOutboundChStatus(recipientAddr string) (string, *big.Int, *big.Int, error)

	// RemoveOutboundCh removes the outbound payment channel status by a given recipient addr.
	RemoveOutboundCh(recipientAddr string) error

	// GetCostToCreate gets the current cost to create a payment channel.
	GetCostToCreate(recipientAddr string, amt *big.Int) (*big.Int, error)

	// CheckRecipientSettlementValidity checks if it is valid for the recipient to settle a selling payment channel.
	CheckRecipientSettlementValidity(recipientAddr string) (bool, error)

	/* For inbound payment */
	// Settle settles a payment channel by given sender addr.
	Settle(senderAddr string) error

	// Receive receives a payment. It returns the amount received and the lane number.
	Receive(senderAddr string, voucher string) (*big.Int, uint64, error)

	// Refund creates a voucher used to refund by given sender addr, given lane and given amount.
	Refund(senderAddr string, lane uint64, amt *big.Int) (string, error)

	// RemoveInboundCh removes the inbound payment channel status by a given sender addr.
	RemoveInboundCh(senderAddr string) error

	// GetInboundChStatus gets the inbound payment channel status by a given sender addr.
	// It returns the payment channel address, balance and redeemed amount.
	GetInboundChStatus(senderAddr string) (string, *big.Int, *big.Int, error)

	// GetCostToSettle gets the current cost to settle a payment channel.
	GetCostToSettle(senderAddr string) (*big.Int, error)

	// CheckSettlementValidity checks if it is valid to settle a payment channel.
	CheckSettlementValidity(senderAddr string) (bool, error)
}
