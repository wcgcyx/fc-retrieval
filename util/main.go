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

// Content
// test1.txt CID-QmUN1ytvX4w2VG5LinWySBriHKjCA6484mdbyMJw36LdHa hash-ab9362d3129b3cf191004454aa447883b13b2f39808384b8a99a514362ac76f7
// 		Before: [GW14, GW29]						After: [GW15, GW29] and GW32
// test2.txt CID-QmYb36f6SPpEN8oeznyxD5qSwztygQrK4jn3JEaCvwxUBx hash-df14ea2bf1745c3bb4c227bc407661eabbe7ac14ff84141dcd241a05ccddb73d
//		Before: [GW20, GW31] and [GW0, GW3]			After: [GW20, GW31] and [GW0, GW3]
// test3.txt CID-QmRLm3zrumrNrWSwy3Awy4MJyK1Ujk2efqfAV8SaHM2z7X hash-494f03ec96e583f1b96778eb40b8ad01038158db2b9821fbd3df1a96cdf4315c
//		Before: [GW2, GW16]							After: [GW2, GW16] and GW32

var localLotusAP = "http://127.0.0.1:1234/rpc/v0"
var gatewaysKeyMap = map[string]string{
	"gateway0":  "5a7c858349ed16806e931b5fbb359e031529faf7218340d3b7e56bd58cfb97e5", // ID - 0523070e9e2dbd47c9c11c6ba861c801cd998d3002f245c3bbec6a0bf0d1fd49 // Before [GW25, GW9]		// After [GW25, GW9]
	"gateway1":  "c9dabac5dec3d1927fd6976011c8f5ef99485631f1099961e686529032da76da", // ID - 0e0282464798944a922f5a6bc8d0e026850b4529797925982af839bde1a43dbd // Before [GW25, GW9]		// After [GW25, GW9]
	"gateway2":  "3b139571edf7816ed090f29a4c2e15c2585728ea245fad2ff64ea6bd1fe8ddbe", // ID - 10e638412342152b9ff317a4fcfd57f215562fe050fe6836c7da3193eddcd29e // Before [GW26, GW10]	// After [GW26, GW10]
	"gateway3":  "9f1d730e7c0199107501cbe299eb63ae8d81e89715213a9b0913d66e439ce9d5", // ID - 160394959c1b6e04be181691c0049965e64cbe2cb06aa30677634ff0da72eba5 // Before [GW26, GW10]	// After [GW26, GW10]
	"gateway4":  "90c1a2dc57d19b58ef2b7fc757356f32c457ba36c13cd0cb99119aca9340881e", // ID - 22850133bd4acd610755550ad9d09bbf9794162fcf81c5ac60acd8f17ceefca4 // Before [GW28, GW12]	// After [GW28, GW12]
	"gateway5":  "33ab67450814288bec155ed1620c1e824f6e19d4ad483a0aa5bc2312ca55bfc5", // ID - 29ea6805c30728655f832e04fde46dc118f09c984e2b1f268cf34f8184d3c90e // Before [GW29, GW13]	// After [GW29, GW13]
	"gateway6":  "38b32ba7f9d7f521fdd00cef9961e244aad4968618b1943c4e70ec9e2bd69389", // ID - 338b7c64c252810dddca85f6014d6daa3a7957a0f974419a9e497f3c95505bcf // Before [GW30, GW14]	// After [GW30, GW14]
	"gateway7":  "5715b6386487c213c3314b9bbac11c53d78be8130d74d21412aeb5f6fdc7f4bd", // ID - 3846722e7739daba07067b3077c1a042f344c1aa47975a68e6e98227f47282cf // Before [GW31, GW15]	// After [GW31, GW15]
	"gateway8":  "2e934326ac30802b3e9c0c0b5adfb846d2eb4eb07092f0b8135dbb855872dc7f", // ID - 43daf05542473ef398c9c876cfc68029bf2719249b09b49a9d672cb06518a492 // Before [GW0, GW16]		// After [GW0, GW15] + GW32
	"gateway9":  "c241b786a28636ca22bb9476a2ff269febc7375452f85654677613bdaf2e66c1", // ID - 468ca79b8de14af50138cb6c264ba7bb227e3097ed7271ce82dede2b7f2da0a7 // Before [GW0, GW16]		// After [GW1, GW16] inc GW32
	"gateway10": "03254c3c191e52558559c2c2324ef49bb689c87ca0ce02e1a7a72f6dbc922223", // ID - 52d0284f80058cad04e5cfb930acf93e1922bf19ba0fcacc62a7e7464055ceb7 // Before [GW2, GW18]		// After [GW3, GW18] inc GW32
	"gateway11": "2fd67c6ffd03695cacb66af2a363954880b21112e02bb79516b8fb45e9821c37", // ID - 5c53dd80f68d319c09ec01505eff45d1cebc90aa0c0b55470d0d9980ee519564 // Before [GW3, GW19]		// After [GW4, GW19] inc GW32
	"gateway12": "9708e98ce259c24515474f9ebd50f80c05354690238775371468a6caa6d83c6d", // ID - 6029da88e8c56eba13f8ce3295452a48571906baa80d2ef774468f506bd79c3f // Before [GW4, GW20]		// After [GW4, GW19] inc GW32
	"gateway13": "1098d37e545472f3bf65ce72dcf69c615527817bc506aeddf0ae00fdde746c9c", // ID - 69eff4c36749fc7108a3f419e5e868a94d73d2fa2fe743e62b1f8b4309a10530 // Before [GW5, GW21]		// After [GW5, GW20] inc GW32
	"gateway14": "dd547100866b46380543fa8bd20cf6c475a49fcf325904d1553d42c844cb9146", // ID - 7128499dd89ffcf278b622190bfb344eefd5fdf33ecc64ce1508197cb22419c8 // Before [GW6, GW22]		// After [GW6, GW21] inc GW32
	"gateway15": "a741d27ec9d30ad8589f8c96d81fdce8108e07e591e6989ed6584daa4e1f2f12", // ID - 7728977ae23b0e0abd494cb33b5fda99bc4d1f66a18bf2eb9056c10ee499450a // Before [GW6, GW22]		// After [GW7, GW22] inc GW32
	"gateway16": "2dacc97bf6d472ff286dfab972c545cbc73c475b0fd5e772d4e2d9f1d40286fb", // ID - 83cae55a61335add6368e63b79471191ba11b3fd43455729efeaf73b1c6de7fa // Before [GW8, GW24]		// After [GW9, GW24] inc GW32
	"gateway17": "ee9ebd34579b551147a96e39812a63a1342011991ab7b6cceacdf128eb949313", // ID - 8eeadada16408ce527bd899dad5fd1650145de8cf337269027fe7c7424fd5991 // Before [GW10, GW26]	// After [GW10, GW25] inc GW32
	"gateway18": "bc7c7c44cf5767908d15783f49883c2df0f7a2c7c5deeb1a89a9b0d8c1ffb2e4", // ID - 925498ad9d0f2aff6a8d8b30e304043eedfcba96c6cd1e64e05d2b12a070cc08 // Before [GW10, GW26]	// After [GW10, GW25] inc GW32
	"gateway19": "9c4a15559755df2f45339f232ede123dfd2019012a2f93d498f18a30e72758a0", // ID - 99ca35e67778bb788ec120b9ddb72e24b641e71c8206f2346ba861d3d5ec8c18 // Before [GW11, GW27]	// After [GW11, GW26] inc GW32
	"gateway20": "7917c714b7781ddf2e0bc49def4a9e9b052542ed72196bdc5df78a4a71b180a4", // ID - a3a2c4a1f7828c3913f6fb2bfa8c3ee311b9f02a9021f6b0f5515adbd9645dd0 // Before [GW12, GW28]	// After [GW13, GW28] inc GW32
	"gateway21": "f9a39836a365aa5868f611a541ffa57f1b98c3bf3e0cadbd5272f50eea396763", // ID - aaa225c3f6cd07fe3c7163e87fcce95265ca02e0e4875ff7362aaa639bb3da99 // Before [GW13, GW29]	// After [GW14, GW29] inc GW32
	"gateway22": "1a1dfacc2a27059dbea77b7a25f71fc0c48bc6fcd5fa55efe8cfadf30bca6118", // ID - b589d9a05f3ee5bef3cd5d0ea2960344d64934d0f8366dfb8d2ec4e97eb05bf9 // Before [GW14, GW30]	// After [GW15, GW30] inc GW32
	"gateway23": "cd97f448fe5ca5e12b10b4bf24887882528af0494d1f3d043aaab843a2d2d6c4", // ID - bee42f182166c41bf1da487183bc6ed1f2ba482906e40e5d1517a0bd2313df01 // Before [GW16, GW0]		// After [GW16, GW31] + GW32
	"gateway24": "2e3d92ec9c7cae61d42ba934ca7a13ee58a54f9dbe59d9226527ee946f0808e7", // ID - c10fc3c35c45ddff235a11ae2c41bcfd345f0414615ad618e149f167ccd94d8f // Before [GW16, GW0]		// After [GW16, GW0]
	"gateway25": "a7164b946e12a907ffbd085385da146cd3a0adf40edba349d1c5bceea6af711e", // ID - ca34fcc6df3ca2cb31073c1ed246c8647276a2903e59fc2f30b6a20358ba710b // Before [GW17, GW1]		// After [GW17, GW1]
	"gateway26": "2dd47d9bc2150583d8d97f1354b7884ecafed863691e9ecf0731895784147299", // ID - d2a03495541ede308f04348a6148c9f753d38201e569c84f06e8f97f673cc9ff // Before [GW18, GW2]		// After [GW18, GW2]
	"gateway27": "288849ae970dce8af8da6fadc7b331bcbfd28705f00ad9ed72605653482d23f7", // ID - d9d97672451730fd8594a2841f3643f829096feee41d382b8b658e9a08d26a67 // Before [GW19, GW3]		// After [GW19, GW3]
	"gateway28": "b6e1350a62be2dfe195119910af6faba7590cac3d15656a59ec7c85d86991932", // ID - e1201cd6f5aa071595d8235b609c2e26886b3ef33119e1c787adb30bdd115737 // Before [GW20, GW4]		// After [GW20, GW4]
	"gateway29": "f22008f10c9b1f86b86a3d9ffc6a71154b15561ce5726ec68328c424365c3ed3", // ID - e9bfa910eed28ace8f5adac05a589283c27d44dcc691674113400c3120ddc876 // Before [GW21, GW5]		// After [GW21, GW5]
	"gateway30": "fdf89b3b7cc26173edfea6614b6a696d89042f7051b3644d8bc469a72961c885", // ID - f7a419c082602154061ed19ee4c0d984963a61b514aa32479956441396e02f3a // Before [GW23, GW7]		// After [GW23, GW7]
	"gateway31": "00b9a3b260e0ebd51f13e7adcc3c4aa7851ca35257cbc160c238c9355534e949", // ID - fe8fbb0e83f243fd7f184c24d7a008d1c89d5d3c137e9158977a78fb19e7869e // Before [GW23, GW7]		// After [GW23, GW7]
	// Gateway 32 is added later. It is between GW15 and GW16. Sort of like GW15.5
	"gateway32": "54b81c2cc94a8a10be4820716ddcddf834b24d1359bdedbf5a5a78251175793b", // ID - 79f1dfc58999bc9a1a3cb9f6cc1b8b3109b6e21350cd85c4641ab9a64907f4b0 // Before []				// [GW7, GW22] inc GW32
}
var providersKeyMap = map[string]string{
	"provider0": "ea9a44d5aa53b4714efb7df4aed727ea0cf68b7ed18ac3d36ac2c90f262daf5f", // ID - 3f3bb8d3768a56b0d718e01f29a491dcdbf91e5fc7193e948689d001a22099b6
	"provider1": "0d90700579ab17bcf579ccf904a9911ff5e6f4b9d5a450d1c1aef41e56736de0", // ID - 56651e4cf52c36b56498df851a52e6d95172f399d86b64d8ac3b69c573087f10
	"provider2": "f2754a52c0fb15e3be023346ccab3919a7f0687d876356810398799410280d57", // ID - f79f39161ed74c86d27ac21d98728ffbfd8ddd7ea5a5c5dbbb411b47162d3494
}

