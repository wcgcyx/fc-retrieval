/*
Package fcroffermgr - offer manager manages all offers stored.
*/
package fcroffermgr

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
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
)

const (
	CID1  = "QmWJi2BHLpKpCnD3sA3jcSWv5M51D6Zf1WY4rN8BrQtCgi"
	CID2  = "baga6ea4seaqesauho7j2thfi4g4u5zbnhn2okd74s2igpvc2lsb7rrsfstoy4by"
	CID3  = "baga6ea4seaqpdbfmh26egh4oqpb4dtr5mjtrktedcp527pfcui5pz63faik3wgq"
	CID4  = "baga6ea4seaqlyo5kq57dbjmw4rcexvtqbr5sdhqdfyezvy6hy3bdezfvu6jeiby"
	CID5  = "QmcVy3EpcDPeVkJExZQxx5ZStaey19min1LLkgwt9cJYYM"
	CID6  = "QmPQ42Rdn4rvsEFab6aoGHg4Cj1kiNh6xpNCpRVq98Qevz"
	CID7  = "QmQ2GoeeLevT6TWk9xt6JLorvHc3AhAvJxqT2bfSUgf63C"
	CID8  = "QmZbmXDNWtDafGzCPgAJwFnMXrtXH6DunKGnd6FNuJVzwS"
	CID9  = "QmdXdUx8VWFQ1uLKukUbyXU2aZy3MQQFESe7427PtXAahQ"
	CID10 = "QmVPhUbiWEoFJ26p4uZveuMhnZvVuFx9Drras6FyD8aw22"
)

