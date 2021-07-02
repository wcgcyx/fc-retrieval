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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
)

// FCRMessage is the message used in communication between filecoin retrieval entities.
type FCRMessage struct {
	messageType byte
	messageBody []byte
	signature   string
}

// fcrMessageJson is used to parse to and from json.
type fcrMessageJson struct {
	MessageType string `json:"message_type"`
	MessageBody string `json:"message_body"`
	Signature   string `json:"message_signature"`
}

// CreateFCRMessage is used to create an unsigned message
func CreateFCRMessage(msgType byte, msgBody []byte) *FCRMessage {
	return &FCRMessage{
		messageType: msgType,
		messageBody: msgBody,
		signature:   "",
	}
}

// GetMessageType is used to get the message type of the message.
func (fcrMsg *FCRMessage) GetMessageType() byte {
	return fcrMsg.messageType
}

// GetMessageBody is used to get the message body.
func (fcrMsg *FCRMessage) GetMessageBody() []byte {
	return fcrMsg.messageBody
}

// GetSignature is used to get the signature.
func (fcrMsg *FCRMessage) GetSignature() string {
	return fcrMsg.signature
}

// Sign is used to sign the message with a given private key and a key version.
func (fcrMsg *FCRMessage) Sign(privKey string, keyVer byte) error {
	data := append([]byte{fcrMsg.messageType}, fcrMsg.messageBody...)
	sig, err := fcrcrypto.Sign(privKey, keyVer, data)
	if err != nil {
		return err
	}
	fcrMsg.signature = sig
	return nil
}

// Verify is used to verify the offer with a given public key.
func (fcrMsg *FCRMessage) Verify(pubKey string, keyVer byte) error {
	data := append([]byte{fcrMsg.messageType}, fcrMsg.messageBody...)
	err := fcrcrypto.Verify(pubKey, keyVer, fcrMsg.signature, data)
	if err != nil {
		return fmt.Errorf("Message fail to verify: %v", err.Error())
	}
	return nil
}

// VerifyByID is used to verify the offer with a given id (hashed public key).
func (fcrMsg *FCRMessage) VerifyByID(id string) error {
	data := append([]byte{fcrMsg.messageType}, fcrMsg.messageBody...)
	err := fcrcrypto.VerifyByID(id, fcrMsg.signature, data)
	if err != nil {
		return fmt.Errorf("Message fail to verify: %v", err.Error())
	}
	return nil
}

// FCRMsgToBytes converts a FCRMessage to bytes
func (fcrMsg *FCRMessage) ToBytes() ([]byte, error) {
	fcrMsgJS := &fcrMessageJson{
		MessageType: hex.EncodeToString([]byte{fcrMsg.messageType}),
		MessageBody: hex.EncodeToString(fcrMsg.messageBody),
		Signature:   fcrMsg.signature,
	}
	return json.Marshal(fcrMsgJS)
}

// FCRMsgFromBytes converts a bytes to FCRMessage
func FromBytes(data []byte) (*FCRMessage, error) {
	res := fcrMessageJson{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	msgType, err := hex.DecodeString(res.MessageType)
	if err != nil {
		return nil, err
	}
	if len(msgType) != 1 {
		return nil, errors.New("Invalid message type length")
	}
	msgBody, err := hex.DecodeString(res.MessageBody)
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		messageType: msgType[0],
		messageBody: msgBody,
		signature:   res.Signature,
	}, nil
}
