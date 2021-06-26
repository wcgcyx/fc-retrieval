package fcroffermgr

import (
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
)

type FCROfferMgr interface {
	Start()

	Shutdown()

	AddOffer(offer *cidoffer.CIDOffer) error

	GetOffers(cID *cid.ContentID) ([]cidoffer.CIDOffer, error)

	// TODO: Paging
	ListOffers() ([]cidoffer.CIDOffer, error)
}
