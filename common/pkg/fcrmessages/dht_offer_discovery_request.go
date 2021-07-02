package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
)

type dhtOfferDiscoveryRequestJson struct {
	NodeID                  string `json:"node_id"`
	PieceCID                string `json:"piece_cid"`
	Nonce                   int64  `json:"nonce"`
	NumDHT                  int64  `json:"num_dht"`
	MaxOfferRequestedPerDHT int64  `json:"max_offer_requested_per_dht"`
	AccountAddr             string `json:"account_addr"`
	Voucher                 string `json:"voucher"`
}

func EncodeDHTOfferDiscoveryRequest(
	NodeID string,
	pieceCID *cid.ContentID,
	nonce int64,
	numDHT int64,
	maxOfferRequestedPerDHT int64,
	accountAddr string,
	voucher string,
) (*FCRMessage, error) {
	body, err := json.Marshal(dhtOfferDiscoveryRequestJson{
		NodeID:                  NodeID,
		PieceCID:                pieceCID.ToString(),
		Nonce:                   nonce,
		NumDHT:                  numDHT,
		MaxOfferRequestedPerDHT: maxOfferRequestedPerDHT,
		AccountAddr:             accountAddr,
		Voucher:                 voucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(DHTOfferDiscoveryRequestType, body), nil
}

func DecodeDHTOfferDiscoveryRequest(fcrMsg *FCRMessage) (
	string,
	*cid.ContentID,
	int64,
	int64,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != DHTOfferDiscoveryRequestType {
		return "", nil, 0, 0, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", DHTOfferDiscoveryRequestType, fcrMsg.GetMessageType())
	}
	msg := dhtOfferDiscoveryRequestJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return "", nil, 0, 0, 0, "", "", err
	}
	pieceCID, err := cid.NewContentID(msg.PieceCID)
	if err != nil {
		return "", nil, 0, 0, 0, "", "", err
	}
	return msg.NodeID, pieceCID, msg.Nonce, msg.NumDHT, msg.MaxOfferRequestedPerDHT, msg.AccountAddr, msg.Voucher, nil
}
