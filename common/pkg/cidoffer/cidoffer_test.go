/*
Package cidoffer - provides functionality like create, verify, sign and get details for CIDOffer and SubCIDOffer structures.

CIDOffer represents an offer from a Storage Provider, explaining on what conditions the client can retrieve a set of uniquely identified files from Filecoin blockchain network.
SubCIDOffer represents an offer from a Storage Provider, just like CIDOffer, but for a single file and includes a merkle proof
*/
package cidoffer

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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
)

const (
	Cid1Str     = "baga6ea4seaqo3bt5lts3745ejrq3sh5jarauuwjn4vlsnc7f63fuxcq3psio2hq"
	Cid2Str     = "QmekfhK273inQzFqx14G1oK2iY2jkuM1q71bjFA47ZieYt"
	Cid3Str     = "mAXCg5AIgxeF8+e09cmz4EVe5eJ9a5hczOe0h+FyAO5Xd5g1dS+E"
	PrivKey     = "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29"
	PubKey      = "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d"
	PubKeyWrong = "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62a"
)

func TestNewCIDOffer(t *testing.T) {
	aCid, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid)
	cids := []cid.ContentID{*aCid}
	price := big.NewInt(100)
	expiry := int64(10)
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	assert.NotEmpty(t, offer)
	assert.Equal(t, "testprovider", offer.GetProviderID())
	offerCIDs := offer.GetCIDs()
	assert.Equal(t, len(cids), len(offerCIDs))
	for i := 0; i < len(cids); i++ {
		assert.Equal(t, cids[i].ToString(), offerCIDs[i].ToString())
	}
	assert.Equal(t, price.String(), offer.GetPrice().String())
	assert.Equal(t, expiry, offer.GetExpiry())
	assert.Equal(t, qos, offer.GetQoS())
}

func TestNewCIDOfferMultipleCIDs(t *testing.T) {
	aCid1, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid1)
	aCid2, err := cid.NewContentID(Cid2Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid2)
	aCid3, err := cid.NewContentID(Cid3Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid3)
	cids := []cid.ContentID{*aCid1, *aCid2, *aCid3}
	price := big.NewInt(100)
	expiry := int64(10)
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	assert.NotEmpty(t, offer)
	assert.Equal(t, "testprovider", offer.GetProviderID())
	offerCIDs := offer.GetCIDs()
	assert.Equal(t, len(cids), len(offerCIDs))
	for i := 0; i < len(cids); i++ {
		assert.Equal(t, cids[i].ToString(), offerCIDs[i].ToString())
	}
	assert.Equal(t, price.String(), offer.GetPrice().String())
	assert.Equal(t, expiry, offer.GetExpiry())
	assert.Equal(t, qos, offer.GetQoS())
}

func TestNewCIDOfferWithError(t *testing.T) {
	cids := []cid.ContentID{}
	price := big.NewInt(100)
	expiry := int64(10)
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.NotEmpty(t, err)
	assert.Empty(t, offer)
}

func TestHasExpired(t *testing.T) {
	aCid, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid)
	cids := []cid.ContentID{*aCid}
	price := big.NewInt(100)
	expiry := time.Now().Add(12 * time.Hour).Unix()
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	assert.NotEmpty(t, offer)
	assert.False(t, offer.HasExpired())
	expiry = time.Now().Add(-12 * time.Hour).Unix()
	offer, err = NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	assert.NotEmpty(t, offer)
	assert.True(t, offer.HasExpired())
}

func TestSigning(t *testing.T) {
	aCid, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid)
	cids := []cid.ContentID{*aCid}
	price := big.NewInt(100)
	expiry := time.Now().Add(12 * time.Hour).Unix()
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)

	// Generate private key
	assert.Empty(t, err)
	err = offer.Sign(PrivKey)
	assert.Empty(t, err)
	assert.Equal(t, "00344497164a52734bb28ef7cb9ce8279ef5018340df53254de26c207d5408aa9f55507e012d235e8f3fc66dcb7603329fbba4d3a5ebdc0465a1a77c314f5becb201", offer.GetSignature())

	err = offer.Verify(PubKey)
	assert.Empty(t, err)

	err = offer.Verify(PubKeyWrong)
	assert.NotEmpty(t, err)

	offer.SetSignature("00344497164a52734bb28ef7cb9ce8279ef5018340df53254de26c207d5408aa9f55507e012d235e8f3fc66dcb7603329fbba4d3a5ebdc0465a1a77c314f5becb202")
	err = offer.Verify(PubKey)
	assert.NotEmpty(t, err)
}

