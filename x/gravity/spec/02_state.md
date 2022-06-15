<!--
order: 2
-->

# State

## Keys

* ValidatorAddressByOrchestratorAddress `0xe  + orchestrator.Bytes()`                                                      -> `validator.Bytes()`
* EthAddressByValidatorKey              `0x1  + validator.Bytes()`                                                         -> `[]byte(ethereumAddress)`
* ValidatorByEthAddressKey              `0x2  + []byte(ethereumAddress)`                                                   -> `validator.Bytes()`
* ValsetRequestKey                      `0x3  + sdk.Uint64ToBigEndian(nonce)`                                              -> `k.cdc.MustMarshal(&Valset)`
* ValsetConfirmKey                      `0x4  + sdk.Uint64ToBigEndian(nonce) + validator.Bytes()`                          -> `k.cdc.MustMarshal(&MsgValsetConfirm)`
* OracleAttestationKey                  `0x5  + sdk.Uint64ToBigEndian(nonce) + claimHash`                                  -> `k.cdc.MustMarshal(Attestation)`
* OutgoingTxPoolKey                     `0x6  + sdk.Uint64ToBigEndian(outgoingTransferTxId)`                               -> `k.cdc.MustMarshal(&OutgoingTransferTx)`
* SecondIndexOutgoingTxFeeKey           `0x7  + []byte(tokenContract) + fee.Amount.BigInt().FillBytes(amount)`             -> `k.cdc.MustMarshal(&IDSet)`
* OutgoingTxBatchKey                    `0x8  + []byte(tokenContract) + sdk.Uint64ToBigEndian(nonce)`                      -> `k.cdc.MustMarshal(&OutgoingTxBatch)`
* OutgoingTxBatchBlockKey               `0x9  + sdk.Uint64ToBigEndian(blockHeight)`                                        -> `k.cdc.MustMarshal(&OutgoingTxBatch)`
* BatchConfirmKey                       `0xa  + []byte(tokenContract) + sdk.Uint64ToBigEndian(nonce) + validator.Bytes()`  -> `k.cdc.MustMarshal(&MsgConfirmBatch)`
* LastEventNonceByValidatorKey          `0xb  + validator.Bytes()`                                                         -> `sdk.Uint64ToBigEndian(nonce)`
* LastObservedEventNonceKey             `0xc`                                                                              -> `sdk.Uint64ToBigEndian(nonce)`
* LastTxPoolIDKey                       `0xd  + []byte("lastTxPoolId")`                                                    -> `sdk.Uint64ToBigEndian(id)`
* LastOutgoingBatchIDKey                `0xd  + []byte("lastBatchId")`                                                     -> `sdk.Uint64ToBigEndian(id)`
* DenomToERC20Key                       `0xf  + []byte(denom)`                                                             -> `[]byte(tokenContract)` 
* ERC20ToDenomKey                       `0x10 + []byte(tokenContract)`                                                     -> `[]byte(denom)` 
* LastSlashedValsetNonce                `0x11 `                                                                            -> `sdk.Uint64ToBigEndian(nonce)` 
* LatestValsetNonce                     `0x12 `                                                                            -> `sdk.Uint64ToBigEndian(nonce)`
* LastSlashedBatchBlock                 `0x13 `                                                                            -> `sdk.Uint64ToBigEndian(blockHeight)`
* LastUnBondingBlockHeight              `0x14 `                                                                            -> `sdk.Uint64ToBigEndian(blockHeight)`
* LastObservedEthereumBlockHeightKey    `0x15 `                                                                            -> `k.cdc.MustMarshal(lastObservedEthereumBlockHeight)`
* LastObservedValsetKey                 `0x16 `                                                                            -> `k.cdc.MustMarshal(&Valset)`
* LastEventBlockHeightByValidatorKey    `0x18 + validator.Bytes()`                                                         -> `sdk.Uint64ToBigEndian(blockHeight)`