func TestAddOffer(t *testing.T) {
	mgr := NewFCROfferMgrImplV1(true)
	err := mgr.Start()
	assert.Empty(t, err)
	defer mgr.Shutdown()

	cid1, err := cid.NewContentID(CID1)
	assert.Empty(t, err)
	cid2, err := cid.NewContentID(CID2)
	assert.Empty(t, err)
	cid3, err := cid.NewContentID(CID3)
	assert.Empty(t, err)
	cid4, err := cid.NewContentID(CID4)
	assert.Empty(t, err)
	cid5, err := cid.NewContentID(CID5)
	assert.Empty(t, err)
	cid6, err := cid.NewContentID(CID6)
	assert.Empty(t, err)
	cid7, err := cid.NewContentID(CID7)
	assert.Empty(t, err)
	cid8, err := cid.NewContentID(CID8)
	assert.Empty(t, err)
	cid9, err := cid.NewContentID(CID9)
	assert.Empty(t, err)
	cid10, err := cid.NewContentID(CID10)
	assert.Empty(t, err)

	offer0, err := cidoffer.NewCIDOffer("testID", []cid.ContentID{*cid1, *cid2, *cid3, *cid4}, big.NewInt(10), 10, 10)
	assert.Empty(t, err)
	offer1, err := cidoffer.NewCIDOffer("testID", []cid.ContentID{*cid2, *cid5, *cid6}, big.NewInt(10), 10, 10)
	assert.Empty(t, err)
	offer2, err := cidoffer.NewCIDOffer("testID", []cid.ContentID{*cid7, *cid8}, big.NewInt(10), 10, 10)
	assert.Empty(t, err)
	offer3, err := cidoffer.NewCIDOffer("testID", []cid.ContentID{*cid7}, big.NewInt(10), 10, 10)
	assert.Empty(t, err)
	offer4, err := cidoffer.NewCIDOffer("testID", []cid.ContentID{*cid8}, big.NewInt(10), 10, 10)
	assert.Empty(t, err)
	offer5, err := cidoffer.NewCIDOffer("testID", []cid.ContentID{*cid8, *cid9}, big.NewInt(10), 10, 10)
	assert.Empty(t, err)

	mgr.AddCIDTag(cid1, "CID1")
	mgr.AddCIDTag(cid2, "CID2")
	mgr.AddCIDTag(cid3, "CID3")
	mgr.AddCIDTag(cid4, "CID4")
	mgr.AddCIDTag(cid5, "CID5")
	mgr.AddCIDTag(cid6, "CID6")
	mgr.AddCIDTag(cid7, "CID7")
	mgr.AddCIDTag(cid8, "CID8")

	testCIDStr := mgr.GetCIDByTag("CID3")
	assert.Equal(t, cid3.ToString(), testCIDStr)

	mgr.AddOffer(offer0)
	mgr.AddOffer(offer0)
	mgr.AddOffer(offer1)
	mgr.AddOffer(offer2)
	mgr.AddOffer(offer3)
	mgr.AddOffer(offer4)
	mgr.AddOffer(offer5)

	res1, res2 := mgr.ListAccessCount(0, 10)
	assert.Equal(t, []string{"QmPQ42Rdn4rvsEFab6aoGHg4Cj1kiNh6xpNCpRVq98Qevz", "QmQ2GoeeLevT6TWk9xt6JLorvHc3AhAvJxqT2bfSUgf63C", "QmWJi2BHLpKpCnD3sA3jcSWv5M51D6Zf1WY4rN8BrQtCgi", "QmZbmXDNWtDafGzCPgAJwFnMXrtXH6DunKGnd6FNuJVzwS", "QmcVy3EpcDPeVkJExZQxx5ZStaey19min1LLkgwt9cJYYM", "baga6ea4seaqesauho7j2thfi4g4u5zbnhn2okd74s2igpvc2lsb7rrsfstoy4by", "baga6ea4seaqlyo5kq57dbjmw4rcexvtqbr5sdhqdfyezvy6hy3bdezfvu6jeiby", "baga6ea4seaqpdbfmh26egh4oqpb4dtr5mjtrktedcp527pfcui5pz63faik3wgq"}, res1)
	assert.Equal(t, []int{0, 0, 0, 0, 0, 0, 0, 0}, res2)

	res := mgr.GetOffers(cid1)
	assert.Equal(t, 1, len(res))
	res = mgr.GetOffers(cid2)
	assert.Equal(t, 2, len(res))
	res = mgr.GetOffers(cid5)
	assert.Equal(t, 1, len(res))
	res = mgr.GetOffers(cid8)
	assert.Equal(t, 3, len(res))
	res = mgr.GetOffers(cid10)
	assert.Equal(t, 0, len(res))

	count := mgr.GetAccessCountByCID(cid1)
	assert.Equal(t, count, 0)
	mgr.IncrementCIDAccessCount(cid1)
	count = mgr.GetAccessCountByCID(cid1)
	assert.Equal(t, count, 1)
	mgr.IncrementCIDAccessCount(cid1)
	count = mgr.GetAccessCountByCID(cid1)
	assert.Equal(t, count, 2)

	res1, res2 = mgr.ListAccessCount(0, 2)
	assert.Equal(t, []string{"QmWJi2BHLpKpCnD3sA3jcSWv5M51D6Zf1WY4rN8BrQtCgi", "QmPQ42Rdn4rvsEFab6aoGHg4Cj1kiNh6xpNCpRVq98Qevz"}, res1)
	assert.Equal(t, []int{2, 0}, res2)

	res = mgr.GetOffersByTag("CID10")
	assert.Equal(t, 0, len(res))
	res = mgr.GetOffersByTag("CID8")
	assert.Equal(t, 3, len(res))

	tag := mgr.GetTagByCID(cid10)
	assert.Equal(t, "", tag)
	tag = mgr.GetTagByCID(cid3)
	assert.Equal(t, "CID3", tag)

	res = mgr.ListOffers(3, 1)
	assert.Equal(t, 0, len(res))

	res = mgr.ListOffers(1, 3)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "1d8b5afd46676b00a4433b313a83e40719fdf7c3b52131b8b09b304c26ec1e82", res[0].GetMessageDigest())
	assert.Equal(t, "697dfe073c9714504b6364e7333feceba4b3bbe64f2104efa5842c1a2331a311", res[1].GetMessageDigest())

	res = mgr.ListOffers(2, 3)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "697dfe073c9714504b6364e7333feceba4b3bbe64f2104efa5842c1a2331a311", res[0].GetMessageDigest())

	res = mgr.ListOffers(2, 5)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "697dfe073c9714504b6364e7333feceba4b3bbe64f2104efa5842c1a2331a311", res[0].GetMessageDigest())
	assert.Equal(t, "9198ee39730bad84b65185b1c306f7b575a0a669ab958f7aba7a35c71f779652", res[1].GetMessageDigest())
	assert.Equal(t, "d746ad9bf2a5deafe1f8848eed376e0c68ccd4c600d8b2c9c5d7b832a729ea21", res[2].GetMessageDigest())

	res = mgr.ListOffers(2, 10)
	assert.Equal(t, 4, len(res))
	assert.Equal(t, "697dfe073c9714504b6364e7333feceba4b3bbe64f2104efa5842c1a2331a311", res[0].GetMessageDigest())
	assert.Equal(t, "9198ee39730bad84b65185b1c306f7b575a0a669ab958f7aba7a35c71f779652", res[1].GetMessageDigest())
	assert.Equal(t, "d746ad9bf2a5deafe1f8848eed376e0c68ccd4c600d8b2c9c5d7b832a729ea21", res[2].GetMessageDigest())
	assert.Equal(t, "fb46952a0a8c2c58d76d3b131099d2cbbfdb0029905efb7a7aad709dd827a9f5", res[3].GetMessageDigest())

	offer := mgr.GetOfferByDigest("a9aac8229414ad4f42e73cf93e79f922ff65d5a6465c83be6070baaeeca988ff")
	assert.Empty(t, offer)
	offer = mgr.GetOfferByDigest("09aac8229414ad4f42e73cf93e79f922ff65d5a6465c83be6070baaeeca988ff")
	assert.NotEmpty(t, offer)

	mgr.RemoveOffer("09aac8229414ad4f42e73cf93e79f922ff65d5a6465c83be6070baaeeca988ff")
	res = mgr.GetOffers(cid5)
	assert.Equal(t, 0, len(res))
}

