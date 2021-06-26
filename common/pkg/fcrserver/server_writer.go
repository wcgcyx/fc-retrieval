package fcrserver

import (
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
)

type FCRServerWriter interface {
	Write(msg *fcrmessages.FCRMessage, timeout time.Duration) error
}
