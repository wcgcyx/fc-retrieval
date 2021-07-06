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
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLotusMgr struct {
	createPaymentChannel func(privKey string, recipientAddr string, amt *big.Int) (string, error)

	topupPaymentChannel func(privKey string, chAddr string, amt *big.Int) error

	settlePaymentChannel func(privKey string, chAddr string, vouchers []string) error

	collectPaymentChannel func(privKey string, chAddr string) error

	checkPaymentChannel func(chAddr string) (bool, *big.Int, string, error)

	getCostToCreate func(privKey string, recipientAddr string, amt *big.Int) (*big.Int, error)

	getCostToSettle func(privKey string, chAddr string, vouchers []string) (*big.Int, error)

	getPaymentChannelCreationBlock func(chAddr string) (*big.Int, error)

	getPaymentChannelSettlementBlock func(chAddr string) (*big.Int, error)
}

func (m *mockLotusMgr) CreatePaymentChannel(privKey string, recipientAddr string, amt *big.Int) (string, error) {
	return m.createPaymentChannel(privKey, recipientAddr, amt)
}

func (m *mockLotusMgr) TopupPaymentChannel(privKey string, chAddr string, amt *big.Int) error {
	return m.topupPaymentChannel(privKey, chAddr, amt)
}

func (m *mockLotusMgr) SettlePaymentChannel(privKey string, chAddr string, vouchers []string) error {
	return m.settlePaymentChannel(privKey, chAddr, vouchers)
}

func (m *mockLotusMgr) CollectPaymentChannel(privKey string, chAddr string) error {
	return m.collectPaymentChannel(privKey, chAddr)
}

func (m *mockLotusMgr) CheckPaymentChannel(chAddr string) (bool, *big.Int, string, error) {
	return m.checkPaymentChannel(chAddr)
}

func (m *mockLotusMgr) GetCostToCreate(privKey string, recipientAddr string, amt *big.Int) (*big.Int, error) {
	return m.getCostToCreate(privKey, recipientAddr, amt)
}

func (m *mockLotusMgr) GetCostToSettle(privKey string, chAddr string, vouchers []string) (*big.Int, error) {
	return m.getCostToSettle(privKey, chAddr, vouchers)
}

func (m *mockLotusMgr) GetPaymentChannelCreationBlock(chAddr string) (*big.Int, error) {
	return m.getPaymentChannelCreationBlock(chAddr)
}

func (m *mockLotusMgr) GetPaymentChannelSettlementBlock(chAddr string) (*big.Int, error) {
	return m.GetPaymentChannelSettlementBlock(chAddr)
}

func TestNewPaymentMgr(t *testing.T) {
	mockLotusMgr := mockLotusMgr{}

	mgr := NewFCRPaymentMgrImplV1("ppp", &mockLotusMgr)
	err := mgr.Start()
	assert.NotEmpty(t, err)

	mgr = NewFCRPaymentMgrImplV1("933dfc0be9ca2d783446fa3fa9ea27bd9cc553ec5131256dd6fddcde3302b9e0", &mockLotusMgr)
	err = mgr.Start()
	assert.Empty(t, err)
	defer mgr.Shutdown()
}

