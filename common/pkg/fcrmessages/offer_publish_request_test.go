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
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
)

func TestOfferPublishRequest(t *testing.T) {
	mockNonce := uint64(100)
	mockNodeID := "mockID"
	mockCID, err := cid.NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	mockOffer, err := cidoffer.NewCIDOffer("testprovider", []cid.ContentID{*mockCID}, big.NewInt(100), 40, 40)
	assert.Empty(t, err)
	msg, err := EncodeOfferPublishRequest(mockNonce, mockNodeID, mockOffer)
	assert.Empty(t, err)
	assert.Equal(t, byte(OfferPublishRequestType), msg.messageType)
	assert.Equal(t, uint64(100), msg.nonce)
	assert.Equal(t, "7b226e6f64655f6964223a226d6f636b4944222c226f66666572223a22376232323730373236663736363936343635373235663639363432323361323237343635373337343730373236663736363936343635373232323263323236333639363437333232336135623232353136643538333535323637333837343339376136383332333634613633363135343662333735363665343435383731373633353533343834383332363235343336343136363635366635343436346335333733373033343634346232323564326332323730373236393633363532323361323233313330333032323263323236353738373036393732373932323361333433303263323237313666373332323361333433303263323237333639363736653631373437353732363532323361323232323764227d", hex.EncodeToString(msg.messageBody))
	assert.Equal(t, "", msg.signature)

	resNonce, resNodeID, resOffer, err := DecodeOfferPublishRequest(msg)
	assert.Empty(t, err)
	assert.Equal(t, mockNonce, resNonce)
	assert.Equal(t, mockNodeID, resNodeID)
	assert.Equal(t, mockOffer.GetMessageDigest(), resOffer.GetMessageDigest())

	msg.messageType = 100
	_, _, _, err = DecodeOfferPublishRequest(msg)
	assert.NotEmpty(t, err)
	msg.messageType = 4

	msg.messageBody = []byte{100, 100, 100}
	_, _, _, err = DecodeOfferPublishRequest(msg)
	assert.NotEmpty(t, err)
}
