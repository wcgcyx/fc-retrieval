/*
Package fcrserver - provides an interface to do networking.
*/
package fcrserver

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
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
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

const (
	PrvKey1 = "308204a30201000282010100b7e77e0288455b0e00016010edb3b4a34cc4d962ea3a6c96c783c0b89ce39c5a82115127eae3fa43150e05a8c47bf0c5c86693a57d43f471f2d0a19de777b7c7a7a43480610193db56e519e16540a210adfd3310a09f7751a518a3dd2bda253bcb5982f757bd722897025428a4b74b5b189c8eff18c7d83d6bb4942501395c65f5c0f9299b15f222bf30be41220c1f1cc1dc20ea0223da9bda8df75f995d1d2b665ed113b22b9ae5ac6668ed155f163e57264f8b8dd7eec51e367f2c88f24f3005c0f7fc08e51a5977f7d80126beea3bf79130cbb94194b7a9c6896e7050c50239cef20fecb0d2d19395ff0bea07797b770f03f2e8e2a4fe599dba70d4b227d90203010001028201002b99f4e440cec0c1d6fa7c7e46fd1e4cc13cc2959316fafbdc9dbe2986f8e7ef057b79944f3a71f149a2a370d9f4d0a6f3d66e1704560234a9ef11025108af47e4d527a5705a6165d57a47e28a91025b9604bc00ab3463a3b5d2dbb6ea58b40f332d2bc1dcc98bb157ec336bd771a5aa1971b4ed82408f62309105b6a84da33c1a34a483a52f2a91c3cedc610d50a42124f3e9b42bb32e892ac153f9d483744962f5d7a9452e83122ca5533fc1cc35a96668d6f19941cb639916f0fe0b3e99117b2e876534e3ad163d29e8513fcbf528b5f7bec4bccc6ea2b24c8c7b4fb7d2ae9a61c3005d2e5d5ae0ed9b5e3973de35eac4c6d8789ee148c8ed66982f825a0902818100e59f0a0a5b96b0ceb59fa62e06f79b0227c30ee2301a448eeb74ee6bc1d4d65ef07211e752f6429a5a8f8fe7674cf78d06aad56a0fe9112616d456b6d8fbf00f7b1f867be480781dadd7c072e8b6cec284841ee031a0b29c9c602e17cf4a3fb8ff560379743298036d1233701518ea5ae554aa36e61e9032e5704af153f0c8bb02818100cd07f4d9d301e96e6cdf48e8ffb3d87695aa3e6a30c02a970b97b8a9e40860f2c33ddcf170d74c1742d255caa5897232b6b1a08317442a9b1987e2d997a23557e5588b91cca9676234f39db708cb8f49af25a3c15bea5fd474739725caa84e4873f779ce4127d275de68b450970c27b232d965acdca6dad6afd10c27fb6cc27b02818011d908c8c151b7307a018cc32b1b77daf5083e51ea774038f3a84517ef1b0206a31ddab2664a69e6e17f232a5367321eae13fd3e9f39f871437901bc78a52c85a7864dc7b77d1cd901b831673d1b687aca1e12e04e3b3566e2e8beec6eda5095aa931ef603c822f4b137a6f3e14fec776037f27b0debf63d5e8419ef241d251d02818100babf8d80bdd616f5728aea10f79eab02500df1adad5bcb2f2aeaf5d320959520693f36b85f6c6aad213b0dd37775baa3808e47c23f75e24cc5336527861ac3f59c3b4b5cf08a3855561fb33e9cef34430c19ff8ec616b354830129e1cd36019fb2a8edb434da7db2c8729c126f922db1fce8d0d8635e43239a9e9130f5ac39730281801e4a6185a28e7aa499a6fbaf286b7d4b333797f1baf9c613e6b96ed0b4883dac3bc56074ba1db0b0400beac47c66e56afc2c57a89ba725fb3443de8309d7281c6e197dba956c621cf34e2bc7f892f861daa1ffc43e472ebbdd7c471d52f6bbdf12b6ceb6e9d6251484dbaf03ee17ffccec3a2399d84417a6a9b4114eee662eee"
	PrvKey2 = "308204a302010002820101009a077b0668f3a98da5d2b51a73fc64a682b14ab51c96dc67c56d43e623ca9ceb383b87da36895b6e1e9d94af310d62cc49ba200c809b3cb437032a548abecf6a5f1e827b34345c915e55d34c5f1933976f9e7b7d6a4ba43b6a49d73d19231301f10ae75c9d8a6c2d1b7eec10a563280dab09fdc0eea3351b4aa00a9c9db60222aee5b4381f220f79e59ca29acd13bdd9a0803407374815d7fa8ff6d37d13bb8ce6d7ddacde2144ee123c85393cc55a651d110abe59a8c9b35f74989fed0e3b56ecc09697a708794605eb873541dfcbf421846014640f3d8f226a23c45c919f5198b0d073deb1893fc371a2cf3eb8be811dd3ee265dfa59d53771c7441756e37d0203010001028201005871592bda11a757055352d828a75127e73d63f750be333a86bb71d470d2c37db0e145e57f912965b6c0a7025d792134ca54cc58417461cbdd16bd34a4226238e2fb42d2f9abe3473952b0ac56a2c2e3fe9c92adf5de0f246aa891a5ac8c5e3aac2ca5a2a1773d1c3d80888e1a59304380e590c63a808e5ae863b31430deb4a44f8bcc79ae61239a28153e5abe01a297dce8182147f8736dad7c6deda1fe194847c71885ad3695147610b2a88a3c623e2b39e41e84a1c775e81bb8988d1bb69d5845b364983ebca3e7e390202458da4a2319859bbeeca04171eb1ad0cd2b7dc51eafd429964bd96e08f7cc30560f8b9b5ec1b960973e50d60035a65f1fe387a102818100c1f9f097d507eb64b16ca7aaa53cd473a208cd38ef5b4e6a40d7ec17d68d565c0598ea8d89071be3c7104a747fafc0eb5eff5306bee5ebe267dfc8ae25cc8edf11c97ce3c278a694f5c6cf0c81715d12e9e2889a35d0dccf8a3e17df62c8e82a5bbe6334c2688eecf5280460651ff05b8e3e7fffde53f6731b6c09a57905df0502818100cb47a4ddf03772ae79b25e541f2f03137374c9fd5284d3f56c6c71fcae123bde169081f529bd1d72a0f68cbc2caa84c1c8310770086c051e6850914d2d667ab2a059e33ad82da333f9791501d7b97c36a8bb95fbe1928acdb058884e79a7fcb1f76f6e5bffd19a34bf7558cc2f216e21926801a299891f20a88f244a60566c1902818050b531c9bab564d7ac8acce84f8013d558e1d8a18bd5adb6bfec172b83f5a2acff1734e056d7425f6f7ff3baad35ef4aff67b49fe5e5bc53a36c950f0063303ed823c176f27f48b049e2c25b2db0814d514b141335b90566c4da390c95098aafb5246e1a9198f77ed8322240095354aa8370b5c93c342b2291924e212f4da611028181009c6beee3913b399624b32a7ed4d81a27d78a20fc3b895688ddfbbce2d117dad594cb7215331f010ff9e87e77366fa8646d25bd316a69a4aeb75a77d4c980b81dc7e223465e9f0f9ca8f59142afbb5d67ba034ef059ada7fd8b1b35181de9343bc5c90b44e3df6827fac3d3a69b05c07738efab82715ee08302f1d2dd20b09fd102818033f6a94f23312cb7f69fa4964fc464e9529d3018e6bdea59ce076cbfb5f7c2736c8e53559753a3a39f168836e788e5c8350a638a5e97dcd5c0515a42bec7a01ff9ede829321167643c0f8b74888b558b155cb4e55fa967a99b755cadab8a9292209d87fcdbad6bf17f23907b7ecc77d54ee89f8cdbd6c7ef2b71f7ea5ea46587"
)

