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
	"errors"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
)

// FCRReqMsg is the request used in communication between filecoin retrieval entities.
type FCRReqMsg struct {
	messageType byte
	nonce       uint64
	messageBody []byte
	signature   string
}

// fcrMessageJson is used to parse to and from json.
type fcrReqMsgJson struct {
	MessageType string `json:"message_type"`
	Nonce       uint64 `json:"nonce"`
	MessageBody string `json:"message_body"`
	Signature   string `json:"message_signature"`
}

// CreateFCRReqMsg is used to create an unsigned message
func CreateFCRReqMsg(msgType byte, nonce uint64, msgBody []byte) *FCRReqMsg {
	return &FCRReqMsg{
		messageType: msgType,
		messageBody: msgBody,
		nonce:       nonce,
		signature:   "",
	}
}

// Type is used to get the message type of the message.
func (fcrMsg *FCRReqMsg) Type() byte {
	return fcrMsg.messageType
}

// Nonce is used to get the nonce of the message.
func (fcrMsg *FCRReqMsg) Nonce() uint64 {
	return fcrMsg.nonce
}

// Body is used to get the message body.
func (fcrMsg *FCRReqMsg) Body() []byte {
	return fcrMsg.messageBody
}

// Signature is used to get the signature.
func (fcrMsg *FCRReqMsg) Signature() string {
	return fcrMsg.signature
}

// Sign is used to sign the message with a given private key and a key version.
func (fcrMsg *FCRReqMsg) Sign(privKey string, keyVer byte) error {
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, fcrMsg.nonce)
	data := append([]byte{fcrMsg.messageType}, nonce...)
	data = append(data, fcrMsg.messageBody...)
	sig, err := fcrcrypto.Sign(privKey, keyVer, data)
	if err != nil {
		return err
	}
	fcrMsg.signature = sig
	return nil
}

// Verify is used to verify the offer with a given public key.
func (fcrMsg *FCRReqMsg) Verify(pubKey string, keyVer byte) error {
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, fcrMsg.nonce)
	data := append([]byte{fcrMsg.messageType}, nonce...)
	data = append(data, fcrMsg.messageBody...)
	err := fcrcrypto.Verify(pubKey, keyVer, fcrMsg.signature, data)
	if err != nil {
		return fmt.Errorf("Message fail to verify: %v", err.Error())
	}
	return nil
}

// VerifyByID is used to verify the offer with a given id (hashed public key).
func (fcrMsg *FCRReqMsg) VerifyByID(id string) error {
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, fcrMsg.nonce)
	data := append([]byte{fcrMsg.messageType}, nonce...)
	data = append(data, fcrMsg.messageBody...)
	err := fcrcrypto.VerifyByID(id, fcrMsg.signature, data)
	if err != nil {
		return fmt.Errorf("Message fail to verify: %v", err.Error())
	}
	return nil
}

// FCRMsgToBytes converts a FCRMessage to bytes
func (fcrMsg *FCRReqMsg) ToBytes() ([]byte, error) {
	fcrMsgJS := &fcrReqMsgJson{
		MessageType: hex.EncodeToString([]byte{fcrMsg.messageType}),
		Nonce:       fcrMsg.nonce,
		MessageBody: hex.EncodeToString(fcrMsg.messageBody),
		Signature:   fcrMsg.signature,
	}
	return json.Marshal(fcrMsgJS)
}

// FCRMsgFromBytes converts a bytes to FCRMessage
func (fcrMsg *FCRReqMsg) FromBytes(data []byte) error {
	res := fcrReqMsgJson{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	msgType, err := hex.DecodeString(res.MessageType)
	if err != nil {
		return err
	}
	if len(msgType) != 1 {
		return errors.New("Invalid message type length")
	}
	msgBody, err := hex.DecodeString(res.MessageBody)
	if err != nil {
		return err
	}
	fcrMsg.messageType = msgType[0]
	fcrMsg.nonce = res.Nonce
	fcrMsg.messageBody = msgBody
	fcrMsg.signature = res.Signature
	return nil
}
