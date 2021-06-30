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

import (
	"encoding/hex"
	"errors"
	"math/big"
	"sync"

	"github.com/filecoin-project/go-address"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrlotusmgr"
)

// FCRPaymentMgrImplV1 implements FCRPaymentMgr, it is an in-memory version.
type FCRPaymentMgrImplV1 struct {
	prvKey string
	addr   string

	lotusMgr fcrlotusmgr.FCRLotusMgr

	// Channel states.
	// map[recipient addr] -> channel state
	outboundChs     map[string]*channelState
	outboundChsLock sync.RWMutex
	// map[sender addr] -> channel state
	inboundChs     map[string]*channelState
	inboundChsLock sync.RWMutex
}

// channelState represents the state of a channel
type channelState struct {
	addr     string
	balance  big.Int
	redeemed big.Int
	lock     sync.RWMutex

	// Lane States.
	// map[lane id] -> lane state
	laneStates map[uint64]*laneState
}

// laneState represents the state of a lane
type laneState struct {
	nonce    uint64
	redeemed big.Int
	vouchers []string
}

func NewFCRPaymentMgrImplV1(prvKey string, lotusMgr fcrlotusmgr.FCRLotusMgr) FCRPaymentMgr {
	return &FCRPaymentMgrImplV1{
		prvKey:          prvKey,
		lotusMgr:        lotusMgr,
		outboundChs:     make(map[string]*channelState),
		outboundChsLock: sync.RWMutex{},
		inboundChs:      make(map[string]*channelState),
		inboundChsLock:  sync.RWMutex{},
	}
}

func (mgr *FCRPaymentMgrImplV1) Start() error {
	pubKey, _, err := fcrcrypto.GetPublicKey(mgr.prvKey)
	if err != nil {
		return err
	}
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return err
	}
	addr, err := address.NewSecp256k1Address(pubKeyBytes)
	if err != nil {
		return err
	}
	mgr.addr = addr.String()
	return nil
}

func (mgr *FCRPaymentMgrImplV1) Shutdown() {
}

func (mgr *FCRPaymentMgrImplV1) Create(recipientAddr string, amt *big.Int) error {
	mgr.outboundChsLock.RLock()
	_, ok := mgr.outboundChs[recipientAddr]
	mgr.outboundChsLock.RUnlock()
	if ok {
		return errors.New("There is an existing channel for given recipient")
	}
	chAddr, err := mgr.lotusMgr.CreatePaymentChannel(mgr.prvKey, recipientAddr, amt)
	if err != nil {
		return err
	}
	mgr.outboundChsLock.Lock()
	defer mgr.outboundChsLock.Unlock()
	mgr.outboundChs[recipientAddr] = &channelState{
		addr:       chAddr,
		balance:    *amt,
		redeemed:   *big.NewInt(0),
		lock:       sync.RWMutex{},
		laneStates: make(map[uint64]*laneState),
	}
	return nil
}

func (mgr *FCRPaymentMgrImplV1) Topup(recipientAddr string, amt *big.Int) error {
	mgr.outboundChsLock.RLock()
	defer mgr.outboundChsLock.RUnlock()
	cs, ok := mgr.outboundChs[recipientAddr]
	if !ok {
		return errors.New("There is no existing channel for given recipient")
	}
	cs.lock.Lock()
	defer cs.lock.Unlock()
	err := mgr.lotusMgr.TopupPaymentChannel(mgr.prvKey, cs.addr, amt)
	if err != nil {
		return err
	}
	// Update channel state
	cs.balance.Add(&cs.balance, amt)
	return nil
}

func (mgr *FCRPaymentMgrImplV1) Pay(recipientAddr string, lane uint64, amt *big.Int) (string, bool, bool, error) {
	mgr.outboundChsLock.RLock()
	defer mgr.outboundChsLock.RUnlock()
	cs, ok := mgr.outboundChs[recipientAddr]
	if !ok {
		return "", true, false, nil
	}
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cNewRedeemed := big.NewInt(0).Add(&cs.redeemed, amt)
	if cs.balance.Cmp(cNewRedeemed) < 0 {
		// Balance not enough
		return "", false, true, nil
	}
	// Blanace is enough
	// Check lane state
	ls, ok := cs.laneStates[lane]
	if !ok {
		// Lane not existed, create a new lane
		ls = &laneState{
			nonce:    0,
			redeemed: *big.NewInt(0),
			vouchers: make([]string, 0),
		}
		cs.laneStates[lane] = ls
	}
	// Create a voucher
	lNewRedeemed := big.NewInt(0).Add(&ls.redeemed, amt)
	voucher, err := mgr.lotusMgr.GenerateVoucher(mgr.prvKey, cs.addr, lane, ls.nonce, lNewRedeemed)
	if err != nil {
		return "", false, false, err
	}
	// Update lane state
	ls.nonce++
	ls.redeemed.Add(&ls.redeemed, amt)
	ls.vouchers = append(ls.vouchers, voucher)
	// Update channel state
	cs.redeemed.Add(&cs.redeemed, amt)
	return voucher, false, false, nil
}

