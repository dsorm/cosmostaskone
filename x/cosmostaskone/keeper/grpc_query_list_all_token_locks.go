package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dsorm/cosmostaskone/x/cosmostaskone/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListAllTokenLocks(goCtx context.Context, req *types.QueryListAllTokenLocksRequest) (*types.QueryListAllTokenLocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// create new store from sdk context
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(keyPrefix))

	var msgK msgServer
	tokenLockCount := msgK.GetTokensLockCount(ctx)

	tokensLockList := make([]*types.TokensLock, 0, tokenLockCount)

	msg := types.MsgAddTokensLock{}
	for i := uint64(1); i <= tokenLockCount; i++ {
		key := keyPrefix + strconv.FormatUint(i, 10)
		bz := store.Get([]byte(key))

		// TODO use TokensLock from tokenslock.proto
		k.cdc.MustUnmarshalBinaryBare(bz, &msg)
		tokensLockList = append(tokensLockList, &types.TokensLock{
			Creator:  msg.Creator,
			Balances: msg.Balances,
		})
	}

	return &types.QueryListAllTokenLocksResponse{TokensLockList: tokensLockList}, nil
}
