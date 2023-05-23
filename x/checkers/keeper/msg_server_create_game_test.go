package keeper_test

import (
	"github.com/peppelan/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateGame(t *testing.T) {
	msgServer, context := setupMsgServer(t)
	createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	require.Nil(t, err)
	require.EqualValues(t, &types.MsgCreateGameResponse{
		GameIndex: "123",
	}, createResponse)
}
