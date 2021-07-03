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
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}

func testHandler(data []byte) (byte, []byte, error) {
	if data[0] != 1 || data[1] != 2 || data[2] != 3 {
		panic("Wrong data")
	}
	return 11, []byte{4, 5, 6}, nil
}

func testWrongHandler(data []byte) (byte, []byte, error) {
	return 0, nil, errors.New("Test error")
}

func TestChat(t *testing.T) {
	// Test New Server
	portReceiver := freePort()
	addr := fmt.Sprintf("localhost:%v", portReceiver)
	receiver := NewFCRAdminServerImplV1(addr, "pppp")
	err := receiver.Start()
	assert.NotEmpty(t, err)

	receiver = NewFCRAdminServerImplV1(addr, "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a2901")
	err = receiver.Start()
	assert.NotEmpty(t, err)

	receiver = NewFCRAdminServerImplV1(addr, "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29")
	receiver.AddHandler(10, testHandler)
	receiver.AddHandler(12, testWrongHandler)
	receiver.Shutdown()
	err = receiver.Start()
	assert.Empty(t, err)
	receiver.AddHandler(10, testHandler)
	defer receiver.Shutdown()
	err = receiver.Start()
	assert.NotEmpty(t, err)

	respType, resp, err := Request(addr, "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29", 10, []byte{1, 2, 3})
	assert.Equal(t, byte(11), respType)
	assert.Equal(t, []byte{4, 5, 6}, resp)
	assert.Empty(t, err)

	client := &http.Client{Timeout: 90 * time.Second}
	req, err := http.NewRequest("POST", "http://"+addr, bytes.NewReader([]byte{0}))
	assert.Empty(t, err)
	_, err = client.Do(req)
	assert.Empty(t, err)

	_, _, err = Request(addr, "ppp", 10, []byte{1, 2, 3})
	assert.NotEmpty(t, err)

	_, _, err = Request(addr, "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a2910", 10, []byte{1, 2, 3})
	assert.NotEmpty(t, err)

	_, _, err = Request(addr, "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29", 11, []byte{1, 2, 3})
	assert.NotEmpty(t, err)

	_, _, err = Request(addr, "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a2a", 10, []byte{1, 2, 3})
	assert.NotEmpty(t, err)

	_, _, err = Request(addr, "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29", 12, []byte{1, 2, 3})
	assert.NotEmpty(t, err)
}
