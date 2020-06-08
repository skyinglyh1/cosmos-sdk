package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/headersync/internal/types"
	polycommon "github.com/cosmos/cosmos-sdk/x/headersync/poly-utils/common"
	vconfig "github.com/cosmos/cosmos-sdk/x/headersync/poly-utils/consensus/vbft/config"
	polysig "github.com/cosmos/cosmos-sdk/x/headersync/poly-utils/core/signature"
	polytype "github.com/cosmos/cosmos-sdk/x/headersync/poly-utils/core/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/pkg/errors"
	"sort"
	"strconv"
)

// Keeper of the mint store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey) Keeper {

	return Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}

func (keeper Keeper) SyncGenesisHeader(ctx sdk.Context, genesisHeaderBytes []byte) error {
	genesisHeader := &polytype.Header{}

	source := polycommon.NewZeroCopySource(genesisHeaderBytes)
	if err := genesisHeader.Deserialization(source); err != nil {
		return types.ErrDeserializeHeader
	}
	if storedHeader, err := keeper.GetHeaderByHeight(ctx, genesisHeader.ChainID, genesisHeader.Height); storedHeader != nil && err == nil {
		return errors.Wrapf(err, "GenesisHeader already synced")
	}
	if err := keeper.SetBlockHeader(ctx, genesisHeader); err != nil {
		return err
	}
	if err := keeper.UpdateConsensusPeer(ctx, genesisHeader); err != nil {
		return err
	}
	return nil
}

func (keeper Keeper) SyncBlockHeaders(ctx sdk.Context, headers [][]byte) error {
	for _, headerBytes := range headers {
		header := &polytype.Header{}
		source := polycommon.NewZeroCopySource(headerBytes)
		if err := header.Deserialization(source); err != nil {
			return types.ErrDeserializeHeader
		}
		h, err := keeper.GetHeaderByHeight(ctx, header.ChainID, header.Height)
		if err != nil {
			return errors.New(fmt.Sprintf("SyncBlockHeader chainId=%d, height=%d, err:%s", header.ChainID, header.Height, err.Error()))
		}

		if h == nil {
			if err := keeper.ProcessHeader(ctx, header); err != nil {
				return errors.New(fmt.Sprintf("SyncBlockHeader error:%s", err.Error()))
			}
		}
	}
	return nil
}

func (keeper Keeper) ProcessHeader(ctx sdk.Context, header *polytype.Header) error {
	if err := keeper.VerifyHeader(ctx, header); err != nil {
		return errors.New(fmt.Sprintf("processHeader, %s", err.Error()))
	}
	if err := keeper.SetBlockHeader(ctx, header); err != nil {
		return errors.New(fmt.Sprintf("processHeader, %s", err.Error()))
	}
	if err := keeper.UpdateConsensusPeer(ctx, header); err != nil {
		return errors.New(fmt.Sprintf("processHeader, %s", err.Error()))
	}
	return nil
}

type HeaderSyncViewKeeper interface {
	GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*polytype.Header, error)
	GetHeaderByHash(ctx sdk.Context, chainId uint64, hash polycommon.Uint256) (*polytype.Header, error)
	GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, error)
	GetConsensusPeers(ctx sdk.Context, chainId uint64, height uint32) (*types.ConsensusPeers, error)
	GetKeyHeights(ctx sdk.Context, chainId uint64) *types.KeyHeights
}

