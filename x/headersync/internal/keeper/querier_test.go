package keeper_test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/headersync"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"

	keep "github.com/cosmos/cosmos-sdk/x/headersync/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/headersync/internal/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
	app, ctx := createTestApp(true)

	h0s, _ := hex.DecodeString(header0)
	err := app.HeaderSyncKeeper.SyncGenesisHeader(ctx, h0s)
	assert.Nil(t, err, "Sync genesis header fail")

	querier := keep.NewQuerier(app.HeaderSyncKeeper)

	// query synced Poly Chain header
	query := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", headersync.StoreKey, types.QueryHeader),
		Data: app.Codec().MustMarshalJSON(types.NewQueryHeaderParams(0, 0)),
	}
	headerBs, err := querier(ctx, []string{types.QueryHeader}, query)
	require.NoError(t, err)
	require.Equal(t, h0s, headerBs, "Synced Header 0 is not equal to the querier result")

	// query synced Poly Chain header height
	query = abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", headersync.StoreKey, types.QueryCurrentHeight),
		Data: app.Codec().MustMarshalJSON(types.NewQueryCurrentHeightParams(0)),
	}
	heightBs, err := querier(ctx, []string{types.QueryCurrentHeight}, query)
	require.NoError(t, err)
	bz, err := codec.MarshalJSONIndent(app.Codec(), uint32(0))
	require.NoError(t, err)
	require.Equal(t, bz, heightBs, "current synced heigh 0 is not equal to the querier result")

	// sync header 1 and header 100 of Poly Chain
	h1s, err := hex.DecodeString(header1)
	assert.Nil(t, err, "Header1 hex to header1 bytes error")
	h100s, err := hex.DecodeString(header100)
	assert.Nil(t, err, "Header100 hex to header100 bytes error")
	err = app.HeaderSyncKeeper.SyncBlockHeaders(ctx, [][]byte{h1s, h100s})
	assert.Nil(t, err, "Sync Poly Chain block headers fail")

	// query synced Poly Chain Key height
	query = abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", headersync.StoreKey, types.QueryKeyHeight),
		Data: app.Codec().MustMarshalJSON(types.NewQueryKeyHeightParams(0, 90)),
	}
	keyHeightBs, err := querier(ctx, []string{types.QueryKeyHeight}, query)
	require.NoError(t, err)
	bz, err = codec.MarshalJSONIndent(app.Codec(), uint32(0))
	require.NoError(t, err)
	require.Equal(t, bz, keyHeightBs, "key height of height:90 should be 0, yet not equal to the querier result")

	// query synced Poly Chain Key heights
	query = abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", headersync.StoreKey, types.QueryKeyHeights),
		Data: app.Codec().MustMarshalJSON(types.NewQueryKeyHeightsParams(0)),
	}
	keyHeightsBs, err := querier(ctx, []string{types.QueryKeyHeights}, query)
	require.NoError(t, err)
	bz, err = codec.MarshalJSONIndent(app.Codec(), types.KeyHeights{[]uint32{0}})
	require.NoError(t, err)
	require.Equal(t, bz, keyHeightsBs, "key heights is not [0]")
}
