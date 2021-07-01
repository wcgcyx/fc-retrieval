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

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmerkletree"
)

// SubCIDOffer represents a sub CID Offer. That is, part of a CID offer.
// It contains one sub cid and a merkle proof showing that this sub cid
// is part of the cid array in the original cid offer.
type SubCIDOffer struct {
	providerID string
	subCID     *cid.ContentID
	price      *big.Int
	expiry     int64
	qos        uint64
	signature  string

	merkleRoot  string
	merkleProof *fcrmerkletree.FCRMerkleProof
}

// subCIDOfferJson is used to parse to and from json.
type subCIDOfferJson struct {
	ProviderID  string `json:"provider_id"`
	SubCID      string `json:"sub_cid"`
	MerkleRoot  string `json:"merkle_root"`
	MerkleProof string `json:"merkle_proof"`
	Price       string `json:"price"`
	Expiry      int64  `json:"expiry"`
	QoS         uint64 `json:"qos"`
	Signature   string `json:"signature"`
}

// subCIDOfferSigning is used to generate and verify signature.
type subCIDOfferSigning struct {
	providerID string
	merkleRoot string
	price      string
	expiry     int64
	qos        uint64
}

// NewSubCIDOffer creates a sub CID Offer.
func NewSubCIDOffer(providerID string, subCID *cid.ContentID, merkleRoot string, merkleProof *fcrmerkletree.FCRMerkleProof, price *big.Int, expiry int64, qos uint64, signature string) *SubCIDOffer {
	return &SubCIDOffer{
		providerID:  providerID,
		subCID:      subCID,
		merkleRoot:  merkleRoot,
		merkleProof: merkleProof,
		price:       price,
		expiry:      expiry,
		qos:         qos,
		signature:   signature,
	}
}

// GetProviderID returns the provider ID of this offer.
func (c *SubCIDOffer) GetProviderID() string {
	return c.providerID
}

// GetSubCID returns the sub cid of this offer.
func (c *SubCIDOffer) GetSubCID() *cid.ContentID {
	return c.subCID
}

// GetMerkleRoot returns the merkle root of this offer.
func (c *SubCIDOffer) GetMerkleRoot() string {
	return c.merkleRoot
}

// GetMerkleProof returns the merkle proof of this offer.
func (c *SubCIDOffer) GetMerkleProof() *fcrmerkletree.FCRMerkleProof {
	return c.merkleProof
}

// GetPrice returns the price of this offer.
func (c *SubCIDOffer) GetPrice() *big.Int {
	return big.NewInt(0).Set(c.price)
}

// GetExpiry returns the expiry of this offer.
func (c *SubCIDOffer) GetExpiry() int64 {
	return c.expiry
}

// GetQoS returns the quality of service of this offer.
func (c *SubCIDOffer) GetQoS() uint64 {
	return c.qos
}

// GetSignature returns the signature of this offer.
func (c *SubCIDOffer) GetSignature() string {
	return c.signature
}

// HasExpired returns true if the offer expiry date is in the past.
func (c *SubCIDOffer) HasExpired() bool {
	expiryTime := time.Unix(c.expiry, 0)
	now := time.Now()
	return expiryTime.Before(now)
}

// Verify is used to verify the offer with a given public key.
func (c *SubCIDOffer) Verify(pubKey string) error {
	data, err := json.Marshal(subCIDOfferSigning{
		providerID: c.providerID,
		merkleRoot: c.merkleRoot,
		price:      c.price.String(),
		expiry:     c.expiry,
		qos:        c.qos,
	})
	if err != nil {
		return err
	}
	return fcrcrypto.Verify(pubKey, 0, c.signature, data)
}

// VerifyMerkleProof is used to verify the sub cid is part of the merkle trie
func (c *SubCIDOffer) VerifyMerkleProof() error {
	if c.merkleProof.VerifyContent(c.subCID, c.merkleRoot) {
		return nil
	}
	return errors.New("Offer does not pass merkle proof verification")
}

// GetMessageDigest calculate the message digest of this sub CID Offer.
// Note that the methodology used here should not be externally visible. The
// message digest should only be used within the system.
func (c *SubCIDOffer) GetMessageDigest() string {
	b := []byte(c.providerID)
	b = append(b, []byte(c.merkleRoot)...)
	b = append(b, []byte(c.subCID.ToString())...)
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
func (c *SubCIDOffer) ToBytes() ([]byte, error) {
	// Merkle proof to string
	proofData, err := c.merkleProof.ToBytes()
	if err != nil {
		return nil, err
	}

	return json.Marshal(subCIDOfferJson{
		ProviderID:  c.providerID,
		SubCID:      c.subCID.ToString(),
		MerkleRoot:  c.merkleRoot,
		MerkleProof: hex.EncodeToString(proofData),
		Price:       c.price.String(),
		Expiry:      c.expiry,
		QoS:         c.qos,
		Signature:   c.signature,
	})
}

// FromBytes is used to turn bytes into offer.
func (c *SubCIDOffer) FromBytes(p []byte) error {
	cJson := subCIDOfferJson{}
	err := json.Unmarshal(p, &cJson)
	if err != nil {
		return err
	}
	c.providerID = cJson.ProviderID
	c.subCID, err = cid.NewContentID(cJson.SubCID)
	if err != nil {
		return err
	}
	c.merkleRoot = cJson.MerkleRoot
	proofData, err := hex.DecodeString(cJson.MerkleProof)
	if err != nil {
		return err
	}
	c.merkleProof = &fcrmerkletree.FCRMerkleProof{}
	err = c.merkleProof.FromBytes(proofData)
	if err != nil {
		return err
	}
	price, good := big.NewInt(0).SetString(cJson.Price, 10)
	if !good {
		return errors.New("Fail to decode price")
	}
	c.price = price
	c.expiry = cJson.Expiry
	c.qos = cJson.QoS
	c.signature = cJson.Signature
	return nil
}

// Copy returns a copy of the Sub CID Offer
func (c *SubCIDOffer) Copy() *SubCIDOffer {
	data, _ := c.ToBytes()
	var copy SubCIDOffer
	copy.FromBytes(data)
	return &copy
}
