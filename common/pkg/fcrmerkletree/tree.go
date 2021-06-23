package fcrmerkletree

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

	"github.com/cbergoon/merkletree"
)

// FCRMerkleTree is used to store a list of CIDs.
type FCRMerkleTree struct {
	tree *merkletree.MerkleTree
}

// CreateMerkleTree creates a merkle tree from a list of cids.
func CreateMerkleTree(contents []merkletree.Content) (*FCRMerkleTree, error) {
	tree, err := merkletree.NewTree(contents)
	if err != nil {
		return nil, err
	}
	return &FCRMerkleTree{tree: tree}, nil
}

// GetMerkleRoot returns the merkle root of the tree.
func (mt *FCRMerkleTree) GetMerkleRoot() string {
	return hex.EncodeToString(mt.tree.MerkleRoot())
}

// GenerateMerkleProof gets the merkle proof for a given cid.
func (mt *FCRMerkleTree) GenerateMerkleProof(content merkletree.Content) (*FCRMerkleProof, error) {
	path, index, err := mt.tree.GetMerklePath(content)
	if err != nil {
		return nil, err
	}
	return &FCRMerkleProof{path: path, index: index}, nil
}
