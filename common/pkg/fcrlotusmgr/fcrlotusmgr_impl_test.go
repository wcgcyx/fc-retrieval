package fcrlotusmgr

import (
	"context"
	"math/big"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
)

const (
	PrvKey       = "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29"
	LotusAPIAddr = "http://127.0.0.1:1234/rpc/v0"
	LotusToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.lYpknouVX5M_BJOEbHZrdcHdxHkfu0ih1W0NCTFlJz0"
)

type mockLotusAPI struct {
	chainReadObj func(ctx context.Context, obj cid.Cid) ([]byte, error)

	gasEstimateFeeCap func(ctx context.Context, msg *types.Message, maxqueueblks int64, tsk types.TipSetKey) (types.BigInt, error)

	gasEstimateGasLimit func(ctx context.Context, msg *types.Message, tsk types.TipSetKey) (int64, error)

	gasEstimateGasPremium func(ctx context.Context, nblocksincl uint64, sender address.Address, gaslimit int64, tsk types.TipSetKey) (types.BigInt, error)

	mpoolGetNonce func(ctx context.Context, addr address.Address) (uint64, error)

	mpoolPush func(ctx context.Context, smsg *types.SignedMessage) (cid.Cid, error)

	stateAccountKey func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error)

	stateGetActor func(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error)

	stateGetReceipt func(ctx context.Context, msg cid.Cid, tsk types.TipSetKey) (*types.MessageReceipt, error)
}

func (m *mockLotusAPI) ChainReadObj(ctx context.Context, obj cid.Cid) ([]byte, error) {
	return m.chainReadObj(ctx, obj)
}

func (m *mockLotusAPI) GasEstimateFeeCap(ctx context.Context, msg *types.Message, maxqueueblks int64, tsk types.TipSetKey) (types.BigInt, error) {
	return m.gasEstimateFeeCap(ctx, msg, maxqueueblks, tsk)
}

func (m *mockLotusAPI) GasEstimateGasLimit(ctx context.Context, msg *types.Message, tsk types.TipSetKey) (int64, error) {
	return m.gasEstimateGasLimit(ctx, msg, tsk)
}

func (m *mockLotusAPI) GasEstimateGasPremium(ctx context.Context, nblocksincl uint64, sender address.Address, gaslimit int64, tsk types.TipSetKey) (types.BigInt, error) {
	return m.gasEstimateGasPremium(ctx, nblocksincl, sender, gaslimit, tsk)
}

func (m *mockLotusAPI) MpoolGetNonce(ctx context.Context, addr address.Address) (uint64, error) {
	return m.mpoolGetNonce(ctx, addr)
}

func (m *mockLotusAPI) MpoolPush(ctx context.Context, smsg *types.SignedMessage) (cid.Cid, error) {
	return m.mpoolPush(ctx, smsg)
}

func (m *mockLotusAPI) StateAccountKey(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return m.stateAccountKey(ctx, addr, tsk)
}

func (m *mockLotusAPI) StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error) {
	return m.stateGetActor(ctx, actor, tsk)
}

func (m *mockLotusAPI) StateGetReceipt(ctx context.Context, msg cid.Cid, tsk types.TipSetKey) (*types.MessageReceipt, error) {
	return m.stateGetReceipt(ctx, msg, tsk)
}

func TestCreate(t *testing.T) {
	mock := mockLotusAPI{
		gasEstimateFeeCap: func(ctx context.Context, msg *types.Message, maxqueueblks int64, tsk types.TipSetKey) (types.BigInt, error) {
			return types.NewInt(100643), nil
		},
		gasEstimateGasPremium: func(ctx context.Context, nblocksincl uint64, sender address.Address, gaslimit int64, tsk types.TipSetKey) (types.BigInt, error) {
			return types.NewInt(99589), nil
		},
		gasEstimateGasLimit: func(ctx context.Context, msg *types.Message, tsk types.TipSetKey) (int64, error) {
			return 3823323, nil
		},
		mpoolGetNonce: func(ctx context.Context, addr address.Address) (uint64, error) {
			return 1, nil
		},
		mpoolPush: func(ctx context.Context, smsg *types.SignedMessage) (cid.Cid, error) {
			return cid.Parse("baga6ea4seaqesauho7j2thfi4g4u5zbnhn2okd74s2igpvc2lsb7rrsfstoy4by")
		},
		stateGetReceipt: func(ctx context.Context, msg cid.Cid, tsk types.TipSetKey) (*types.MessageReceipt, error) {
			return &types.MessageReceipt{
				ExitCode: 0,
				Return:   []byte{130, 67, 0, 236, 7, 85, 2, 111, 159, 23, 63, 130, 221, 152, 12, 77, 31, 30, 188, 113, 217, 40, 153, 105, 233, 205, 79},
				GasUsed:  3823323,
			}, nil
		},
	}

	mgr := NewFCRLotusMgrImpl(LotusAPIAddr, LotusToken, func(authToken, lotusAPIAddr string) (LotusAPI, jsonrpc.ClientCloser, error) {
		return &mock, nil, nil
	})

	address, err := mgr.CreatePaymentChannel(PrvKey, "t1hn3o5excejl2uyea7efs3licozuycghzpdiikjy", big.NewInt(1000000))
	assert.Empty(t, err)
	assert.Equal(t, "f2n6prop4c3wmayti7d26hdwjitfu6ttkp5qhu6ni", address)
}

