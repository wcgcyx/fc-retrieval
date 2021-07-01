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

	mgr.AddOffer(offer0)
	mgr.AddOffer(offer0)
	mgr.AddOffer(offer1)
	mgr.AddOffer(offer2)
	mgr.AddOffer(offer3)
	mgr.AddOffer(offer4)
	mgr.AddOfferWithTag(offer5, "testtag")

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

	res = mgr.GetOffersByTag("testtag2")
	assert.Equal(t, 0, len(res))
	res = mgr.GetOffersByTag("testtag")
	assert.Equal(t, 1, len(res))

	res = mgr.ListOffers(3, 1)
	assert.Equal(t, 0, len(res))

	res = mgr.ListOffers(1, 3)
	assert.Equal(t, 2, len(res))

	res = mgr.ListOffers(2, 3)
	assert.Equal(t, 1, len(res))

	res = mgr.ListOffers(2, 5)
	assert.Equal(t, 3, len(res))

	res = mgr.ListOffers(2, 10)
	assert.Equal(t, 4, len(res))

	res, tags := mgr.ListOffersWithTag(3, 2)
	assert.Equal(t, 0, len(res))
	assert.Equal(t, 0, len(tags))

	res, tags = mgr.ListOffersWithTag(2, 3)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, 1, len(tags))
}