func (mgr *FCRPaymentMgrImplV1) ReceiveRefund(recipientAddr string, voucher string) (*big.Int, error) {
	senderVAddr, chAddr, lane, nonce, newRedeemed, err := mgr.lotusMgr.VerifyVoucher(voucher)
	if err != nil {
		return nil, err
	}
	if recipientAddr != senderVAddr {
		return nil, errors.New("Refund sender address mismtach")
	}
	mgr.outboundChsLock.RLock()
	defer mgr.outboundChsLock.RUnlock()
	cs, ok := mgr.outboundChs[recipientAddr]
	if !ok {
		return nil, errors.New("Channel not found")
	}
	cs.lock.Lock()
	defer cs.lock.Unlock()
	if chAddr != cs.addr {
		return nil, errors.New("Refund channel address mismatch")
	}
	ls, ok := cs.laneStates[lane]
	if !ok {
		return nil, errors.New("Refund lane not existed")
	}
	if ls.nonce > nonce {
		return nil, errors.New("Refund nonce is not valid")
	}
	diff := big.NewInt(0).Sub(&ls.redeemed, newRedeemed)
	if diff.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.New("Refund value is not positive")
	}
	// Refund is valid, update lane state and channel state
	ls.nonce = nonce
	ls.redeemed.Sub(&ls.redeemed, diff)
	ls.vouchers = append(ls.vouchers, voucher)
	cs.redeemed.Sub(&cs.redeemed, diff)
	return diff, nil
}

func (mgr *FCRPaymentMgrImplV1) GetOutboundChStatus(recipientAddr string) (string, *big.Int, *big.Int, error) {
	mgr.outboundChsLock.RLock()
	defer mgr.outboundChsLock.RUnlock()
	cs, ok := mgr.outboundChs[recipientAddr]
	if !ok {
		return "", nil, nil, errors.New("Channel not found")
	}
	cs.lock.Lock()
	defer cs.lock.Unlock()
	return cs.addr, big.NewInt(0).Set(&cs.balance), big.NewInt(0).Set(&cs.redeemed), nil
}

func (mgr *FCRPaymentMgrImplV1) RemoveOutboundCh(recipientAddr string) error {
	mgr.outboundChsLock.RLock()
	_, ok := mgr.outboundChs[recipientAddr]
	mgr.outboundChsLock.RUnlock()
	if !ok {
		return nil
	}
	mgr.outboundChsLock.Lock()
	defer mgr.outboundChsLock.Unlock()
	delete(mgr.outboundChs, recipientAddr)
	return nil
}

func (mgr *FCRPaymentMgrImplV1) GetCostToCreate(recipientAddr string, amt *big.Int) (*big.Int, error) {
	return mgr.lotusMgr.GetCostToCreate(mgr.prvKey, recipientAddr, amt)
}

func (mgr *FCRPaymentMgrImplV1) CheckRecipientSettlementValidity(recipientAddr string) (bool, error) {
	return false, errors.New("No implementation")
}

func (mgr *FCRPaymentMgrImplV1) Settle(senderAddr string) error {
	mgr.inboundChsLock.Lock()
	defer mgr.inboundChsLock.Unlock()
	cs, ok := mgr.inboundChs[senderAddr]
	if !ok {
		return errors.New("Channel not found")
	}
	vouchers := make([]string, 0)
	for _, ls := range cs.laneStates {
		vouchers = append(vouchers, ls.vouchers[len(ls.vouchers)-1])
	}
	err := mgr.lotusMgr.SettlePaymentChannel(mgr.prvKey, cs.addr, vouchers)
	if err != nil {
		return err
	}
	delete(mgr.inboundChs, senderAddr)
	return nil
}

