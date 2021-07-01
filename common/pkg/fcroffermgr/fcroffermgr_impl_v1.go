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
)

// FCROfferMgrImplV1 implements FCROfferMgr interface, it is an in-memory version.
type FCROfferMgrImplV1 struct {
	// tracking indicates whether to track count
	tracking bool

	// lock is the lock to offer storage
	lock sync.RWMutex

	// dMap is the digest offer map, map from digest string -> offer entry
	dMap map[string]*offerStorageEntry

	// cidMap is a map from cid string -> (map from digest string -> true)
	cidMap map[string]map[string]bool

	// tagMap is a map from tag string -> (map from digest string -> true)
	tagMap map[string]map[string]bool

	// countMap is a map from count -> (map from digest string -> true)
	countMap map[int]map[string]bool

	// dsMap is the digest sub offer map, map from digest string -> sub offer
	dsMap map[string]*cidoffer.SubCIDOffer

	// cidsMap is a map from cid string -> (map from digest string -> true)
	cidsMap map[string]map[string]bool
}

// offerStorageEntry is an entry storing cid offer
type offerStorageEntry struct {
	offer *cidoffer.CIDOffer
	tag   string
	count int
}

func NewFCROfferMgrImplV1(tracking bool) FCROfferMgr {
	return &FCROfferMgrImplV1{
		tracking: tracking,
		lock:     sync.RWMutex{},
		dMap:     make(map[string]*offerStorageEntry),
		cidMap:   make(map[string]map[string]bool),
		tagMap:   make(map[string]map[string]bool),
		countMap: make(map[int]map[string]bool),
		dsMap:    make(map[string]*cidoffer.SubCIDOffer),
		cidsMap:  make(map[string]map[string]bool),
	}
}

func (mgr *FCROfferMgrImplV1) Start() error {
	return nil
}

func (mgr *FCROfferMgrImplV1) Shutdown() {
}

func (mgr *FCROfferMgrImplV1) AddOffer(offer *cidoffer.CIDOffer) {
	mgr.AddOfferWithTag(offer, "")
}

func (mgr *FCROfferMgrImplV1) AddOfferWithTag(offer *cidoffer.CIDOffer, tag string) {
	digest := offer.GetMessageDigest()
	mgr.lock.RLock()
	_, ok := mgr.dMap[digest]
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
		panic("Fail to get an offer copy")
	}
	mgr.dMap[digest] = &offerStorageEntry{
		offer: copy,
		tag:   tag,
		count: 0,
	}
	for _, cid := range copy.GetCIDs() {
		_, ok = mgr.cidMap[cid.ToString()]
		if !ok {
			mgr.cidMap[cid.ToString()] = make(map[string]bool)
		}
		mgr.cidMap[cid.ToString()][digest] = true
	}
	_, ok = mgr.tagMap[tag]
	if !ok {
		mgr.tagMap[tag] = make(map[string]bool)
	}
	mgr.tagMap[tag][digest] = true
	_, ok = mgr.countMap[0]
	if !ok {
		mgr.countMap[0] = map[string]bool{}
	}
	mgr.countMap[0][digest] = true
}

func (mgr *FCROfferMgrImplV1) GetOffers(cID *cid.ContentID) []cidoffer.CIDOffer {
	res := make([]cidoffer.CIDOffer, 0)
	mgr.lock.RLock()
	digests, ok := mgr.cidMap[cID.ToString()]
	if !ok {
		mgr.lock.RUnlock()
		return res
	}
	for digest := range digests {
		copy := mgr.dMap[digest].offer.Copy()
		if copy == nil {
			panic("Fail to get an offer copy")
		}
		res = append(res, *copy)
		if mgr.tracking {
			mgr.lock.RUnlock()
			mgr.lock.Lock()
			// Update count
			prvCount := mgr.dMap[digest].count
			mgr.dMap[digest].count++
			delete(mgr.countMap[prvCount], digest)
			if len(mgr.countMap[prvCount]) == 0 {
				delete(mgr.countMap, prvCount)
			}
			newCount := prvCount + 1
			_, ok = mgr.countMap[newCount]
			if !ok {
				mgr.countMap[newCount] = make(map[string]bool)
			}
			mgr.countMap[newCount][digest] = true
			mgr.lock.Unlock()
			mgr.lock.RLock()
		}
	}
	mgr.lock.RUnlock()
	return res
}

