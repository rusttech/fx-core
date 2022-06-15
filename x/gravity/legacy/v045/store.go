package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"

	v042gravity "github.com/functionx/fx-core/x/gravity/legacy/v042"
)

// MigrateStore performs in-place store migrations from v0.42 to v0.45.
// migrate data from gravity module
func MigrateStore(ctx sdk.Context, cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {

	// gravity 0x1 -> eth ? -
	// gravity ? -> eth 0x12 *

	// gravity 0x2 -> eth 0x13
	migratePrefix(gravityStore, ethStore, v042gravity.ValidatorByEthAddressKey, crosschaintypes.OracleAddressByExternalKey)

	// gravity 0xe -> eth 0x14
	migratePrefix(gravityStore, ethStore, v042gravity.ValidatorAddressByOrchestratorAddress, crosschaintypes.OracleAddressByBridgerKey)

	// gravity 0x3 -> eth 0x15
	migratePrefix(gravityStore, ethStore, v042gravity.ValsetRequestKey, crosschaintypes.OracleSetRequestKey)

	// gravity 0x4 -> eth 0x16
	migratePrefix(gravityStore, ethStore, v042gravity.ValsetConfirmKey, crosschaintypes.OracleSetConfirmKey)

	// gravity 0x5 -> eth 0x17
	migratePrefix(gravityStore, ethStore, v042gravity.OracleAttestationKey, crosschaintypes.OracleAttestationKey)

	// gravity 0x6 and 0x7 -> eth 0x18
	migrateOutgoingTxPool(cdc, gravityStore, ethStore)

	// gravity 0x8 -> eth 0x20
	migratePrefix(gravityStore, ethStore, v042gravity.OutgoingTxBatchKey, crosschaintypes.OutgoingTxBatchKey)

	// gravity 0x9 -> eth 0x21
	migratePrefix(gravityStore, ethStore, v042gravity.OutgoingTxBatchBlockKey, crosschaintypes.OutgoingTxBatchBlockKey)

	// gravity 0xa -> eth 0x22
	migratePrefix(gravityStore, ethStore, v042gravity.BatchConfirmKey, crosschaintypes.BatchConfirmKey)

	// gravity 0xb -> eth 0x23
	migratePrefix(gravityStore, ethStore, v042gravity.LastEventNonceByValidatorKey, crosschaintypes.LastEventNonceByValidatorKey)

	// gravity 0xc -> eth 0x24
	migratePrefix(gravityStore, ethStore, v042gravity.LastObservedEventNonceKey, crosschaintypes.LastObservedEventNonceKey)

	// gravity 0xd+"lastTxPoolId" -> eth 0x25+"lastTxPoolId"
	// gravity 0xd+"lastBatchId" -> eth 0x25+"lastBatchId"
	migratePrefix(gravityStore, ethStore, v042gravity.SequenceKeyPrefix, crosschaintypes.SequenceKeyPrefix)

	// gravity 0xf -> eth 0x26
	migratePrefix(gravityStore, ethStore, v042gravity.DenomToERC20Key, crosschaintypes.DenomToTokenKey)

	// gravity 0x10 -> eth 0x27
	migratePrefix(gravityStore, ethStore, v042gravity.ERC20ToDenomKey, crosschaintypes.TokenToDenomKey)

	// gravity 0x11 -> eth 0x28
	migratePrefix(gravityStore, ethStore, v042gravity.LastSlashedValsetNonce, crosschaintypes.LastSlashedOracleSetNonce)

	// gravity 0x12 -> eth 0x29
	migratePrefix(gravityStore, ethStore, v042gravity.LatestValsetNonce, crosschaintypes.LatestOracleSetNonce)

	// gravity 0x13 -> eth 0x30
	migratePrefix(gravityStore, ethStore, v042gravity.LastSlashedBatchBlock, crosschaintypes.LastSlashedBatchBlock)

	// gravity 0x14 -> eth 0x31 -

	// gravity 0x15 -> eth 0x32
	migratePrefix(gravityStore, ethStore, v042gravity.LastObservedEthereumBlockHeightKey, crosschaintypes.LastObservedBlockHeightKey)

	// gravity 0x16 -> eth 0x33
	migratePrefix(gravityStore, ethStore, v042gravity.LastObservedValsetKey, crosschaintypes.LastObservedOracleSetKey)

	// gravity 0x17 -> eth 0x34
	//migratePrefix(gravityStore, ethStore, v042gravity.IbcSequenceHeightKey, crosschaintypes.KeyIbcSequenceHeight)

	// gravity 0x18 -> eth 0x35
	migratePrefix(gravityStore, ethStore, v042gravity.LastEventBlockHeightByValidatorKey, crosschaintypes.LastEventBlockHeightByOracleKey)

	// gravity ? -> eth 0x36 -
	// gravity ? -> eth 0x37 -
	// gravity ? -> eth 0x38 *
	// gravity ? -> eth 0x39 *
}

func migratePrefix(gravityStore, ethStore sdk.KVStore, oldPrefix, newPrefix []byte) {
	oldStore := prefix.NewStore(gravityStore, oldPrefix)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()
	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		newStoreKey := append(newPrefix, oldStoreIter.Key()[len(oldPrefix):]...)
		ethStore.Set(newStoreKey, oldStoreIter.Value())
		oldStore.Delete(oldStoreIter.Key())
	}
}

