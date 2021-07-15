/*
Package e2e - end to end testing.
*/
package e2e

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
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/client/pkg/client"
	"github.com/wcgcyx/fc-retrieval/gateway-admin/pkg/gatewayadmin"
	"github.com/wcgcyx/fc-retrieval/itest/pkg/util"
	"github.com/wcgcyx/fc-retrieval/provider-admin/pkg/provideradmin"
)

var gatewayKeys = []string{
	/* gateway0 */ "5a7c858349ed16806e931b5fbb359e031529faf7218340d3b7e56bd58cfb97e5", // ID - 0523070e9e2dbd47c9c11c6ba861c801cd998d3002f245c3bbec6a0bf0d1fd49 // Before [GW25, GW9]		// After [GW25, GW9]
	/* gateway1 */ "c9dabac5dec3d1927fd6976011c8f5ef99485631f1099961e686529032da76da", // ID - 0e0282464798944a922f5a6bc8d0e026850b4529797925982af839bde1a43dbd // Before [GW25, GW9]		// After [GW25, GW9]
	/* gateway2 */ "3b139571edf7816ed090f29a4c2e15c2585728ea245fad2ff64ea6bd1fe8ddbe", // ID - 10e638412342152b9ff317a4fcfd57f215562fe050fe6836c7da3193eddcd29e // Before [GW26, GW10]		// After [GW26, GW10]
	/* gateway3 */ "9f1d730e7c0199107501cbe299eb63ae8d81e89715213a9b0913d66e439ce9d5", // ID - 160394959c1b6e04be181691c0049965e64cbe2cb06aa30677634ff0da72eba5 // Before [GW26, GW10]		// After [GW26, GW10]
	/* gateway4 */ "90c1a2dc57d19b58ef2b7fc757356f32c457ba36c13cd0cb99119aca9340881e", // ID - 22850133bd4acd610755550ad9d09bbf9794162fcf81c5ac60acd8f17ceefca4 // Before [GW28, GW12]		// After [GW28, GW12]
	/* gateway5 */ "33ab67450814288bec155ed1620c1e824f6e19d4ad483a0aa5bc2312ca55bfc5", // ID - 29ea6805c30728655f832e04fde46dc118f09c984e2b1f268cf34f8184d3c90e // Before [GW29, GW13]		// After [GW29, GW13]
	/* gateway6 */ "38b32ba7f9d7f521fdd00cef9961e244aad4968618b1943c4e70ec9e2bd69389", // ID - 338b7c64c252810dddca85f6014d6daa3a7957a0f974419a9e497f3c95505bcf // Before [GW30, GW14]		// After [GW30, GW14]
	/* gateway7 */ "5715b6386487c213c3314b9bbac11c53d78be8130d74d21412aeb5f6fdc7f4bd", // ID - 3846722e7739daba07067b3077c1a042f344c1aa47975a68e6e98227f47282cf // Before [GW31, GW15]		// After [GW31, GW15]
	/* gateway8 */ "2e934326ac30802b3e9c0c0b5adfb846d2eb4eb07092f0b8135dbb855872dc7f", // ID - 43daf05542473ef398c9c876cfc68029bf2719249b09b49a9d672cb06518a492 // Before [GW0, GW16]		// After [GW0, GW15] + GW32
	/* gateway9 */ "c241b786a28636ca22bb9476a2ff269febc7375452f85654677613bdaf2e66c1", // ID - 468ca79b8de14af50138cb6c264ba7bb227e3097ed7271ce82dede2b7f2da0a7 // Before [GW0, GW16]		// After [GW1, GW16] inc GW32
	/* gateway10 */ "03254c3c191e52558559c2c2324ef49bb689c87ca0ce02e1a7a72f6dbc922223", // ID - 52d0284f80058cad04e5cfb930acf93e1922bf19ba0fcacc62a7e7464055ceb7 // Before [GW2, GW18]		// After [GW3, GW18] inc GW32
	/* gateway11 */ "2fd67c6ffd03695cacb66af2a363954880b21112e02bb79516b8fb45e9821c37", // ID - 5c53dd80f68d319c09ec01505eff45d1cebc90aa0c0b55470d0d9980ee519564 // Before [GW3, GW19]		// After [GW4, GW19] inc GW32
	/* gateway12 */ "9708e98ce259c24515474f9ebd50f80c05354690238775371468a6caa6d83c6d", // ID - 6029da88e8c56eba13f8ce3295452a48571906baa80d2ef774468f506bd79c3f // Before [GW4, GW20]		// After [GW4, GW19] inc GW32
	/* gateway13 */ "1098d37e545472f3bf65ce72dcf69c615527817bc506aeddf0ae00fdde746c9c", // ID - 69eff4c36749fc7108a3f419e5e868a94d73d2fa2fe743e62b1f8b4309a10530 // Before [GW5, GW21]		// After [GW5, GW20] inc GW32
	/* gateway14 */ "dd547100866b46380543fa8bd20cf6c475a49fcf325904d1553d42c844cb9146", // ID - 7128499dd89ffcf278b622190bfb344eefd5fdf33ecc64ce1508197cb22419c8 // Before [GW6, GW22]		// After [GW6, GW21] inc GW32
	/* gateway15 */ "a741d27ec9d30ad8589f8c96d81fdce8108e07e591e6989ed6584daa4e1f2f12", // ID - 7728977ae23b0e0abd494cb33b5fda99bc4d1f66a18bf2eb9056c10ee499450a // Before [GW6, GW22]		// After [GW7, GW22] inc GW32
	/* gateway16 */ "2dacc97bf6d472ff286dfab972c545cbc73c475b0fd5e772d4e2d9f1d40286fb", // ID - 83cae55a61335add6368e63b79471191ba11b3fd43455729efeaf73b1c6de7fa // Before [GW8, GW24]		// After [GW9, GW24] inc GW32
	/* gateway17 */ "ee9ebd34579b551147a96e39812a63a1342011991ab7b6cceacdf128eb949313", // ID - 8eeadada16408ce527bd899dad5fd1650145de8cf337269027fe7c7424fd5991 // Before [GW10, GW26]		// After [GW10, GW25] inc GW32
	/* gateway18 */ "bc7c7c44cf5767908d15783f49883c2df0f7a2c7c5deeb1a89a9b0d8c1ffb2e4", // ID - 925498ad9d0f2aff6a8d8b30e304043eedfcba96c6cd1e64e05d2b12a070cc08 // Before [GW10, GW26]		// After [GW10, GW25] inc GW32
	/* gateway19 */ "9c4a15559755df2f45339f232ede123dfd2019012a2f93d498f18a30e72758a0", // ID - 99ca35e67778bb788ec120b9ddb72e24b641e71c8206f2346ba861d3d5ec8c18 // Before [GW11, GW27]		// After [GW11, GW26] inc GW32
	/* gateway20 */ "7917c714b7781ddf2e0bc49def4a9e9b052542ed72196bdc5df78a4a71b180a4", // ID - a3a2c4a1f7828c3913f6fb2bfa8c3ee311b9f02a9021f6b0f5515adbd9645dd0 // Before [GW12, GW28]		// After [GW13, GW28] inc GW32
	/* gateway21 */ "f9a39836a365aa5868f611a541ffa57f1b98c3bf3e0cadbd5272f50eea396763", // ID - aaa225c3f6cd07fe3c7163e87fcce95265ca02e0e4875ff7362aaa639bb3da99 // Before [GW13, GW29]		// After [GW14, GW29] inc GW32
	/* gateway22 */ "1a1dfacc2a27059dbea77b7a25f71fc0c48bc6fcd5fa55efe8cfadf30bca6118", // ID - b589d9a05f3ee5bef3cd5d0ea2960344d64934d0f8366dfb8d2ec4e97eb05bf9 // Before [GW14, GW30]		// After [GW15, GW30] inc GW32
	/* gateway23 */ "cd97f448fe5ca5e12b10b4bf24887882528af0494d1f3d043aaab843a2d2d6c4", // ID - bee42f182166c41bf1da487183bc6ed1f2ba482906e40e5d1517a0bd2313df01 // Before [GW16, GW0]		// After [GW16, GW31] + GW32
	/* gateway24 */ "2e3d92ec9c7cae61d42ba934ca7a13ee58a54f9dbe59d9226527ee946f0808e7", // ID - c10fc3c35c45ddff235a11ae2c41bcfd345f0414615ad618e149f167ccd94d8f // Before [GW16, GW0]		// After [GW16, GW0]
	/* gateway25 */ "a7164b946e12a907ffbd085385da146cd3a0adf40edba349d1c5bceea6af711e", // ID - ca34fcc6df3ca2cb31073c1ed246c8647276a2903e59fc2f30b6a20358ba710b // Before [GW17, GW1]		// After [GW17, GW1]
	/* gateway26 */ "2dd47d9bc2150583d8d97f1354b7884ecafed863691e9ecf0731895784147299", // ID - d2a03495541ede308f04348a6148c9f753d38201e569c84f06e8f97f673cc9ff // Before [GW18, GW2]		// After [GW18, GW2]
	/* gateway27 */ "288849ae970dce8af8da6fadc7b331bcbfd28705f00ad9ed72605653482d23f7", // ID - d9d97672451730fd8594a2841f3643f829096feee41d382b8b658e9a08d26a67 // Before [GW19, GW3]		// After [GW19, GW3]
	/* gateway28 */ "b6e1350a62be2dfe195119910af6faba7590cac3d15656a59ec7c85d86991932", // ID - e1201cd6f5aa071595d8235b609c2e26886b3ef33119e1c787adb30bdd115737 // Before [GW20, GW4]		// After [GW20, GW4]
	/* gateway29 */ "f22008f10c9b1f86b86a3d9ffc6a71154b15561ce5726ec68328c424365c3ed3", // ID - e9bfa910eed28ace8f5adac05a589283c27d44dcc691674113400c3120ddc876 // Before [GW21, GW5]		// After [GW21, GW5]
	/* gateway30 */ "fdf89b3b7cc26173edfea6614b6a696d89042f7051b3644d8bc469a72961c885", // ID - f7a419c082602154061ed19ee4c0d984963a61b514aa32479956441396e02f3a // Before [GW23, GW7]		// After [GW23, GW7]
	/* gateway31 */ "00b9a3b260e0ebd51f13e7adcc3c4aa7851ca35257cbc160c238c9355534e949", // ID - fe8fbb0e83f243fd7f184c24d7a008d1c89d5d3c137e9158977a78fb19e7869e // Before [GW23, GW7]		// After [GW23, GW7]
	// Gateway 32 is added later. It is between GW15 and GW16. Sort of like GW15.5
	/* gateway32 */ "54b81c2cc94a8a10be4820716ddcddf834b24d1359bdedbf5a5a78251175793b", // ID - 79f1dfc58999bc9a1a3cb9f6cc1b8b3109b6e21350cd85c4641ab9a64907f4b0 // Before []				// [GW7, GW22] inc GW32
}

