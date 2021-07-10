package main

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

var localLotusAP = "http://127.0.0.1:1234/rpc/v0"

func main() {
	// Get lotus token
	token, acct := getLotusToken()
	res := token
	if os.Args[1] == "gw" {
		cmd := exec.Command("docker", "ps", "--filter", "ancestor=wcgcyx/fc-retrieval/gateway", "--format", "{{.Names}}")
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		gws := strings.Split(string(stdout[:len(stdout)-1]), "\n")
		privKeys, _, err := generateAccount(localLotusAP, token, acct, len(gws))
		if err != nil {
			panic(err)
		}
		for i, gw := range gws {
			cmd = exec.Command("docker", "run", "--net", "shared", "--rm", "giantswarm/tiny-tools", "dig", "+short", gw)
			stdout, err = cmd.Output()
			if err != nil {
				panic(err)
			}
			ip := string(stdout[:len(stdout)-1])
			res = fmt.Sprintf("%v;%v,%v,%v,%v", res, gw, "6465616135313132656636333864653962393132306366363537333664656465", privKeys[i], ip)
		}
		fmt.Println(res)
	} else if os.Args[1] == "pvd" {
		cmd := exec.Command("docker", "ps", "--filter", "ancestor=wcgcyx/fc-retrieval/provider", "--format", "{{.Names}}")
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		gws := strings.Split(string(stdout[:len(stdout)-1]), "\n")
		privKeys, _, err := generateAccount(localLotusAP, token, acct, len(gws))
		if err != nil {
			panic(err)
		}
		for i, gw := range gws {
			cmd = exec.Command("docker", "run", "--net", "shared", "--rm", "giantswarm/tiny-tools", "dig", "+short", gw)
			stdout, err = cmd.Output()
			if err != nil {
				panic(err)
			}
			ip := string(stdout[:len(stdout)-1])
			res = fmt.Sprintf("%v;%v,%v,%v,%v", res, gw, "6465616135313132656636333864653962393132306366363537333664656465", privKeys[i], ip)
		}
		fmt.Println(res)
	} else {
		// Must be client
		privKeys, _, err := generateAccount(localLotusAP, token, acct, 1)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v;%v\n", token, privKeys[0])
	}
}

func getLotusToken() (string, string) {
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
}

// The following helper method is used to generate a new filecoin account with 10 filecoins of balance
func generateAccount(localLotusAP string, token string, superAcct string, num int) ([]string, []string, error) {
	// Get API
	var api v0api.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + token}}
	closer, err := jsonrpc.NewMergeClient(context.Background(), localLotusAP, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		return nil, nil, err
	}
	defer closer()

	mainAddress, err := address.NewFromString(superAcct)
	if err != nil {
		return nil, nil, err
	}

	privateKeys := make([]string, 0)
	addresses := make([]string, 0)
	cids := make([]cid2.Cid, 0)

	// Send messages
	for i := 0; i < num; i++ {
		privKey, pubKey, err := generateKeyPair()
		if err != nil {
			return nil, nil, err
		}
		privKeyStr := hex.EncodeToString(privKey)

		address1, err := address.NewSecp256k1Address(pubKey)
		if err != nil {
			return nil, nil, err
		}

		// Get amount
		amt, err := types.ParseFIL("100")
		if err != nil {
			return nil, nil, err
		}

		msg := &types.Message{
			To:     address1,
			From:   mainAddress,
			Value:  types.BigInt(amt),
			Method: 0,
		}
		signedMsg, err := fillMsg(mainAddress, &api, msg)
		if err != nil {
			return nil, nil, err
		}

		// Send request to lotus
		cid, err := api.MpoolPush(context.Background(), signedMsg)
		if err != nil {
			return nil, nil, err
		}
		cids = append(cids, cid)

		// Add to result
		privateKeys = append(privateKeys, privKeyStr)
		addresses = append(addresses, address1.String())
	}

	// Finally check receipts
	for _, cid := range cids {
		receipt := waitReceipt(&cid, &api)
		if receipt.ExitCode != 0 {
			return nil, nil, errors.New("Transaction fail to execute")
		}
	}

	return privateKeys, addresses, nil
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

func generateKeyPair() ([]byte, []byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	pubKey := crypto.PublicKey(privKey)
	return privKey, pubKey, nil
}
