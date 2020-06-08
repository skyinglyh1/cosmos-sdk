package types

const (
	QueryHeader        = "header"
	QueryCurrentHeight = "current_height"
	QueryKeyHeights    = "key_heights"
	QueryKeyHeight     = "key_height"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryHeaderParams struct {
	ChainId uint64
	Height  uint32
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryHeaderParams(chainId uint64, height uint32) QueryHeaderParams {
	return QueryHeaderParams{ChainId: chainId, Height: height}
}

type QueryCurrentHeightParams struct {
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryCurrentHeightParams(chainId uint64) QueryCurrentHeightParams {
	return QueryCurrentHeightParams{ChainId: chainId}
}

type QueryKeyHeightsParams struct {
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryKeyHeightsParams(chainId uint64) QueryKeyHeightsParams {
	return QueryKeyHeightsParams{ChainId: chainId}
}

type QueryKeyHeightParams struct {
	ChainId uint64
	Height  uint32
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryKeyHeightParams(chainId uint64, height uint32) QueryKeyHeightParams {
	return QueryKeyHeightParams{ChainId: chainId, Height: height}
}
