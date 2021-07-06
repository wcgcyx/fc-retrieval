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
)

func TestDHTOfferDiscoveryResponse(t *testing.T) {
	mockNonce := uint64(100)
	mockContacted := make(map[string]*FCRACKMsg)
	mockContacted["01"] = CreateFCRACKMsg(1, []byte{1, 2, 3})
	mockContacted["02"] = CreateFCRACKMsg(1, []byte{2, 3, 4})
	mockVoucher := "mockVoucher"

	msg, err := EncodeDHTOfferDiscoveryResponse(mockNonce, mockContacted, mockVoucher)
	assert.Empty(t, err)
	assert.Equal(t, true, msg.ack)
	assert.Equal(t, uint64(100), mockNonce)
	assert.Equal(t, "7b22636f6e746163746564223a5b223031222c223032225d2c22726573706f6e736573223a5b223762323236313633366232323361373437323735363532633232373536393665373433363334323233613331326332323664363537333733363136373635356636323666363437393232336132323330333133303332333033333232326332323664363537333733363136373635356637333639363736653631373437353732363532323361323232323764222c223762323236313633366232323361373437323735363532633232373536393665373433363334323233613331326332323664363537333733363136373635356636323666363437393232336132323330333233303333333033343232326332323664363537333733363136373635356637333639363736653631373437353732363532323361323232323764225d2c22726566756e645f766f7563686572223a226d6f636b566f7563686572227d", hex.EncodeToString(msg.messageBody))
	assert.Equal(t, "", msg.signature)

	resNonce, resContacted, resVoucher, err := DecodeDHTOfferDiscoveryResponse(msg)
	assert.Empty(t, err)
	assert.Equal(t, 2, len(resContacted))
	for key, val := range resContacted {
		assert.Equal(t, mockContacted[key].ACK(), val.ACK())
		assert.Equal(t, mockContacted[key].Body(), val.Body())
		assert.Equal(t, mockContacted[key].Signature(), val.Signature())
	}
	assert.Equal(t, mockNonce, resNonce)
	assert.Equal(t, mockVoucher, resVoucher)

	msg.ack = false
	_, _, _, err = DecodeDHTOfferDiscoveryResponse(msg)
	assert.NotEmpty(t, err)
	msg.ack = true

	msg.messageBody = []byte{100, 100, 100}
	_, _, _, err = DecodeDHTOfferDiscoveryResponse(msg)
	assert.NotEmpty(t, err)
}
