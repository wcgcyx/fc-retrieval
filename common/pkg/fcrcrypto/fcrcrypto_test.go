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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	PrvKey      = "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29"
	PrvKeyWrong = "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd16480"
	PubKey      = "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d"
	PubKeyWrong = "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62a"
	ID          = "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe"
)

func TestGenerateKey(t *testing.T) {
	prvKey, pubKey, id, err := GenerateRetrievalKeyPair()
	assert.Empty(t, err)
	assert.NotEmpty(t, prvKey)
	assert.NotEmpty(t, pubKey)
	assert.NotEmpty(t, id)
}

func TestGetPublicKey(t *testing.T) {
	pubKey, id, err := GetPublicKey(PrvKey)
	assert.Empty(t, err)
	assert.Equal(t, PubKey, pubKey)
	assert.Equal(t, ID, id)
}

func TestSign(t *testing.T) {
	sig, err := Sign(PrvKey, 0, []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.Empty(t, err)
	assert.Equal(t, "006e9654ac82348a7ff3ff5e0bf906a34c799f3841e0119ed32a64d32ba92258f2735c90af4295684485735c63d85514beda21037cdb2b501735cdccec6d3a625301", sig)

	sig, err = Sign("abcedfg", 0, []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.NotEmpty(t, err)
	assert.Empty(t, sig)

	sig, err = Sign(PrvKeyWrong, 0, []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.NotEmpty(t, err)
	assert.Empty(t, sig)
}

func TestVerify(t *testing.T) {
	err := Verify(PubKey, 0, "006e9654ac82348a7ff3ff5e0bf906a34c799f3841e0119ed32a64d32ba92258f2735c90af4295684485735c63d85514beda21037cdb2b501735cdccec6d3a625301", []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.Empty(t, err)

	err = Verify("abcdefg", 0, "006e9654ac82348a7ff3ff5e0bf906a34c799f3841e0119ed32a64d32ba92258f2735c90af4295684485735c63d85514beda21037cdb2b501735cdccec6d3a625301", []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.NotEmpty(t, err)

	err = Verify(PubKey, 0, "abcdefg", []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.NotEmpty(t, err)

	err = Verify(PubKey, 1, "006e9654ac82348a7ff3ff5e0bf906a34c799f3841e0119ed32a64d32ba92258f2735c90af4295684485735c63d85514beda21037cdb2b501735cdccec6d3a625301", []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.NotEmpty(t, err)

	err = Verify(PubKeyWrong[2:], 0, "006e9654ac82348a7ff3ff5e0bf906a34c799f3841e0119ed32a64d32ba92258f2735c90af4295684485735c63d85514beda21037cdb2b501735cdccec6d3a6253", []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.NotEmpty(t, err)

	err = Verify(PubKeyWrong[2:], 0, "006e9654ac82348a7ff3ff5e0bf906a34c799f3841e0119ed32a64d32ba92258f2735c90af4295684485735c63d85514beda21037cdb2b501735cdccec6d3a625301", []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.NotEmpty(t, err)

	err = Verify(PubKeyWrong, 0, "006e9654ac82348a7ff3ff5e0bf906a34c799f3841e0119ed32a64d32ba92258f2735c90af4295684485735c63d85514beda21037cdb2b501735cdccec6d3a625301", []byte{0x00, 0x01, 0x02, 0x03, 0x04})
	assert.NotEmpty(t, err)
}
