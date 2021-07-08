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

// initialisationRequestJson represents the request to initialise.
type initialisationRequestJson struct {
	P2PPrivKey        string `json:"p2p_private_key"`
	P2PPort           int    `json:"p2p_port"`
	NetworkAddr       string `json:"network_addr"`
	RootPrivKey       string `json:"root_private_key"`
	LotusAPIAddr      string `json:"lotus_api_addr"`
	LotusAuthToken    string `json:"lotus_auth_token"`
	RegisterPrivKey   string `json:"register_private_key"`
	RegisterAPIAddr   string `json:"register_api_addr"`
	RegisterAuthToken string `json:"register_auth_token"`
	RegionCode        string `json:"region_code"`
}

// EncodeInitialisationRequest is used to get the byte array of initialisationRequestJson
func EncodeInitialisationRequest(
	p2pPrivKey string,
	p2pPort int,
	networkAddr string,
	rootPrivKey string,
	lotusAPIAddr string,
	lotusAuthToken string,
	registerPrivKey string,
	registerAPIAddr string,
	registerAuthToken string,
	regionCode string,
) ([]byte, error) {
	return json.Marshal(&initialisationRequestJson{
		P2PPrivKey:        p2pPrivKey,
		P2PPort:           p2pPort,
		NetworkAddr:       networkAddr,
		RootPrivKey:       rootPrivKey,
		LotusAPIAddr:      lotusAPIAddr,
		LotusAuthToken:    lotusAuthToken,
		RegisterPrivKey:   registerPrivKey,
		RegisterAPIAddr:   registerAPIAddr,
		RegisterAuthToken: registerAuthToken,
		RegionCode:        regionCode,
	})
}

// DecodeInitialisationRequest is used to get the fields from byte array of initialisationRequestJson
func DecodeInitialisationRequest(data []byte) (
	string, // p2p private key
	int, // p2p port
	string, // network addr
	string, // root priv key
	string, // lotus api addr
	string, // lotus auth token
	string, // register private key
	string, // register api addr
	string, // register auth token
	string, // region code
	error, // error
) {
	msg := initialisationRequestJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return "", 0, "", "", "", "", "", "", "", "", err
	}
	return msg.P2PPrivKey, msg.P2PPort, msg.NetworkAddr, msg.RootPrivKey, msg.LotusAPIAddr, msg.LotusAuthToken, msg.RegisterPrivKey, msg.RegisterAPIAddr, msg.RegisterAuthToken, msg.RegionCode, nil
}
