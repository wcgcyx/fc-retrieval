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

// listFilesResponseJson represents the response of listing files.
type listFilesResponseJson struct {
	Files     []string `json:"files"`
	CIDs      []string `json:"cids"`
	Sizes     []int64  `json:"sizes"`
	Published []bool   `json:"published"`
	Frequency []int    `json:"frequency"`
}

// EncodeListFilesResponse is used to get the byte array of listFilesResponseJson
func EncodeListFilesResponse(
	files []string,
	cids []string,
	sizes []int64,
	published []bool,
	frequency []int,
) ([]byte, error) {
	return json.Marshal(&listFilesResponseJson{
		Files:     files,
		CIDs:      cids,
		Sizes:     sizes,
		Published: published,
		Frequency: frequency,
	})
}

// DecodeListFilesResponse is used to get the fields from byte array of listFilesResponseJson
func DecodeListFilesResponse(data []byte) (
	[]string, // files
	[]string, // cids
	[]int64, // sizes
	[]bool, // published
	[]int, // frequency
	error, // error
) {
	msg := listFilesResponseJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return msg.Files, msg.CIDs, msg.Sizes, msg.Published, msg.Frequency, nil
}
