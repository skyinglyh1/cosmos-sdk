package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/headersync/internal/types"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryHeader(cliCtx context.CLIContext, queryRoute string, chainId uint64, height uint32) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryHeader),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryHeaderParams(chainId, height)),
	)
	return res, err
}

func QueryCurrentHeaderHeight(cliCtx context.CLIContext, queryRoute string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryCurrentHeight),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryCurrentHeightParams(chainId)),
	)

	return res, err
}
func QueryKeyHeights(cliCtx context.CLIContext, queryRoute string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryKeyHeights),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryKeyHeightsParams(chainId)),
	)

	return res, err
}

func QueryKeyHeight(cliCtx context.CLIContext, queryRoute string, chainId uint64, height uint32) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryKeyHeight),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryKeyHeightParams(chainId, height)),
	)

	return res, err
}
