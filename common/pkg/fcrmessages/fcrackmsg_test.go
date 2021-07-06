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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestACKGetter(t *testing.T) {
	msg := CreateFCRACKMsg(100, []byte{1, 2, 3, 4})
	msg.signature = "testsignature"
	assert.Equal(t, true, msg.ACK())
	assert.Equal(t, uint64(100), msg.Nonce())
	assert.Equal(t, []byte{1, 2, 3, 4}, msg.Body())
	assert.Equal(t, "testsignature", msg.Signature())

	msg = CreateFCRACKErrorMsg(100, errors.New("Test error"))
	msg.signature = "testsignature2"
	assert.Equal(t, false, msg.ACK())
	assert.Equal(t, uint64(100), msg.Nonce())
	assert.Equal(t, "Test error", msg.Error())
	assert.Equal(t, "testsignature2", msg.Signature())
}

func TestACKParse(t *testing.T) {
	msg := CreateFCRACKMsg(100, []byte{1, 2, 3, 4})
	msg.signature = "testsignature"
	data, err := msg.ToBytes()
	assert.Empty(t, err)
	msg2 := FCRACKMsg{}
	err = msg2.FromBytes(data)
	assert.Empty(t, err)
	assert.Equal(t, true, msg2.ACK())
	assert.Equal(t, uint64(100), msg2.Nonce())
	assert.Equal(t, []byte{1, 2, 3, 4}, msg2.Body())
	assert.Equal(t, "testsignature", msg2.Signature())

	err = msg2.FromBytes([]byte{111, 111, 111})
	assert.NotEmpty(t, err)
}

func TestACKSigning(t *testing.T) {
	msg := CreateFCRACKMsg(100, []byte{1, 2, 3, 4})
	err := msg.Sign("wrongkey", 0)
	assert.NotEmpty(t, err)
	err = msg.Sign(PrivKey, 0)
	assert.Empty(t, err)
	assert.Equal(t, "0023cc79a46dd760a3ef2f6245cfdf028ed582b97a72cf9991cf8ff47e20f4ccdc3708a6a27e5f8057516b1322958cbc10fc012c828450fc45f245154c7488b41b00", msg.Signature())

	err = msg.Verify("wrongkey", 0)
	assert.NotEmpty(t, err)
	err = msg.Verify(PubKey, 1)
	assert.NotEmpty(t, err)
	err = msg.Verify(PubKey, 0)
	assert.Empty(t, err)
}
