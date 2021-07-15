/*
Package util - common functions used in end-to-end and integration testing. Allowing to start different types of
Retrieval network nodes for testing.
*/
package util

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
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/lotus/chain/types"
	cid2 "github.com/ipfs/go-cid"
)

func GetLotusAPI() string {
	// If running from container
	env := os.Getenv("DEV_TEST")
	if env == "" {
		cmd := exec.Command("docker", "ps", "--filter", "ancestor=wcgcyx/lotusfull", "--format", "{{.ID}}")
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		id := string(stdout[:len(stdout)-1])
		cmd = exec.Command("docker", "inspect", id, "--format", "{{.NetworkSettings.Networks.shared.IPAddress}}")
		stdout, err = cmd.Output()
		if err != nil {
			panic(err)
		}
		ip := string(stdout[:len(stdout)-1])
		return fmt.Sprintf("http://%v:1234/rpc/v0", ip)
	} else {
		vars := strings.Split(env, ";")
		return vars[0]
	}
}

func GetRegisterAPI() string {
	// If running from container
	env := os.Getenv("DEV_TEST")
	if env == "" {
		cmd := exec.Command("docker", "ps", "--filter", "ancestor=wcgcyx/fc-retrieval/register", "--format", "{{.ID}}")
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		id := string(stdout[:len(stdout)-1])
		cmd = exec.Command("docker", "inspect", id, "--format", "{{.NetworkSettings.Networks.shared.IPAddress}}")
		stdout, err = cmd.Output()
		if err != nil {
			panic(err)
		}
		ip := string(stdout[:len(stdout)-1])
		return fmt.Sprintf("%v:9020", ip)
	} else {
		vars := strings.Split(env, ";")
		return vars[1]
	}
}

func GetLotusToken() (string, string) {
	env := os.Getenv("DEV_TEST")
	if env == "" {
		cmd := exec.Command("docker", "ps", "--filter", "ancestor=wcgcyx/lotusfull", "--format", "{{.ID}}")
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		id := string(stdout[:len(stdout)-1])
		cmd = exec.Command("docker", "exec", id, "bash", "-c", "cd ~/.lotus; cat token")
		stdout, err = cmd.Output()
		if err != nil {
			panic(err)
		}
		token := string(stdout)

		cmd = exec.Command("docker", "exec", id, "bash", "-c", "./lotus wallet default")
		stdout, err = cmd.Output()
		if err != nil {
			panic(err)
		}
		acct := string(stdout[:len(stdout)-1])
		return token, acct
	} else {
		vars := strings.Split(env, ";")
		return vars[2], vars[3]
	}
}

func GetContainerInfo(pvd bool) []string {
	env := os.Getenv("DEV_TEST")
	if env == "" {
		var image string
		ip := make([]string, 0)
		if pvd {
			image = "ancestor=wcgcyx/fc-retrieval/provider"
		} else {
			image = "ancestor=wcgcyx/fc-retrieval/gateway"
		}
		cmd := exec.Command("docker", "ps", "--filter", image, "--format", "{{.Names}}")
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		nodes := strings.Split(string(stdout[:len(stdout)-1]), "\n")
		// Get ips
		for _, node := range nodes {
			cmd = exec.Command("docker", "inspect", node, "--format", "{{.NetworkSettings.Networks.shared.IPAddress}}")
			stdout, err = cmd.Output()
			if err != nil {
				panic(err)
			}
			ip = append(ip, string(stdout[:len(stdout)-1]))
		}
		return ip
	} else {
		vars := strings.Split(env, ";")
		if pvd {
			return vars[37:]
		} else {
			return vars[33:37]
		}
	}
}

func Topup(lotusAPI string, token string, superAcct string, privKeyStrs []string) {
	// Get API
	var api v0api.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + token}}
	closer, err := jsonrpc.NewMergeClient(context.Background(), lotusAPI, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		panic(err)
	}
	defer closer()

	mainAddress, err := address.NewFromString(superAcct)
	if err != nil {
		panic(err)
	}

	cids := make([]cid2.Cid, 0)
	// Send messages
	for _, privKeyStr := range privKeyStrs {
		privKey, err := hex.DecodeString(privKeyStr)
		if err != nil {
			panic(err)
		}
		pubKey := crypto.PublicKey(privKey)
		address1, err := address.NewSecp256k1Address(pubKey)
		if err != nil {
			panic(err)
		}
		// Get amount
		amt, err := types.ParseFIL("100")
		if err != nil {
			panic(err)
		}
		msg := &types.Message{
			To:     address1,
			From:   mainAddress,
			Value:  types.BigInt(amt),
			Method: 0,
		}
		signedMsg, err := fillMsg(mainAddress, &api, msg)
		if err != nil {
			panic(err)
		}
		// Send request to lotus
		cid, err := api.MpoolPush(context.Background(), signedMsg)
		if err != nil {
			panic(err)
		}
		cids = append(cids, cid)
	}

	// Finally check receipts
	for _, cid := range cids {
		receipt := waitReceipt(&cid, &api)
		if receipt.ExitCode != 0 {
			panic(errors.New("Transaction fail to execute"))
		}
	}
}

// fillMsg will fill the gas and sign a given message
func fillMsg(fromAddress address.Address, api *v0api.FullNodeStruct, msg *types.Message) (*types.SignedMessage, error) {
	// Get nonce
	nonce, err := api.MpoolGetNonce(context.Background(), msg.From)
	if err != nil {
		return nil, err
	}
	msg.Nonce = nonce

	// Calculate gas
	limit, err := api.GasEstimateGasLimit(context.Background(), msg, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasLimit = int64(float64(limit) * 1.25)

	premium, err := api.GasEstimateGasPremium(context.Background(), 10, msg.From, msg.GasLimit, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasPremium = premium

	feeCap, err := api.GasEstimateFeeCap(context.Background(), msg, 20, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	msg.GasFeeCap = feeCap

	// Sign message
	return api.WalletSignMessage(context.Background(), fromAddress, msg)
}

// wait receipt will wait until receipt is received for a given cid
func waitReceipt(cid *cid2.Cid, api *v0api.FullNodeStruct) *types.MessageReceipt {
	// Return until recipient is returned (transaction is processed)
	var receipt *types.MessageReceipt
	var err error
	for {
		receipt, err = api.StateGetReceipt(context.Background(), *cid, types.EmptyTSK)
		if err != nil {
			fmt.Printf("Payment manager has error getting recipient of cid: %s\n", cid.String())
		}
		if receipt != nil {
			break
		}
		// TODO, Make the interval configurable
		time.Sleep(1 * time.Second)
	}
	return receipt
}
