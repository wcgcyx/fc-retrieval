package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
)

type standardOfferDiscoveryRequestJson struct {
	NodeID            string `json:"node_id"`
	PieceCID          string `json:"piece_cid"`
	Nonce             int64  `json:"nonce"`
	MaxOfferRequested int64  `json:"max_offer_requested"`
	AccountAddr       string `json:"account_addr"`
	Voucher           string `json:"voucher"`
}

func EncodeStandardOfferDiscoveryRequest(
	NodeID string,
	pieceCID *cid.ContentID,
	nonce int64,
	maxOfferRequested int64,
	accountAddr string,
	voucher string,
) (*FCRMessage, error) {
	body, err := json.Marshal(standardOfferDiscoveryRequestJson{
		NodeID:            NodeID,
		PieceCID:          pieceCID.ToString(),
		Nonce:             nonce,
		MaxOfferRequested: maxOfferRequested,
		AccountAddr:       accountAddr,
		Voucher:           voucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(StandardOfferDiscoveryRequestType, body), nil
}

func DecodeStandardOfferDiscoveryRequest(fcrMsg *FCRMessage) (
	string,
	*cid.ContentID,
	int64,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != StandardOfferDiscoveryRequestType {
		return "", nil, 0, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", StandardOfferDiscoveryRequestType, fcrMsg.GetMessageType())
	}
	msg := standardOfferDiscoveryRequestJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return "", nil, 0, 0, "", "", err
	}
	pieceCID, err := cid.NewContentID(msg.PieceCID)
	if err != nil {
		return "", nil, 0, 0, "", "", err
	}
	return msg.NodeID, pieceCID, msg.Nonce, msg.MaxOfferRequested, msg.AccountAddr, msg.Voucher, nil
}
