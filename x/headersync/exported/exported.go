package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	polycommon "github.com/cosmos/cosmos-sdk/x/headersync/poly-utils/common"
	polytype "github.com/cosmos/cosmos-sdk/x/headersync/poly-utils/core/types"
)

// DelegationI delegation bond for a delegated proof of stake system
type HeaderSyncKeeper interface {
	ProcessHeader(ctx sdk.Context, header *polytype.Header) error
	GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*polytype.Header, error)
	GetHeaderByHash(ctx sdk.Context, chainId uint64, hash polycommon.Uint256) (*polytype.Header, error)
	GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, error)
}
