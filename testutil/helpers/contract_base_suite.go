package helpers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

type ContractBaseSuite struct {
	require  *require.Assertions
	signer   *Signer
	contract common.Address
}

func NewContractBaseSuite(require *require.Assertions, signer *Signer) *ContractBaseSuite {
	return &ContractBaseSuite{
		require:  require,
		signer:   signer,
		contract: common.Address{},
	}
}

func (s *ContractBaseSuite) WithContract(addr common.Address) {
	s.contract = addr
}

func (s *ContractBaseSuite) WithSigner(signer *Signer) {
	s.signer = signer
}

func (s *ContractBaseSuite) HexAddress() common.Address {
	return s.signer.Address()
}

func (s *ContractBaseSuite) AccAddress() sdk.AccAddress {
	return s.signer.AccAddress()
}