func (mgr *FCRPaymentMgrImplV1) Receive(senderAddr string, voucher string) (*big.Int, uint64, error) {
	senderVAddr, chAddr, lane, nonce, newRedeemed, err := mgr.lotusMgr.VerifyVoucher(voucher)
	if err != nil {
		return nil, 0, err
	}
	if senderAddr != senderVAddr {
		return nil, 0, errors.New("Receive sender address mismtach")
	}
	mgr.inboundChsLock.RLock()
	cs, ok := mgr.inboundChs[senderAddr]
	if !ok {
		// Need to create a new entry
		// Get channel address
		_, balance, recipientAddr, err := mgr.lotusMgr.CheckPaymentChannel(chAddr)
		if err != nil {
			return nil, 0, err
		}
		if recipientAddr != mgr.addr {
			return nil, 0, errors.New("Receive sender address mismatch")
		}
		mgr.inboundChsLock.RUnlock()
		mgr.inboundChsLock.Lock()
		cs = &channelState{
			addr:       chAddr,
			balance:    *balance,
			redeemed:   *big.NewInt(0),
			lock:       sync.RWMutex{},
			laneStates: make(map[uint64]*laneState),
		}
		mgr.inboundChs[senderAddr] = cs
		mgr.inboundChsLock.Unlock()
		mgr.inboundChsLock.RLock()
	}
	defer mgr.inboundChsLock.RUnlock()
	cs.lock.Lock()
	defer cs.lock.Unlock()
	if chAddr != cs.addr {
		return nil, 0, errors.New("Receive channel address mismatch")
	}
	ls, ok := cs.laneStates[lane]
	if !ok {
		// Need to create a new lane
		ls = &laneState{
			nonce:    0,
			redeemed: *big.NewInt(0),
			vouchers: make([]string, 0),
		}
		cs.laneStates[lane] = ls
	}
	if ls.nonce > nonce {
		// Nonce not match
		return nil, 0, errors.New("Receive nonce is not valid")
	}
	paymentValue := big.NewInt(0).Sub(newRedeemed, &ls.redeemed)
	if paymentValue.Cmp(big.NewInt(0)) <= 0 {
		return nil, 0, errors.New("Receive has bad amount")
	}
	// Verify channel balance
	csNewRedeemed := big.NewInt(0).Add(&cs.redeemed, paymentValue)
	if cs.balance.Cmp(csNewRedeemed) < 0 {
		// Channel balance is not enough
		// Update one time.
		_, balance, _, err := mgr.lotusMgr.CheckPaymentChannel(chAddr)
		if err != nil {
			return nil, 0, err
		}
		cs.balance = *balance
		if cs.balance.Cmp(csNewRedeemed) < 0 {
			return nil, 0, errors.New("Receive not enough channel balance")
		}
	}
	// Voucher is valid, update lane state and channel state
	ls.nonce = nonce + 1
	ls.redeemed = *newRedeemed
	ls.vouchers = append(ls.vouchers, voucher)
	// Update channel state
	cs.redeemed.Add(&cs.redeemed, paymentValue)
	return paymentValue, lane, nil
}

func (mgr *FCRPaymentMgrImplV1) Refund(senderAddr string, lane uint64, amt *big.Int) (string, error) {
	mgr.inboundChsLock.RLock()
	defer mgr.inboundChsLock.RUnlock()
	cs, ok := mgr.inboundChs[senderAddr]
	if !ok {
		return "", errors.New("Channel not found")
	}
	cs.lock.Lock()
	defer cs.lock.Unlock()
	// Check lane state
	ls, ok := cs.laneStates[lane]
	if !ok {
		return "", errors.New("Refund lane not existed")
	}
	lNewRedeemed := big.NewInt(0).Sub(&ls.redeemed, amt)
	if lNewRedeemed.Cmp(big.NewInt(0)) < 0 {
		// Refund too much
		return "", errors.New("")
	}
	// Create a voucher
	voucher, err := mgr.lotusMgr.GenerateVoucher(mgr.prvKey, cs.addr, lane, ls.nonce, lNewRedeemed)
	if err != nil {
		return "", err
	}
	// Update lane state
	ls.nonce++
	ls.redeemed.Sub(&ls.redeemed, amt)
	ls.vouchers = append(ls.vouchers, voucher)
	cs.redeemed.Sub(&cs.redeemed, amt)
	return voucher, nil
}

func (mgr *FCRPaymentMgrImplV1) GetInboundChStatus(senderAddr string) (string, *big.Int, *big.Int, error) {
	mgr.inboundChsLock.RLock()
	defer mgr.inboundChsLock.RUnlock()
	cs, ok := mgr.inboundChs[senderAddr]
	if !ok {
		return "", nil, nil, errors.New("Channel not found")
	}
	cs.lock.Lock()
	defer cs.lock.Unlock()
	return cs.addr, big.NewInt(0).Set(&cs.balance), big.NewInt(0).Set(&cs.redeemed), nil
}

func (mgr *FCRPaymentMgrImplV1) RemoveInboundCh(senderAddr string) error {
	mgr.inboundChsLock.RLock()
	_, ok := mgr.inboundChs[senderAddr]
	mgr.inboundChsLock.RUnlock()
	if !ok {
		return nil
	}
	mgr.inboundChsLock.Lock()
	defer mgr.inboundChsLock.Unlock()
	delete(mgr.inboundChs, senderAddr)
	return nil
}

func (mgr *FCRPaymentMgrImplV1) GetCostToSettle(senderAddr string) (*big.Int, error) {
	return nil, errors.New("No implementation")
}

func (mgr *FCRPaymentMgrImplV1) CheckSettlementValidity(senderAddr string) (bool, error) {
	return false, errors.New("No implementation")
}
