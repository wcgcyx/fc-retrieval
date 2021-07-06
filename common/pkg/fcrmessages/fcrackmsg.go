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
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
)

// FCRACKMsg is the response used in communication between filecoin retrieval entities.
type FCRACKMsg struct {
	ack         bool
	nonce       uint64
	messageBody []byte
	signature   string
}

// fcrACKMsgJson is used to parse to and from json.
type fcrACKMsgJson struct {
	ACK         bool   `json:"ack"`
	Nonce       uint64 `json:"uint64"`
	MessageBody string `json:"message_body"`
	Signature   string `json:"message_signature"`
}

// CreateFCRACKMsg is used to create an unsigned ack message
func CreateFCRACKMsg(nonce uint64, msgBody []byte) *FCRACKMsg {
	return &FCRACKMsg{
		ack:         true,
		nonce:       nonce,
		messageBody: msgBody,
		signature:   "",
	}
}

// CreateFCRACKErrorMsg is used to create an unsigned error message
func CreateFCRACKErrorMsg(nonce uint64, err error) *FCRACKMsg {
	return &FCRACKMsg{
		ack:         false,
		nonce:       nonce,
		messageBody: []byte(err.Error()),
		signature:   "",
	}
}

// Type is used to get the message type of the message.
func (fcrMsg *FCRACKMsg) ACK() bool {
	return fcrMsg.ack
}

// Nonce is used to get the nonce of the message.
func (fcrMsg *FCRACKMsg) Nonce() uint64 {
	return fcrMsg.nonce
}

// Body is used to get the message body.
func (fcrMsg *FCRACKMsg) Body() []byte {
	return fcrMsg.messageBody
}

// Error is used to get the error.
func (fcrMsg *FCRACKMsg) Error() string {
	return string(fcrMsg.Body())
}

// Signature is used to get the signature.
func (fcrMsg *FCRACKMsg) Signature() string {
	return fcrMsg.signature
}

// Sign is used to sign the message with a given private key and a key version.
func (fcrMsg *FCRACKMsg) Sign(privKey string, keyVer byte) error {
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, fcrMsg.nonce)
	var data []byte
	if fcrMsg.ack {
		data = append([]byte{0}, nonce...)
	} else {
		data = append([]byte{1}, nonce...)
	}
	data = append(data, fcrMsg.messageBody...)
	sig, err := fcrcrypto.Sign(privKey, keyVer, data)
	if err != nil {
		return err
	}
	fcrMsg.signature = sig
	return nil
}

// Verify is used to verify the offer with a given public key.
func (fcrMsg *FCRACKMsg) Verify(pubKey string, keyVer byte) error {
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, fcrMsg.nonce)
	var data []byte
	if fcrMsg.ack {
		data = append([]byte{0}, nonce...)
	} else {
		data = append([]byte{1}, nonce...)
	}
	data = append(data, fcrMsg.messageBody...)
	err := fcrcrypto.Verify(pubKey, keyVer, fcrMsg.signature, data)
	if err != nil {
		return fmt.Errorf("Message fail to verify: %v", err.Error())
	}
	return nil
}

// FCRMsgToBytes converts a FCRMessage to bytes
func (fcrMsg *FCRACKMsg) ToBytes() ([]byte, error) {
	fcrMsgJS := &fcrACKMsgJson{
		ACK:         fcrMsg.ack,
		Nonce:       fcrMsg.nonce,
		MessageBody: hex.EncodeToString(fcrMsg.messageBody),
		Signature:   fcrMsg.signature,
	}
	return json.Marshal(fcrMsgJS)
}

// FCRMsgFromBytes converts a bytes to FCRMessage
func (fcrMsg *FCRACKMsg) FromBytes(data []byte) error {
	res := fcrACKMsgJson{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	msgBody, err := hex.DecodeString(res.MessageBody)
	if err != nil {
		return err
	}
	fcrMsg.ack = res.ACK
	fcrMsg.nonce = res.Nonce
	fcrMsg.messageBody = msgBody
	fcrMsg.signature = res.Signature
	return nil
}
