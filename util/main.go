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
var gatewaysKeyMap = map[string]string{
	"gateway0":  "dd34d72efd4730e73d7825b6138d2986d1c1e8722d3a4d59e247a4217fec16ee", // ID - 007160c839e60ebbd6a104a87a20e366277c5d00cd3c24ca453716b8a6ed889c
	"gateway1":  "a12a33f3d89d4b884138d5ee39c1e7d7ebded826453abca76e88b348bf4968b5", // ID - 07cd6375b4cd3eb4a45d2896ba1eaae9ea786382bd32dce0335607409a9b156d
	"gateway2":  "61dc2992a36571627f393158f9fee7e197bee997a22c1b72cac80788b1abf1ac", // ID - 1191275e532d85d0d3ac760e999f6c09aa9bba0ccd7db263b2ba0cc76d034460
	"gateway3":  "4ac6e88c8cc058558fda8fcc79b9000b7c5f1cbb734669960f2e26732936099e", // ID - 19f81363bbe352407e3a5060d4b72d710947f400a3185848bcfaf568dcee38a9
	"gateway4":  "30db4070a659c16f2b097f0946d58ec6cc1408cd2f28db34fcb647abbf5fa3db", // ID - 214f442cc624005321158c19d793fd63207441cd2fd79be79229d0ef971cbe6b
	"gateway5":  "7351931bbfd68c8c4bada67f3c9fafe621fff79f8e2e2eba23ad450a168ed70d", // ID - 28b172c99067f28fbd3e065fdc43ca0f238b3c4ee4b5513faa821d8c70fcacf6
	"gateway6":  "dca5e724b372cf6ec8e8b4efcece55daedf217d6bd586fa7250dcdec314f22ef", // ID - 333841f38ab3385e5794ccbc036abd6fae7120a84b1cafc59ef9ce85921fb9dd
	"gateway7":  "56a45ed8d18ea777f099aef2e9486c709ca61a21af701794c6d2217134c366ec", // ID - 3888cdb4af23ee3348f073ad999127788015cc4b24bc57e8ba51b0dda1606fe7
	"gateway8":  "86952e14f8f4cc3c4c3b1d4fe6a7e06d2ef21c7d7d5a77f48d6ca147e3f8a857", // ID - 46561045d56d241a3215d1268ac07603b5461c2e8a8e5ee972acbcbd0e18de08
	"gateway9":  "9fab680ac419e35b4d389c3c27bbfa865bc0d69969862555b9b1be606b7616d8", // ID - 468627830cf8d00ab92cd8e7097babf93499d7ce6e46ef1eeef5b156d8361379
	"gateway10": "9869cb2d369c4e55a8131afd5e193abedbba54c097f920a283d2d1924dd21bfc", // ID - 5397ea26fc5521fcedafa4ee98a4934c0717f08014d3c429a708aa997f599a47
	"gateway11": "85744beca232ba5203aa85cb44f1b95a5eedc567c4dce4c64626d04f3c7a2739", // ID - 585e6baf1763b06eef32b59c4f24675481807051011333a97d6becf3b8291fc9
	"gateway12": "c86bff3fe95d3ea327e09b42eb748da0fcae2cccad698742363423b8ce07b13e", // ID - 64cd860b7335189b2841f8e230f2a051fef12f6379722e8d459b124c38e775d4
	"gateway13": "68993f9eae34e02b588d109c16580df76ff29d577075e26cf1a8bb689692f073", // ID - 6ac71d2abe0ad5dc79045244a6cbc768179b30b5c276f819403c6dbe58093dd3
	"gateway14": "930ef1d6c65dbd3ee4fef31a95a9ad7545bd37724e230f92f4d9735526c101d1", // ID - 74a00891168be328e57928d13c46d43e9ebbc935f0357a13a995a4fab4c4f9d9
	"gateway15": "40bd41fb09ce7a1129e45cf55800c5c0e72004c0cdcde5e33226314313516957", // ID - 790ae917a8fc54d2bda5705b81e372ce4d217428957bb9f4e8ea06691a8ce558
	"gateway16": "7d26c07fb28722a5299920f5c37f5c2e4e55f6b2faca49574ac3700c9b806a1e", // ID - 8383406aeee226a460e98226f84fd134b40b8029cc1b45ea27c00697e4a583b8
	"gateway17": "dada231f41f7f0bb27a6bcf2987c6fbd0ebe5bbd5bd471c5bcdbc85ac83ed2e7", // ID - 8d5253e264c74e35dce9b389738686478c061e1fe2fca6f264d2e887b4d5b972
	"gateway18": "165ad659ee6b98acdd7c4901b9ac8ec12f826e83191c17a04034c2fa214a27d1", // ID - 92ccd60c9bbe23973e3fc6cda3b4fafd520d5af9f13f6eac1ac9d82e263599bc
	"gateway19": "8ffc4780945d2e7066a68afa79d85595fc9161fab5260d19f8a0aea5f8b7a178", // ID - 9838d0415c77f03769a862526029010ee72866e64f1230459612100dac0c8b68
	"gateway20": "ed0969fa664b9e07e1d193821d99281d0d48793072bc2824a493ba6ae0da682c", // ID - a10dcc06cd77f901d4385a1f639a3cba957dd82b83f216ff833a9c20946e48a2
	"gateway21": "1b4a9e8f9a39ad2afd22d608267a5b60420f279ec3bcdf334bce82afd493bfec", // ID - acc46adb9f3ee347098a8186761f105f08d07e2b8d5c824ec19c31d49e940123
	"gateway22": "936e6a2613066e49b9c5a56a2bc76063079fa1d405b25c3656dc12db9f652ea9", // ID - b046737e63934838cd323bf77535f8a5767dbd723c894e2d7610855c1ad6e970
	"gateway23": "c092df233858e2125add9241e037beaf4220b15600d72b4f4d4b30ccbd933861", // ID - b9ac3a21004be7a217f28d648ea1f1d43142b89c4ba6d632487ab58a1d6a8263
	"gateway24": "87dce26dccdac09842301fc56f0bd3e3cf9129e15e90199d28987675363540c5", // ID - c140dae4df5bf865665b5a6ccba9a9767a49f9ec6c78210beac93f755eb12334
	"gateway25": "1bbd21cd8738c5be459b3b0a37a251ed8ca56a24858a9ae6f431032ec0c34d8f", // ID - cc4698d1ee3eef7adc05750083d5701dc8b012ad137190746a9a21206b290e9d
	"gateway26": "8d5f51ade31ff752f81156cae4b44db885a3acf008b1c9acf23e75cf5bc3ba0e", // ID - d52d129d09b331e99397f2236669c14a0b9b294b2e0bc594dd492773cd249ac1
	"gateway27": "d2b783c6cca688a41c664fee322cb4685b32083937def36deb8df8a9c796670c", // ID - d845d7bb84aa044f13d2a13007667aec581d11beae49a029e0aaabb845132f84
	"gateway28": "c0f0d4971ac06bf2c8dc538ed967278e6f5494936a9c5dbcd3ca6960d90765b9", // ID - e468682606494e0b9a12f1cb088d2cec8b6c5442c5aba3f1f1ce3e6a495585fe
	"gateway29": "5a9fe074cc5260dba62c812772f8f828efb0cfbd09270c8ae53fa289c3eb59d2", // ID - eaed3aa6fefdbbab8b13191ef09fcc296c3a48ab17beb573a06cd2a9def0cb41
	"gateway30": "264da2220cfc3b1455d761921d39694ed035866558085b207c624607586da936", // ID - f453ebc11de95e36852b86cf8832555982ee3e6423eabd9a6ec19e0ca4e31cca
	"gateway31": "32529e528a62f73ba78a7fb33e1b2fab4002228abdde1c08f7af4ae8701b0b0e", // ID - fdb31714a5726e5027df9c5fb9c3c947bad50977948324c9f4c209f24e0dfa0d
	"gateway32": "f166353658d0c08a41d6d230246312a9fed5076dfae228c56ed392e6a31b413e", // ID - 810ec37dc67672e38dc4de8ce7078a448b3177ac06d124b3c8e44564c1acd559
}
var providersKeyMap = map[string]string{
	"provider0": "d4e1eb66860a21521ca32a9fff6caf802ac457b0df1553069639441959032230", // ID - 3e299f2dc239bd63dc3fbd96326d7e03d43146a0356d45bc131cb787e42303d5
	"provider1": "4b8420bc090bc07a7e2b1510ccf2bb097d69d7c0c2753881f456f56f52ff50fc", // ID - 510f7ea057c25e86661646b3fcbcb5b2863b8799401bf38fc2ed094f3e177471
	"provider2": "6835c1e1f7fa86f33b44b4943a53283a927b7ee23af614a850be1aafc95599ca", // ID - f3b5b3be53fab51dc02474e28359cb7c8cd7a6fa36c7bac65d27c1b7d392f4bc
}

const clientKey = "ce83e7fcadaab477a750bcb3d900e07a977193be70f0682591e2f80a812e61bf" // ID - 755d637b4b01d44505b7f826da71a8b0ce99c52380ff1a79c7d4084b6f76ce55

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
