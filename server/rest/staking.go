package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"
)

// RegisterStakingRESTRoutes
// Deprecated
func RegisterStakingRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := WithHTTPDeprecationHeaders(rtr)
	registerStakingQueryRoutes(clientCtx, r)
}

func registerStakingQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/delegations",
		delegatorDelegationsHandlerFn(clientCtx),
	).Methods("GET")

	// Get all unbonding delegations from a delegator
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/unbonding_delegations",
		delegatorUnbondingDelegationsHandlerFn(clientCtx),
	).Methods("GET")

	// Get all staking txs (i.e msgs) from a delegator
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/txs",
		delegatorTxsHandlerFn(clientCtx),
	).Methods("GET")

	// Query all validators that a delegator is bonded to
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/validators",
		delegatorValidatorsHandlerFn(clientCtx),
	).Methods("GET")

	// Query a validator that a delegator is bonded to
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/validators/{validatorAddr}",
		delegatorValidatorHandlerFn(clientCtx),
	).Methods("GET")

	// Query a delegation between a delegator and a validator
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/delegations/{validatorAddr}",
		delegationHandlerFn(clientCtx),
	).Methods("GET")

	// Query all unbonding delegations between a delegator and a validator
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}",
		unbondingDelegationHandlerFn(clientCtx),
	).Methods("GET")

	// Query redelegations (filters in query params)
	r.HandleFunc(
		"/staking/redelegations",
		redelegationsHandlerFn(clientCtx),
	).Methods("GET")

	// Get all validators
	r.HandleFunc(
		"/staking/validators",
		validatorsHandlerFn(clientCtx),
	).Methods("GET")

	// Get a single validator info
	r.HandleFunc(
		"/staking/validators/{validatorAddr}",
		validatorHandlerFn(clientCtx),
	).Methods("GET")

	// Get all delegations to a validator
	r.HandleFunc(
		"/staking/validators/{validatorAddr}/delegations",
		validatorDelegationsHandlerFn(clientCtx),
	).Methods("GET")

	// Get all unbonding delegations from a validator
	r.HandleFunc(
		"/staking/validators/{validatorAddr}/unbonding_delegations",
		validatorUnbondingDelegationsHandlerFn(clientCtx),
	).Methods("GET")

	// Get HistoricalInfo at a given height
	r.HandleFunc(
		"/staking/historical_info/{height}",
		historicalInfoHandlerFn(clientCtx),
	).Methods("GET")

	// Get the current state of the staking pool
	r.HandleFunc(
		"/staking/pool",
		poolHandlerFn(clientCtx),
	).Methods("GET")

	// Get the current staking parameter values
	r.HandleFunc(
		"/staking/parameters",
		stakingParamsHandlerFn(clientCtx),
	).Methods("GET")
}

// HTTP request handler to query a delegator delegations
func delegatorDelegationsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return queryDelegator(clientCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDelegatorDelegations))
}

// HTTP request handler to query a delegator unbonding delegations
func delegatorUnbondingDelegationsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return queryDelegator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDelegatorUnbondingDelegations))
}

// HTTP request handler to query all staking txs (msgs) from a delegator
func delegatorTxsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var typesQuerySlice []string

		vars := mux.Vars(r)
		delegatorAddr := vars["delegatorAddr"]

		if _, err := sdk.AccAddressFromBech32(delegatorAddr); CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		typesQuery := r.URL.Query().Get("type")
		trimmedQuery := strings.TrimSpace(typesQuery)

		if len(trimmedQuery) != 0 {
			typesQuerySlice = strings.Split(trimmedQuery, " ")
		}

		noQuery := len(typesQuerySlice) == 0
		isBondTx := contains(typesQuerySlice, "bond")
		isUnbondTx := contains(typesQuerySlice, "unbond")
		isRedTx := contains(typesQuerySlice, "redelegate")

		var (
			txs     []*sdk.SearchTxsResult
			actions []string
		)

		// For each case, we search txs for both:
		// - legacy messages: their Type() is a custom string, e.g. "delegate"
		// - service Msgs: their Type() is their FQ method name, e.g. "/cosmos.staking.v1beta1.MsgDelegate"
		// and we combine the results.
		switch {
		case isBondTx:
			actions = append(actions, types.TypeMsgDelegate)

		case isUnbondTx:
			actions = append(actions, types.TypeMsgUndelegate)

		case isRedTx:
			actions = append(actions, types.TypeMsgBeginRedelegate)

		case noQuery:
			actions = append(actions, types.TypeMsgDelegate)
			actions = append(actions, types.TypeMsgUndelegate)
			actions = append(actions, types.TypeMsgBeginRedelegate)

		default:
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for _, action := range actions {
			foundTxs, errQuery := queryTxs(clientCtx, action, delegatorAddr)
			if CheckInternalServerError(w, errQuery) {
				return
			}

			txs = append(txs, foundTxs)
		}

		res, err := clientCtx.LegacyAmino.MarshalJSON(txs)
		if CheckInternalServerError(w, err) {
			return
		}

		PostProcessResponseBare(w, clientCtx, res)
	}
}

