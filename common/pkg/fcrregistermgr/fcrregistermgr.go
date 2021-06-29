/*
Package fcrregistermgr - register manager handles the interaction with the register.
*/
package fcrregistermgr

import "github.com/wcgcyx/fc-retrieval/common/pkg/register"

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

// FCRRegisterMgr represents the manager that interacts with the register.
type FCRRegisterMgr interface {
	// GetHeight gets the current height of the register.
	GetHeight() (uint64, error)

	// GetMaxPage gets the maximum page of the register at given height.
	GetMaxPage(height uint64) (uint64, error)

	// RegisterGateway registers a given gateway.
	RegisterGateway(id string, gwInfo *register.GatewayRegisteredInfo) error

	// UpdateGateway updates the given gateway's register information (usually for rolling msg signing key).
	UpdateGateway(id string, gwInfo *register.GatewayRegisteredInfo) error

	// RequestDeregisterGateway requests the deregistration of a given gateway. You need to be the owner.
	RequestDeregisterGateway(id string) error

	// DeregisterGateway removes the gateway entry. It will only be successful 24 hours (5760 blocks for rinkeby) after a request is sent.
	DeregisterGateway(id string) error

	// RegisterProvider registers a given provider.
	RegisterProvider(id string, pvdInfo *register.ProviderRegisteredInfo) error

	// UpdateProvider updates the given provider's register information (usually for rolling msg signing key).
	UpdateProvider(id string, pvdInfo *register.ProviderRegisteredInfo) error

	// RequestDeregisterProvider requests the deregistration of a given provider. You need to be the owner.
	RequestDeregisterProvider(id string) error

	// DeregisterProvider removes the provider entry. It will only be successful 24 hours (5760 blocks for rinkeby) after a request is sent.
	DeregisterProvider(id string) error

	// GetAllRegisteredGateway gets the registered gateways' information at given page at a given height.
	GetAllRegisteredGateway(height uint64, page uint64) ([]register.GatewayRegisteredInfo, error)

	// GetAllRegisteredProvider gets the registered providers' information at given page at a given height.
	GetAllRegisteredProvider(height uint64, page uint64) ([]register.ProviderRegisteredInfo, error)

	// GetRegisteredGatewayByID gets the gateway's information by a given ID.
	GetRegisteredGatewayByID(id string) (*register.GatewayRegisteredInfo, error)

	// GetRegisteredProviderByID gets the provider's information by a given ID.
	GetRegisteredProviderByID(id string) (*register.ProviderRegisteredInfo, error)

	// GetRegisteredGatewaysByRegion gets the registered gateways' information at given page at a given height with a given region code.
	GetRegisteredGatewaysByRegion(height uint64, region string, page uint64) ([]register.GatewayRegisteredInfo, error)

	// GetRegisteredProvidersByRegion gets the registered providers' information at given page at a given height with a given region code.
	GetRegisteredProvidersByRegion(height uint64, region string, page uint64) ([]register.ProviderRegisteredInfo, error)
}
