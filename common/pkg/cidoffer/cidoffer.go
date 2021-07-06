/*
Package cidoffer - provides functionality like create, verify, sign and get details for CIDOffer and SubCIDOffer structures.

CIDOffer represents an offer from a Storage Provider, explaining on what conditions the client can retrieve a set of uniquely identified files from Filecoin blockchain network.
SubCIDOffer represents an offer from a Storage Provider, just like CIDOffer, but for a single file and includes a merkle proof
*/
package cidoffer

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
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/cbergoon/merkletree"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmerkletree"
)

// CIDOffer represents a CID Offer. That is, an offer to deliver content
// for Piece CID(s) at a certain price.
type CIDOffer struct {
	providerID string
	cids       []cid.ContentID
	price      *big.Int
	expiry     int64
	qos        uint64 // TODO: Represent Quality of Service (QoS)
	signature  string

	merkleRoot string
	merkleTree *fcrmerkletree.FCRMerkleTree
}

// cidOfferJson is used to parse to and from json.
type cidOfferJson struct {
	ProviderID string   `json:"provider_id"`
	CIDs       []string `json:"cids"`
	Price      string   `json:"price"`
	Expiry     int64    `json:"expiry"`
	QoS        uint64   `json:"qos"`
	Signature  string   `json:"signature"`
}

// cidOfferSigning is used to generate and verify signature.
type cidOfferSigning struct {
	providerID string
	merkleRoot string
	price      string
	expiry     int64
	qos        uint64
}

// NewCidOffer creates an unsigned CID Offer.
func NewCIDOffer(providerID string, cids []cid.ContentID, price *big.Int, expiry int64, qos uint64) (*CIDOffer, error) {
	if len(cids) < 1 {
		return nil, errors.New("Group CID Offer: need to provide at least 1 CID")
	}
	// TODO: Check that the expiry is in the future (are there scenarios where an expired offer should be loadable?)
	var c = CIDOffer{
		providerID: providerID,
		cids:       cids,
		price:      price,
		expiry:     expiry,
		qos:        qos,
	}

	// Create merkle tree & merkle root
	list := make([]merkletree.Content, len(cids))
	for i := 0; i < len(cids); i++ {
		list[i] = (cids)[i]
	}
	var err error
	c.merkleTree, err = fcrmerkletree.CreateMerkleTree(list)
	if err != nil {
		return nil, err
	}
	c.merkleRoot = c.merkleTree.GetMerkleRoot()

	return &c, nil
}

// GetProviderID returns the provider ID of this offer.
func (c *CIDOffer) GetProviderID() string {
	return c.providerID
}

// GetCIDs returns the cids of this offer.
func (c *CIDOffer) GetCIDs() []cid.ContentID {
	return c.cids
}

// GetPrice returns the price of this offer.
func (c *CIDOffer) GetPrice() *big.Int {
	return big.NewInt(0).Set(c.price)
}

// GetExpiry returns the expiry of this offer.
func (c *CIDOffer) GetExpiry() int64 {
	return c.expiry
}

// GetQoS returns the quality of service of this offer.
func (c *CIDOffer) GetQoS() uint64 {
	return c.qos
}

// GetSignature returns the signature of this offer.
func (c *CIDOffer) GetSignature() string {
	return c.signature
}

// SetSignature sets the signature of this offer.
func (c *CIDOffer) SetSignature(sig string) {
	c.signature = sig
}

// HasExpired returns true if the offer expiry date is in the past.
func (c *CIDOffer) HasExpired() bool {
	expiryTime := time.Unix(c.expiry, 0)
	now := time.Now()
	return expiryTime.Before(now)
}

// Sign is used to sign the offer with a given private key.
func (c *CIDOffer) Sign(privKey string) error {
	// Get msg bytes array
	data, err := json.Marshal(cidOfferSigning{
		providerID: c.providerID,
		merkleRoot: c.merkleRoot,
		price:      c.price.String(),
		expiry:     c.expiry,
		qos:        c.qos,
	})
	if err != nil {
		return err
	}
	// Sign
	sig, err := fcrcrypto.Sign(privKey, 0, data)
	if err != nil {
		return err
	}
	c.signature = sig
	return nil
}

