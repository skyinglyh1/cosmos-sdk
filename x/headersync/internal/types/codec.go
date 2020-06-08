package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSyncGenesisParam{}, ModuleName+"/MsgSyncGenesisParam", nil)
	cdc.RegisterConcrete(MsgSyncHeadersParam{}, ModuleName+"/MsgSyncHeadersParam", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
