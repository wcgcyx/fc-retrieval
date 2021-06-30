/*
Package cid - provides methods for ContentID struct.

ContentID is wrapper over cid of a file stored in the system.
*/
package cid

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
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"io/ioutil"

	"github.com/cbergoon/merkletree"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

// ContentID represents a CID.
type ContentID struct {
	id string
}

// NewContentID creates a ContentID object from a cid string.
func NewContentID(cidStr string) (*ContentID, error) {
	_, err := cid.Parse(cidStr)
	if err != nil {
		return nil, err
	}
	return &ContentID{cidStr}, nil
}

// NewContentIDFromFile creates a ContentID object from a given file.
func NewContentIDFromFile(file fs.File) (*ContentID, error) {
	// TODO:
	// 1. Might have problem for large file.
	// 2. For now, it just uses the plain bytes, we need to generate a merkle DAG
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	pref := cid.Prefix{
		Version:  0,
		Codec:    cid.DagProtobuf,
		MhType:   multihash.SHA2_256,
		MhLength: -1, // Default length
	}
	id, err := pref.Sum(data)
	if err != nil {
		return nil, err
	}
	return &ContentID{id.String()}, nil
}

// NewRandomContentID creates a random ContentID object.
func NewRandomContentID() *ContentID {
	// We generate a version 0 cid
	pref := cid.Prefix{
		Version:  0,
		Codec:    cid.DagProtobuf,
		MhType:   multihash.SHA2_256,
		MhLength: -1, // Default length
	}

	seed := make([]byte, 32)
	rand.Read(seed)

	id, err := pref.Sum(seed)
	if err != nil {
		panic(err) // This should never happen.
	}
	return &ContentID{id.String()}
}

// GetHashID gets a 32 bytes hash of this cid string in hex string format
func (n *ContentID) GetHashID() (string, error) {
	temp, err := n.CalculateHash()
	if err != nil || len(temp) != 32 {
		return "", err
	}
	return hex.EncodeToString(temp), nil
}

// ToString returns a string for the ContentID.
func (n *ContentID) ToString() string {
	return n.id
}

// ToBytes is used to turn CID into bytes.
func (n *ContentID) ToBytes() ([]byte, error) {
	return json.Marshal(n.id)
}

// FromBytes is used to turn bytes into ContentID.
func (n *ContentID) FromBytes(p []byte) error {
	return json.Unmarshal(p, &n.id)
}

// CalculateHash hashes the values of a ContentID.
func (n ContentID) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(n.id)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Equals tests for equality of two ContentIDs.
func (n ContentID) Equals(other merkletree.Content) (bool, error) {
	return n.id == other.(*ContentID).id, nil
}
