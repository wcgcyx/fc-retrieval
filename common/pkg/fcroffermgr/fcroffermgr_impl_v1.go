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
	"sort"
	"sync"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
)

// FCROfferMgrImplV1 implements FCROfferMgr interface, it is an in-memory version.
type FCROfferMgrImplV1 struct {

	// lock is the lock to offer storage
	lock sync.RWMutex

	// cidTagMap is a map from cid string -> tag string
	cidTagMap map[string]string

	// cidCountMap is a map from cid -> count
	cidCountMap map[string]int

	// cidDigestMap is a map from cid string -> (map digest string -> true)
	cidDigestMap map[string]map[string]bool

	// tagDigestMap is a map from tag string -> (map digest string -> true)
	tagDigestMap map[string]map[string]bool

	// digestOfferMap is a map from digest string -> offer
	digestOfferMap map[string]*cidoffer.CIDOffer

	// digestOfferMapS is the digest sub offer map, map from digest string -> sub offer
	digestOfferMapS map[string]*cidoffer.SubCIDOffer

	// cidDigestMapS is a map from cid string -> (map from digest string -> true)
	cidDigestMapS map[string]map[string]bool
}

func NewFCROfferMgrImplV1(tracking bool) FCROfferMgr {
	return &FCROfferMgrImplV1{
		lock:            sync.RWMutex{},
		cidTagMap:       make(map[string]string),
		cidCountMap:     make(map[string]int),
		cidDigestMap:    make(map[string]map[string]bool),
		tagDigestMap:    make(map[string]map[string]bool),
		digestOfferMap:  make(map[string]*cidoffer.CIDOffer),
		digestOfferMapS: make(map[string]*cidoffer.SubCIDOffer),
		cidDigestMapS:   make(map[string]map[string]bool),
	}
}

func (mgr *FCROfferMgrImplV1) Start() error {
	return nil
}

func (mgr *FCROfferMgrImplV1) Shutdown() {
}

func (mgr *FCROfferMgrImplV1) AddCIDTag(cid *cid.ContentID, tag string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.cidTagMap[cid.ToString()] = tag
}

func (mgr *FCROfferMgrImplV1) GetTagByCID(cid *cid.ContentID) string {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	return mgr.cidTagMap[cid.ToString()]
}

func (mgr *FCROfferMgrImplV1) IncrementCIDAccessCount(cid *cid.ContentID) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	_, ok := mgr.cidCountMap[cid.ToString()]
	if !ok {
		mgr.cidCountMap[cid.ToString()] = 1
	} else {
		mgr.cidCountMap[cid.ToString()]++
	}
}

func (mgr *FCROfferMgrImplV1) GetAccessCountByCID(cid *cid.ContentID) int {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	return mgr.cidCountMap[cid.ToString()]
}

