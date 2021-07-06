package p2papi

import (
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
)

// OfferQueryRequester sends an offer query request.
func OfferQueryRequester(reader fcrserver.FCRServerResponseReader, writer fcrserver.FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error) {
	return nil, nil
}