func TestPayAndReceive(t *testing.T) {
	mockLotusMgr := mockLotusMgr{
		createPaymentChannel: func(privKey string, recipientAddr string, amt *big.Int) (string, error) {
			return "f12yybez3cfe2yb2nsartagpwkk23q5hmmiluqafi", nil
		},
		topupPaymentChannel: func(privKey string, chAddr string, amt *big.Int) error {
			return nil
		},
		checkPaymentChannel: func(chAddr string) (bool, *big.Int, string, error) {
			return false, big.NewInt(100000000), "f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", nil
		},
	}
	mgr1 := NewFCRPaymentMgrImplV1("933dfc0be9ca2d783446fa3fa9ea27bd9cc553ec5131256dd6fddcde3302b9e0", &mockLotusMgr)
	err := mgr1.Start()
	assert.Empty(t, err)
	defer mgr1.Shutdown()
	mgr2 := NewFCRPaymentMgrImplV1("8495f24f3bfab01404671400d876d2887314086d4fd73792e52c46386039ec32", &mockLotusMgr)
	err = mgr2.Start()
	assert.Empty(t, err)
	defer mgr2.Shutdown()

	_, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(-50000000))
	assert.NotEmpty(t, err)

	mgr1.RevertPay("f2wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0)
	mgr1.RevertPay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0)

	_, create, _, err := mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(50000000))
	assert.Empty(t, err)
	assert.True(t, create)

	err = mgr1.Topup("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", big.NewInt(100000000))
	assert.NotEmpty(t, err)

	err = mgr1.Create("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", big.NewInt(100000000))
	assert.Empty(t, err)

	err = mgr1.Create("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", big.NewInt(100000000))
	assert.NotEmpty(t, err)

	_, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(50000000))
	assert.Empty(t, err)
	mgr1.RevertPay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 1)
	mgr1.RevertPay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0)
	voucher, _, _, err := mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(50000000))
	assert.Empty(t, err)

	received, lane, err := mgr2.Receive("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", voucher)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), lane)
	assert.Equal(t, "50000000", received.String())

	_, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(10000000))
	assert.Empty(t, err)
	mgr1.RevertPay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0)
	voucher, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(10000000))
	assert.Empty(t, err)

	received, lane, err = mgr2.Receive("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", voucher)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), lane)
	assert.Equal(t, "10000000", received.String())

	voucher, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(10000000))
	assert.Empty(t, err)
	received, lane, err = mgr2.Receive("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", voucher)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), lane)
	assert.Equal(t, "10000000", received.String())

	voucher, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(10000000))
	assert.Empty(t, err)
	received, lane, err = mgr2.Receive("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", voucher)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), lane)
	assert.Equal(t, "10000000", received.String())

	voucher, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(10000000))
	assert.Empty(t, err)
	received, lane, err = mgr2.Receive("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", voucher)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), lane)
	assert.Equal(t, "10000000", received.String())

	voucher, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(10000000))
	assert.Empty(t, err)
	received, lane, err = mgr2.Receive("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", voucher)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), lane)
	assert.Equal(t, "10000000", received.String())

	_, _, topup, err := mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(10000000))
	assert.Empty(t, err)
	assert.True(t, topup)

	err = mgr1.Topup("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", big.NewInt(100000000))
	assert.Empty(t, err)

	mockLotusMgr.checkPaymentChannel = func(chAddr string) (bool, *big.Int, string, error) {
		return false, big.NewInt(200000000), "f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", nil
	}

	voucher, _, _, err = mgr1.Pay("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", 0, big.NewInt(10000000))
	assert.Empty(t, err)
	received, lane, err = mgr2.Receive("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", voucher)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), lane)
	assert.Equal(t, "10000000", received.String())

	// Test refund
	_, err = mgr2.Refund("f2qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", 0, big.NewInt(10000000))
	assert.NotEmpty(t, err)

	_, err = mgr2.Refund("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", 1, big.NewInt(10000000))
	assert.NotEmpty(t, err)

	_, err = mgr2.Refund("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", 0, big.NewInt(1000000000000))
	assert.NotEmpty(t, err)

	voucher, err = mgr2.Refund("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy", 0, big.NewInt(10000000))
	assert.Empty(t, err)

	_, err = mgr1.ReceiveRefund("f2wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", voucher)
	assert.NotEmpty(t, err)

	received, err = mgr1.ReceiveRefund("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi", voucher)
	assert.Empty(t, err)
	assert.Equal(t, "10000000", received.String())

	_, _, _, err = mgr1.GetOutboundChStatus("f2wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi")
	assert.NotEmpty(t, err)

	paychAddr, balance, redeemed, err := mgr1.GetOutboundChStatus("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi")
	assert.Empty(t, err)
	assert.Equal(t, "f12yybez3cfe2yb2nsartagpwkk23q5hmmiluqafi", paychAddr)
	assert.Equal(t, "200000000", balance.String())
	assert.Equal(t, "100000000", redeemed.String())

	_, _, _, err = mgr2.GetInboundChStatus("f2qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy")
	assert.NotEmpty(t, err)

	paychAddr, balance, redeemed, err = mgr2.GetInboundChStatus("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy")
	assert.Empty(t, err)
	assert.Equal(t, "f12yybez3cfe2yb2nsartagpwkk23q5hmmiluqafi", paychAddr)
	assert.Equal(t, "200000000", balance.String())
	assert.Equal(t, "100000000", redeemed.String())

	err = mgr1.RemoveOutboundCh("f2wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi")
	assert.NotEmpty(t, err)

	err = mgr1.RemoveOutboundCh("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi")
	assert.Empty(t, err)

	_, _, _, err = mgr1.GetOutboundChStatus("f1wcl5t2jld4iqtthqmj4ef4xvx7jy64eqvyvkchi")
	assert.NotEmpty(t, err)

	err = mgr2.RemoveInboundCh("f2qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy")
	assert.NotEmpty(t, err)

	err = mgr2.RemoveInboundCh("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy")
	assert.Empty(t, err)

	_, _, _, err = mgr2.GetInboundChStatus("f1qsbhbdqnbzjxqmz3fchnodr5vfae2twuwstoxuy")
	assert.NotEmpty(t, err)
}

func TestUnimplemented(t *testing.T) {
	mockLotusMgr := mockLotusMgr{}
	mgr := NewFCRPaymentMgrImplV1("933dfc0be9ca2d783446fa3fa9ea27bd9cc553ec5131256dd6fddcde3302b9e0", &mockLotusMgr)
	err := mgr.Start()
	assert.Empty(t, err)

	_, err = mgr.GetCostToCreate("", nil)
	assert.NotEmpty(t, err)

	_, err = mgr.CheckRecipientSettlementValidity("")
	assert.NotEmpty(t, err)

	err = mgr.Settle("")
	assert.NotEmpty(t, err)

	_, err = mgr.GetCostToSettle("")
	assert.NotEmpty(t, err)

	_, err = mgr.CheckSettlementValidity("")
	assert.NotEmpty(t, err)
}
