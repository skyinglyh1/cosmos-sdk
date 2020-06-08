//nolint
package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"reflect"
)

var (
	ErrDeserializeHeader         = sdkerrors.Register(ModuleName, 1, "Header deserialization error")
	ErrSupplyKeeperMintCoinsFail = sdkerrors.Register(ModuleName, 2, "supplyKeeper mint coins failed ")
)

func ErrMarshalSpecificTypeFail(o interface{}, err error) error {
	return sdkerrors.Wrap(err, fmt.Sprintf("marshal type: %s, error: %s", reflect.TypeOf(o).String(), err.Error()))
}

func ErrUnmarshalSpecificTypeFail(o interface{}, err error) error {
	return sdkerrors.Wrap(err, fmt.Sprintf("marshal type: %s, error: %s", reflect.TypeOf(o).String(), err.Error()))
}

func ErrFindKeyHeight(height uint32, chainId uint64) error {
	return errors.New(fmt.Sprintf("findKeyHeight error: can not find key height with height:%d and chainId:%d", height, chainId))
}

func ErrGetConsensusPeers(height uint32, chainId uint64) error {
	return errors.New(fmt.Sprintf("get consensus peers empty error: chainId: %d, height: %d", height, chainId))
}

func ErrBookKeeperNum(headerBookKeeperNum int, consensusNodeNum int) error {
	return errors.New(fmt.Sprintf("header Bookkeepers number:%d must more than 2/3 consensus node number:%d", headerBookKeeperNum, consensusNodeNum))
}

func ErrInvalidPublicKey(pubkey string) error {
	return errors.New(fmt.Sprintf("invalid pubkey error:%s", pubkey))
}

func ErrVerifyMultiSignatureFail(err error, height uint32) error {
	return sdkerrors.Wrap(err, fmt.Sprintf("verify multi signature error:%s, height:%d", err.Error(), height))
}

func ErrInvalidChainId(chainId uint64) error {
	return errors.New(fmt.Sprintf("unknown chainId with id %d", chainId))
}

func ErrEmptyTargetHash(targetHashStr string) error {
	return errors.New(fmt.Sprintf("empty target asset hash %s", targetHashStr))
}
func ErrBelowCrossedLimit(limit sdk.Int, storedCrossedLimit sdk.Int) error {
	return errors.New(fmt.Sprintf("new Limit:%s should be greater than stored Limit:%s", limit.String(), storedCrossedLimit.String()))
}

func ErrCrossedAmountOverLimit(newCrossedAmount sdk.Int, crossedLimit sdk.Int) error {
	return errors.New(fmt.Sprintf("new crossedAmount:%s is over the crossedLimit:%s", newCrossedAmount.String(), crossedLimit.String()))
}

func ErrCrossedAmountOverflow(newCrossedAmount sdk.Int, storedCrossedAmount sdk.Int) error {
	return errors.New(fmt.Sprintf("new crossedAmount:%s is not greater than stored crossed amount:%s", newCrossedAmount.String(), storedCrossedAmount.String()))
}