var providerKeys = []string{
	/* provider0 */ "ea9a44d5aa53b4714efb7df4aed727ea0cf68b7ed18ac3d36ac2c90f262daf5f", // ID - 3f3bb8d3768a56b0d718e01f29a491dcdbf91e5fc7193e948689d001a22099b6
	/* provider1 */ "0d90700579ab17bcf579ccf904a9911ff5e6f4b9d5a450d1c1aef41e56736de0", // ID - 56651e4cf52c36b56498df851a52e6d95172f399d86b64d8ac3b69c573087f10
	/* provider2 */ "f2754a52c0fb15e3be023346ccab3919a7f0687d876356810398799410280d57", // ID - f79f39161ed74c86d27ac21d98728ffbfd8ddd7ea5a5c5dbbb411b47162d3494
}

const clientKey = "72dd0be8b35fac690d0e763ce13326d9512c81c664b2a0a143bfb87bde5fc195" // ID - 87cd9ced77cb602a83b80a883f4f52d14901279fefb0ff40e7816f960e083f66

const adminKey = "6465616135313132656636333864653962393132306366363537333664656465"

// test1.txt - QmZcKGSc63SnDfLoeiHcU9B1qZ53xBDEr7Y4a7LUSi9d5T
// test2.txt - QmTQ5MBuph9v4ggsWRK9rUBrUJqRoC6HG7MvTQusYGaBfy
// test3.txt - QmU83BhFDvqT2BR1yvm8EuexeYLvMoAnPXRvmF7bq1gj4x