func (mgr *FCROfferMgrImplV1) GetOffersByTag(tag string) []cidoffer.CIDOffer {
	res := make([]cidoffer.CIDOffer, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	digests, ok := mgr.tagMap[tag]
	if !ok {
		return res
	}
	for digest := range digests {
		copy := mgr.dMap[digest].offer.Copy()
		if copy == nil {
			panic("Fail to get an offer copy")
		}
		res = append(res, *copy)
	}
	return res
}

func (mgr *FCROfferMgrImplV1) ListOffers(from uint, to uint) []cidoffer.CIDOffer {
	res := make([]cidoffer.CIDOffer, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	if from >= to || from >= uint(len(mgr.dMap)) {
		return res
	}
	i := uint(0)
	for _, val := range mgr.dMap {
		if i >= from && i < to {
			copy := val.offer.Copy()
			if copy == nil {
				panic("Fail to get an offer copy")
			}
			res = append(res, *copy)
		}
		i++
	}
	return res
}

func (mgr *FCROfferMgrImplV1) ListOffersWithTag(from uint, to uint) ([]cidoffer.CIDOffer, []string) {
	res1 := make([]cidoffer.CIDOffer, 0)
	res2 := make([]string, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	if from >= to || from >= uint(len(mgr.dMap)) {
		return res1, res2
	}
	i := uint(0)
	for _, val := range mgr.dMap {
		if i >= from && i < to {
			copy := val.offer.Copy()
			if copy == nil {
				panic("Fail to get an offer copy")
			}
			res1 = append(res1, *copy)
			res2 = append(res2, val.tag)
		}
		i++
	}
	return res1, res2
}

func (mgr *FCROfferMgrImplV1) ListOffersWithAccessCount(from uint, to uint) ([]cidoffer.CIDOffer, []int) {
	res1 := make([]cidoffer.CIDOffer, 0)
	res2 := make([]int, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	if from >= to || from >= uint(len(mgr.dMap)) {
		return res1, res2
	}
	keys := make([]int, len(mgr.countMap))
	i := 0
	for key := range mgr.countMap {
		keys[i] = key
		i++
	}
	sort.Ints(keys)

	index := uint(0)
	for key := range keys {
		for digest := range mgr.countMap[key] {
			if index >= from {
				if index < to {
					copy := mgr.dMap[digest].offer.Copy()
					if copy == nil {
						panic("Fail to get an offer copy")
					}
					if key != mgr.dMap[digest].count {
						panic("Offer access count mismatch")
					}
					res1 = append(res1, *copy)
					res2 = append(res2, key)
				} else {
					return res1, res2
				}
			}
			index++
		}
	}
	return res1, res2
}

func (mgr *FCROfferMgrImplV1) RemoveOffer(digest string) {
	mgr.lock.RLock()
	_, ok := mgr.dMap[digest]
	mgr.lock.RUnlock()
	if !ok {
		// Offer not existed
		return
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	cids := mgr.dMap[digest].offer.GetCIDs()
	count := mgr.dMap[digest].count
	tag := mgr.dMap[digest].tag
	delete(mgr.dMap, digest)
	// Update cid map
	for _, cid := range cids {
		delete(mgr.cidMap[cid.ToString()], digest)
		if len(mgr.cidMap[cid.ToString()]) == 0 {
			delete(mgr.cidMap, cid.ToString())
		}
	}
	// Update count map
	delete(mgr.countMap[count], digest)
	if len(mgr.countMap[count]) == 0 {
		delete(mgr.countMap, count)
	}
	// Update tag map
	delete(mgr.tagMap[tag], digest)
	if len(mgr.tagMap[tag]) == 0 {
		delete(mgr.tagMap, tag)
	}
}

func (mgr *FCROfferMgrImplV1) AddSubOffer(offer *cidoffer.SubCIDOffer) {
	digest := offer.GetMessageDigest()
	mgr.lock.RLock()
	_, ok := mgr.dsMap[digest]
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
		panic("Fail to get an offer copy")
	}
	mgr.dsMap[digest] = copy
	subCID := offer.GetSubCID().ToString()
	_, ok = mgr.cidsMap[subCID]
	if !ok {
		mgr.cidsMap[subCID] = map[string]bool{}
	}
	mgr.cidsMap[subCID][digest] = true
}

func (mgr *FCROfferMgrImplV1) GetSubOffers(cID *cid.ContentID) []cidoffer.SubCIDOffer {
	res := make([]cidoffer.SubCIDOffer, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	digests, ok := mgr.cidsMap[cID.ToString()]
	if !ok {
		return res
	}
	for digest := range digests {
		copy := mgr.dsMap[digest].Copy()
		if copy == nil {
			panic("Fail to get an offer copy")
		}
		res = append(res, *copy)
	}
	return res
}

func (mgr *FCROfferMgrImplV1) ListSubOffers(from uint, to uint) []cidoffer.SubCIDOffer {
	res := make([]cidoffer.SubCIDOffer, 0)
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	if from >= to || from >= uint(len(mgr.dMap)) {
		return res
	}
	i := uint(0)
	for _, val := range mgr.dsMap {
		if i >= from && i < to {
			copy := val.Copy()
			if copy == nil {
				panic("Fail to get an offer copy")
			}
			res = append(res, *copy)
		}
		i++
	}
	return res
}

func (mgr *FCROfferMgrImplV1) RemoveSubOffer(digest string) {
	mgr.lock.RLock()
	_, ok := mgr.dsMap[digest]
	mgr.lock.RUnlock()
	if !ok {
		// Offer not existed
		return
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	cid := mgr.dsMap[digest].GetSubCID()
	delete(mgr.dsMap, digest)
	// Update cid map
	delete(mgr.cidsMap[cid.ToString()], digest)
	if len(mgr.cidsMap[cid.ToString()]) == 0 {
		delete(mgr.cidsMap, cid.ToString())
	}
}
