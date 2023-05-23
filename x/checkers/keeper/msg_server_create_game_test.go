package keeper_test

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/peppelan/checkers/testutil/keeper"
	"github.com/peppelan/checkers/x/checkers"
	"github.com/peppelan/checkers/x/checkers/keeper"
	"github.com/peppelan/checkers/x/checkers/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupMsgServerCreateGame(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	return keeper.NewMsgServerImpl(*k), *k, sdk.WrapSDKContext(ctx)
}

func TestCreateGame(t *testing.T) {
	msgServer, _, ctx := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(ctx, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	require.Nil(t, err)
	require.EqualValues(t, &types.MsgCreateGameResponse{
		GameIndex: "1",
	}, createResponse)
}

func TestCreate1GameHasSaved(t *testing.T) {
	// arrange
	msgServer, k, ctx := setupMsgServerCreateGame(t)

	// act
	createResponse, err := msgServer.CreateGame(ctx, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})

	// assert
	require.Nil(t, err)
	assert.EqualValues(t, &types.MsgCreateGameResponse{
		GameIndex: "1",
	}, createResponse)

	assert.EqualValues(t, []types.StoredGame{
		{
			Index: "1",
			Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
			Turn:  "b",
			Black: bob,
			Red:   carol,
		},
	}, k.GetAllStoredGame(sdk.UnwrapSDKContext(ctx)))

	systemInfo, found := k.GetSystemInfo(sdk.UnwrapSDKContext(ctx))
	require.True(t, found)
	assert.EqualValues(t, types.SystemInfo{NextId: 2}, systemInfo)
}