func TestTopup(t *testing.T) {
	mock := mockLotusAPI{
		gasEstimateFeeCap: func(ctx context.Context, msg *types.Message, maxqueueblks int64, tsk types.TipSetKey) (types.BigInt, error) {
			return types.NewInt(101130), nil
		},
		gasEstimateGasPremium: func(ctx context.Context, nblocksincl uint64, sender address.Address, gaslimit int64, tsk types.TipSetKey) (types.BigInt, error) {
			return types.NewInt(100076), nil
		},
		gasEstimateGasLimit: func(ctx context.Context, msg *types.Message, tsk types.TipSetKey) (int64, error) {
			return 481468, nil
		},
		mpoolGetNonce: func(ctx context.Context, addr address.Address) (uint64, error) {
			return 1, nil
		},
		mpoolPush: func(ctx context.Context, smsg *types.SignedMessage) (cid.Cid, error) {
			return cid.Parse("baga6ea4seaqesauho7j2thfi4g4u5zbnhn2okd74s2igpvc2lsb7rrsfstoy4by")
		},
		stateGetReceipt: func(ctx context.Context, msg cid.Cid, tsk types.TipSetKey) (*types.MessageReceipt, error) {
			return &types.MessageReceipt{
				ExitCode: 0,
				Return:   []byte{},
				GasUsed:  481468,
			}, nil
		},
	}

	mgr := NewFCRLotusMgrImpl(LotusAPIAddr, LotusToken, func(authToken, lotusAPIAddr string) (LotusAPI, jsonrpc.ClientCloser, error) {
		return &mock, nil, nil
	})
	err := mgr.TopupPaymentChannel(PrvKey, "t1hn3o5excejl2uyea7efs3licozuycghzpdiikjy", big.NewInt(1000000))
	assert.Empty(t, err)
}

func TestCheck(t *testing.T) {
	codeCID, err := cid.Parse("bafkqafdgnfwc6nbpobqxs3lfnz2gg2dbnzxgk3a")
	assert.Empty(t, err)
	headCID, err := cid.Parse("bafy2bzaceazcxcw4ggk66uh76z7k2owst7d6xwkmrycz3cvnjf6g5havlo5lw")
	assert.Empty(t, err)

	mock := mockLotusAPI{
		stateGetActor: func(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error) {
			return &types.Actor{
				Code:    codeCID,
				Head:    headCID,
				Nonce:   0,
				Balance: types.NewInt(1000000),
			}, nil
		},
		chainReadObj: func(ctx context.Context, obj cid.Cid) ([]byte, error) {
			return []byte{134, 67, 0, 233, 7, 67, 0, 235, 7, 64, 0, 0, 216, 42, 88, 39, 0, 1,
				113, 160, 228, 2, 32, 208, 155, 127, 152, 162, 62, 233, 213, 222, 219, 74, 189,
				156, 112, 213, 71, 154, 99, 246, 190, 166, 215, 206, 196, 31, 178, 197, 242, 181,
				207, 236, 66}, nil
		},
		stateAccountKey: func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
			return address.NewFromString("f1hn3o5excejl2uyea7efs3licozuycghzpdiikjy")
		},
	}

	mgr := NewFCRLotusMgrImpl(LotusAPIAddr, LotusToken, func(authToken, lotusAPIAddr string) (LotusAPI, jsonrpc.ClientCloser, error) {
		return &mock, nil, nil
	})
	settling, balance, recipient, err := mgr.CheckPaymentChannel("f2n6prop4c3wmayti7d26hdwjitfu6ttkp5qhu6ni")
	assert.Empty(t, err)
	assert.False(t, settling)
	assert.Equal(t, big.NewInt(1000000), balance)
	assert.Equal(t, "f1hn3o5excejl2uyea7efs3licozuycghzpdiikjy", recipient)
}

func TestUnimplemented(t *testing.T) {
	mock := mockLotusAPI{}
	mgr := NewFCRLotusMgrImpl(LotusAPIAddr, LotusToken, func(authToken, lotusAPIAddr string) (LotusAPI, jsonrpc.ClientCloser, error) {
		return &mock, nil, nil
	})
	err := mgr.SettlePaymentChannel("", "", "")
	assert.NotEmpty(t, err)
	err = mgr.CollectPaymentChannel("", "")
	assert.NotEmpty(t, err)
	_, err = mgr.GetCostToCreate("", "", nil)
	assert.NotEmpty(t, err)
	_, err = mgr.GetCostToSettle("", "")
	assert.NotEmpty(t, err)
	_, err = mgr.GetPaymentChannelCreationBlock("")
	assert.NotEmpty(t, err)
	_, err = mgr.GetPaymentChannelSettlementBlock("")
	assert.NotEmpty(t, err)
}
