package fcrserver

import (
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
)

type FCRServerReader interface {
	Read(timeout time.Duration) (*fcrmessages.FCRMessage, error)
}
