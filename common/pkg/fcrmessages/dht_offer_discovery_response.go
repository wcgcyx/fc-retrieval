package fcrmessages

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type dhtOfferDiscoveryResponseJson struct {
	Contacted         []string `json:"contacted"`
	Responses         []string `json:"responses"`
	Nonce             int64    `json:"nonce"`
	RefundAccountAddr string   `json:"refund_account_addr"`
	RefundVoucher     string   `json:"refund_voucher"`
}

func EncodeDHTOfferDiscoveryResponse(
	contacted map[string]FCRMessage,
	nonce int64,
	refundAccountAddr string,
	refundVoucher string,
) (*FCRMessage, error) {
	contactedStr := make([]string, 0)
	responsesStr := make([]string, 0)
	for id, resp := range contacted {
		contactedStr = append(contactedStr, id)
		data, err := resp.ToBytes()
		if err != nil {
			return nil, err
		}
		responsesStr = append(responsesStr, hex.EncodeToString(data))
	}
	body, err := json.Marshal(dhtOfferDiscoveryResponseJson{
		Contacted:         contactedStr,
		Responses:         responsesStr,
		Nonce:             nonce,
		RefundAccountAddr: refundAccountAddr,
		RefundVoucher:     refundVoucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(DHTOfferDiscoveryResponseType, body), nil
}

func DecodeDHTOfferDiscoveryResponse(fcrMsg *FCRMessage) (
	map[string]FCRMessage,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != DHTOfferDiscoveryResponseType {
		return nil, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", DHTOfferDiscoveryResponseType, fcrMsg.GetMessageType())
	}
	msg := dhtOfferDiscoveryResponseJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, "", "", err
	}
	if len(msg.Contacted) != len(msg.Responses) {
		return nil, 0, "", "", fmt.Errorf("Contacted length %v mismatches response length %v", len(msg.Contacted), len(msg.Responses))
	}
	contacted := make(map[string]FCRMessage)
	for i := 0; i < len(msg.Contacted); i++ {
		data, err := hex.DecodeString(msg.Responses[i])
		if err != nil {
			return nil, 0, "", "", err
		}
		resp, err := FromBytes(data)
		if err != nil {
			return nil, 0, "", "", err
		}
		_, ok := contacted[msg.Contacted[i]]
		if ok {
			return nil, 0, "", "", fmt.Errorf("Node %v appears at least twice in the response", msg.Contacted[i])
		}
		contacted[msg.Contacted[i]] = *resp
	}
	return contacted, msg.Nonce, msg.RefundAccountAddr, msg.RefundVoucher, nil
}
