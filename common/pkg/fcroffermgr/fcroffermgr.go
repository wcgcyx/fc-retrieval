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

	/* CID Offer related functions */
	// AddOffer adds an cid offer to the storage.
	AddOffer(offer *cidoffer.CIDOffer)

	// AddOffer adds an offer to the storage with a tag.
	// Should be called by a provider to link offer with filename.
	AddOfferWithTag(offer *cidoffer.CIDOffer, tag string)

	// GetOffers gets offers containing given cid.
	GetOffers(cID *cid.ContentID) []cidoffer.CIDOffer

	// GetOffersByTag gets offers by given tag.
	// Should be called by a provider to show offers linked with filename.
	GetOffersByTag(tag string) []cidoffer.CIDOffer

	// ListOffers gets a list of offers from given index to given index.
	ListOffers(from uint, to uint) []cidoffer.CIDOffer

	// ListOffersWithTag gets a list of offers and their tag from given index to given index.
	ListOffersWithTag(from uint, to uint) ([]cidoffer.CIDOffer, []string)

	// ListOffersWithAccessCount gets a list of offers and their access count from given index to given index.
	// From most frequently accessed offer to least frequently accessed offer.
	// It is used by gateway to list offers.
	ListOffersWithAccessCount(from uint, to uint) ([]cidoffer.CIDOffer, []int)

	// RemoveOffer removes an offer by digest
	RemoveOffer(digest string)

	/* SubCID Offer related functions */
	// AddSubOffer adds an cid offer to the storage.
	AddSubOffer(offer *cidoffer.SubCIDOffer)

	// GetSubOffers gets offers containing given cid.
	GetSubOffers(cID *cid.ContentID) []cidoffer.SubCIDOffer

	// ListSubOffers gets a list of offers from given index to given index.
	ListSubOffers(from uint, to uint) []cidoffer.SubCIDOffer

	// RemoveSubOffer removes an offer by digest
	RemoveSubOffer(digest string)
}
