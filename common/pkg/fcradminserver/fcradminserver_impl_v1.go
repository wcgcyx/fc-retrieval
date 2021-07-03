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
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
)

// FCRAdminServerImplV1 implements the FCRAdminServer
type FCRAdminServerImplV1 struct {
	// 32 bytes hex string
	keyStr     string
	key        []byte
	listenAddr string
	start      bool

	server *http.Server

	// handlers for different message type
	handlers map[byte]func(data []byte) (byte, []byte, error)
}

func NewFCRAdminServerImplV1(listenAddr string, keyStr string) FCRAdminServer {
	return &FCRAdminServerImplV1{
		keyStr:     keyStr,
		listenAddr: listenAddr,
		start:      false,
		handlers:   make(map[byte]func(data []byte) (byte, []byte, error)),
	}
}

func (s *FCRAdminServerImplV1) Start() error {
	if s.start {
		return errors.New("Admin server has been started already")
	}
	key, err := hex.DecodeString(s.keyStr)
	if err != nil {
		return err
	}
	if len(key) != 32 {
		return fmt.Errorf("Wrong key size, expect 32, got: %v", len(key))
	}
	s.key = key
	s.server = &http.Server{
		Addr:           s.listenAddr,
		Handler:        s,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		s.server.ListenAndServe()
	}()
	s.start = true
	return nil
}

func (s *FCRAdminServerImplV1) Shutdown() {
	if !s.start {
		return
	}
	s.server.Shutdown(context.Background())
	s.start = false
}

func (s *FCRAdminServerImplV1) AddHandler(msgType byte, handler func(data []byte) (byte, []byte, error)) FCRAdminServer {
	if s.start {
		return s
	}
	s.handlers[msgType] = handler
	return s
}

func (s *FCRAdminServerImplV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if closeErr := r.Body.Close(); closeErr != nil {
		logging.Warn("HTTP can't close request body")
	}
	if err != nil {
		logging.Error("Error reading request %v", err.Error())
		WriteError(w, 400, fmt.Sprintf("Invalid Request: %v.", err.Error()))
		return
	}
	if len(content) <= 1 {
		logging.Error("Received content with empty request %v", content)
		WriteError(w, 400, "Content body is empty")
		return
	}
	plain, err := decrypt(content, s.key)
	if err != nil {
		logging.Error("Request fails to verify")
		WriteError(w, 400, "Request fails to verify")
		return
	}
	msgType := plain[0]
	msgData := plain[1:]
	handler, ok := s.handlers[msgType]
	if !ok {
		logging.Error("Unsupported message type: %v", msgType)
		WriteError(w, 400, fmt.Sprintf("Unsupported method: %v", msgType))
		return
	}
	respType, respData, err := handler(msgData)
	if err != nil {
		logging.Error("Error handling request: %v", err.Error())
		WriteError(w, 400, fmt.Sprintf("Error handling request: %v", err.Error()))
		return
	}
	respEnc, err := encrypt(append([]byte{respType}, respData...), s.key)
	if err != nil {
		logging.Error("Internal Error in encryption: %v", err.Error())
		WriteError(w, 500, "Internal Error")
		return
	}
	w.WriteHeader(200)
	_, err = w.Write(respEnc)
	if err != nil {
		logging.Error("Error responding to client: %v", err.Error())
	}
}

func WriteError(w http.ResponseWriter, header int, msg string) {
	w.WriteHeader(header)
	resp, err := json.Marshal(map[string]string{"Error": msg})
	if err == nil {
		w.Write(resp)
	}
}

func encrypt(plain []byte, key []byte) ([]byte, error) {
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	enc := aesGCM.Seal(nonce, nonce, plain, nil)
	return enc, nil
}

func decrypt(enc []byte, key []byte) ([]byte, error) {
	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// Decrypt the data
	plain, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}

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
