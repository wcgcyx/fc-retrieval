/*
Package fcrcrypto - location for cryptographic tools to perform common operations on hashes, keys and signatures
*/
package fcrcrypto

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
	"errors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
	"github.com/minio/blake2b-simd"
)

// GenerateRetrievalKeyPair generates a new key,
// returns the private key, its associated public key, address and error.
func GenerateRetrievalKeyPair() (string, string, string, error) {
	prvKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", err
	}
	prvKeyStr := hex.EncodeToString(prvKey)
	pubKeyStr, addr, err := GetPublicKey(prvKeyStr)
	if err != nil {
		return "", "", "", err
	}
	return prvKeyStr, pubKeyStr, addr, nil
}

func GetPublicKey(prvKeyStr string) (string, string, error) {
	prvKey, err := hex.DecodeString(prvKeyStr)
	if err != nil {
		return "", "", err
	}
	pubKey := crypto.PublicKey(prvKey)
	addr, err := address.NewSecp256k1Address(pubKey)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(pubKey), addr.String(), nil
}

// Sign signs given bytes using given private key and given version,
// returns the signature in bytes and error.
func Sign(prvKeyStr string, ver byte, data []byte) (string, error) {
	prvKey, err := hex.DecodeString(prvKeyStr)
	if err != nil {
		return "", err
	}
	b2sum := blake2b.Sum256(data)
	sig, err := crypto.Sign(prvKey, b2sum[:])
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(append([]byte{ver}, sig...)), nil
}

// Verify verifies the given msg and its signature against the public key and key version,
// returns error.
func Verify(pubKeyStr string, ver byte, sigStr string, data []byte) error {
	pubKey, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return err
	}
	sig, err := hex.DecodeString(sigStr)
	if err != nil {
		return err
	}
	// First to check key version
	if sig[0] != ver {
		return errors.New("Key version mismatch")
	}
	sig = sig[1:]
	b2sum := blake2b.Sum256(data)
	pubk, err := crypto.EcRecover(b2sum[:], sig)
	if err != nil {
		return err
	}
	// Check public key
	if len(pubKey) != len(pubk) {
		return errors.New("Public key length mismatch")
	}
	for i := range pubKey {
		if pubKey[i] != pubk[i] {
			return errors.New("Public key mismatch")
		}
	}
	return nil
}
