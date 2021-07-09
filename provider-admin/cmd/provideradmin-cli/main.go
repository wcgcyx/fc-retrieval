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
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/wcgcyx/fc-retrieval/provider-admin/pkg/provideradmin"
)

// ProviderAdminCLI stores the provider admin struct for api calls
type ProviderAdminCLI struct {
	admin *provideradmin.FilecoinRetrievalProviderAdmin
}

// Start Client CLI
func main() {
	c := ProviderAdminCLI{
		provideradmin.NewFilecoinRetrievalProviderAdmin(),
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
		{Text: "init-provider", Description: "Initialise given provider"},
		{Text: "ls-providers", Description: "List providers this admin is administering"},
		{Text: "list-provider-files", Description: "List files given provider is monitoring"},
		{Text: "get-providers-offers-by-cid", Description: "Get offers by given cid from a given administered provider"},
		{Text: "upload-provider-file", Description: "Upload a file to a given administered provider (max 25MB)"},
		{Text: "publish-provider-offer", Description: "Ask a administered provider to publish an offer"},
		{Text: "fast-publish-provider-offer", Description: "Upload a given file to a administered provider and ask it to publish an offer"},
		{Text: "exit", Description: "Exit the program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// executor executes the command
func (c *ProviderAdminCLI) executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "init-dev":
		// Note: this is a hidden command, used by developers to test
		c.initDev()
	case "init-provider":
		if len(blocks) != 13 {
			fmt.Println("Usage: init-provider ${adminURL} ${adminKey} ${p2pPort} ${providerIP} ${rootPrivKey} ${lotusAPIAddr} {lotusAuthToken} {registerPrivKey} {registerAPIAddr} {registerAuthToken} {regionCode} {alias}")
			return
		}
		port, err := strconv.ParseInt(blocks[3], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing unit int %v: %v\n", blocks[3], err.Error())
			return
		}
		err = c.admin.InitialiseProvider(blocks[1], blocks[2], int(port), blocks[4], blocks[5], blocks[6], blocks[7], blocks[8], blocks[9], blocks[10], blocks[11], blocks[12])
		if err != nil {
			fmt.Printf("Error in initialising given provider: %v\n", err.Error())
			return
		}
		fmt.Printf("Provider has been initialised\n")
	case "ls-providers":
		ids, regions, aliases := c.admin.ListProviders()
		fmt.Println("Managed providers:")
		for i, id := range ids {
			fmt.Printf("Provider %v:\tid-%v\tregion-%v\talias-%v\n", i, id, regions[i], aliases[i])
		}
	case "list-provider-files":
		if len(blocks) != 2 {
			fmt.Println("Usage: ls-provider-files ${targetID}")
			return
		}
		files, cids, sizes, published, frequency, err := c.admin.ListFiles(blocks[1])
		if err != nil {
			fmt.Printf("Error in listing files for given provider: %v\n", err.Error())
			return
		}
		fmt.Println("Listing files:")
		for i, file := range files {
			fmt.Printf("Name: %v\n", file)
			fmt.Printf("\tCID: %v\t Size: %v\t Published: %t\t Frequency: %v\n", cids[i], sizes[i], published[i], frequency[i])
		}
	case "get-provider-offers-by-cid":
		if len(blocks) != 3 {
			fmt.Println("Usage: get-provider-offers-by-cid ${targetID} ${cid}")
			return
		}
		digests, providers, prices, expriy, qos, err := c.admin.GetOfferByCID(blocks[1], blocks[2])
		if err != nil {
			fmt.Printf("Error in get provider offers by cid: %v\n", err.Error())
			return
		}
		fmt.Printf("Offers containing cid %v:\n", blocks[2])
		for i, digest := range digests {
			fmt.Printf("Offer %v: provider-%v price-%v expiry-%v qos-%v\n", digest, providers[i], prices[i], expriy[i], qos[i])
		}
	case "upload-provider-file":
		if len(blocks) != 4 {
			fmt.Println("Usage: upload-provider-file ${targetID} ${local-file} ${remote-filename}")
			return
		}
		ok, msg, err := c.admin.UploadFile(blocks[1], blocks[2], blocks[3])
		if err != nil {
			fmt.Printf("Error in uploading file to provider: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to upload file to given provider: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "publish-provider-offer":
		if len(blocks) < 6 {
			fmt.Println("Usage: publish-provider-offer ${targetID} [${file}...] ${price} ${expiry} ${qos}")
			return
		}
		price, ok := big.NewInt(0).SetString(blocks[len(blocks)-3], 10)
		if !ok {
			fmt.Println("Error parsing price")
			return
		}
		period, err := time.ParseDuration(blocks[len(blocks)-2])
		if err != nil {
			fmt.Printf("Error parsing period: %v\n", err.Error())
			return
		}
		if period <= time.Hour*12 {
			fmt.Printf("Too short period: %v, need to be at least 12 hours\n", period)
			return
		}
		qos, err := strconv.ParseUint(blocks[len(blocks)-1], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing qos: %v\n", err.Error())
			return
		}
		ok, msg, err := c.admin.PublishOffer(blocks[1], blocks[2:len(blocks)-3], price, time.Now().Add(period).Unix(), qos)
		if err != nil {
			fmt.Printf("Error in publishing offer from provider: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to publish offer from provider: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "fast-publish-provider-offer":
		if len(blocks) != 7 {
			fmt.Println("Usage: publish-provider-offer ${targetID} ${local-file} ${remote-filename} ${price} ${expiry} ${qos}")
			return
		}
		price, ok := big.NewInt(0).SetString(blocks[len(blocks)-3], 10)
		if !ok {
			fmt.Println("Error parsing price")
			return
		}
		period, err := time.ParseDuration(blocks[len(blocks)-2])
		if err != nil {
			fmt.Printf("Error parsing period: %v\n", err.Error())
			return
		}
		if period <= time.Hour*12 {
			fmt.Printf("Too short period: %v, need to be at least 12 hours\n", period)
			return
		}
		qos, err := strconv.ParseUint(blocks[len(blocks)-1], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing qos: %v\n", err.Error())
			return
		}
		ok, msg, err := c.admin.UploadFile(blocks[1], blocks[2], blocks[3])
		if err != nil {
			fmt.Printf("Error in uploading file to provider: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to upload file to given provider: %v\n", msg)
			return
		}
		ok, msg, err = c.admin.PublishOffer(blocks[1], []string{blocks[3]}, price, time.Now().Add(period).Unix(), qos)
		if err != nil {
			fmt.Printf("Error in publishing offer from provider: %v\n", err.Error())
			return
		}
		if !ok {
			fmt.Printf("Fail to publish offer from provider: %v\n", msg)
			return
		}
		fmt.Println("Done")
	case "exit":
		fmt.Println("Shutdown provider admin...")
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
func (c *ProviderAdminCLI) initDev() {
	env := os.Getenv("DEVINIT")
	vars := strings.Split(env, ";")
	lotusAuthToken := vars[0]
	for i := 1; i < len(vars); i++ {
		info := strings.Split(vars[i], ",")
		adminURL := fmt.Sprintf("%v:9010", info[0])
		adminKey := info[1]
		providerIP := info[3]
		rootPrivKey := info[2]
		lotusAPIAddr := "http://lotus:1234/rpc/v0"
		registerPrivKey := "_"
		registerAPIAddr := "register:9020"
		registerAuthToken := "_"
		regionCode := "au"
		alias := info[0]
		err := c.admin.InitialiseProvider(adminURL, adminKey, 9011, providerIP, rootPrivKey, lotusAPIAddr, lotusAuthToken, registerPrivKey, registerAPIAddr, registerAuthToken, regionCode, alias)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("All providers are initialised.")
}
