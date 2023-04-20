package server_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v4/server"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

func TestDatabase_GetChainId(t *testing.T) {
	database, err := server.NewDatabase(fxtypes.GetDefaultNodeHome(), dbm.GoLevelDBBackend)
	require.NoError(t, err)
	if database == nil {
		return
	}
	id, err := database.GetChainId()
	require.NoError(t, err)
	require.NotNil(t, id)
	height, err := database.GetBlockHeight()
	require.NoError(t, err)
	require.NotNil(t, height)
	nodeInfo, err := database.GetNodeInfo()
	require.NoError(t, err)
	require.NotNil(t, nodeInfo)
	_, err = database.CurrentPlan()
	require.NoError(t, err)
	validators, err := database.GetValidators()
	require.NoError(t, err)
	require.NotNil(t, validators)
	database.Close()
}