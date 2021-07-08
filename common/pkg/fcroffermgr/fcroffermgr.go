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
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
)

// FCROfferMgr represents the manager that manages all stored offers.
type FCROfferMgr interface {
	// Start starts the manager's routine.
	Start() error

	// Shutdown ends the manager's routine safely.
	Shutdown()

	/* CID related functions */
	// AddCIDTag adds a cid to system and its tag.
	// If cid already existed, it fails silently
	AddCIDTag(cid *cid.ContentID, tag string)

	// GetTagByCID gets the tag associated by the given cid.
	GetTagByCID(cid *cid.ContentID) string

	// GetCIDByTag gets the cid string associated by the given tag.
	GetCIDByTag(tag string) string

	// IncrementCIDAccessCount increments the access count for a given cid.
	IncrementCIDAccessCount(cid *cid.ContentID)

	// GetAccessCountByCID gets the access count of a given cid.
	GetAccessCountByCID(cid *cid.ContentID) int

	// ListAccessCount lists the cid and access count from most accessed to least accessed
	ListAccessCount(from uint, to uint) ([]string, []int)

	/* CID Offer related functions */
	// AddOffer adds an cid offer to the storage.
	// If calling from provider, needs to first call add cid tag to track tag.
	AddOffer(offer *cidoffer.CIDOffer)

	// GetOffers gets offers containing given cid.
	GetOffers(cID *cid.ContentID) []cidoffer.CIDOffer

	// GetOffersByTag gets offers by given tag.
	// Should be called by a provider to show offers linked with filename.
	GetOffersByTag(tag string) []cidoffer.CIDOffer

	// ListOffers gets a list of offers from given index to given index.
	ListOffers(from uint, to uint) []cidoffer.CIDOffer

	// GetOfferByDigest
	GetOfferByDigest(digest string) *cidoffer.CIDOffer

	// RemoveOffer removes an offer by digest
	RemoveOffer(digest string)

	/* SubCID Offer related functions */
	// AddSubOffer adds an cid offer to the storage.
	AddSubOffer(offer *cidoffer.SubCIDOffer)

	// GetSubOffers gets offers containing given cid.
	GetSubOffers(cID *cid.ContentID) []cidoffer.SubCIDOffer

	// ListSubOffers gets a list of offers from given index to given index.
	ListSubOffers(from uint, to uint) []cidoffer.SubCIDOffer

	// GetSubOfferByDigest
	GetSubOfferByDigest(digest string) *cidoffer.SubCIDOffer

	// RemoveSubOffer removes an offer by digest
	RemoveSubOffer(digest string)
}