// Testing scenario
// Step 1:
// Start 33 Gateways (GW0 - GW32)
// Start 3 providers (PVD0 - PVD2)
// Start 1 client, 1 gateway admin and 1 provider admin.
// Step 2:
// Initialise client, initialise gateway admin and provider admin.
// Initialise 32 gateways (GW0 - GW31) and 3 providers
// Force every gateway and provider to sync.
// Step 3:
// Client add one active gateway. (GW14)
// Step 4:
// PVD0 publishes offer for test1.txt
// PVD1 publishes offer for test2.txt
// Step 5:
// Standard discovery for test1.txt (Found 1 offer)
// Step 6:
// Standard discovery for test2.txt (Found 0 offer)
// Step 7:
// DHT discovery for test2.txt (Found 1 offer)
// Step 8:
// Initialise GW32.
// Force every gateway and provider to sync.
// Step 9:
// PVD2 publishes offer for test1.txt + test2.txt
// Step 10:
// Standard discovery for test1.txt (Found 1 offer)
// Client add one active gateway. (GW32)
// Standard discovery for text1.txt (Found 2 offers)

var lotusAPI string
var lotusToken string
var registerAPI string
var gateway32IP string
var gwAdmin *gatewayadmin.FilecoinRetrievalGatewayAdmin
var pvdAdmin *provideradmin.FilecoinRetrievalProviderAdmin
var fcrClient *client.FilecoinRetrievalClient

