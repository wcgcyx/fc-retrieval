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
	"encoding/hex"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestNewContentID(t *testing.T) {
	// Test V0
	id, err := NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	assert.Equal(t, "QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK", id.ToString())

	// Test V1
	id, err = NewContentID("baga6ea4seaqesauho7j2thfi4g4u5zbnhn2okd74s2igpvc2lsb7rrsfstoy4by")
	assert.Empty(t, err)
	assert.Equal(t, "baga6ea4seaqesauho7j2thfi4g4u5zbnhn2okd74s2igpvc2lsb7rrsfstoy4by", id.ToString())

	id, err = NewContentID("mAXCg5AIgiDG5WwMyC/w5d7U0XXkzok3wYOIejTx9IuLCaF7PNZI")
	assert.Empty(t, err)
	assert.Equal(t, "mAXCg5AIgiDG5WwMyC/w5d7U0XXkzok3wYOIejTx9IuLCaF7PNZI", id.ToString())

	// Test Fail
	id, err = NewContentID("QmAXCg5AIgiDG5WwMyC/w5d7U0XXkzok3wYOIejTx9IuLCaF7PNZI")
	assert.NotEmpty(t, err)
	assert.Empty(t, id)
}

func TestNewContentIDFromFile(t *testing.T) {
	m := fstest.MapFS{
		"test.txt": {
			Data: []byte("test, test, test"),
		},
	}
	file, err := m.Open("test.txt")
	assert.Empty(t, err)
	id, err := NewContentIDFromFile(file)
	assert.Empty(t, err)
	assert.Equal(t, "QmYiTrMsaCM9LRRybYzm3DkR93dt25Va75m5BbeggVNFRB", id.ToString())
}

func TestRandomContentID(t *testing.T) {
	id := NewRandomContentID()
	assert.NotEmpty(t, id)
}

func TestSerialization(t *testing.T) {
	cid1, err := NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	p, err := cid1.ToBytes()
	assert.Empty(t, err)
	cid2 := ContentID{}
	err = cid2.FromBytes(p)
	assert.Empty(t, err)
	assert.Equal(t, cid1.ToString(), cid2.ToString())
}

func TestCalculateHash(t *testing.T) {
	cid, err := NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	hash, err := cid.CalculateHash()
	assert.Empty(t, err)
	assert.Equal(t, "d85899bc6b64f09d68fd7b8d0fcea30c42b756250ef89f2cb12d03f907baacd7", hex.EncodeToString(hash))
}

func TestEqual(t *testing.T) {
	cid1, err := NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	cid2, err := NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	cid3, err := NewContentID("baga6ea4seaqesauho7j2thfi4g4u5zbnhn2okd74s2igpvc2lsb7rrsfstoy4by")
	assert.Empty(t, err)
	res1, err := cid1.Equals(cid2)
	assert.Empty(t, err)
	assert.True(t, res1)
	res2, err := cid1.Equals(cid3)
	assert.Empty(t, err)
	assert.False(t, res2)
}

func TestGetHashID(t *testing.T) {
	cid1, err := NewContentID("QmX5Rg8t9zh26JcaTk7VnDXqv5SHH2bT6AfeoTFLSsp4dK")
	assert.Empty(t, err)
	res, err := cid1.GetHashID()
	assert.Empty(t, err)
	assert.Equal(t, "d85899bc6b64f09d68fd7b8d0fcea30c42b756250ef89f2cb12d03f907baacd7", res)
}
