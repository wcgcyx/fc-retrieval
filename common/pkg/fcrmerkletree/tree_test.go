package fcrmerkletree

/*
 * Copyright 2021 ConsenSys Software Inc.
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
	"testing"

	"github.com/cbergoon/merkletree"
	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
)

func TestCreateTree(t *testing.T) {
	cid1, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	cid2, err := cid.NewContentID(Cid2Str)
	assert.Empty(t, err)
	cid3, err := cid.NewContentID(Cid3Str)
	assert.Empty(t, err)
	cid4, err := cid.NewContentID(Cid4Str)
	assert.Empty(t, err)
	cid5, err := cid.NewContentID(Cid5Str)
	assert.Empty(t, err)
	tree, err := CreateMerkleTree([]merkletree.Content{cid1, cid2, cid3, cid4, cid5})
	assert.Empty(t, err)
	assert.NotEmpty(t, tree)
	_, err = CreateMerkleTree([]merkletree.Content{})
	assert.NotEmpty(t, err)
	assert.Equal(t, "500e4c003e478a3b91857a6ad419d39b18d7da2f39b5fd694c0fbf0c2f7e783b", tree.GetMerkleRoot())
}

func TestCreateTreeOneElement(t *testing.T) {
	cid1, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	tree, err := CreateMerkleTree([]merkletree.Content{cid1})
	assert.Empty(t, err)
	assert.NotEmpty(t, tree)

	assert.Equal(t, "8e5d1faa539d8677e618eea34ed2fa7ae75c25bfe89c5caee2846554739557a6", tree.GetMerkleRoot())
}
