package p2papi

import (
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
)

// EstablishmentRequester sends an establishment request.
func EstablishmentRequester(reader fcrserver.FCRServerReader, writer fcrserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	return nil, nil
}
