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
	"math/big"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/wcgcyx/fc-retrieval/client/pkg/client"
)

// ClientCLI stores the client struct for api calls
type ClientCLI struct {
	// Boolean indicates if client has been initialised
	initialised bool
	// client struct
	client *client.FilecoinRetrievalClient
}

// Start Client CLI
func main() {
	c := ClientCLI{
		initialised: false,
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
		{Text: "init", Description: "Initialise the client by given key and service API Addr"},
		{Text: "search", Description: "Search gateways by given location"},
		{Text: "add-active-gateway", Description: "Add active gateway"},
		{Text: "ls-active-gateways", Description: "List active gateways"},
		{Text: "inspect-active-gateway", Description: "Inspect given active gateway"},
		{Text: "add-active-provider", Description: "Add active provider"},
		{Text: "ls-active-providers", Description: "List active providers"},
		{Text: "inspect-active-provider", Description: "Inspect given active provider"},
		{Text: "block-gateway", Description: "Block given gateway"},
		{Text: "unblock-gateway", Description: "Unblock given gateway"},
		{Text: "resume-gateway", Description: "Resume given gateway"},
		{Text: "block-provider", Description: "Block given provider"},
		{Text: "unblock-provider", Description: "Unblock given provider"},
		{Text: "resume-provider", Description: "Resume given provider"},
		{Text: "find-offer", Description: "Find offers for given cid"},
		{Text: "find-offer-dht", Description: "Find offers for given cid using DHT discovery"},
		{Text: "list-offers", Description: "List obtained offers for given cid"},
		{Text: "retrieve", Description: "Retrieve data using an offer by given offer digest"},
		{Text: "retrieve-fast", Description: "Fast-retrieve data by given cid (automated offer discovery, selection and data retrieval)"},
		{Text: "exit", Description: "Exit the program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// executor executes the command
func (c *ClientCLI) executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "init-dev":
		// Note: this is a hidden command, used by developers to test
		c.initDev()
	case "init":
		if c.initialised {
			fmt.Println("Client has already been initialised")
			return
		}
		if len(blocks) != 7 {
			fmt.Println("Usage: init ${walletPrivKey} ${lotusAPIAddr} ${lotusAuthToken} ${registerPrivKey} ${registerAPIAddr} ${registerAuthToken}")
			return
		}
		var err error
		c.client, err = client.NewFilecoinRetrievalClient(blocks[1], blocks[2], blocks[3], blocks[4], blocks[5], blocks[6])
		if err != nil {
			fmt.Printf("Error in initialising the client: %v\n", err.Error())
			return
		}
		fmt.Println("Client has been initialised successfully")
		c.initialised = true
	case "search":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: search ${location}")
			return
		}
		gws, err := c.client.Search(blocks[1])
		if err != nil {
			fmt.Printf("Error in searching for gateways in location %v: %v\n", blocks[1], err.Error())
			return
		}
		fmt.Printf("Find gateways in given location %v:\n", blocks[1])
		for _, gw := range gws {
			fmt.Printf("ID: %v\n", gw)
		}
	case "add-active-gateway":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: add-active-gateway ${targetID}")
			return
		}
		err := c.client.AddActiveGW(blocks[1])
		if err != nil {
			fmt.Printf("Error in adding active gateway for %v: %v\n", blocks[1], err.Error())
			return
		}
		fmt.Println("Added.")
	case "ls-active-gateways":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		gws := c.client.ListActiveGWS()
		fmt.Println("Current active gateways:")
		for _, gw := range gws {
			score, pending, blocked, err := c.client.GetGWReputaion(gw)
			history := ""
			temp := c.client.GetGWHistory(gw, 0, 1)
			if len(temp) == 1 {
				history = temp[0]
			}
			if err != nil {
				fmt.Printf("ID: %v\t Error loading reputation details: %v\n", gw, err.Error())
			} else {
				fmt.Printf("ID: %v\tReputation score: %v\tPending: %t\tBlocked: %t\tRecent: %v\n", gw, score, pending, blocked, history)
			}
		}
	case "inspect-active-gateway":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: inspect-active-gateway ${gatewayID}")
			return
		}
		score, pending, blocked, err := c.client.GetGWReputaion(blocks[1])
		if err != nil {
			fmt.Printf("ID: %v\t Error loading reputation details: %v\n", blocks[1], err.Error())
		} else {
			fmt.Printf("ID: %v\tReputation score: %v\tPending: %t\tBlocked: %t\n", blocks[1], score, pending, blocked)
		}
		history := c.client.GetGWHistory(blocks[1], 0, 10)
		fmt.Println("Recent 10 activites:")
		for index, entry := range history {
			fmt.Printf("Activity %v: %v\n", index, entry)
		}
	case "add-active-provider":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: add-active-provider ${targetID}")
			return
		}
		err := c.client.AddActivePVD(blocks[1])
		if err != nil {
			fmt.Printf("Error in adding active provider for %v: %v\n", blocks[1], err.Error())
			return
		}
		fmt.Println("Added.")
	case "ls-active-providers":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		pvds := c.client.ListActivePVDS()
		fmt.Println("Current active providers:")
		for _, pvd := range pvds {
			score, pending, blocked, err := c.client.GetPVDReputaion(pvd)
			history := ""
			temp := c.client.GetPVDHistory(pvd, 0, 1)
			if len(temp) == 1 {
				history = temp[0]
			}
			if err != nil {
				fmt.Printf("ID: %v\t Error loading reputation details: %v\n", pvd, err.Error())
			} else {
				fmt.Printf("ID: %v\tReputation score: %v\tPending: %t\tBlocked: %t\tRecent: %v\n", pvd, score, pending, blocked, history)
			}
		}
	case "inspect-active-provider":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: inspect-active-provider ${providerID}")
			return
		}
		score, pending, blocked, err := c.client.GetPVDReputaion(blocks[1])
		if err != nil {
			fmt.Printf("ID: %v\t Error loading reputation details: %v\n", blocks[1], err.Error())
		} else {
			fmt.Printf("ID: %v\tReputation score: %v\tPending: %t\tBlocked: %t\n", blocks[1], score, pending, blocked)
		}
		history := c.client.GetPVDHistory(blocks[1], 0, 10)
		fmt.Println("Recent 10 activites:")
		for index, entry := range history {
			fmt.Printf("Activity %v: %v\n", index, entry)
		}
	case "block-gateway":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: block-gateway ${gatewayID}")
			return
		}
		c.client.BlockGW(blocks[1])
	case "unblock-gateway":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: unblock-gateway ${gatewayID}")
			return
		}
		c.client.UnblockGW(blocks[1])
	case "resume-gateway":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: resume-gateway ${gatewayID}")
			return
		}
		c.client.ResumeGW(blocks[1])
	case "block-provider":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: block-provider ${providerID}")
			return
		}
		c.client.BlockPVD(blocks[1])
	case "unblock-provider":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: unblock-provider ${providerID}")
			return
		}
		c.client.UnblockPVD(blocks[1])
	case "resume-provider":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: resume-provider ${providerID}")
			return
		}
		c.client.ResumePVD(blocks[1])
	case "find-offer":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: find-offer ${contentID}")
			return
		}
		toContact := make(map[string]uint32)
		for _, gw := range c.client.ListActiveGWS() {
			toContact[gw] = 1
		}
		offers, err := c.client.StandardDiscovery(blocks[1], toContact)
		if err != nil {
			fmt.Printf("Error doing standard discovery for %v: %v\n", blocks[1], err.Error())
			return
		}
		fmt.Println("Find offers: ")
		for _, offer := range offers {
			fmt.Printf("Offer %v: provider-%v, price-%v, expiry-%v, qos-%v\n", offer.GetMessageDigest(), offer.GetProviderID(), offer.GetPrice().String(), offer.GetExpiry(), offer.GetQoS())
		}
	case "find-offer-dht":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 3 {
			fmt.Println("Usage: find-offer ${contentID} ${gatewayID}")
			return
		}
		offers, err := c.client.DHTDiscovery(blocks[1], blocks[2], 8, 1)
		if err != nil {
			fmt.Printf("Error doing dht discovery for %v by %v: %v\n", blocks[1], blocks[2], err.Error())
			return
		}
		fmt.Println("Find offers: ")
		for _, offer := range offers {
			fmt.Printf("Offer %v: provider-%v, price-%v, expiry-%v, qos-%v\n", offer.GetMessageDigest(), offer.GetProviderID(), offer.GetPrice().String(), offer.GetExpiry(), offer.GetQoS())
		}
	case "list-offers":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 2 {
			fmt.Println("Usage: list-offers ${contentID}")
			return
		}
		offers, err := c.client.ListOffers(blocks[1])
		if err != nil {
			fmt.Printf("Error doing list offers for %v: %v\n", blocks[1], err.Error())
			return
		}
		fmt.Println("Find offers: ")
		for _, offer := range offers {
			fmt.Printf("Offer %v: provider-%v, price-%v, expiry-%v, qos-%v\n", offer.GetMessageDigest(), offer.GetProviderID(), offer.GetPrice().String(), offer.GetExpiry(), offer.GetQoS())
		}
	case "retrieve":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 3 {
			fmt.Println("Usage: retrieve ${offerDigest} ${outputName}")
			return
		}
		err := c.client.Retrieve(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error retrieval of offer %v to %v: %v\n", blocks[1], blocks[2], err.Error())
			return
		}
		fmt.Printf("Success, file saved to %v\n", blocks[2])
	case "retrieve-fast":
		if !c.initialised {
			fmt.Println("Client has not been initialised yet")
			return
		}
		if len(blocks) != 4 {
			fmt.Println("Usage: retrieve-fast ${contentID} ${outputName} ${maxPrice}")
			return
		}
		maxPrice, ok := big.NewInt(0).SetString(blocks[3], 10)
		if !ok {
			fmt.Printf("Error parsing bigInt from %v\n", blocks[3])
			return
		}
		err := c.client.FastRetrieve(blocks[1], blocks[2], maxPrice)
		if err != nil {
			fmt.Printf("Error retrieval of offer %v to %v: %v\n", blocks[1], blocks[2], err.Error())
			return
		}
		fmt.Printf("Success, file saved to %v\n", blocks[2])
	case "exit":
		fmt.Println("Shutdown client...")
		if c.client != nil {
			c.client.Shutdown()
		}
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
func (c *ClientCLI) initDev() {
	if c.initialised {
		fmt.Println("Client has already been initialised")
		return
	}

	env := os.Getenv("DEVINIT")
	vars := strings.Split(env, ";")
	lotusAuthToken := vars[0]
	walletPrivKey := vars[1]
	lotusAPIAddr := "http://lotus:1234/rpc/v0"
	registerPrivKey := "_"
	registerAPIAddr := "register:9020"
	registerAuthToken := "_"

	var err error
	c.client, err = client.NewFilecoinRetrievalClient(walletPrivKey, lotusAPIAddr, lotusAuthToken, registerPrivKey, registerAPIAddr, registerAuthToken)
	if err != nil {
		fmt.Printf("Error in initialising the client: %v\n", err.Error())
		return
	}
	fmt.Println("Client has been initialised successfully")
	c.initialised = true
}