// HTTP request handler to query an unbonding-delegation
func unbondingDelegationHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return queryBonds(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryUnbondingDelegation))
}

// HTTP request handler to query redelegations
func redelegationsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params types.QueryRedelegationParams

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		bechDelegatorAddr := r.URL.Query().Get("delegator")
		bechSrcValidatorAddr := r.URL.Query().Get("validator_from")
		bechDstValidatorAddr := r.URL.Query().Get("validator_to")

		if len(bechDelegatorAddr) != 0 {
			delegatorAddr, err := sdk.AccAddressFromBech32(bechDelegatorAddr)
			if CheckBadRequestError(w, err) {
				return
			}

			params.DelegatorAddr = delegatorAddr
		}

		if len(bechSrcValidatorAddr) != 0 {
			srcValidatorAddr, err := sdk.ValAddressFromBech32(bechSrcValidatorAddr)
			if CheckBadRequestError(w, err) {
				return
			}

			params.SrcValidatorAddr = srcValidatorAddr
		}

		if len(bechDstValidatorAddr) != 0 {
			dstValidatorAddr, err := sdk.ValAddressFromBech32(bechDstValidatorAddr)
			if CheckBadRequestError(w, err) {
				return
			}

			params.DstValidatorAddr = dstValidatorAddr
		}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRedelegations), bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

// HTTP request handler to query a delegation
func delegationHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return queryBonds(clientCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDelegation))
}

// HTTP request handler to query all delegator bonded validators
func delegatorValidatorsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return queryDelegator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDelegatorValidators))
}

// HTTP request handler to get information from a currently bonded validator
func delegatorValidatorHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return queryBonds(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDelegatorValidator))
}

// HTTP request handler to query list of validators
func validatorsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := ParseHTTPArgsWithLimit(r, 0)
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		status := r.FormValue("status")
		// These are query params that were available in =<0.39. We show a nice
		// error message for this breaking change.
		if status == "bonded" || status == "unbonding" || status == "unbonded" {
			err := fmt.Errorf("cosmos sdk v0.40 introduces a breaking change on this endpoint:"+
				" instead of querying using `?status=%s`, please use `status=BOND_STATUS_%s`. For more"+
				" info, please see our REST endpoint migration guide at %s", status, strings.ToUpper(status), DeprecationURL)

			if CheckBadRequestError(w, err) {
				return
			}

		}

		if status == "" {
			status = types.BondStatusBonded
		}

		params := types.NewQueryValidatorsParams(page, limit, status)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidators)

		res, height, err := clientCtx.QueryWithData(route, bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

// HTTP request handler to query the validator information from a given validator address
func validatorHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return queryValidator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidator))
}

// HTTP request handler to query all unbonding delegations from a validator
func validatorDelegationsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return queryValidator(clientCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorDelegations))
}

// HTTP request handler to query all unbonding delegations from a validator
func validatorUnbondingDelegationsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return queryValidator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorUnbondingDelegations))
}

// HTTP request handler to query historical info at a given height
func historicalInfoHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		heightStr := vars["height"]

		height, err := strconv.ParseInt(heightStr, 10, 64)
		if err != nil || height < 0 {
			WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Must provide non-negative integer for height: %v", err))
			return
		}

		params := types.QueryHistoricalInfoRequest{Height: height}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckInternalServerError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryHistoricalInfo), bz)
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

// HTTP request handler to query the pool information
func poolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPool), nil)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

// HTTP request handler to query the staking params values
func stakingParamsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

func queryDelegator(clientCtx client.Context, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryDelegatorParams(delegatorAddr)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(endpoint, bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

// contains checks if the a given query contains one of the tx types
func contains(stringSlice []string, txType string) bool {
	for _, word := range stringSlice {
		if word == txType {
			return true
		}
	}

	return false
}

// queries staking txs
func queryTxs(clientCtx client.Context, action string, delegatorAddr string) (*sdk.SearchTxsResult, error) {
	page := 1
	limit := 100
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, action),
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, delegatorAddr),
	}

	return authtx.QueryTxsByEvents(clientCtx, events, page, limit, "")
}

func queryBonds(clientCtx client.Context, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]
		bech32validator := vars["validatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if CheckBadRequestError(w, err) {
			return
		}

		validatorAddr, err := sdk.ValAddressFromBech32(bech32validator)
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.QueryDelegatorValidatorRequest{DelegatorAddr: delegatorAddr.String(), ValidatorAddr: validatorAddr.String()}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(endpoint, bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}

func queryValidator(clientCtx client.Context, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32validatorAddr := vars["validatorAddr"]

		_, page, limit, err := ParseHTTPArgsWithLimit(r, 0)
		if CheckBadRequestError(w, err) {
			return
		}

		validatorAddr, err := sdk.ValAddressFromBech32(bech32validatorAddr)
		if CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryValidatorParams(validatorAddr, page, limit)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(endpoint, bz)
		if CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		PostProcessResponse(w, clientCtx, res)
	}
}
