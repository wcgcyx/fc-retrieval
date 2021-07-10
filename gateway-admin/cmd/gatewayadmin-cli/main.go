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
	defaultGW string
	admin     *gatewayadmin.FilecoinRetrievalGatewayAdmin
}

// Start Client CLI
func main() {
	c := GatewayAdminCLI{
		defaultGW: "",
		admin:     gatewayadmin.NewFilecoinRetrievalGatewayAdmin(),
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
		{Text: "init", Description: "Initialise a given gateway"},
		{Text: "set-default", Description: "Set the default gateway"},
		{Text: "sync", Description: "Force the default gateway to sync"},
		{Text: "ls", Description: "List gateways this admin is administering"},
		{Text: "ls-peers", Description: "List the peers of the default gateway"},
		{Text: "inspect-gwpeer", Description: "Inspect a given gateway peer of the default gateway"},
		{Text: "block-gwpeer", Description: "Block a given gateway peer of the default gateway"},
		{Text: "unblock-gwpeer", Description: "Unblock a given gateway peer of the default gateway"},
		{Text: "resume-gwpeer", Description: "Resume a given gateway peer of the default gateway"},
		{Text: "inspect-pvdpeer", Description: "Inspect a given provider peer of the default gateway"},
		{Text: "block-pvdpeer", Description: "Block a given provider peer of the default gateway"},
		{Text: "unblock-pvdpeer", Description: "Unblock a given provider peer of the default gateway"},
		{Text: "resume-pvdpeer", Description: "Resume a given provider peer of the default gateway"},
		{Text: "list-cids", Description: "List the cid access frequency of the default gateway"},
		{Text: "get-offers", Description: "Get offers by given cid from the default gateway"},
		{Text: "cache-content", Description: "Cache offer by given offer digest and and a given sub cid using the default gateway"},
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
	case "sync-dev":
		// Note: this is a hidden command, used by developers to test
		c.syncDev()
	case "init":
		if len(blocks) != 13 {
			fmt.Println("Usage: init ${adminURL} ${adminKey} ${p2pPort} ${gatewayIP} ${rootPrivKey} ${lotusAPIAddr} {lotusAuthToken} {registerPrivKey} {registerAPIAddr} {registerAuthToken} {regionCode} {alias}")
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
		if c.defaultGW == "" {
			ids, _, _ := c.admin.ListGateways()
			c.defaultGW = ids[0]
		}
		fmt.Printf("Gateway has been initialised\n")
	case "sync":
		err := c.admin.ForceSync(c.defaultGW)
		if err != nil {
			fmt.Printf("Error in force syncing the given gateway: %v\n", err.Error())
			return
		}
		fmt.Println("Done")
	case "set-default":
		if len(blocks) != 2 {
			fmt.Println("Usage: set-default ${gatewayID}")
			return
		}
		c.defaultGW = blocks[1]
		fmt.Println("Done")
	case "ls":
		ids, regions, aliases := c.admin.ListGateways()
		fmt.Println("Managed gateways:")
		for i, id := range ids {
			fmt.Printf("Gateway %v:\tid-%v\tregion-%v\talias-%v", i, id, regions[i], aliases[i])
			if id == c.defaultGW {
				fmt.Printf("\t(default)\n")
			} else {
				fmt.Printf("\n")
			}
		}
	case "ls-peers":
		gwIDs, gwScore, gwPending, gwBlocked, gwRecent, pvdIDs, pvdScore, pvdPending, pvdBlocked, pvdRecent, err := c.admin.ListPeers(c.defaultGW)
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
	case "inspect-gwpeer":
		if len(blocks) != 2 {
			fmt.Println("Usage: inspect-gwpeer ${peerID}")
			return
		}
		score, pending, blocked, history, err := c.admin.InspectGateway(c.defaultGW, blocks[1])
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
	case "block-gwpeer":
		if len(blocks) != 2 {
			fmt.Println("Usage: block-gwpeer ${peerID}")
			return
		}
		ok, msg, err := c.admin.BlockGateway(c.defaultGW, blocks[1])
		if err != nil {
			fmt.Printf("Error in blocking peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to block peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "unblock-gwpeer":
		if len(blocks) != 2 {
			fmt.Println("Usage: unblock-gwpeer ${peerID}")
			return
		}
		ok, msg, err := c.admin.UnblockGateway(c.defaultGW, blocks[1])
		if err != nil {
			fmt.Printf("Error in unblocking peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to unblock peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "resume-gwpeer":
		if len(blocks) != 2 {
			fmt.Println("Usage: resume-gwpeer ${peerID}")
			return
		}
		ok, msg, err := c.admin.ResumeGateway(c.defaultGW, blocks[1])
		if err != nil {
			fmt.Printf("Error in resuming peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to resuming peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "inspect-pvdpeer":
		if len(blocks) != 2 {
			fmt.Println("Usage: inspect-pvdpeer ${peerID}")
			return
		}
		score, pending, blocked, history, err := c.admin.InspectProvider(c.defaultGW, blocks[1])
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
	case "block-pvdpeer":
		if len(blocks) != 2 {
			fmt.Println("Usage: block-pvdpeer ${peerID}")
			return
		}
		ok, msg, err := c.admin.BlockProvider(c.defaultGW, blocks[1])
		if err != nil {
			fmt.Printf("Error in blocking peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to block peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "unblock-pvdpeer":
		if len(blocks) != 2 {
			fmt.Println("Usage: unblock-pvdpeer ${peerID}")
			return
		}
		ok, msg, err := c.admin.UnblockProvider(c.defaultGW, blocks[1])
		if err != nil {
			fmt.Printf("Error in unblocking peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to unblock peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "resume-pvdpeer":
		if len(blocks) != 2 {
			fmt.Println("Usage: resume-pvdpeer ${peerID}")
			return
		}
		ok, msg, err := c.admin.ResumeProvider(c.defaultGW, blocks[1])
		if err != nil {
			fmt.Printf("Error in resuming peer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to resuming peer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "list-cids":
		if len(blocks) != 2 {
			fmt.Println("Usage: list-cids ${page}")
			return
		}
		page, err := strconv.ParseUint(blocks[1], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing unit int %v: %v\n", blocks[1], err.Error())
			return
		}
		cids, counts, err := c.admin.ListCIDFrequency(c.defaultGW, uint(page))
		if err != nil {
			fmt.Printf("Error in listing cids for given gateway: %v\n", err.Error())
			return
		}
		fmt.Println("CID access frequency:")
		for i, cid := range cids {
			fmt.Printf("Access count: %v\t\tCID: %v\n", counts[i], cid)
		}
	case "get-offers":
		if len(blocks) != 2 {
			fmt.Println("Usage: get-offers ${cid}")
			return
		}
		digests, providers, prices, expriy, qos, err := c.admin.GetOfferByCID(c.defaultGW, blocks[1])
		if err != nil {
			fmt.Printf("Error in get gateway offers by cid: %v\n", err.Error())
			return
		}
		fmt.Printf("Offers containing cid %v:\n", blocks[1])
		for i, digest := range digests {
			fmt.Printf("Offer %v: provider-%v price-%v expiry-%v qos-%v\n", digest, providers[i], prices[i], expriy[i], qos[i])
		}
	case "cache-content":
		if len(blocks) != 3 {
			fmt.Println("Usage: cache-content ${digest} ${cid}")
			return
		}
		ok, msg, err := c.admin.CacheOfferByDigest(c.defaultGW, blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in cache offer for given gateway: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to cache offer for given gateway: %v\n", msg)
			return
		}
		fmt.Println("Content cached")
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
		gatewayIP := info[3]
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
		if c.defaultGW == "" {
			ids, _, _ := c.admin.ListGateways()
			c.defaultGW = ids[0]
		}
	}
	fmt.Println("All gateways are initialised.")
}

// syncDev is only used by developers to test, its hard-coded
func (c *GatewayAdminCLI) syncDev() {
	ids, _, _ := c.admin.ListGateways()
	for _, id := range ids {
		err := c.admin.ForceSync(id)
		if err != nil {
			panic(err)
		}
	}
}