func testHandler(reader FCRServerReader, writer FCRServerWriter, req *fcrmessages.FCRMessage) error {
	if req.GetMessageType() != 1 {
		panic("wrong msg")
	}
	body := req.GetMessageBody()
	if len(body) != 3 {
		panic("wrong msg")
	}
	if body[0] != 1 || body[1] != 0 || body[2] != 1 {
		panic("wrong msg")
	}
	resp := fcrmessages.CreateFCRMessage(0, []byte{0, 1, 0})
	err := writer.Write(resp, time.Minute)
	if err != nil {
		panic(err)
	}
	return nil
}

func testHandlerV1(reader FCRServerReader, writer FCRServerWriter, req *fcrmessages.FCRMessage) error {
	return errors.New("Test error")
}

func testRequester(reader FCRServerReader, writer FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	req := fcrmessages.CreateFCRMessage(1, []byte{1, 0, 1})
	err := writer.Write(req, time.Minute)
	if err != nil {
		panic(err)
	}
	resp, err := reader.Read(time.Minute)
	if err != nil {
		panic(err)
	}
	if resp.GetMessageType() != 0 {
		panic("wrong msg")
	}
	body := resp.GetMessageBody()
	if len(body) != 3 {
		panic("wrong msg")
	}
	if body[0] != 0 || body[1] != 1 || body[2] != 0 {
		panic("wrong msg")
	}
	return resp, nil
}

