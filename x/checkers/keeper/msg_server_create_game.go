package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
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

	black, err := sdk.AccAddressFromBech32(msg.Black)
	if err != nil {
		return nil, errors.Wrapf(err, types.ErrInvalidBlack.Error(), msg.Black)
	}

	red, err := sdk.AccAddressFromBech32(msg.Red)
	if err != nil {
		return nil, errors.Wrapf(err, types.ErrInvalidRed.Error(), msg.Red)
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(err, types.ErrGameInvalidCreator.Error(), creator)
	}

	newId := si.NextId
	si.NextId += 1

	game := rules.New()
	sg := types.StoredGame{
		Index: fmt.Sprintf("%d", newId),
		Board: game.String(),
		Turn:  "b",
		Red:   red.String(),
		Black: black.String(),
	}

	if err := sg.Validate(); err != nil {
		return nil, err
	}

	k.Keeper.SetSystemInfo(ctx, si)
	k.Keeper.SetStoredGame(ctx, sg)

	return &types.MsgCreateGameResponse{GameIndex: sg.Index}, nil
}
