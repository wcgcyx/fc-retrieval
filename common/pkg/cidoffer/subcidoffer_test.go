package cidoffer

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
)

func TestNewSubCIDOfferWithGet(t *testing.T) {
	aCid1, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid1)
	aCid2, err := cid.NewContentID(Cid2Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid2)
	aCid3, err := cid.NewContentID(Cid3Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid3)
	cids := []cid.ContentID{*aCid1, *aCid2, *aCid3}
	price := big.NewInt(100)
	expiry := int64(10)
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	subOffer, err := offer.GenerateSubCIDOffer(aCid1)
	assert.Empty(t, err)
	assert.NotEmpty(t, subOffer)

	assert.Equal(t, "testprovider", subOffer.GetProviderID())
	assert.Equal(t, aCid1, subOffer.GetSubCID())
	assert.Equal(t, offer.merkleRoot, subOffer.GetMerkleRoot())
	assert.Equal(t, price, subOffer.GetPrice())
	assert.Equal(t, expiry, subOffer.GetExpiry())
	assert.Equal(t, qos, subOffer.GetQoS())
	assert.Equal(t, offer.GetSignature(), subOffer.GetSignature())
	p, err := subOffer.GetMerkleProof().ToBytes()
	assert.Empty(t, err)
	assert.Equal(t, []byte{0x22, 0x41, 0x41, 0x41, 0x41, 0x58,
		0x31, 0x73, 0x69, 0x52, 0x6c, 0x63, 0x76, 0x56, 0x48,
		0x63, 0x33, 0x55, 0x6c, 0x6c, 0x71, 0x5a, 0x45, 0x63,
		0x31, 0x63, 0x56, 0x46, 0x72, 0x64, 0x31, 0x4d, 0x72,
		0x59, 0x6d, 0x5a, 0x48, 0x4e, 0x56, 0x42, 0x71, 0x51,
		0x32, 0x6c, 0x34, 0x61, 0x58, 0x6c, 0x31, 0x59, 0x54,
		0x59, 0x79, 0x64, 0x47, 0x67, 0x77, 0x53, 0x57, 0x70,
		0x6f, 0x4e, 0x44, 0x63, 0x76, 0x63, 0x7a, 0x30, 0x69,
		0x4c, 0x43, 0x4a, 0x74, 0x51, 0x6b, 0x4d, 0x78, 0x59,
		0x33, 0x56, 0x46, 0x63, 0x6c, 0x64, 0x57, 0x57, 0x48,
		0x52, 0x56, 0x64, 0x30, 0x45, 0x76, 0x55, 0x55, 0x52,
		0x48, 0x55, 0x30, 0x5a, 0x68, 0x54, 0x56, 0x52, 0x34,
		0x56, 0x6a, 0x4e, 0x46, 0x64, 0x32, 0x46, 0x57, 0x64,
		0x6d, 0x78, 0x51, 0x4f, 0x44, 0x68, 0x6f, 0x4f, 0x56,
		0x46, 0x6a, 0x63, 0x32, 0x31, 0x6a, 0x50, 0x53, 0x4a,
		0x64, 0x41, 0x41, 0x41, 0x41, 0x42, 0x56, 0x73, 0x78,
		0x4c, 0x44, 0x46, 0x64, 0x22}, p)
}

func TestSubOfferHasExpired(t *testing.T) {
	aCid1, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid1)
	aCid2, err := cid.NewContentID(Cid2Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid2)
	aCid3, err := cid.NewContentID(Cid3Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid3)
	cids := []cid.ContentID{*aCid1, *aCid2, *aCid3}
	price := big.NewInt(100)
	expiry := time.Now().Add(12 * time.Hour).Unix()
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	subOffer, err := offer.GenerateSubCIDOffer(aCid1)
	assert.Empty(t, err)
	assert.NotEmpty(t, subOffer)

	assert.False(t, subOffer.HasExpired())
	subOffer.expiry = time.Now().Add(-12 * time.Hour).Unix()
	assert.True(t, subOffer.HasExpired())
}

