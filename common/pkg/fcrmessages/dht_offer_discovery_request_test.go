/*
Package fcrmessages - stores all the p2p messages.
*/
package fcrmessages

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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
)

func TestDHTOfferDiscoveryRequest(t *testing.T) {
	mockNonce := uint64(100)
	mockNodeID := "mockID"
	mockCID, err := cid.NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	mockNumDHT := uint32(10)
	mockMaxOfferRequestedPerDHT := uint32(10)
	mockAccountAddr := "mockAddr"
	mockVoucher := "mockVoucher"

	msg, err := EncodeDHTOfferDiscoveryRequest(mockNonce, mockNodeID, mockCID, mockNumDHT, mockMaxOfferRequestedPerDHT, mockAccountAddr, mockVoucher)
	assert.Empty(t, err)
	assert.Equal(t, byte(DHTOfferDiscoveryRequestType), msg.messageType)
	assert.Equal(t, uint64(100), msg.nonce)
	assert.Equal(t, "7b226e6f64655f6964223a226d6f636b4944222c2270696563655f636964223a22516d583552673874397a6832364a6361546b37566e4458717635534848326254364166656f54464c53737034644b222c226e756d5f646874223a31302c226d61785f6f666665725f7265717565737465645f7065725f646874223a31302c226163636f756e745f61646472223a226d6f636b41646472222c22766f7563686572223a226d6f636b566f7563686572227d", hex.EncodeToString(msg.messageBody))
	assert.Equal(t, "", msg.signature)

	resNonce, resNodeID, resCID, resNumDHT, resMaxOfferRequestedPerDHT, resAcountAddr, resVoucher, err := DecodeDHTOfferDiscoveryRequest(msg)
	assert.Empty(t, err)
	assert.Equal(t, mockNonce, resNonce)
	assert.Equal(t, mockNodeID, resNodeID)
	assert.Equal(t, mockCID.ToString(), resCID.ToString())
	assert.Equal(t, mockNumDHT, resNumDHT)
	assert.Equal(t, mockMaxOfferRequestedPerDHT, resMaxOfferRequestedPerDHT)
	assert.Equal(t, mockAccountAddr, resAcountAddr)
	assert.Equal(t, mockVoucher, resVoucher)

	msg.messageType = 100
	_, _, _, _, _, _, _, err = DecodeDHTOfferDiscoveryRequest(msg)
	assert.NotEmpty(t, err)
	msg.messageType = 2

	msg.messageBody = []byte{100, 100, 100}
	_, _, _, _, _, _, _, err = DecodeDHTOfferDiscoveryRequest(msg)
	assert.NotEmpty(t, err)
}