func TestInitialisation(t *testing.T) {
	var err error
	// Get lotusAPI
	lotusAPI = util.GetLotusAPI()
	// Get registerAPI
	registerAPI = util.GetRegisterAPI()
	// Get lotus token, super Acct
	token, acct := util.GetLotusToken()
	lotusToken = token
	// Topup all accounts
	util.Topup(lotusAPI, token, acct, append(append(gatewayKeys, providerKeys...), clientKey))

	// Test gateway initialisation
	// Get admin AP and IPs for all gateways
	ips := util.GetContainerInfo(false)
	gwAdmin = gatewayadmin.NewFilecoinRetrievalGatewayAdmin()
	for i, ip := range ips {
		if i == 32 {
			// Initialise gateway-32 later
			gateway32IP = ip
			break
		}
		err = gwAdmin.InitialiseGateway(
			fmt.Sprintf("%v:9010", ip),
			adminKey,
			9011,
			ips[i],
			gatewayKeys[i],
			lotusAPI,
			token,
			"",
			registerAPI,
			"",
			"au",
			fmt.Sprintf("gateway-%v", i))
		if !assert.Empty(t, err) {
			panic(fmt.Errorf("Error in initialisating gateway: %v", err.Error()))
		}
	}

	// Test provider initialisation
	// Get admin AP and IPs for all providers
	ips = util.GetContainerInfo(true)
	pvdAdmin = provideradmin.NewFilecoinRetrievalProviderAdmin()
	for i, ip := range ips {
		err = pvdAdmin.InitialiseProvider(
			fmt.Sprintf("%v:9010", ip),
			adminKey,
			9011,
			ips[i],
			providerKeys[i],
			lotusAPI,
			token,
			"",
			registerAPI,
			"",
			"au",
			fmt.Sprintf("provider-%v", i))
		if !assert.Empty(t, err) {
			panic(fmt.Errorf("Error in initialisating provider: %v", err.Error()))
		}
	}

	// Initialise client
	fcrClient, err = client.NewFilecoinRetrievalClient(
		clientKey,
		lotusAPI,
		token,
		"",
		registerAPI,
		"",
	)
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in initialisating client: %v", err.Error()))
	}
}

func TestSync(t *testing.T) {
	var err error
	// Sync provider
	ids, _, _ := pvdAdmin.ListProviders()
	for _, id := range ids {
		err = pvdAdmin.ForceSync(id)
		if !assert.Empty(t, err) {
			panic(fmt.Errorf("Error in foce syncing provider: %v", err.Error()))
		}
	}
	// Sync gateawy
	ids, _, _ = gwAdmin.ListGateways()
	for _, id := range ids {
		err = gwAdmin.ForceSync(id)
		if !assert.Empty(t, err) {
			panic(fmt.Errorf("Error in foce syncing gateway: %v", err.Error()))
		}
	}
}

func TestClientAddActiveGateway(t *testing.T) {
	err := fcrClient.AddActivePeer("7128499dd89ffcf278b622190bfb344eefd5fdf33ecc64ce1508197cb22419c8")
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in adding active gateway: %v", err.Error()))
	}
}

func TestPublishOffer(t *testing.T) {
	ok, msg, err := pvdAdmin.PublishOffer("ea9a44d5aa53b4714efb7df4aed727ea0cf68b7ed18ac3d36ac2c90f262daf5f", []string{"test1.txt"}, big.NewInt(1000000), time.Now().Add(time.Hour*24).Unix(), 24)
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in publishing offer: %v", err.Error()))
	}
	if !assert.True(t, ok) {
		panic(fmt.Errorf("Fail to publish offer: %v", msg))
	}
	ok, msg, err = pvdAdmin.PublishOffer("0d90700579ab17bcf579ccf904a9911ff5e6f4b9d5a450d1c1aef41e56736de0", []string{"test2.txt"}, big.NewInt(1000000), time.Now().Add(time.Hour*24).Unix(), 24)
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in publishing offer: %v", err.Error()))
	}
	if !assert.True(t, ok) {
		panic(fmt.Errorf("Fail to publish offer: %v", msg))
	}
}

