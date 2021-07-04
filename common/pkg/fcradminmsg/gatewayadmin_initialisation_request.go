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

type gatewayAdminInitialisationRequestJson struct {
	P2PPrvKey         string `json:"p2p_private_key"`
	P2PPort           int    `json:"p2p_port"`
	NetworkAddr       string `json:"network_addr"`
	RootPrvKey        string `json:"root_private_key"`
	LotusAPIAddr      string `json:"lotus_api_addr"`
	LotusAuthToken    string `json:"lotus_auth_token"`
	RegisterPrvKey    string `json:"register_private_key"`
	RegisterAPIAddr   string `json:"register_api_addr"`
	RegisterAuthToken string `json:"register_auth_token"`
	RegionCode        string `json:"region_code"`
}

func EncodeGatewayAdminInitialisationRequest(
	p2pPrvKey string,
	p2pPort int,
	networkAddr string,
	rootPrvKey string,
	lotusAPIAddr string,
	lotusAuthToken string,
	registerPrvKey string,
	registerAPIAddr string,
	registerAuthToken string,
	regionCode string,
) ([]byte, error) {
	return json.Marshal(&gatewayAdminInitialisationRequestJson{
		P2PPrvKey:         p2pPrvKey,
		P2PPort:           p2pPort,
		NetworkAddr:       networkAddr,
		RootPrvKey:        rootPrvKey,
		LotusAPIAddr:      lotusAPIAddr,
		LotusAuthToken:    lotusAuthToken,
		RegisterPrvKey:    registerPrvKey,
		RegisterAPIAddr:   registerAPIAddr,
		RegisterAuthToken: registerAuthToken,
		RegionCode:        regionCode,
	})
}

func DecodeGatewayAdminInitialisationRequest(data []byte) (
	string,
	int,
	string,
	string,
	string,
	string,
	string,
	string,
	string,
	string,
	error,
) {
	msg := gatewayAdminInitialisationRequestJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return "", 0, "", "", "", "", "", "", "", "", err
	}
	return msg.P2PPrvKey, msg.P2PPort, msg.NetworkAddr, msg.RootPrvKey, msg.LotusAPIAddr, msg.LotusAPIAddr, msg.RegisterPrvKey, msg.RegisterAPIAddr, msg.RegisterAuthToken, msg.RegionCode, nil
}
