package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peppelan/checkers/x/checkers/rules"
	"strconv"
)

const TypeMsgPlayMove = "play_move"

var _ sdk.Msg = &MsgPlayMove{}

func NewMsgPlayMove(creator string, gameIndex string, fromX uint64, fromY uint64, toX uint64, toY uint64) *MsgPlayMove {
	return &MsgPlayMove{
		Creator:   creator,
		GameIndex: gameIndex,
		FromX:     fromX,
		FromY:     fromY,
		ToX:       toX,
		ToY:       toY,
	}
}

func (msg *MsgPlayMove) Route() string {
	return RouterKey
}

func (msg *MsgPlayMove) Type() string {
	return TypeMsgPlayMove
}

func (msg *MsgPlayMove) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPlayMove) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPlayMove) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.GameIndex == "" {
		return sdkerrors.Wrapf(ErrInvalidGameIndex, "gameIndex is missing")
	}

	gameIndex, err := strconv.ParseUint(msg.GameIndex, 10, 64)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidGameIndex, "gameIndex is invalid (%s)", err)
	}

	if gameIndex < DefaultIndex {
		return sdkerrors.Wrapf(ErrInvalidGameIndex, "gameIndex (%s) is below the default (%d)", gameIndex, DefaultIndex)
	}

	if msg.FromX >= rules.BOARD_DIM {
		return sdkerrors.Wrapf(ErrInvalidPositionIndex, "fromX is out of range (%d)", msg.FromX)
	}

	if msg.FromY >= rules.BOARD_DIM {
		return sdkerrors.Wrapf(ErrInvalidPositionIndex, "fromY is out of range (%d)", msg.FromY)
	}

	if msg.ToX >= rules.BOARD_DIM {
		return sdkerrors.Wrapf(ErrInvalidPositionIndex, "toX is out of range (%d)", msg.ToX)
	}

	if msg.ToY >= rules.BOARD_DIM {
		return sdkerrors.Wrapf(ErrInvalidPositionIndex, "toY is out of range (%d)", msg.ToY)
	}

	if msg.FromX == msg.ToX && msg.FromY == msg.ToY {
		return sdkerrors.Wrapf(ErrMoveAbsent, "x (%d) and y (%d)", msg.FromX, msg.FromY)
	}

	return nil
}
