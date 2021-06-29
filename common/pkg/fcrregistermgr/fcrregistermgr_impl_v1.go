/*
Package fcrregistermgr - register manager handles the interaction with the register.
*/
package fcrregistermgr

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/wcgcyx/fc-retrieval/common/pkg/register"
)

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

// FCRRegisterMgrImplV1 implements FCRRegisterMgr, it interacts with a mocked register (http server).
type FCRRegisterMgrImplV1 struct {
	registerAPI string
	client      *http.Client
	lock        sync.RWMutex
}

func NewFCRRegisterMgrImplV1(registerAPI string, client *http.Client) FCRRegisterMgr {
	return &FCRRegisterMgrImplV1{registerAPI: registerAPI, client: client, lock: sync.RWMutex{}}
}

func (mgr *FCRRegisterMgrImplV1) GetHeight() (uint64, error) {
	return 0, nil
}

func (mgr *FCRRegisterMgrImplV1) GetMaxPage(height uint64) (uint64, error) {
	return 0, nil
}

func (mgr *FCRRegisterMgrImplV1) RegisterGateway(id string, gwInfo *register.GatewayRegisteredInfo) error {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	url := mgr.registerAPI + "/registers/gateway"
	return SendJSON(url, mgr.client, gwInfo)
}

func (mgr *FCRRegisterMgrImplV1) UpdateGateway(id string, gwInfo *register.GatewayRegisteredInfo) error {
	// Not implemented
	return errors.New("No implementation")
}

func (mgr *FCRRegisterMgrImplV1) RequestDeregisterGateway(id string) error {
	// Not implemented
	return errors.New("No implementation")
}

func (mgr *FCRRegisterMgrImplV1) DeregisterGateway(id string) error {
	// Not implemented
	return errors.New("No implementation")
}

func (mgr *FCRRegisterMgrImplV1) RegisterProvider(id string, pvdInfo *register.ProviderRegisteredInfo) error {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	url := mgr.registerAPI + "/registers/provider"
	return SendJSON(url, mgr.client, pvdInfo)
}

func (mgr *FCRRegisterMgrImplV1) UpdateProvider(id string, gwInfo *register.ProviderRegisteredInfo) error {
	// Not implemented
	return errors.New("No implementation")
}

func (mgr *FCRRegisterMgrImplV1) RequestDeregisterProvider(id string) error {
	// Not implemented
	return errors.New("No implementation")
}

func (mgr *FCRRegisterMgrImplV1) DeregisterProvider(id string) error {
	// Not implemented
	return errors.New("No implementation")
}

func (mgr *FCRRegisterMgrImplV1) GetAllRegisteredGateway(height uint64, page uint64) ([]register.GatewayRegisteredInfo, error) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	url := mgr.registerAPI + "/registers/gateway"
	var gateways []register.GatewayRegisteredInfo
	err := GetJSON(url, mgr.client, &gateways)
	if err != nil {
		return gateways, err
	}
	return gateways, nil
}

func (mgr *FCRRegisterMgrImplV1) GetAllRegisteredProvider(height uint64, page uint64) ([]register.ProviderRegisteredInfo, error) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	url := mgr.registerAPI + "/registers/provider"
	var providers []register.ProviderRegisteredInfo
	err := GetJSON(url, mgr.client, &providers)
	if err != nil {
		return providers, err
	}
	return providers, nil
}

func (mgr *FCRRegisterMgrImplV1) GetRegisteredGatewayByID(id string) (*register.GatewayRegisteredInfo, error) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	url := mgr.registerAPI + "/registers/gateway/" + id
	gateway := register.GatewayRegisteredInfo{}
	err := GetJSON(url, mgr.client, &gateway)
	if err != nil {
		return &gateway, err
	}
	return &gateway, nil
}

func (mgr *FCRRegisterMgrImplV1) GetRegisteredProviderByID(id string) (*register.ProviderRegisteredInfo, error) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	url := mgr.registerAPI + "/registers/provider/" + id
	provider := register.ProviderRegisteredInfo{}
	err := GetJSON(url, mgr.client, &provider)
	if err != nil {
		return &provider, err
	}
	return &provider, nil
}

func (mgr *FCRRegisterMgrImplV1) GetRegisteredGatewaysByRegion(height uint64, region string, page uint64) ([]register.GatewayRegisteredInfo, error) {
	// Not implemented
	return nil, errors.New("No implementation")
}

func (mgr *FCRRegisterMgrImplV1) GetRegisteredProvidersByRegion(height uint64, region string, page uint64) ([]register.ProviderRegisteredInfo, error) {
	// Not implemented
	return nil, errors.New("No implementation")
}

// GetJSON request Get JSON
func GetJSON(url string, client *http.Client, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	if decodeErr := json.NewDecoder(r.Body).Decode(target); decodeErr != nil {
		return decodeErr
	}
	if closeErr := r.Body.Close(); closeErr != nil {
		return closeErr
	}
	return nil
}

// SendJSON request Send JSON
func SendJSON(url string, client *http.Client, data interface{}) error {
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if req == nil {
		return errors.New("SendJSON error, can't create request")
	}
	req.Header.Set("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	return nil
}
