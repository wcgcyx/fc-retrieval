/*
Package fcrmessages - stores all the messages.
*/
package fcrmessages

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
	PrivKey = "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29"
	PubKey  = "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d"
	ID      = "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe"
)

func TestReqGetter(t *testing.T) {
	msg := CreateFCRReqMsg(1, 100, []byte{1, 2, 3, 4})
	msg.signature = "testsignature"
	assert.Equal(t, byte(1), msg.Type())
	assert.Equal(t, uint64(100), msg.Nonce())
	assert.Equal(t, []byte{1, 2, 3, 4}, msg.Body())
	assert.Equal(t, "testsignature", msg.Signature())
}

func TestReqParse(t *testing.T) {
	msg := CreateFCRReqMsg(1, 100, []byte{1, 2, 3, 4})
	msg.signature = "testsignature"
	data, err := msg.ToBytes()
	assert.Empty(t, err)
	msg2 := FCRReqMsg{}
	err = msg2.FromBytes(data)
	assert.Empty(t, err)
	assert.Equal(t, byte(1), msg2.Type())
	assert.Equal(t, uint64(100), msg2.Nonce())
	assert.Equal(t, []byte{1, 2, 3, 4}, msg2.Body())
	assert.Equal(t, "testsignature", msg2.Signature())

	err = msg2.FromBytes([]byte{111, 111, 111})
	assert.NotEmpty(t, err)
}

func TestReqSigning(t *testing.T) {
	msg := CreateFCRReqMsg(1, 100, []byte{1, 2, 3, 4})
	err := msg.Sign("wrongkey", 0)
	assert.NotEmpty(t, err)
	err = msg.Sign(PrivKey, 0)
	assert.Empty(t, err)
	assert.Equal(t, "00c224806c95a99cd31bce68f0d502c5abc23da173debea0cef3f61d25a6d0649b2dd1940131468b971992be78d51907cb9454cb05b23b42f5fc861255c12b24fb00", msg.Signature())

	err = msg.Verify("wrongkey", 0)
	assert.NotEmpty(t, err)
	err = msg.Verify(PubKey, 1)
	assert.NotEmpty(t, err)
	err = msg.Verify(PubKey, 0)
	assert.Empty(t, err)

	// Test Verify By ID
	err = msg.VerifyByID("wrong id")
	assert.NotEmpty(t, err)

	err = msg.VerifyByID(ID)
	assert.Empty(t, err)
}