func TestStandardDiscovery(t *testing.T) {
	res, err := fcrClient.StandardDiscovery("QmZcKGSc63SnDfLoeiHcU9B1qZ53xBDEr7Y4a7LUSi9d5T")
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in standard discovery: %v", err.Error()))
	}
	if !assert.Equal(t, 1, len(res)) {
		panic(fmt.Errorf("Should find 1 offer by standard discovery but find %v offers", len(res)))
	}
}

func TestDHTDiscovery(t *testing.T) {
	res, err := fcrClient.StandardDiscovery("QmTQ5MBuph9v4ggsWRK9rUBrUJqRoC6HG7MvTQusYGaBfy")
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in standard discovery: %v", err.Error()))
	}
	if !assert.Equal(t, 0, len(res)) {
		panic(fmt.Errorf("Should find 0 offer by standard discovery but find %v offers", len(res)))
	}
	res, err = fcrClient.DHTDiscovery("QmTQ5MBuph9v4ggsWRK9rUBrUJqRoC6HG7MvTQusYGaBfy", "7128499dd89ffcf278b622190bfb344eefd5fdf33ecc64ce1508197cb22419c8")
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in dht discovery: %v", err.Error()))
	}
	if !assert.Equal(t, 1, len(res)) {
		panic(fmt.Errorf("Should find 1 offer by dht discovery but find %v offers", len(res)))
	}
}

func TestNewGateway(t *testing.T) {
	err := gwAdmin.InitialiseGateway(
		fmt.Sprintf("%v:9010", gateway32IP),
		adminKey,
		9011,
		gateway32IP,
		gatewayKeys[32],
		lotusAPI,
		lotusToken,
		"",
		registerAPI,
		"",
		"au",
		"gateway-32")
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in initialisating gateway: %v", err.Error()))
	}
	// Sync provider
	ids, _, _ := pvdAdmin.ListProviders()
	for _, id := range ids {
		err = pvdAdmin.ForceSync(id)
		if !assert.Empty(t, err) {
			panic(fmt.Errorf("Error in foce syncing provider: %v", err.Error()))
		}
	}
	// Sync gateawy
	ids, _, _ = gwAdmin.ListGateways()
	for _, id := range ids {
		err = gwAdmin.ForceSync(id)
		if !assert.Empty(t, err) {
			panic(fmt.Errorf("Error in foce syncing gateway: %v", err.Error()))
		}
	}
}

func TestRingUpdate(t *testing.T) {
	ok, msg, err := pvdAdmin.PublishOffer("f2754a52c0fb15e3be023346ccab3919a7f0687d876356810398799410280d57", []string{"test1.txt", "test2.txt"}, big.NewInt(1000000), time.Now().Add(time.Hour*24).Unix(), 24)
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in publishing offer: %v", err.Error()))
	}
	if !assert.True(t, ok) {
		panic(fmt.Errorf("Fail to publish offer: %v", msg))
	}
	res, err := fcrClient.StandardDiscovery("QmZcKGSc63SnDfLoeiHcU9B1qZ53xBDEr7Y4a7LUSi9d5T")
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in standard discovery: %v", err.Error()))
	}
	if !assert.Equal(t, 1, len(res)) {
		panic(fmt.Errorf("Should find 1 offer by standard discovery but find %v offers", len(res)))
	}
	err = fcrClient.AddActivePeer("54b81c2cc94a8a10be4820716ddcddf834b24d1359bdedbf5a5a78251175793b")
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in adding active gateway: %v", err.Error()))
	}
	res, err = fcrClient.StandardDiscovery("QmZcKGSc63SnDfLoeiHcU9B1qZ53xBDEr7Y4a7LUSi9d5T")
	if !assert.Empty(t, err) {
		panic(fmt.Errorf("Error in standard discovery: %v", err.Error()))
	}
	if !assert.Equal(t, 2, len(res)) {
		panic(fmt.Errorf("Should find 2 offer by standard discovery but find %v offers", len(res)))
	}
}
