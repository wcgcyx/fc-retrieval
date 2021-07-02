package fcrmessages

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
)

type offerPublishRequestJson struct {
	NodeID string `json:"node_id"`
	Nonce  int64  `json:"nonce"`
	Offer  string `json:"offer"`
}

func EncodeOfferPublishRequest(
	nodeID string,
	nonce int64,
	offer *cidoffer.CIDOffer,
) (*FCRMessage, error) {
	data, err := offer.ToBytes()
	if err != nil {
		return nil, err
	}
	offerStr := hex.EncodeToString(data)
	body, err := json.Marshal(offerPublishRequestJson{
		NodeID: nodeID,
		Nonce:  nonce,
		Offer:  offerStr,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(OfferPublishRequestType, body), nil
}

func DecodeOfferPublishRequest(fcrMsg *FCRMessage) (
	string,
	int64,
	*cidoffer.CIDOffer,
	error,
) {
	if fcrMsg.GetMessageType() != OfferPublishRequestType {
		return "", 0, nil, fmt.Errorf("Message type mismatch, expect %v, got %v", OfferPublishRequestType, fcrMsg.GetMessageType())
	}
	msg := offerPublishRequestJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return "", 0, nil, err
	}
	data, err := hex.DecodeString(msg.Offer)
	if err != nil {
		return "", 0, nil, err
	}
	offer := cidoffer.CIDOffer{}
	err = offer.FromBytes(data)
	if err != nil {
		return "", 0, nil, err
	}
	return msg.NodeID, msg.Nonce, &offer, nil
}
