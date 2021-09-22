package keeper

import (
	"context"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dsorm/cosmostaskone/x/cosmostaskone/types"
)

func (k msgServer) AddTokensLock(goCtx context.Context, msg *types.MsgAddTokensLock) (*types.MsgAddTokensLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// create new store from sdk context
	store := ctx.KVStore(k.storeKey)

	// check if the account has the balance specified
	creatorAddress, err := cosmosTypes.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	for _, v := range msg.Balances {
		if !k.bankKeeper.HasBalance(ctx, creatorAddress, *v) {
			return nil,
				sdkerrors.Wrapf(
					sdkerrors.ErrInsufficientFunds,
					"The account (%v) doesn't have %v %v\n",
					msg.Creator,
					v.Amount,
					v.Denom)
		}
	}

	// convert []*cosmosTypes.Coin to []cosmosTypes.Coin
	//coins := types.DereferenceCoinSlice(msg.Balances)

	// temporary workaround for problem that causes the transaction to fail when multiple coins are entered at once
	tempCoins := make(cosmosTypes.Coins, 1, 1)
	for _, v := range msg.Balances {
		tempCoins[0] = *v
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddress, types.ModuleName, tempCoins)
		if err != nil {
			return nil, err
		}
	}

	// err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddress, types.ModuleName, cosmosTypes.Coins(coins))

	if err != nil {
		return nil, err
	}

	tokenLock := types.TokenLockInternal{
		Creator:  msg.Creator,
		Balances: msg.Balances,
		Disabled: false,
	}
	tokenLock.GenerateKeyForTokenLock(store)
	tokenLock.Save(store, k.cdc)

	return &types.MsgAddTokensLockResponse{}, nil
}