func (keeper Keeper) SetBlockHeader(ctx sdk.Context, blockHeader *polytype.Header) error {

	store := ctx.KVStore(keeper.storeKey)
	blockHash := blockHeader.Hash()
	sink := polycommon.NewZeroCopySink(nil)
	if err := blockHeader.Serialization(sink); err != nil {
		return types.ErrDeserializeHeader
	}
	store.Set(GetBlockHeaderKey(blockHeader.ChainID, blockHash.ToArray()), sink.Bytes())
	store.Set(GetBlockHashKey(blockHeader.ChainID, blockHeader.Height), types.ModuleCdc.MustMarshalJSON(blockHash))
	store.Set(GetBlockCurHeightKey(blockHeader.ChainID), types.ModuleCdc.MustMarshalJSON(blockHeader.Height))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSyncHeader,
			sdk.NewAttribute(types.AttributeKeyChainId, strconv.FormatUint(blockHeader.ChainID, 10)),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatUint(uint64(blockHeader.Height), 10)),
			sdk.NewAttribute(types.AttributeKeyBlockHash, hex.EncodeToString(blockHash[:])),
			sdk.NewAttribute(types.AttributeKeyNativeChainHeight, strconv.FormatUint(uint64(ctx.BlockHeight()), 10)),
		),
	})
	return nil
}
func (keeper Keeper) GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, error) {
	store := ctx.KVStore(keeper.storeKey)
	heightBs := store.Get(GetBlockCurHeightKey(chainId))
	if heightBs == nil {
		return 0, nil
	}
	var height uint32
	if err := types.ModuleCdc.UnmarshalJSON(heightBs, &height); err != nil {
		return 0, types.ErrUnmarshalSpecificTypeFail(height, err)
	}
	return height, nil

}

func (keeper Keeper) GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*polytype.Header, error) {
	store := ctx.KVStore(keeper.storeKey)
	hashBytes := store.Get(GetBlockHashKey(chainId, height))
	if hashBytes == nil {
		return nil, nil
	}
	blockHash := new(polycommon.Uint256)
	types.ModuleCdc.MustUnmarshalJSON(hashBytes, blockHash)
	headerBytes := store.Get(GetBlockHeaderKey(chainId, blockHash.ToArray()))
	if headerBytes == nil {
		return nil, nil
	}
	header := new(polytype.Header)
	source := polycommon.NewZeroCopySource(headerBytes)
	if err := header.Deserialization(source); err != nil {
		return nil, types.ErrDeserializeHeader
	}
	return header, nil

}
func (keeper Keeper) GetHeaderByHash(ctx sdk.Context, chainId uint64, hash polycommon.Uint256) (*polytype.Header, error) {
	store := ctx.KVStore(keeper.storeKey)
	headerBytes := store.Get(GetBlockHeaderKey(chainId, hash.ToArray()))
	if headerBytes == nil {
		return nil, errors.New(fmt.Sprintf("get block header error: chainid = %d, hash = %s", chainId, hex.EncodeToString(hash.ToArray())))
	}
	header := new(polytype.Header)
	source := polycommon.NewZeroCopySource(headerBytes)
	if err := header.Deserialization(source); err != nil {
		return nil, types.ErrDeserializeHeader
	}
	return header, nil

}

func (keeper Keeper) UpdateConsensusPeer(ctx sdk.Context, blockHeader *polytype.Header) error {

	blkInfo := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(blockHeader.ConsensusPayload, blkInfo); err != nil {
		return types.ErrUnmarshalSpecificTypeFail(blkInfo, err)
	}
	if blkInfo.NewChainConfig != nil {
		consensusPeers := &types.ConsensusPeers{
			ChainID: blockHeader.ChainID,
			Height:  blockHeader.Height,
			PeerMap: make(map[string]*types.Peer),
		}
		for _, p := range blkInfo.NewChainConfig.Peers {
			consensusPeers.PeerMap[p.ID] = &types.Peer{Index: p.Index, PeerPubkey: p.ID}
		}
		//TODO: check consensus peer to avoid duplicate
		if err := keeper.SetConsensusPeers(ctx, *consensusPeers); err != nil {
			return errors.New(fmt.Sprintf("updateConsensusPeer, set ConsensusPeer error: %s", err.Error()))
		}
	}

	return nil
}

