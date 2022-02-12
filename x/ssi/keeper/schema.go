package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hypersign-protocol/hid-node/utils"
	"github.com/hypersign-protocol/hid-node/x/ssi/types"
)

func (k Keeper) GetSchemaCount(ctx sdk.Context) uint64 {
	// Get the store using storeKey (which is "blog") and PostCountKey (which is "Post-count-")
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SchemaCountKey))
	// Convert the PostCountKey to bytes
	byteKey := []byte(types.SchemaCountKey)
	// Get the value of the count
	bz := store.Get(byteKey)
	// Return zero if the count value is not found (for example, it's the first post)
	if bz == nil {
		return 0
	}
	// Convert the count into a uint64
	return binary.BigEndian.Uint64(bz)
}

// Check whether the given DID is already present in the store
func (k Keeper) HasSchema(ctx sdk.Context, id string) bool {
	extractedID := utils.ExtractIDFromSchema(id)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	return store.Has(utils.UnsafeStrToBytes(extractedID))
}

// Retrieves the DID from the store
// func (k Keeper) GetDid(ctx *sdk.Context, id string) (*types.DidDocStructCreateDID, error) {
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))

// 	if !k.HasDid(*ctx, id) {
// 		return nil, sdkerrors.ErrNotFound
// 	}

// 	// TODO: Currently, we can interchange DidDocStructCreateDID with DidDocStructUpdateDID
// 	// as they are similar. They need to have a similar Interface
// 	var value types.DidDocStructCreateDID
// 	var bytes = store.Get([]byte(id))
// 	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
// 	}

// 	return &value, nil
// }

func (k Keeper) SetSchemaCount(ctx sdk.Context, count uint64) {
	// Get the store using storeKey (which is "blog") and PostCountKey (which is "Post-count-")
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SchemaCountKey))
	// Convert the PostCountKey to bytes
	byteKey := []byte(types.SchemaCountKey)
	// Convert count from uint64 to string and get bytes
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	// Set the value of Post-count- to count
	store.Set(byteKey, bz)
}

// SetDid set a specific did in the store
// func (k Keeper) SetDid(ctx sdk.Context, did types.Did) error {
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
// 	b := k.cdc.MustMarshal(&did)
// 	store.Set([]byte(did.Id), b)
// 	return nil
// }

func (k Keeper) AppendSchema(ctx sdk.Context, schema types.Schema) uint64 {
	// Get the current number of posts in the store
	count := k.GetDidCount(ctx)
	// Get the store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SchemaKey))
	// Convert the post ID into bytes
	// Marshal the post into bytes
	schemaBytes := k.cdc.MustMarshal(&schema)
	// Strip the id part from schema.ID
	schemaID := utils.ExtractIDFromSchema(schema.Id)
	store.Set(utils.UnsafeStrToBytes(schemaID), schemaBytes)
	// Update the post count
	k.SetDidCount(ctx, count+1)
	return count
}