// Verify is used to verify the offer with a given public key.
func (c *CIDOffer) Verify(pubKey string) error {
	// Get msg bytes array
	data, err := json.Marshal(cidOfferSigning{
		providerID: c.providerID,
		merkleRoot: c.merkleRoot,
		price:      c.price.String(),
		expiry:     c.expiry,
		qos:        c.qos,
	})
	if err != nil {
		return err
	}
	// Verify
	return fcrcrypto.Verify(pubKey, 0, c.signature, data)
}

// GenerateSubCIDOffer is used to generate a sub cid offer with proof for a given cid.
func (c *CIDOffer) GenerateSubCIDOffer(cid *cid.ContentID) (*SubCIDOffer, error) {
	proof, err := c.merkleTree.GenerateMerkleProof(cid)
	if err != nil {
		return nil, err
	}
	return NewSubCIDOffer(c.providerID, cid, c.merkleRoot, proof, c.price, c.expiry, c.qos, c.signature), nil
}

// GetMessageDigest calculate the message digest of this CID Group Offer.
// Note that the methodology used here should not be externally visible. The
// message digest should only be used within the system.
func (c *CIDOffer) GetMessageDigest() string {
	b := []byte(c.providerID)
	for _, id := range c.cids {
		b = append(b, []byte(id.ToString())...)
	}
	b = append(b, []byte(c.price.String())...)
	bExpiry := make([]byte, 8)
	binary.BigEndian.PutUint64(bExpiry, uint64(c.expiry))
	b = append(b, bExpiry...)
	bQoS := make([]byte, 8)
	binary.BigEndian.PutUint64(bQoS, uint64(c.qos))
	b = append(b, bQoS...)
	res := sha512.Sum512_256(b)
	return hex.EncodeToString(res[:])
}

// ToBytes is used to turn offer into bytes.
func (c *CIDOffer) ToBytes() ([]byte, error) {
	cidStrs := make([]string, 0)
	for _, id := range c.cids {
		cidStrs = append(cidStrs, id.ToString())
	}

	return json.Marshal(cidOfferJson{
		ProviderID: c.providerID,
		CIDs:       cidStrs,
		Price:      c.price.String(),
		Expiry:     c.expiry,
		QoS:        c.qos,
		Signature:  c.signature,
	})
}

// FromBytes is used to turn bytes into offer.
func (c *CIDOffer) FromBytes(p []byte) error {
	cJson := cidOfferJson{}
	err := json.Unmarshal(p, &cJson)
	if err != nil {
		return err
	}
	c.providerID = cJson.ProviderID
	cids := make([]cid.ContentID, 0)
	for _, cidStr := range cJson.CIDs {
		cid, err := cid.NewContentID(cidStr)
		if err != nil {
			return err
		}
		cids = append(cids, *cid)
	}
	c.cids = cids
	price, good := big.NewInt(0).SetString(cJson.Price, 10)
	if !good {
		return errors.New("Fail to decode price")
	}
	c.price = price
	c.expiry = cJson.Expiry
	c.qos = cJson.QoS
	c.signature = cJson.Signature
	// Reconstrct the merkle trie
	list := make([]merkletree.Content, len(c.cids))
	for i := 0; i < len(c.cids); i++ {
		list[i] = (c.cids)[i]
	}
	c.merkleTree, err = fcrmerkletree.CreateMerkleTree(list)
	if err != nil {
		return err
	}
	c.merkleRoot = c.merkleTree.GetMerkleRoot()
	return nil
}

// Copy returns a copy of the CID Offer
func (c *CIDOffer) Copy() *CIDOffer {
	data, _ := c.ToBytes()
	var copy CIDOffer
	copy.FromBytes(data)
	return &copy
}
