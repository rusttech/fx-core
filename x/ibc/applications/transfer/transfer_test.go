package transfer_test

import (
	"github.com/functionx/fx-core/app/fxcore"
	ibctesting "github.com/functionx/fx-core/x/ibc/testing"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	"github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

var (
	zeroFee  = sdk.NewCoin(fxcore.MintDenom, sdk.ZeroInt())
	noRouter = ""
	noFeeStr = ""
)

type TransferTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	// setup between chainA and chainB
	clientA, clientB, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateTransferChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.UNORDERED)
	//	originalBalance := suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
	timeoutHeight := clienttypes.NewHeight(0, 110)

	coinToSendToB := sdk.NewCoin(fxcore.MintDenom, sdk.NewInt(100))

	// send from chainA to chainB
	msg := types.NewMsgTransfer(channelA.PortID, channelA.ID, coinToSendToB, suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, noRouter, zeroFee)

	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, clientB, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	fungibleTokenPacket := types.NewFungibleTokenPacketData(coinToSendToB.Denom, coinToSendToB.Amount.String(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), noRouter, noFeeStr)
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, timeoutHeight, 0)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainA, suite.chainB, clientA, clientB, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain B
	voucherDenomTrace := types.ParseDenomTrace(types.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), fxcore.MintDenom))
	balance := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())

	coinSentFromAToB := types.GetTransferCoin(channelB.PortID, channelB.ID, fxcore.MintDenom, 100)
	suite.Require().Equal(coinSentFromAToB, balance)

	// setup between chainB to chainC
	clientOnBForC, clientOnCForB, connOnBForC, connOnCForB := suite.coordinator.SetupClientConnections(suite.chainB, suite.chainC, exported.Tendermint)
	channelOnBForC, channelOnCForB := suite.coordinator.CreateTransferChannels(suite.chainB, suite.chainC, connOnBForC, connOnCForB, channeltypes.UNORDERED)

	// send from chainB to chainC
	msg = types.NewMsgTransfer(channelOnBForC.PortID, channelOnBForC.ID, coinSentFromAToB, suite.chainB.SenderAccount.GetAddress(), suite.chainC.SenderAccount.GetAddress().String(), timeoutHeight, 0, noRouter, zeroFee)

	err = suite.coordinator.SendMsg(suite.chainB, suite.chainC, clientOnCForB, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
	fullDenomPath := types.GetPrefixedDenom(channelOnCForB.PortID, channelOnCForB.ID, voucherDenomTrace.GetFullDenomPath())
	fungibleTokenPacket = types.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), coinSentFromAToB.Amount.String(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainC.SenderAccount.GetAddress().String(), noRouter, noFeeStr)
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnBForC.PortID, channelOnBForC.ID, channelOnCForB.PortID, channelOnCForB.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainB, suite.chainC, clientOnBForC, clientOnCForB, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	coinSentFromBToC := sdk.NewInt64Coin(types.ParseDenomTrace(fullDenomPath).IBCDenom(), 100)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromBToC.Denom)

	// check that the balance is updated on chainC
	suite.Require().Equal(coinSentFromBToC, balance)

	// check that balance on chain B is empty
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())

	// send from chainC back to chainB
	msg = types.NewMsgTransfer(channelOnCForB.PortID, channelOnCForB.ID, coinSentFromBToC, suite.chainC.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, noRouter, zeroFee)

	err = suite.coordinator.SendMsg(suite.chainC, suite.chainB, clientOnBForC, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
	fungibleTokenPacket = types.NewFungibleTokenPacketData(fullDenomPath, coinSentFromBToC.Amount.String(), suite.chainC.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), noRouter, noFeeStr)
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnCForB.PortID, channelOnCForB.ID, channelOnBForC.PortID, channelOnBForC.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainC, suite.chainB, clientOnCForB, clientOnBForC, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)

	// check that the balance on chainA returned back to the original state
	suite.Require().Equal(coinSentFromAToB, balance)

	// check that module account escrow address is empty
	escrowAddress := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), escrowAddress, fxcore.MintDenom)
	suite.Require().Equal(sdk.NewCoin(fxcore.MintDenom, sdk.ZeroInt()), balance)

	// check that balance on chain B is empty
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}