func TestSubOfferVerify(t *testing.T) {
	aCid1, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid1)
	aCid2, err := cid.NewContentID(Cid2Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid2)
	aCid3, err := cid.NewContentID(Cid3Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid3)
	cids := []cid.ContentID{*aCid1, *aCid2, *aCid3}
	price := big.NewInt(100)
	expiry := time.Now().Add(12 * time.Hour).Unix()
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	err = offer.Sign(PrivKey)
	assert.Empty(t, err)

	subOffer, err := offer.GenerateSubCIDOffer(aCid1)
	assert.Empty(t, err)
	assert.NotEmpty(t, subOffer)

	err = subOffer.Verify(PubKey)
	assert.Empty(t, err)

	err = subOffer.Verify(PubKeyWrong)
	assert.NotEmpty(t, err)

	assert.Empty(t, subOffer.VerifyMerkleProof())
	subOffer.subCID = aCid2
	assert.NotEmpty(t, subOffer.VerifyMerkleProof())
}

func TestSerializationSubOffer(t *testing.T) {
	aCid1, err := cid.NewContentID(Cid1Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid1)
	aCid2, err := cid.NewContentID(Cid2Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid2)
	aCid3, err := cid.NewContentID(Cid3Str)
	assert.Empty(t, err)
	assert.NotEmpty(t, aCid3)
	cids := []cid.ContentID{*aCid1, *aCid2, *aCid3}
	price := big.NewInt(100)
	expiry := int64(10)
	qos := uint64(5)
	offer, err := NewCIDOffer("testprovider", cids, price, expiry, qos)
	assert.Empty(t, err)
	subOffer, err := offer.GenerateSubCIDOffer(aCid1)
	assert.Empty(t, err)
	assert.NotEmpty(t, subOffer)
	p, err := subOffer.ToBytes()
	assert.Empty(t, err)
	assert.Equal(t, []byte{0x7b, 0x22, 0x70, 0x72, 0x6f, 0x76,
		0x69, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x3a,
		0x22, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x76,
		0x69, 0x64, 0x65, 0x72, 0x22, 0x2c, 0x22, 0x73, 0x75,
		0x62, 0x5f, 0x63, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x62,
		0x61, 0x67, 0x61, 0x36, 0x65, 0x61, 0x34, 0x73, 0x65,
		0x61, 0x71, 0x6f, 0x33, 0x62, 0x74, 0x35, 0x6c, 0x74,
		0x73, 0x33, 0x37, 0x34, 0x35, 0x65, 0x6a, 0x72, 0x71,
		0x33, 0x73, 0x68, 0x35, 0x6a, 0x61, 0x72, 0x61, 0x75,
		0x75, 0x77, 0x6a, 0x6e, 0x34, 0x76, 0x6c, 0x73, 0x6e,
		0x63, 0x37, 0x66, 0x36, 0x33, 0x66, 0x75, 0x78, 0x63,
		0x71, 0x33, 0x70, 0x73, 0x69, 0x6f, 0x32, 0x68, 0x71,
		0x22, 0x2c, 0x22, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65,
		0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x22, 0x3a, 0x22, 0x39,
		0x35, 0x38, 0x35, 0x61, 0x34, 0x32, 0x63, 0x33, 0x31,
		0x31, 0x61, 0x61, 0x31, 0x62, 0x37, 0x37, 0x38, 0x31,
		0x38, 0x66, 0x39, 0x33, 0x64, 0x63, 0x64, 0x31, 0x66,
		0x62, 0x39, 0x63, 0x63, 0x38, 0x66, 0x36, 0x38, 0x38,
		0x36, 0x35, 0x34, 0x37, 0x34, 0x36, 0x34, 0x37, 0x30,
		0x63, 0x61, 0x37, 0x33, 0x38, 0x62, 0x36, 0x31, 0x36,
		0x37, 0x34, 0x35, 0x63, 0x35, 0x35, 0x31, 0x31, 0x33,
		0x22, 0x2c, 0x22, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65,
		0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x22, 0x3a, 0x22,
		0x32, 0x32, 0x34, 0x31, 0x34, 0x31, 0x34, 0x31, 0x34,
		0x31, 0x35, 0x38, 0x33, 0x31, 0x37, 0x33, 0x36, 0x39,
		0x35, 0x32, 0x36, 0x63, 0x36, 0x33, 0x37, 0x36, 0x35,
		0x36, 0x34, 0x38, 0x36, 0x33, 0x33, 0x33, 0x35, 0x35,
		0x36, 0x63, 0x36, 0x63, 0x37, 0x31, 0x35, 0x61, 0x34,
		0x35, 0x36, 0x33, 0x33, 0x31, 0x36, 0x33, 0x35, 0x36,
		0x34, 0x36, 0x37, 0x32, 0x36, 0x34, 0x33, 0x31, 0x34,
		0x64, 0x37, 0x32, 0x35, 0x39, 0x36, 0x64, 0x35, 0x61,
		0x34, 0x38, 0x34, 0x65, 0x35, 0x36, 0x34, 0x32, 0x37,
		0x31, 0x35, 0x31, 0x33, 0x32, 0x36, 0x63, 0x33, 0x34,
		0x36, 0x31, 0x35, 0x38, 0x36, 0x63, 0x33, 0x31, 0x35,
		0x39, 0x35, 0x34, 0x35, 0x39, 0x37, 0x39, 0x36, 0x34,
		0x34, 0x37, 0x36, 0x37, 0x37, 0x37, 0x35, 0x33, 0x35,
		0x37, 0x37, 0x30, 0x36, 0x66, 0x34, 0x65, 0x34, 0x34,
		0x36, 0x33, 0x37, 0x36, 0x36, 0x33, 0x37, 0x61, 0x33,
		0x30, 0x36, 0x39, 0x34, 0x63, 0x34, 0x33, 0x34, 0x61,
		0x37, 0x34, 0x35, 0x31, 0x36, 0x62, 0x34, 0x64, 0x37,
		0x38, 0x35, 0x39, 0x33, 0x33, 0x35, 0x36, 0x34, 0x36,
		0x36, 0x33, 0x36, 0x63, 0x36, 0x34, 0x35, 0x37, 0x35,
		0x37, 0x34, 0x38, 0x35, 0x32, 0x35, 0x36, 0x36, 0x34,
		0x33, 0x30, 0x34, 0x35, 0x37, 0x36, 0x35, 0x35, 0x35,
		0x35, 0x35, 0x32, 0x34, 0x38, 0x35, 0x35, 0x33, 0x30,
		0x35, 0x61, 0x36, 0x38, 0x35, 0x34, 0x35, 0x36, 0x35,
		0x32, 0x33, 0x34, 0x35, 0x36, 0x36, 0x61, 0x34, 0x65,
		0x34, 0x36, 0x36, 0x34, 0x33, 0x32, 0x34, 0x36, 0x35,
		0x37, 0x36, 0x34, 0x36, 0x64, 0x37, 0x38, 0x35, 0x31,
		0x34, 0x66, 0x34, 0x34, 0x36, 0x38, 0x36, 0x66, 0x34,
		0x66, 0x35, 0x36, 0x34, 0x36, 0x36, 0x61, 0x36, 0x33,
		0x33, 0x32, 0x33, 0x31, 0x36, 0x61, 0x35, 0x30, 0x35,
		0x33, 0x34, 0x61, 0x36, 0x34, 0x34, 0x31, 0x34, 0x31,
		0x34, 0x31, 0x34, 0x31, 0x34, 0x32, 0x35, 0x36, 0x37,
		0x33, 0x37, 0x38, 0x34, 0x63, 0x34, 0x34, 0x34, 0x36,
		0x36, 0x34, 0x32, 0x32, 0x22, 0x2c, 0x22, 0x70, 0x72,
		0x69, 0x63, 0x65, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30,
		0x22, 0x2c, 0x22, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79,
		0x22, 0x3a, 0x31, 0x30, 0x2c, 0x22, 0x71, 0x6f, 0x73,
		0x22, 0x3a, 0x35, 0x2c, 0x22, 0x73, 0x69, 0x67, 0x6e,
		0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x3a, 0x22, 0x22,
		0x7d}, p)
	subOffer2 := SubCIDOffer{}
	err = subOffer2.FromBytes(p)
	assert.Empty(t, err)
	assert.Equal(t, subOffer.providerID, subOffer2.providerID)
	assert.Equal(t, subOffer.subCID, subOffer2.subCID)
	assert.Equal(t, subOffer.merkleRoot, subOffer2.merkleRoot)
	assert.Equal(t, subOffer.price, subOffer2.price)
	assert.Equal(t, subOffer.expiry, subOffer2.expiry)
	assert.Equal(t, subOffer.qos, subOffer2.qos)
	assert.Equal(t, subOffer.signature, subOffer2.signature)
	p1, err := subOffer.merkleProof.ToBytes()
	assert.Empty(t, err)
	p2, err := subOffer2.merkleProof.ToBytes()
	assert.Empty(t, err)
	assert.Equal(t, p1, p2)
	err = subOffer2.FromBytes([]byte{})
	assert.NotEmpty(t, err)
}
