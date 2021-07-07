/*
Package fcradminserver - provides an interface to do admin networking.
*/
package fcradminserver

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
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// FCRAdminServer represents a server handling admin requests.
type FCRAdminServer interface {
	// Start starts the server.
	Start() error

	// Shutdown stops the server.
	Shutdown()

	// AddHandler adds a handler to the server, which handles a given message type (POST).
	AddHandler(msgType byte, handler func(data []byte) (byte, []byte, error)) FCRAdminServer
}

// Request sends a request to given addr with given key, msg type and data.
func Request(addr string, keyStr string, msgType byte, data []byte) (byte, []byte, error) {
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}
	key, err := hex.DecodeString(keyStr)
	if err != nil {
		return 0, nil, err
	}
	if len(key) != 32 {
		return 0, nil, fmt.Errorf("Wrong key size, expect 32, got: %v", len(key))
	}
	enc, err := encrypt(append([]byte{msgType}, data...), key)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", addr, bytes.NewReader(enc))
	if err != nil {
		return 0, nil, err
	}
	client := &http.Client{Timeout: 90 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	content, err := ioutil.ReadAll(r.Body)
	if closeErr := r.Body.Close(); closeErr != nil {
		return 0, nil, closeErr
	}
	if len(content) <= 1 {
		return 0, nil, fmt.Errorf("Received content with empty request %v", content)
	}
	plain, err := decrypt(content, key)
	if err != nil {
		return 0, nil, err
	}
	return plain[0], plain[1:], nil
}
