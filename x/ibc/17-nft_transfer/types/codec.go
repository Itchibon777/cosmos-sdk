package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTransfer{}, "ibc/nft_transfer/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgRecvPacket{}, "ibc/nft_transfer/MsgRecvPacket", nil)
	cdc.RegisterConcrete(PacketData{}, "ibc/nft_transfer/PacketData", nil)
}

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	channel.RegisterCodec(ModuleCdc)
	commitment.RegisterCodec(ModuleCdc)
}
