package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	"github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

// nolint: unused
func (k Keeper) OnChanOpenInit(
	ctx sdk.Context,
	order channel.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channel.Counterparty,
	version string,
) error {
	if order != channel.UNORDERED {
		return channel.ErrInvalidChannel(k.codespace, "channel must be UNORDERED")
	}

	// NOTE: here the capability key name defines the port ID of the counterparty
	if counterparty.PortID != k.boundedCapability.Name() {
		return port.ErrInvalidPort(
			k.codespace,
			fmt.Sprintf("counterparty port ID doesn't match the capability key (%s ≠ %s)", counterparty.PortID, k.boundedCapability.Name()))
	}

	if strings.TrimSpace(version) != "" {
		return ibctypes.ErrInvalidVersion(k.codespace, "version must be blank")
	}

	// NOTE: as the escrow address is generated from both the port and channel IDs
	// there's no need to store it on a map.
	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenTry(
	ctx sdk.Context,
	order channel.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channel.Counterparty,
	version string,
	counterpartyVersion string,
) error {
	if order != channel.UNORDERED {
		return channel.ErrInvalidChannel(k.codespace, "channel must be UNORDERED")
	}

	// NOTE: here the capability key name defines the port ID of the counterparty
	if counterparty.PortID != k.boundedCapability.Name() {
		return port.ErrInvalidPort(
			k.codespace,
			fmt.Sprintf("counterparty port ID doesn't match the capability key (%s ≠ %s)", counterparty.PortID, k.boundedCapability.Name()))
	}

	if strings.TrimSpace(version) != "" {
		return ibctypes.ErrInvalidVersion(k.codespace, "version must be blank")
	}

	if strings.TrimSpace(counterpartyVersion) != "" {
		return ibctypes.ErrInvalidVersion(k.codespace, "counterparty version must be blank")
	}

	// NOTE: as the escrow address is generated from both the port and channel IDs
	// there's no need to store it on a map.
	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	version string,
) error {
	if strings.TrimSpace(version) != "" {
		return ibctypes.ErrInvalidVersion(k.codespace, "version must be blank")
	}

	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// no-op
	return nil
}

// onRecvPacket is called when an FTTransfer packet is received
// nolint: unused
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channelexported.Packet,
	data types.PacketDataTransfer,
) error {
	return k.ReceiveTransfer(
		ctx, packet.GetSourcePort(), packet.GetSourceChannel(),
		packet.GetDestPort(), packet.GetDestChannel(), data,
	)
}

// nolint: unused
func (k Keeper) OnAcknowledgePacket(
	ctx sdk.Context,
	packet channelexported.Packet,
	acknowledgement []byte,
) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnTimeoutPacket(
	ctx sdk.Context,
	packet channelexported.Packet,
	data types.PacketDataTransfer,
) error {
	// check the denom prefix
	prefix := types.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
	coins := make(sdk.Coins, len(data.Amount))
	for i, coin := range data.Amount {
		coin := coin
		if !strings.HasPrefix(coin.Denom, prefix) {
			return sdk.ErrInvalidCoins(fmt.Sprintf("%s doesn't contain the prefix '%s'", coin.Denom, prefix))
		}
		coins[i] = sdk.NewCoin(coin.Denom[len(prefix):], coin.Amount)
	}

	if data.Source {
		escrowAddress := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
		return k.bankKeeper.SendCoins(ctx, escrowAddress, data.Sender, coins)
	}

	// mint from supply
	err := k.supplyKeeper.MintCoins(ctx, types.GetModuleAccountName(), data.Amount)
	if err != nil {
		return err
	}

	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.GetModuleAccountName(), data.Sender, data.Amount)
}

// nolint: unused
func (k Keeper) OnTimeoutPacketClose(_ sdk.Context, _ channelexported.Packet) {
	panic("can't happen, only unordered channels allowed")
}