func (mgr *FCROfferMgrImplV1) AddOffer(offer *cidoffer.CIDOffer) {
	// Need to update cid -> digest map, tag -> digest map and digest -> offer map
	digest := offer.GetMessageDigest()
	mgr.lock.RLock()
	_, ok := mgr.digestOfferMap[digest]
	mgr.lock.RUnlock()
	if ok {
		// Offer existed
		return
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	// Add a copy
	copy := offer.Copy()
	if copy == nil {
		logging.Error("Fail to obtain a copy of the offer when adding to storage.")
		return
	}
	// Update digest -> offer map
	mgr.digestOfferMap[digest] = copy

	for _, cid := range copy.GetCIDs() {
		cidStr := cid.ToString()
		// Update cid -> digest map
		_, ok = mgr.cidDigestMap[cidStr]
		if !ok {
			mgr.cidDigestMap[cidStr] = make(map[string]bool)
		}
		mgr.cidDigestMap[cidStr][digest] = true
		// Update tag -> digest map
		tag := mgr.cidTagMap[cidStr]
		_, ok = mgr.tagDigestMap[tag]
		if !ok {
			mgr.tagDigestMap[tag] = make(map[string]bool)
		}
		mgr.tagDigestMap[tag][digest] = true
	}
}

func (mgr *FCROfferMgrImplV1) GetOffers(cid *cid.ContentID) []cidoffer.CIDOffer {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	cidStr := cid.ToString()
	res := make([]cidoffer.CIDOffer, 0)
	digests, ok := mgr.cidDigestMap[cidStr]
	if !ok {
		return res
	}
	for digest := range digests {
		copy := mgr.digestOfferMap[digest].Copy()
		if copy == nil {
			logging.Error("Fail to obtain a copy of the offer when getting offers by cid.")
			continue
		}
		res = append(res, *copy)
	}
	return res
}

func (mgr *FCROfferMgrImplV1) GetOffersByTag(tag string) []cidoffer.CIDOffer {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	res := make([]cidoffer.CIDOffer, 0)
	digests, ok := mgr.tagDigestMap[tag]
	if !ok {
		return res
	}
	for digest := range digests {
		copy := mgr.digestOfferMap[digest].Copy()
		if copy == nil {
			logging.Error("Fail to obtain a copy of the offer when getting offers by tag.")
			continue
		}
		res = append(res, *copy)
	}
	return res
}

func (mgr *FCROfferMgrImplV1) ListOffers(from uint, to uint) []cidoffer.CIDOffer {
	res := make([]cidoffer.CIDOffer, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	if from >= to || from >= uint(len(mgr.digestOfferMap)) {
		return res
	}
	keys := make([]string, len(mgr.digestOfferMap))
	i := 0
	for key := range mgr.digestOfferMap {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	index := uint(0)
	for _, key := range keys {
		if index >= from {
			if index < to {
				copy := mgr.digestOfferMap[key].Copy()
				if copy == nil {
					logging.Error("Fail to obtain a copy of the offer when listing offers.")
					continue
				}
				res = append(res, *copy)
			} else {
				return res
			}
		}
		index++
	}
	return res
}

func (mgr *FCROfferMgrImplV1) GetOfferByDigest(digest string) *cidoffer.CIDOffer {
	// Return a copy
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	res := mgr.digestOfferMap[digest]
	if res != nil {
		res = res.Copy()
	}
	return res
}

func (mgr *FCROfferMgrImplV1) RemoveOffer(digest string) {
	// Need to update cid -> digest map, tag -> digest map and digest -> offer map
	mgr.lock.RLock()
	_, ok := mgr.digestOfferMap[digest]
	mgr.lock.RUnlock()
	if !ok {
		// Offer not existed
		return
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	cids := mgr.digestOfferMap[digest].GetCIDs()
	delete(mgr.digestOfferMap, digest)

	for _, cid := range cids {
		cidStr := cid.ToString()
		// Update cid map
		delete(mgr.cidDigestMap[cidStr], digest)
		if len(mgr.cidDigestMap[cidStr]) == 0 {
			delete(mgr.cidDigestMap, cidStr)
		}
		// Update tag map
		tag := mgr.cidTagMap[cidStr]
		delete(mgr.tagDigestMap, digest)
		if len(mgr.tagDigestMap[tag]) == 0 {
			delete(mgr.tagDigestMap, tag)
		}
	}
}

func (mgr *FCROfferMgrImplV1) AddSubOffer(offer *cidoffer.SubCIDOffer) {
	digest := offer.GetMessageDigest()
	mgr.lock.RLock()
	_, ok := mgr.digestOfferMapS[digest]
	mgr.lock.RUnlock()
	if ok {
		// Offer existed
		return
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	// Add a copy
	copy := offer.Copy()
	if copy == nil {
		logging.Error("Fail to obtain a copy of the sub offer when adding to storage.")
		return
	}
	mgr.digestOfferMapS[digest] = copy

	subCIDStr := offer.GetSubCID().ToString()
	_, ok = mgr.cidDigestMapS[subCIDStr]
	if !ok {
		mgr.cidDigestMapS[subCIDStr] = map[string]bool{}
	}
	mgr.cidDigestMapS[subCIDStr][digest] = true
}

func (mgr *FCROfferMgrImplV1) GetSubOffers(cID *cid.ContentID) []cidoffer.SubCIDOffer {
	res := make([]cidoffer.SubCIDOffer, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	digests, ok := mgr.cidDigestMapS[cID.ToString()]
	if !ok {
		return res
	}
	for digest := range digests {
		copy := mgr.digestOfferMapS[digest].Copy()
		if copy == nil {
			logging.Error("Fail to obtain a copy of the offer when getting sub offers by cid.")
			continue
		}
		res = append(res, *copy)
	}
	return res
}

func (mgr *FCROfferMgrImplV1) ListSubOffers(from uint, to uint) []cidoffer.SubCIDOffer {
	res := make([]cidoffer.SubCIDOffer, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	if from >= to || from >= uint(len(mgr.digestOfferMapS)) {
		return res
	}
	keys := make([]string, len(mgr.digestOfferMapS))
	i := 0
	for key := range mgr.digestOfferMapS {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	index := uint(0)
	for _, key := range keys {
		if index >= from {
			if index < to {
				copy := mgr.digestOfferMapS[key].Copy()
				if copy == nil {
					logging.Error("Fail to obtain a copy of the offer when listing sub offers.")
					continue
				}
				res = append(res, *copy)
			} else {
				return res
			}
		}
		index++
	}
	return res
}

func (mgr *FCROfferMgrImplV1) GetSubOfferByDigest(digest string) *cidoffer.SubCIDOffer {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	res := mgr.digestOfferMapS[digest]
	if res != nil {
		res = res.Copy()
	}
	return res
}

func (mgr *FCROfferMgrImplV1) RemoveSubOffer(digest string) {
	mgr.lock.RLock()
	_, ok := mgr.digestOfferMapS[digest]
	mgr.lock.RUnlock()
	if !ok {
		// Offer not existed
		return
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	subCIDStr := mgr.digestOfferMapS[digest].GetSubCID().ToString()
	delete(mgr.digestOfferMapS, digest)

	// Update cid map
	delete(mgr.cidDigestMapS[subCIDStr], digest)
	if len(mgr.cidDigestMapS[subCIDStr]) == 0 {
		delete(mgr.cidDigestMapS, subCIDStr)
	}
}
