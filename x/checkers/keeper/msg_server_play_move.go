package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peppelan/checkers/x/checkers/rules"
	"github.com/peppelan/checkers/x/checkers/types"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the game
	sg, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
		return nil, types.ErrGameNotFound
	}

	game, err := sg.ParseGame()
	if err != nil {
		panic(err.Error())
	}

	// am I a player?
	if msg.Creator != sg.Red && msg.Creator != sg.Black {
		return nil, types.ErrCreatorNotPlayer
	}

	// is it my turn?
	var turnAddr string
	if game.Turn.Color == rules.BLACK {
		turnAddr = sg.Black
	} else {
		turnAddr = sg.Red
	}

	if turnAddr != msg.Creator {
		return nil, types.ErrNotPlayerTurn
	}

	captPos, err := game.Move(rules.Pos{
		X: int(msg.FromX),
		Y: int(msg.FromY),
	}, rules.Pos{
		X: int(msg.ToX),
		Y: int(msg.ToY),
	})

	if err != nil {
		return nil, types.ErrWrongMove
	}

	sg.Board = game.String()
	sg.Turn = rules.PieceStrings[game.Turn]
	k.SetStoredGame(ctx, sg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.MovePlayedEventType,
		sdk.NewAttribute(types.MovePlayedEventCreator, msg.Creator),
		sdk.NewAttribute(types.MovePlayedEventGameIndex, sg.Index),
		sdk.NewAttribute(types.MovePlayedEventCapturedX, fmt.Sprintf("%d", captPos.X)),
		sdk.NewAttribute(types.MovePlayedEventCapturedY, fmt.Sprintf("%d", captPos.Y)),
		sdk.NewAttribute(types.MovePlayedEventWinner, rules.PieceStrings[game.Winner()]),
	))

	return &types.MsgPlayMoveResponse{
		CapturedX: int32(captPos.X),
		CapturedY: int32(captPos.Y),
		Winner:    rules.PieceStrings[game.Winner()],
	}, nil
}
