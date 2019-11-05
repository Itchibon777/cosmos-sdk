package types

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

const (
	// SubModuleName defines the IBC nft_transfer name
	SubModuleName = "nft_transfer"

	// StoreKey is the store key string for IBC nft_transfer
	StoreKey = SubModuleName

	// RouterKey is the message route for IBC nft_transfer
	RouterKey = SubModuleName

	// QuerierRoute is the querier route for IBC nft_transfer
	QuerierRoute = SubModuleName

	// BoundPortID defines the name of the capability key
	BoundPortID = "nftnftnft"
)

// GetEscrowAddress returns the escrow address for the specified channel
//
// CONTRACT: this assumes that there's only one nft bridge module that owns the
// port associated with the channel ID so that the address created is actually
// unique.
func GetEscrowAddress(portID, channelID string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(portID + channelID)))
}

// GetDenomPrefix returns the receiving denomination prefix
func GetDenomPrefix(portID, channelID string) string {
	return fmt.Sprintf("%s/%s/", portID, channelID)
}

// GetModuleAccountName returns the IBC nft_transfer module account name for supply
func GetModuleAccountName() string {
	return fmt.Sprintf("%s/%s", ibctypes.ModuleName, SubModuleName)
}
