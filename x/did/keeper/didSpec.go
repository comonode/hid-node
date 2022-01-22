package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hypersign-protocol/hid-node/utils"
	"github.com/hypersign-protocol/hid-node/x/did/types"
)

func (k Keeper) GetDidCount(ctx sdk.Context) uint64 {
	// Get the store using storeKey (which is "blog") and PostCountKey (which is "Post-count-")
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DidCountKey))
	// Convert the PostCountKey to bytes
	byteKey := []byte(types.DidCountKey)
	// Get the value of the count
	bz := store.Get(byteKey)
	// Return zero if the count value is not found (for example, it's the first post)
	if bz == nil {
		return 0
	}
	// Convert the count into a uint64
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetDidCount(ctx sdk.Context, count uint64) {
	// Get the store using storeKey (which is "blog") and PostCountKey (which is "Post-count-")
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DidCountKey))
	// Convert the PostCountKey to bytes
	byteKey := []byte(types.DidCountKey)
	// Convert count from uint64 to string and get bytes
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	// Set the value of Post-count- to count
	store.Set(byteKey, bz)
}

func (k Keeper) AppendDID(ctx sdk.Context, didSpec types.DidSpec) uint64 {
	// Get the current number of posts in the store
	count := k.GetDidCount(ctx)
	// Assign an ID to the post based on the number of posts in the store
	didSpec.Id = count
	// Get the store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.DidKey))
	// Convert the post ID into bytes
	byteKey := make([]byte, 8)
	binary.BigEndian.PutUint64(byteKey, didSpec.Id)
	// Marshal the post into bytes
	appendedValue := k.cdc.MustMarshal(&didSpec)
	// Insert the post bytes using post ID as a key

	store.Set(utils.UnsafeStrToBytes(didSpec.Did), appendedValue)
	// Update the post count
	k.SetDidCount(ctx, count+1)
	return count
}