func TestDigest(t *testing.T) {
	aCid, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	cids := []cid.ContentID{*aCid}
	price := big.NewInt(100)
	expiry := int64(10)
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	err = offer.Sign(PrivKey)
	assert.Empty(t, err)
	assert.Equal(t, [32]byte{0xbb, 0x6d, 0x76, 0x96, 0xcf,
		0xb2, 0xbf, 0xb6, 0xaa, 0x5e, 0x5, 0xe4, 0x81,
		0x7b, 0x11, 0xcb, 0x8f, 0xa7, 0x99, 0x11, 0x94,
		0xb3, 0x78, 0x1f, 0x4a, 0xf4, 0xca, 0x6, 0x71,
		0x58, 0x81, 0x82}, offer.GetMessageDigest())
}

func TestJSON(t *testing.T) {
	aCid, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	cids := []cid.ContentID{*aCid}
	price := big.NewInt(100)
	expiry := int64(1)
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	p, err := offer.MarshalJSON()
	assert.Empty(t, err)
	assert.Equal(t, []byte{0x7b, 0x22, 0x70, 0x72, 0x6f,
		0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64,
		0x22, 0x3a, 0x22, 0x74, 0x65, 0x73, 0x74, 0x70,
		0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x22,
		0x2c, 0x22, 0x63, 0x69, 0x64, 0x73, 0x22, 0x3a,
		0x5b, 0x22, 0x62, 0x61, 0x67, 0x61, 0x36, 0x65,
		0x61, 0x34, 0x73, 0x65, 0x61, 0x71, 0x6f, 0x33,
		0x62, 0x74, 0x35, 0x6c, 0x74, 0x73, 0x33, 0x37,
		0x34, 0x35, 0x65, 0x6a, 0x72, 0x71, 0x33, 0x73,
		0x68, 0x35, 0x6a, 0x61, 0x72, 0x61, 0x75, 0x75,
		0x77, 0x6a, 0x6e, 0x34, 0x76, 0x6c, 0x73, 0x6e,
		0x63, 0x37, 0x66, 0x36, 0x33, 0x66, 0x75, 0x78,
		0x63, 0x71, 0x33, 0x70, 0x73, 0x69, 0x6f, 0x32,
		0x68, 0x71, 0x22, 0x5d, 0x2c, 0x22, 0x70, 0x72,
		0x69, 0x63, 0x65, 0x22, 0x3a, 0x22, 0x31, 0x30,
		0x30, 0x22, 0x2c, 0x22, 0x65, 0x78, 0x70, 0x69,
		0x72, 0x79, 0x22, 0x3a, 0x31, 0x2c, 0x22, 0x71,
		0x6f, 0x73, 0x22, 0x3a, 0x35, 0x2c, 0x22, 0x73,
		0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
		0x22, 0x3a, 0x22, 0x22, 0x7d}, p)
	offer2 := CIDOffer{}
	err = offer2.UnmarshalJSON(p)
	assert.Empty(t, err)
	assert.Equal(t, offer.GetProviderID(), offer2.GetProviderID())
	assert.Equal(t, offer.GetCIDs(), offer2.GetCIDs())
	assert.Equal(t, offer.GetPrice(), offer2.GetPrice())
	assert.Equal(t, offer.GetExpiry(), offer2.GetExpiry())
	assert.Equal(t, offer.GetQoS(), offer2.GetQoS())
	assert.Equal(t, offer.GetSignature(), offer2.GetSignature())
	assert.Equal(t, offer.merkleRoot, offer2.merkleRoot)
	err = offer2.UnmarshalJSON([]byte{})
	assert.NotEmpty(t, err)
}
