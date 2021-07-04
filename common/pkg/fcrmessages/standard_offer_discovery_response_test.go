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

func TestStandardOfferDiscoveryResponse(t *testing.T) {
	mockCID1, err := cid.NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	mockCID2, err := cid.NewContentID("baga6ea4seaqesauho7j2thfi4g4u5zbnhn2okd74s2igpvc2lsb7rrsfstoy4by")
	assert.Empty(t, err)
	mockOffer, err := cidoffer.NewCIDOffer("testprovider", []cid.ContentID{*mockCID1, *mockCID2}, big.NewInt(40), 40, 101)
	assert.Empty(t, err)
	mockSubOffer, err := mockOffer.GenerateSubCIDOffer(mockCID1)
	assert.Empty(t, err)
	mockNonce := int64(42)
	mockVoucher := "mockVoucher"

	msg, err := EncodeStandardOfferDiscoveryResponse([]cidoffer.SubCIDOffer{*mockSubOffer}, mockNonce, mockVoucher)
	assert.Empty(t, err)
	assert.Equal(t, byte(StandardOfferDiscoveryResponseType), msg.messageType)
	assert.Equal(t, "7b226f6666657273223a5b22376232323730373236663736363936343635373235663639363432323361323237343635373337343730373236663736363936343635373232323263323237333735363235663633363936343232336132323531366435383335353236373338373433393761363833323336346136333631353436623337353636653434353837313736333535333438343833323632353433363431363636353666353434363463353337333730333436343462323232633232366436353732366236633635356637323666366637343232336132323338333133343636333536353334333433383636363536323631363133323338333636313636333733313330333533323333363633353339363536333333363133363335333933363335333733313336363233323332333036313331363233373339363336363633333633393330333033323633333333353333363136313333323232633232366436353732366236633635356637303732366636663636323233613232333233323334333133343331333433313334333133343634333433363337333333363339333533323336363233373338333436353334363433363632333533363336333133363335333636333334363533373339333533353336363433333335333436353335333933333332333733343335333433353631333433383335363133353336333633323336363433353631333433353335363133353336333433323334333533363332333433383335363133353333333436343336363133363334333533313335333533353337333333313337333733363335333533373336363333363634333633323335333833353332333433353334363233333330333733303336333133343636333433343333333033363339333533383335333133343331333433313334333133343331333436353336333233343634333533363333333033333634333233323232326332323730373236393633363532323361323233343330323232633232363537383730363937323739323233613334333032633232373136663733323233613331333033313263323237333639363736653631373437353732363532323361323232323764225d2c226e6f6e6365223a34322c22726566756e645f766f7563686572223a226d6f636b566f7563686572227d", hex.EncodeToString(msg.messageBody))
	assert.Equal(t, "", msg.signature)

	resOffers, resNonce, resVoucher, err := DecodeStandardOfferDiscoveryResponse(msg)
	assert.Empty(t, err)
	assert.Equal(t, 1, len(resOffers))
	assert.Equal(t, mockSubOffer.GetMessageDigest(), resOffers[0].GetMessageDigest())
	assert.Equal(t, mockNonce, resNonce)
	assert.Equal(t, mockVoucher, resVoucher)

	msg.messageType = 100
	_, _, _, err = DecodeStandardOfferDiscoveryResponse(msg)
	assert.NotEmpty(t, err)
	msg.messageType = 1

	msg.messageBody = []byte{100, 100, 100}
	_, _, _, err = DecodeStandardOfferDiscoveryResponse(msg)
	assert.NotEmpty(t, err)
}
