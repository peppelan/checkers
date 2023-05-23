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
	uCtx := sdk.UnwrapSDKContext(ctx)

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
	}, k.GetAllStoredGame(uCtx))

	systemInfo, found := k.GetSystemInfo(uCtx)
	require.True(t, found)
	assert.EqualValues(t, types.SystemInfo{NextId: 2}, systemInfo)

	events := sdk.StringifyEvents(uCtx.EventManager().ABCIEvents())
	require.Equal(t, 1, len(events))

	assert.EqualValues(t, sdk.StringEvent{
		Type: "new-game-created",
		Attributes: []sdk.Attribute{
			{Key: "creator", Value: alice},
			{Key: "game-index", Value: "1"},
			{Key: "black", Value: bob},
			{Key: "red", Value: carol},
		},
	}, events[0])
}