const clientKey = "72dd0be8b35fac690d0e763ce13326d9512c81c664b2a0a143bfb87bde5fc195" // ID - 87cd9ced77cb602a83b80a883f4f52d14901279fefb0ff40e7816f960e083f66

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
		// Topup
		privKeys := make([]string, 0)
		for _, gw := range gws {
			privKey, ok := gatewaysKeyMap[gw]
			if !ok {
				panic("gateway not found in local map")
			}
			privKeys = append(privKeys, privKey)
		}
		err = topup(localLotusAP, token, acct, privKeys)
		if err != nil {
			panic(err)
		}
		// Get ips
		for i, gw := range gws {
			cmd = exec.Command("docker", "inspect", gw, "--format", "{{.NetworkSettings.Networks.shared.IPAddress}}")
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
		pvds := strings.Split(string(stdout[:len(stdout)-1]), "\n")
		// topup
		privKeys := make([]string, 0)
		for _, pvd := range pvds {
			privKey, ok := providersKeyMap[pvd]
			if !ok {
				panic("provider not found in local map")
			}
			privKeys = append(privKeys, privKey)
		}
		err = topup(localLotusAP, token, acct, privKeys)
		if err != nil {
			panic(err)
		}
		// Get ips
		for i, pvd := range pvds {
			cmd = exec.Command("docker", "inspect", pvd, "--format", "{{.NetworkSettings.Networks.shared.IPAddress}}")
			stdout, err = cmd.Output()
			if err != nil {
				panic(err)
			}
			ip := string(stdout[:len(stdout)-1])
			res = fmt.Sprintf("%v;%v,%v,%v,%v", res, pvd, "6465616135313132656636333864653962393132306366363537333664656465", privKeys[i], ip)
		}
		fmt.Println(res)
	} else {
		// Must be client
		err := topup(localLotusAP, token, acct, []string{clientKey})
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v;%v\n", token, clientKey)
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

// The following helper method is used to topup a new filecoin account with 100 filecoins of balance
func topup(localLotusAP string, token string, superAcct string, privKeyStrs []string) error {
	// Get API
	var api v0api.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + token}}
	closer, err := jsonrpc.NewMergeClient(context.Background(), localLotusAP, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		return err
	}
	defer closer()

	mainAddress, err := address.NewFromString(superAcct)
	if err != nil {
		return err
	}

	cids := make([]cid2.Cid, 0)
	// Send messages
	for _, privKeyStr := range privKeyStrs {
		privKey, err := hex.DecodeString(privKeyStr)
		if err != nil {
			return err
		}
		pubKey := crypto.PublicKey(privKey)
		address1, err := address.NewSecp256k1Address(pubKey)
		if err != nil {
			return err
		}
		// Get amount
		amt, err := types.ParseFIL("100")
		if err != nil {
			return err
		}
		msg := &types.Message{
			To:     address1,
			From:   mainAddress,
			Value:  types.BigInt(amt),
			Method: 0,
		}
		signedMsg, err := fillMsg(mainAddress, &api, msg)
		if err != nil {
			return err
		}
		// Send request to lotus
		cid, err := api.MpoolPush(context.Background(), signedMsg)
		if err != nil {
			return err
		}
		cids = append(cids, cid)
	}

	// Finally check receipts
	for _, cid := range cids {
		receipt := waitReceipt(&cid, &api)
		if receipt.ExitCode != 0 {
			return errors.New("Transaction fail to execute")
		}
	}
	return nil
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
