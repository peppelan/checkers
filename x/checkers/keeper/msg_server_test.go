package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/peppelan/checkers/testutil/keeper"
	"github.com/peppelan/checkers/x/checkers/keeper"
	"github.com/peppelan/checkers/x/checkers/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	k.SetSystemInfo(ctx, types.SystemInfo{NextId: 123})
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