func testRequesterV2(reader FCRServerReader, writer FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	req := fcrmessages.CreateFCRMessage(10, []byte{1, 0, 1})
	err := writer.Write(req, time.Minute)
	if err != nil {
		panic(err)
	}
	return nil, nil
}

func testRequesterV3(reader FCRServerReader, writer FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	req := fcrmessages.CreateFCRMessage(3, []byte{1, 0, 1})
	err := writer.Write(req, time.Minute)
	if !IsTimeoutError(err) && err != nil {
		panic(err)
	}
	return nil, nil
}

func TestChat(t *testing.T) {
	portSender := freePort()
	sender := NewFCRServerImplV1(PrvKey1, uint(portSender), time.Minute)
	sender.AddRequester(1, testRequester)
	sender.AddRequester(2, testRequesterV2)
	sender.AddRequester(3, testRequesterV3)
	sender.Shutdown()
	err := sender.Start()
	assert.Empty(t, err)
	defer sender.Shutdown()
	err = sender.Start()
	assert.NotEmpty(t, err)

	portReceiver := freePort()
	receiver := NewFCRServerImplV1(PrvKey2, uint(portReceiver), time.Minute)
	receiver.AddHandler(1, testHandler)
	receiver.AddHandler(2, testHandler)
	receiver.AddHandler(3, testHandlerV1)
	receiver.Shutdown()
	err = receiver.Start()
	assert.Empty(t, err)
	defer receiver.Shutdown()
	err = receiver.Start()
	assert.NotEmpty(t, err)

	receiverAddr, err := GetMultiAddr(PrvKey2, "127.0.0.1", uint(portReceiver))
	assert.Empty(t, err)

	sender.Request(receiverAddr, 1, 0)
	time.Sleep(1 * time.Second)

	sender.Request(receiverAddr, 2, 0)
	time.Sleep(1 * time.Second)

	sender.Request(receiverAddr, 3, 0)
	time.Sleep(1 * time.Second)

	// Test invalid message
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h, err := libp2p.New(
		ctx,
	)
	assert.Empty(t, err)
	// Get multiaddr
	maddr, err := multiaddr.NewMultiaddr(receiverAddr)
	assert.Empty(t, err)
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	assert.Empty(t, err)
	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	conn, err := h.NewStream(context.Background(), info.ID, "/fc-retrieval/0.0.1")
	assert.Empty(t, err)
	defer conn.Close()
	conn.Write([]byte{192, 192, 192})
	time.Sleep(1 * time.Second)
}
