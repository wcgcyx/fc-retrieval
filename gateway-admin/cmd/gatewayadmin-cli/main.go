/*
Package main - program entry point for a Retrieval Client cli.
*/
package main

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
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/wcgcyx/fc-retrieval/gateway-admin/pkg/gatewayadmin"
)

// GatewayAdminCLI stores the gateway admin struct for api calls
type GatewayAdminCLI struct {
	admin *gatewayadmin.FilecoinRetrievalGatewayAdmin
}

// Start Client CLI
func main() {
	c := GatewayAdminCLI{
		gatewayadmin.NewFilecoinRetrievalGatewayAdmin(),
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
		}
		handleExit()
	}()
	p := prompt.New(
		c.executor,
		completer,
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}

// completer completes the input
func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "init-gateway", Description: "Initialise given gateway"},
		{Text: "ls-gateways", Description: "List gateways this admin is administering"},
		{Text: "ls-gateway-peers", Description: "List the peers of a given administered gateway"},
		{Text: "inspect-gateway-gwpeer", Description: "Inspect a given gateway peer of a given administered gateway"},
		{Text: "block-gateway-gwpeer", Description: "Block a given gateway peer of a given administered gateway"},
		{Text: "unblock-gateway-gwpeer", Description: "Unblock a given gateway peer of a given administered gateway"},
		{Text: "resume-gateway-gwpeer", Description: "Resume a given gateway peer of a given administered gateway"},
		{Text: "inspect-gateway-pvdpeer", Description: "Inspect a given provider peer of a given administered gateway"},
		{Text: "block-gateway-pvdpeer", Description: "Block a given provider peer of a given administered gateway"},
		{Text: "unblock-gateway-pvdpeer", Description: "Unblock a given provider peer of a given administered gateway"},
		{Text: "resume-gateway-pvdpeer", Description: "Resume a given provider peer of a given administered gateway"},
		{Text: "list-gateway-cids", Description: "List the cid access frequency of a given administered gateway"},
		{Text: "get-gateway-offers-by-cid", Description: "Get offers by given cid from a given administered gateway"},
		{Text: "cache-gateway-offer-by-digest", Description: "Cache offer by given digest by a given administered gateway"},
		{Text: "exit", Description: "Exit the program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// executor executes the command
func (c *GatewayAdminCLI) executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "init-dev":
		// Note: this is a hidden command, used by developers to test
		c.initDev()
	case "init-gateway":
		if len(blocks) != 13 {
			fmt.Println("Usage: init-gateway ${adminURL} ${adminKey} ${p2pPort} ${gatewayIP} ${rootPrivKey} ${lotusAPIAddr} {lotusAuthToken} {registerPrivKey} {registerAPIAddr} {registerAuthToken} {regionCode} {alias}")
			return
		}
		port, err := strconv.ParseInt(blocks[3], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing unit int %v: %v\n", blocks[3], err.Error())
			return
		}
		err = c.admin.InitialiseGateway(blocks[1], blocks[2], int(port), blocks[4], blocks[5], blocks[6], blocks[7], blocks[8], blocks[9], blocks[10], blocks[11], blocks[12])
		if err != nil {
			fmt.Printf("Error in initialising given gateway: %v\n", err.Error())
			return
		}
		fmt.Printf("Gateway has been initialised\n")
	case "ls-gateways":
		ids, regions, aliases := c.admin.ListGateways()
		fmt.Println("Managed gateways:")
		for i, id := range ids {
			fmt.Printf("Gateway %v:\tid-%v\tregion-%v\talias-%v\n", i, id, regions[i], aliases[i])
		}
	case "ls-gateway-peers":
		if len(blocks) != 2 {
			fmt.Println("Usage: ls-gateway-peers ${targetID}")
			return
		}
		gwIDs, gwScore, gwPending, gwBlocked, gwRecent, pvdIDs, pvdScore, pvdPending, pvdBlocked, pvdRecent, err := c.admin.ListPeers(blocks[1])
		if err != nil {
			fmt.Printf("Error in listing peers for given gateway: %v\n", err.Error())
			return
		}
		fmt.Println("Peer gateways:")
		for i, gwID := range gwIDs {
			fmt.Printf("%v:\tid-%v\tscore-%v\tpending-%t\tblocked-%t\trecent-%v", i, gwID, gwScore[i], gwPending[i], gwBlocked[i], gwRecent[i])
		}
		fmt.Println("Peer providers:")
		for i, pvdID := range pvdIDs {
			fmt.Printf("%v:\tid-%v\tscore-%v\tpending-%t\tblocked-%t\trecent-%v", i, pvdID, pvdScore[i], pvdPending[i], pvdBlocked[i], pvdRecent[i])
		}
	case "inspect-gateway-gwpeer":
		if len(blocks) != 3 {
			fmt.Println("Usage: inspect-gateway-gwpeer ${targetID} ${peerID}")
			return
		}
		score, pending, blocked, history, err := c.admin.InspectGateway(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in inspecting peer for given gateway: %v\n", err.Error())
			return
		}
		fmt.Printf("Gateway peer %v:\n", blocks[2])
		fmt.Printf("Reputation score: %v\n", score)
		fmt.Printf("Pending: %t\n", pending)
		fmt.Printf("Blocked: %t\n", blocked)
		fmt.Println("Recent history:")
		for i, entry := range history {
			fmt.Printf("History %v - %v\n", i, entry)
		}
	case "block-gateway-gwpeer":
		if len(blocks) != 3 {
			fmt.Println("Usage: block-gateway-gwpeer ${targetID} ${peerID}")
			return
		}
		ok, msg, err := c.admin.BlockGateway(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in blocking peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to block peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "unblock-gateway-gwpeer":
		if len(blocks) != 3 {
			fmt.Println("Usage: unblock-gateway-gwpeer ${targetID} ${peerID}")
			return
		}
		ok, msg, err := c.admin.UnblockGateway(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in unblocking peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to unblock peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "resume-gateway-gwpeer":
		if len(blocks) != 3 {
			fmt.Println("Usage: resume-gateway-gwpeer ${targetID} ${peerID}")
			return
		}
		ok, msg, err := c.admin.ResumeGateway(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in resuming peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to resuming peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "inspect-gateway-pvdpeer":
		if len(blocks) != 3 {
			fmt.Println("Usage: inspect-gateway-pvdpeer ${targetID} ${peerID}")
			return
		}
		score, pending, blocked, history, err := c.admin.InspectProvider(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in inspecting peer for given gateway: %v\n", err.Error())
			return
		}
		fmt.Printf("Provider peer %v:\n", blocks[2])
		fmt.Printf("Reputation score: %v\n", score)
		fmt.Printf("Pending: %t\n", pending)
		fmt.Printf("Blocked: %t\n", blocked)
		fmt.Println("Recent history:")
		for i, entry := range history {
			fmt.Printf("History %v - %v\n", i, entry)
		}
	case "block-gateway-pvdpeer":
		if len(blocks) != 3 {
			fmt.Println("Usage: block-gateway-pvdpeer ${targetID} ${peerID}")
			return
		}
		ok, msg, err := c.admin.BlockProvider(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in blocking peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to block peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "unblock-gateway-pvdpeer":
		if len(blocks) != 3 {
			fmt.Println("Usage: unblock-gateway-pvdpeer ${targetID} ${peerID}")
			return
		}
		ok, msg, err := c.admin.UnblockProvider(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in unblocking peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to unblock peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "resume-gateway-pvdpeer":
		if len(blocks) != 3 {
			fmt.Println("Usage: resume-gateway-pvdpeer ${targetID} ${peerID}")
			return
		}
		ok, msg, err := c.admin.ResumeProvider(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in resuming peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to resuming peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "list-gateway-cids":
		if len(blocks) != 3 {
			fmt.Println("Usage: list-gateway-cids ${targetID} ${page}")
			return
		}
		page, err := strconv.ParseUint(blocks[2], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing unit int %v: %v\n", blocks[3], err.Error())
			return
		}
		cids, counts, err := c.admin.ListCIDFrequency(blocks[1], uint(page))
		if err != nil {
			fmt.Printf("Error in listing cids for given gateway: %v\n", err.Error())
			return
		}
		fmt.Println("CID access frequency:")
		for i, cid := range cids {
			fmt.Printf("Access count: %v\t\tCID: %v\n", counts[i], cid)
		}
	case "get-gateway-offers-by-cid":
		if len(blocks) != 3 {
			fmt.Println("Usage: get-gateway-offers-by-cid ${targetID} ${cid}")
			return
		}
		digests, providers, prices, expriy, qos, err := c.admin.GetOfferByCID(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in get gateway offers by cid: %v\n", err.Error())
			return
		}
		fmt.Printf("Offers containing cid %v:\n", blocks[2])
		for i, digest := range digests {
			fmt.Printf("Offer %v: provider-%v price-%v expiry-%v qos-%v\n", digest, providers[i], prices[i], expriy[i], qos[i])
		}
	case "cache-gateway-offer-by-digest":
		if len(blocks) != 3 {
			fmt.Println("Usage: cache-gateway-offer-by-digest ${targetID} ${digest} ${cid}")
			return
		}
		ok, msg, err := c.admin.CacheOfferByDigest(blocks[1], blocks[2], blocks[3])
		if err != nil {
			fmt.Printf("Error in cache offer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to cache offer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Offer cached")
	case "exit":
		fmt.Println("Shutdown gateway admin...")
		fmt.Println("Bye!")
		os.Exit(0)
	}
}

// handleExit fixes the problem of broken terminal when exit in Linux
// ref: https://www.gitmemory.com/issue/c-bata/go-prompt/228/820639887
func handleExit() {
	if _, err := os.Stat("/bin/stty"); os.IsNotExist(err) {
		return
	}
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

// initDev is only used by developers to test, its hard-coded
func (c *GatewayAdminCLI) initDev() {
	env := os.Getenv("DEVINIT")
	vars := strings.Split(env, ";")
	lotusAuthToken := vars[0]
	for i := 1; i < len(vars); i++ {
		info := strings.Split(vars[i], ",")
		adminURL := fmt.Sprintf("%v:9010", info[0])
		adminKey := info[1]
		gatewayIP := info[0]
		rootPrivKey := info[2]
		lotusAPIAddr := "http://lotus:1234/rpc/v0"
		registerPrivKey := "_"
		registerAPIAddr := "register:9020"
		registerAuthToken := "_"
		regionCode := "au"
		alias := info[0]
		err := c.admin.InitialiseGateway(adminURL, adminKey, 9011, gatewayIP, rootPrivKey, lotusAPIAddr, lotusAuthToken, registerPrivKey, registerAPIAddr, registerAuthToken, regionCode, alias)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("All gateways are initialised.")
}