func (keeper Keeper) SetConsensusPeers(ctx sdk.Context, consensusPeers types.ConsensusPeers) error {
	store := ctx.KVStore(keeper.storeKey)

	sink := polycommon.NewZeroCopySink(nil)
	consensusPeers.Serialization(sink)

	store.Set(GetConsensusPeerKey(consensusPeers.ChainID, consensusPeers.Height), sink.Bytes())
	// update key heights
	keyHeights := keeper.GetKeyHeights(ctx, consensusPeers.ChainID)

	keyHeights.HeightList = append(keyHeights.HeightList, consensusPeers.Height)

	if e := keeper.SetKeyHeights(ctx, consensusPeers.ChainID, keyHeights); e != nil {
		return errors.New(fmt.Sprintf("SetConsensusPeers, set KeyHeights error: %s", e.Error()))
	}
	return nil
}

func (keeper Keeper) GetConsensusPeers(ctx sdk.Context, chainId uint64, height uint32) (*types.ConsensusPeers, error) {
	store := ctx.KVStore(keeper.storeKey)

	consensusPeerBytes := store.Get(GetConsensusPeerKey(chainId, height))
	if consensusPeerBytes == nil {
		return nil, types.ErrGetConsensusPeers(height, chainId)
	}
	consensusPeers := new(types.ConsensusPeers)
	if err := consensusPeers.Deserialization(polycommon.NewZeroCopySource(consensusPeerBytes)); err != nil {
		return nil, errors.New(fmt.Sprintf("GetConsensusPeers, Deserialization Error:%v", err))
	}
	return consensusPeers, nil
}

func (keeper Keeper) SetKeyHeights(ctx sdk.Context, chainId uint64, keyHeights *types.KeyHeights) error {
	//first sort the list  (big -> small)
	sort.SliceStable(keyHeights.HeightList, func(i, j int) bool {
		return keyHeights.HeightList[i] > keyHeights.HeightList[j]
	})
	store := ctx.KVStore(keeper.storeKey)
	bz, err := types.ModuleCdc.MarshalBinaryLengthPrefixed(keyHeights)

	if err != nil {
		return types.ErrMarshalSpecificTypeFail(keyHeights, err)
	}
	store.Set(GetKeyHeightsKey(chainId), bz)
	return nil
}

func (keeper Keeper) GetKeyHeights(ctx sdk.Context, chainId uint64) *types.KeyHeights {
	store := ctx.KVStore(keeper.storeKey)
	keyHeightBytes := store.Get(GetKeyHeightsKey(chainId))
	keyHeights := new(types.KeyHeights)
	if keyHeightBytes == nil {
		return keyHeights
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(keyHeightBytes, keyHeights)
	return keyHeights
}

func (keeper Keeper) VerifyHeader(ctx sdk.Context, header *polytype.Header) error {
	height := header.Height
	keyHeight, err := keeper.FindKeyHeight(ctx, height, header.ChainID)
	if err != nil {
		return err
	}
	consensusPeer, err := keeper.GetConsensusPeers(ctx, header.ChainID, keyHeight)
	if err != nil {
		return err
	}
	if len(header.Bookkeepers)*3 < len(consensusPeer.PeerMap)*2 {
		return types.ErrBookKeeperNum(len(header.Bookkeepers), len(consensusPeer.PeerMap))
	}
	for _, bookkeeper := range header.Bookkeepers {
		pubkey := vconfig.PubkeyID(bookkeeper)
		_, present := consensusPeer.PeerMap[pubkey]
		if !present {
			return types.ErrInvalidPublicKey(pubkey)
		}
	}
	hash := header.Hash()
	if e := polysig.VerifyMultiSignature(hash[:], header.Bookkeepers, len(header.Bookkeepers), header.SigData); e != nil {
		return types.ErrVerifyMultiSignatureFail(err, header.Height)
	}
	return nil
}

func (keeper Keeper) FindKeyHeight(ctx sdk.Context, height uint32, chainId uint64) (uint32, error) {
	keyHeights := keeper.GetKeyHeights(ctx, chainId)
	for _, v := range keyHeights.HeightList {
		if (height - v) > 0 {
			return v, nil
		}
	}
	return 0, types.ErrFindKeyHeight(height, chainId)
}
