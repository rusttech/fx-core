package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

var _ sdk.Msg = &MsgUpdateParams{}

const (
	TypeMsgUpdateParams = "update_params"
)

// Route returns the MsgUpdateParams message route.
func (m *MsgUpdateParams) Route() string { return types.ModuleName }

// Type returns the MsgUpdateParams message type.
func (m *MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errortypes.ErrUnauthorized.Wrapf("invalid authority address: %s", err.Error())
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	return nil
}
