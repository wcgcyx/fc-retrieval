/*
Package dhtring - provides operations like find a closest node, add new and remove for a Distributed Hash Table Ring data structure
*/
package dhtring

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
)

// ringNode is a node inside the ring
type ringNode struct {
	prv     *ringNode
	distPrv *big.Int

	key *big.Int
	val string

	distNext *big.Int
	next     *ringNode
}

// Ring is a struct to store the DHT Ring to store 32-bytes hex
type Ring struct {
	entry *ringNode
	size  int
}

// CreateRing creates a new ring data structure
func CreateRing() *Ring {
	return &Ring{
		entry: nil,
		size:  0,
	}
}

// Insert inserts a hex string into the ring
func (r *Ring) Insert(hex string) {
	if !validateInput(hex) {
		return
	}
	// Construct new node
	hexKey, _ := new(big.Int).SetString(hex, 16)
	newNode := &ringNode{
		prv:      nil,
		distPrv:  nil,
		key:      hexKey,
		val:      hex,
		distNext: nil,
		next:     nil,
	}
	// If size is 0
	if r.size == 0 {
		r.entry = newNode
		r.size++
		return
	}
	// If size is 1
	if r.size == 1 {
		cmp := r.entry.key.Cmp(newNode.key)
		if cmp == 0 {
			return
		}
		// Connect entry and new node
		newNode.prv = r.entry
		newNode.distPrv = getDist(r.entry.key, newNode.key)
		r.entry.next = newNode
		r.entry.distNext = getDist(r.entry.key, newNode.key)

		newNode.next = r.entry
		newNode.distNext = getDist(newNode.key, r.entry.key)
		r.entry.prv = newNode
		r.entry.distPrv = getDist(newNode.key, r.entry.key)
		r.size++
		return
	}
	// r.size >= 2
	prv := r.entry.prv
	current := r.entry
	for ok := true; ok; ok = current.val != r.entry.val {
		if current.val == hex {
			return
		}
		if between(prv.key, current.key, newNode.key) {
			// Put as prv -> newNode -> current
			newNode.prv = prv
			newNode.distPrv = getDist(prv.key, newNode.key)
			prv.next = newNode
			prv.distNext = getDist(prv.key, newNode.key)

			newNode.next = current
			newNode.distNext = getDist(newNode.key, current.key)
			current.prv = newNode
			current.distPrv = getDist(newNode.key, current.key)
			r.size++
			return
		}
		current = current.next
		prv = current.prv
	}
}

// Remove inserts a given hex string out of the ring
func (r *Ring) Remove(hex string) {
	if !validateInput(hex) {
		return
	}
	hexKey, _ := new(big.Int).SetString(hex, 16)
	// If size is 0
	if r.size == 0 {
		return
	}
	// If size is 1
	if r.size == 1 {
		if r.entry.val == hex {
			r.entry = nil
			r.size = 0
		}
		return
	}
	// If size is 2
	if r.size == 2 {
		node1 := r.entry
		node2 := r.entry.next
		if node1.val == hex {
			// Remove node1
			node2.prv = nil
			node2.distPrv = nil

			node2.next = nil
			node2.distNext = nil
			r.entry = node2
			r.size = 1
			return
		}
		if node2.val == hex {
			// Remove node2
			node1.prv = nil
			node1.distPrv = nil

			node1.next = nil
			node1.distNext = nil
			r.entry = node1
			r.size = 1
			return
		}
		return
	}
	// If size >= 3
	var toRemove *ringNode
	prv := r.entry.prv
	current := r.entry
	for ok := true; ok; ok = current.val != r.entry.val {
		if current.val == hex {
			toRemove = current
			break
		}
		if between(prv.key, current.key, hexKey) {
			break
		}
		current = current.next
	}
	if toRemove != nil {
		// Change from prv -> toRemove -> next
		// to prv -> next
		prv := toRemove.prv
		next := toRemove.next

		prv.next = next
		prv.distNext = getDist(prv.key, next.key)
		next.prv = prv
		next.distPrv = getDist(prv.key, next.key)
		// Change to next if to remove is the entry
		if toRemove.val == r.entry.val {
			r.entry = next
		}
		r.size--
	}
}

// GetClosest gets the closest hexes close to the given hex
func (r *Ring) GetClosest(hex string, num int, exclude string) []string {
	if !validateInput(hex) || (exclude != "" && !validateInput(exclude)) {
		return nil
	}
	// First consider exclusion
	if exclude != "" {
		before := r.size
		r.Remove(exclude)
		if r.size != before {
			defer r.Insert(exclude)
		}
	}
	res := make([]string, 0)
	if r.size == 0 || num == 0 {
		return res
	}
	if num >= r.size {
		// Add everything in the ring
		current := r.entry
		for ok := true; ok; ok = current != nil && current != r.entry {
			res = append(res, current.val)
			current = current.next
		}
		return res
	}
	// Return partial result
	// Insert hex -> search hex -> get result -> remove hex if need
	before := r.size
	r.Insert(hex)
	if r.size != before {
		defer r.Remove(hex)
	} else {
		// This already exists
		res = append(res, hex)
	}
	// Now search
	var anc *ringNode
	current := r.entry
	for ok := true; ok; ok = current != nil && current != r.entry {
		if current.val == hex {
			anc = current
			break
		}
		current = current.next
	}
	// Search from anchor
	prv := anc.prv
	distToPrv := big.NewInt(0)
	distToPrv.Add(distToPrv, anc.distPrv)
	next := anc.next
	distToNext := big.NewInt(0)
	distToNext.Add(distToNext, anc.distNext)
	for len(res) < num {
		cmp := distToPrv.Cmp(distToNext)
		if cmp <= 0 {
			res = append([]string{prv.val}, res...)
			distToPrv.Add(distToPrv, prv.distPrv)
			prv = prv.prv
		} else {
			res = append(res, next.val)
			distToNext.Add(distToNext, next.distNext)
			next = next.next
		}
	}
	return res
}

// Size gets the size of the ring
func (r *Ring) Size() int {
	return r.size
}

// getDist gets the distance from one to another, clockwise
func getDist(from *big.Int, to *big.Int) *big.Int {
	// So from is always smaller than to
	if from.Cmp(to) < 0 {
		return big.NewInt(0).Sub(to, from)
	} else {
		// It has across the max/min boundary
		max, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
		min, _ := new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
		dist1 := big.NewInt(0).Sub(max, from)
		dist2 := big.NewInt(0).Sub(to, min)
		sum := big.NewInt(0).Add(dist1, dist2)
		return sum.Add(sum, big.NewInt(1))
	}
}

// between checks if new is between prv and current clockwise
func between(prv *big.Int, current *big.Int, new *big.Int) bool {
	// so prv is always smaller than current
	if prv.Cmp(current) < 0 {
		// check if new is bigger than prv and smaller than current
		return new.Cmp(prv) > 0 && new.Cmp(current) < 0
	} else {
		// It has across the max/min boundary
		// check if new is bigger than prv or smaller than current
		return new.Cmp(prv) > 0 || new.Cmp(current) < 0
	}
}

// validateInput makes sure the given hex string is 32 bytes hex string
func validateInput(hex string) bool {
	if len(hex) != 64 {
		return false
	}
	for _, char := range hex {
		if (char < '0' || char > '9') && (char < 'A' || char > 'F') && (char < 'a' || char > 'f') {
			return false
		}
	}
	return true
}
