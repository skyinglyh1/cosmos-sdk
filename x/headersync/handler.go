package headersync

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/headersync/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/headersync/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSyncGenesisParam:
			return handleMsgGenesisHeader(ctx, k, msg)
		case types.MsgSyncHeadersParam:
			return handleMsgBlockHeaders(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

func handleMsgGenesisHeader(ctx sdk.Context, k keeper.Keeper, msg types.MsgSyncGenesisParam) (*sdk.Result, error) {

	//err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	//ctx.BlockHeader()
	err := k.SyncGenesisHeader(ctx, msg.GenesisHeader)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// Handle MsgMultiSend.
func handleMsgBlockHeaders(ctx sdk.Context, k keeper.Keeper, msg types.MsgSyncHeadersParam) (*sdk.Result, error) {
	err := k.SyncBlockHeaders(ctx, msg.Headers)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
