/*
Package fcradminmsg - stores all the admin messages.
*/
package fcradminmsg

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

import "encoding/json"

// uploadFileStartRequestJson represents the request to upload a file.
type uploadFileStartRequestJson struct {
	Tag  string `json:"tag"`
	Data []byte `json:"data"`
}

// EncodeUploadFileStartRequest is used to get the byte array of uploadFileStartRequestJson
func EncodeUploadFileStartRequest(
	tag string,
	data []byte,
) ([]byte, error) {
	return json.Marshal(&uploadFileStartRequestJson{
		Tag:  tag,
		Data: data,
	})
}

// DecodeUploadFileStartRequest is used to get the fields from byte array of uploadFileStartRequestJson
func DecodeUploadFileStartRequest(data []byte) (
	string, // tag
	[]byte, // data
	error, // error
) {
	msg := uploadFileStartRequestJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return "", nil, err
	}
	return msg.Tag, msg.Data, nil
}