func MigrateValidatorToOracle(ctx sdk.Context, cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore, stakingKeeper StakingKeeper) {

	chainOracle := new(crosschaintypes.ProposalOracle)
	totalPower := sdk.ZeroInt()

	oldStore := prefix.NewStore(gravityStore, v042gravity.ValidatorAddressByOrchestratorAddress)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()
	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		bridgerAddr := sdk.AccAddress(oldStoreIter.Key()[len(v042gravity.ValidatorAddressByOrchestratorAddress):])
		oracleAddress := sdk.AccAddress(oldStoreIter.Value())
		externalAddress := sdk.AccAddress(oldStore.Get(append(v042gravity.EthAddressByValidatorKey, oracleAddress.Bytes()...)))
		validator, found := stakingKeeper.GetValidator(ctx, oldStoreIter.Value())
		if !found {
			continue
		}
		oracle := crosschaintypes.Oracle{
			OracleAddress:     oracleAddress.String(),
			BridgerAddress:    bridgerAddr.String(),
			ExternalAddress:   externalAddress.String(),
			DelegateAmount:    sdk.NewInt(validator.ConsensusPower(sdk.DefaultPowerReduction)),
			StartHeight:       0,
			Jailed:            validator.Jailed,
			JailedHeight:      validator.UnbondingHeight,
			DelegateValidator: oracleAddress.String(),
			IsValidator:       true,
		}
		if !oracle.Jailed {
			totalPower = totalPower.Add(oracle.GetPower())
		}
		// SetOracle
		ethStore.Set(crosschaintypes.GetOracleKey(oracle.GetOracle()), cdc.MustMarshal(&oracle))
		oldStore.Delete(oldStoreIter.Key())

		chainOracle.Oracles = append(chainOracle.Oracles, oracleAddress.String())
	}

	// SetProposalOracle
	if len(chainOracle.Oracles) > 0 {
		ethStore.Set(crosschaintypes.ProposalOracleKey, cdc.MustMarshal(chainOracle))
	}
	// setLastTotalPower
	ethStore.Set(types.LastTotalPowerKey, cdc.MustMarshal(&sdk.IntProto{Int: totalPower}))
}

func migrateOutgoingTxPool(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {

	oldStore := prefix.NewStore(gravityStore, v042gravity.OutgoingTxPoolKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var transact crosschaintypes.OutgoingTransferTx
		cdc.MustUnmarshal(oldStoreIter.Value(), &transact)

		ethStore.Set(crosschaintypes.GetOutgoingTxPoolKey(transact.Fee, transact.Id), oldStoreIter.Value())
		oldStore.Delete(oldStoreIter.Key())
	}

	oldStore2 := prefix.NewStore(gravityStore, v042gravity.SecondIndexOutgoingTxFeeKey)

	oldStoreIter2 := oldStore2.Iterator(nil, nil)
	defer oldStoreIter2.Close()
	for ; oldStoreIter2.Valid(); oldStoreIter2.Next() {
		oldStore2.Delete(oldStoreIter2.Key())
	}
}