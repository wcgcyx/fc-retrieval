package fcrmessages

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
)

type standardOfferDiscoveryResponseJson struct {
	Offers            []string `json:"offers"`
	Nonce             int64    `json:"nonce"`
	RefundAccountAddr string   `json:"refund_account_addr"`
	RefundVoucher     string   `json:"refund_voucher"`
}

func EncodeStandardOfferDiscoveryResponse(
	offers []cidoffer.SubCIDOffer,
	nonce int64,
	refundAccountAddr string,
	refundVoucher string,
) (*FCRMessage, error) {
	offersStr := make([]string, 0)
	for _, offer := range offers {
		data, err := offer.ToBytes()
		if err != nil {
			return nil, err
		}
		offersStr = append(offersStr, hex.EncodeToString(data))
	}
	body, err := json.Marshal(standardOfferDiscoveryResponseJson{
		Offers:            offersStr,
		Nonce:             nonce,
		RefundAccountAddr: refundAccountAddr,
		RefundVoucher:     refundVoucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(StandardOfferDiscoveryResponseType, body), nil
}

func DecodeStandardOfferDiscoveryResponse(fcrMsg *FCRMessage) (
	[]cidoffer.SubCIDOffer,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != StandardOfferDiscoveryResponseType {
		return nil, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", StandardOfferDiscoveryResponseType, fcrMsg.GetMessageType())
	}
	msg := standardOfferDiscoveryResponseJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, "", "", err
	}
	offers := make([]cidoffer.SubCIDOffer, 0)
	offersStr := msg.Offers
	for _, offerStr := range offersStr {
		data, err := hex.DecodeString(offerStr)
		if err != nil {
			return nil, 0, "", "", err
		}
		offer := cidoffer.SubCIDOffer{}
		err = offer.FromBytes(data)
		if err != nil {
			return nil, 0, "", "", err
		}
		offers = append(offers, offer)
	}
	return offers, msg.Nonce, msg.RefundAccountAddr, msg.RefundVoucher, nil
}
