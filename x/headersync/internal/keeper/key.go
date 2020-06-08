package keeper

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/x/headersync/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

var (
	BlockHeaderPrefix   = []byte{0x01}
	BlockHashPrefix     = []byte{0x02}
	ConsensusPeerPrefix = []byte{0x03}
	KeyHeightPrefix     = []byte{0x04}

	BlockCurrentHeightKey = []byte("currentHeight")
)

func GetBlockHeaderKey(chainId uint64, blockHash []byte) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(append(BlockHeaderPrefix, b...), blockHash...)
}

func GetBlockHashKey(chainId uint64, height uint32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	bh := make([]byte, 4)
	binary.LittleEndian.PutUint32(bh, height)
	return append(append(BlockHashPrefix, b...), bh...)
}
func GetBlockCurHeightKey(chainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(BlockCurrentHeightKey, b...)
}

func GetConsensusPeerKey(chainId uint64, height uint32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	bh := make([]byte, 4)
	binary.LittleEndian.PutUint32(bh, height)
	return append(append(ConsensusPeerPrefix, b...), bh...)
}

func GetKeyHeightsKey(chainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(append(KeyHeightPrefix, b...), b...)
}
