package fcrlotusmgr

import "math/big"

type FCRLotusMgr interface {
	Start()

	Shutdown()

	CreatePaymentChannel(prvKey string, toID string, amt *big.Int) (string, error)

	TopupPaymentChannel(prvKey string, chAddr string, amt *big.Int) error

	SettlePaymentChannel(prvKey string, chAddr string, voucher string) error

	// If settling or settled, balance, recipient address, error
	CheckPaymentChannel(chAddr string) (bool, *big.Int, string, error)

	GetCostToCreate(fromAddr string, toAddr string, amt *big.Int) (*big.Int, error)

	GetCostToSettle(fromAddr string, chAddr string) (*big.Int, error)

	GetPaymentChannelCreationBlock(chAddr string) (*big.Int, error)

	GetPaymentChannelSettlementBlock(chAddr string) (*big.Int, error)
}
