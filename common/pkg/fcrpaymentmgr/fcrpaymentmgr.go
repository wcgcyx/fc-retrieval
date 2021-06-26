package fcrpaymentmgr

import "math/big"

type FCRPaymentMgr interface {
	Start()

	Shutdown()

	// For outbound payment

	Create(recipientID string, amount *big.Int) error

	Topup(recipientID string, amount *big.Int) error

	Pay(recipientID string, lane uint64, amount *big.Int) (string, bool, bool, error)

	ReceiveRefund(recipientID string, voucher string) (*big.Int, error)

	GetOutboundChStatus(recipientID string) (string, *big.Int, *big.Int, error)

	GetCostToCreate(recipientID string) (*big.Int, error)

	CheckRecipientSettlementValidity(recipientID string) (bool, error)

	// For inbound payment

	Settle(senderID string) error

	Receive(senderID string, voucher string) (*big.Int, uint64, error)

	Refund(senderID string, lane uint64, amount *big.Int) (string, error)

	GetInboundChStatus(senderID string) (string, *big.Int, *big.Int, error)

	GetCostToSettle(senderID string) (*big.Int, error)

	CheckSettlementValidity(senderID string) (bool, error)
}
