package keeper

import (
	"context"
	"fmt"
	"github.com/peppelan/checkers/x/checkers/rules"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peppelan/checkers/x/checkers/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	si, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}

	newId := si.NextId
	si.NextId += 1

	game := rules.New()
	sg := types.StoredGame{
		Index:  fmt.Sprintf("%d", newId),
		Board:  game.String(),
		Turn:   rules.PieceStrings[game.Turn],
		Red:    msg.Red,
		Black:  msg.Black,
		Winner: rules.PieceStrings[rules.NO_PLAYER],
	}

	if err := sg.Validate(); err != nil {
		return nil, err
	}

	k.Keeper.SetSystemInfo(ctx, si)
	k.Keeper.SetStoredGame(ctx, sg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.GameCreatedEventType,
		sdk.NewAttribute(types.GameCreatedEventCreator, msg.Creator),
		sdk.NewAttribute(types.GameCreatedEventGameIndex, sg.Index),
		sdk.NewAttribute(types.GameCreatedEventBlack, sg.Black),
		sdk.NewAttribute(types.GameCreatedEventRed, sg.Red),
	))

	return &types.MsgCreateGameResponse{GameIndex: sg.Index}, nil
}
