package fcrmessages

import (
	"encoding/json"
	"fmt"
)

type ackJson struct {
	ACK   bool   `json:"ack"`
	Nonce int64  `json:"nonce"`
	Data  string `json:"data"`
}

func EncodeACK(
	ack bool,
	nonce int64,
	data string,
) (*FCRMessage, error) {
	body, err := json.Marshal(ackJson{
		ACK:   ack,
		Nonce: nonce,
		Data:  data,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(ACKType, body), nil
}

func DecodeACK(fcrMsg *FCRMessage) (
	bool,
	int64,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != ACKType {
		return false, 0, "", fmt.Errorf("Message type mismatch, expect %v, got %v", ACKType, fcrMsg.GetMessageType())
	}
	msg := ackJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return false, 0, "", err
	}
	return msg.ACK, msg.Nonce, msg.Data, nil
}