func TestSubOffer(t *testing.T) {
	mgr := NewFCROfferMgrImplV1(false)
	err := mgr.Start()
	assert.Empty(t, err)
	defer mgr.Shutdown()

	cid1, err := cid.NewContentID(CID1)
	assert.Empty(t, err)
	cid2, err := cid.NewContentID(CID2)
	assert.Empty(t, err)
	cid3, err := cid.NewContentID(CID3)
	assert.Empty(t, err)
	cid4, err := cid.NewContentID(CID4)
	assert.Empty(t, err)
	cid5, err := cid.NewContentID(CID5)
	assert.Empty(t, err)
	cid6, err := cid.NewContentID(CID6)
	assert.Empty(t, err)
	cid7, err := cid.NewContentID(CID7)
	assert.Empty(t, err)

	offer, err := cidoffer.NewCIDOffer("testID", []cid.ContentID{*cid1, *cid2, *cid3, *cid4, *cid5, *cid6}, big.NewInt(10), 10, 10)
	assert.Empty(t, err)

	offer2, err := cidoffer.NewCIDOffer("testID", []cid.ContentID{*cid1, *cid2}, big.NewInt(10), 10, 10)
	assert.Empty(t, err)

	subOffer0, err := offer.GenerateSubCIDOffer(cid1)
	assert.Empty(t, err)

	subOffer1, err := offer.GenerateSubCIDOffer(cid2)
	assert.Empty(t, err)

	subOffer2, err := offer.GenerateSubCIDOffer(cid3)
	assert.Empty(t, err)

	subOffer3, err := offer.GenerateSubCIDOffer(cid4)
	assert.Empty(t, err)

	subOffer4, err := offer.GenerateSubCIDOffer(cid5)
	assert.Empty(t, err)

	subOffer5, err := offer.GenerateSubCIDOffer(cid6)
	assert.Empty(t, err)

	subOffer6, err := offer2.GenerateSubCIDOffer(cid1)
	assert.Empty(t, err)

	subOffer7, err := offer2.GenerateSubCIDOffer(cid2)
	assert.Empty(t, err)

	mgr.AddSubOffer(subOffer0)
	mgr.AddSubOffer(subOffer0)
	mgr.AddSubOffer(subOffer1)
	mgr.AddSubOffer(subOffer2)
	mgr.AddSubOffer(subOffer3)
	mgr.AddSubOffer(subOffer4)
	mgr.AddSubOffer(subOffer5)
	mgr.AddSubOffer(subOffer6)
	mgr.AddSubOffer(subOffer7)

	res := mgr.GetSubOffers(cid1)
	assert.Equal(t, 2, len(res))
	res = mgr.GetSubOffers(cid2)
	assert.Equal(t, 2, len(res))
	res = mgr.GetSubOffers(cid3)
	assert.Equal(t, 1, len(res))
	res = mgr.GetSubOffers(cid4)
	assert.Equal(t, 1, len(res))
	res = mgr.GetSubOffers(cid7)
	assert.Equal(t, 0, len(res))

	res = mgr.ListSubOffers(3, 1)
	assert.Equal(t, 0, len(res))

	res = mgr.ListSubOffers(1, 3)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "254c88ee494aa0b61fa394d23ed493ffae277f95f36166c4d8e2de03a9708faf", res[0].GetMessageDigest())
	assert.Equal(t, "6cafb9da087efacf5325fe8def74e0af9eeb3b070137a4a4cb56fc86969f54e0", res[1].GetMessageDigest())

	res = mgr.ListSubOffers(2, 3)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "6cafb9da087efacf5325fe8def74e0af9eeb3b070137a4a4cb56fc86969f54e0", res[0].GetMessageDigest())

	res = mgr.ListSubOffers(2, 5)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "6cafb9da087efacf5325fe8def74e0af9eeb3b070137a4a4cb56fc86969f54e0", res[0].GetMessageDigest())
	assert.Equal(t, "975e4e63d8e71bceeff50bfcdf408861bb656d97d86339c33d85e715da94963a", res[1].GetMessageDigest())
	assert.Equal(t, "9761a57643bc9a1c7fa57b7860a7705b8bac7dd8c6d7933718f195005ac6950f", res[2].GetMessageDigest())

	res = mgr.ListSubOffers(2, 10)
	assert.Equal(t, 6, len(res))
	assert.Equal(t, "6cafb9da087efacf5325fe8def74e0af9eeb3b070137a4a4cb56fc86969f54e0", res[0].GetMessageDigest())
	assert.Equal(t, "975e4e63d8e71bceeff50bfcdf408861bb656d97d86339c33d85e715da94963a", res[1].GetMessageDigest())
	assert.Equal(t, "9761a57643bc9a1c7fa57b7860a7705b8bac7dd8c6d7933718f195005ac6950f", res[2].GetMessageDigest())
	assert.Equal(t, "ae291b1b04d853387cedbb5d4a05578b516e3838b5ce5e206a340dacaef7387e", res[3].GetMessageDigest())
	assert.Equal(t, "d1b1a90b7430bef3113868b13a2ea71fecd06b3ba98950afa54ba5ef31af3382", res[4].GetMessageDigest())
	assert.Equal(t, "e567ac85122552402b0e142925998fdbebf086441d24d31a91e5022105a84ee9", res[5].GetMessageDigest())

	subOffer := mgr.GetSubOfferByDigest("8cafb9da087efacf5325fe8def74e0af9eeb3b070137a4a4cb56fc86969f54e0")
	assert.Empty(t, subOffer)
	subOffer = mgr.GetSubOfferByDigest("6cafb9da087efacf5325fe8def74e0af9eeb3b070137a4a4cb56fc86969f54e0")
	assert.NotEmpty(t, subOffer)

	mgr.RemoveSubOffer("8cafb9da087efacf5325fe8def74e0af9eeb3b070137a4a4cb56fc86969f54e0")
	mgr.RemoveSubOffer("6cafb9da087efacf5325fe8def74e0af9eeb3b070137a4a4cb56fc86969f54e0")
	res = mgr.ListSubOffers(2, 4)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "975e4e63d8e71bceeff50bfcdf408861bb656d97d86339c33d85e715da94963a", res[0].GetMessageDigest())
	assert.Equal(t, "9761a57643bc9a1c7fa57b7860a7705b8bac7dd8c6d7933718f195005ac6950f", res[1].GetMessageDigest())

	mgr.RemoveSubOffer("e567ac85122552402b0e142925998fdbebf086441d24d31a91e5022105a84ee9")
	res = mgr.ListSubOffers(5, 6)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "d1b1a90b7430bef3113868b13a2ea71fecd06b3ba98950afa54ba5ef31af3382", res[0].GetMessageDigest())
}
