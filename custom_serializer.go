package znp

import (
	"bytes"

	"github.com/dyrkin/composer"
)

func readNetwork(buf *bytes.Buffer) *Network {
	c := composer.NewWithRW(buf)
	neighborPanID, _ := c.ReadUint16le()
	logicalChannel, _ := c.ReadUint8()
	b, _ := c.ReadUint8()
	stackProfile := b & 0x0F
	zigbeeVersion := (b & 0xF0) >> 4
	b, _ = c.ReadUint8()
	beaconOrder := b & 0x0F
	superFrameOrder := (b & 0xF0) >> 4
	permitJoin, _ := c.ReadUint8()
	return &Network{
		NeighborPanID:   neighborPanID,
		LogicalChannel:  logicalChannel,
		StackProfile:    stackProfile,
		ZigbeeVersion:   zigbeeVersion,
		BeaconOrder:     beaconOrder,
		SuperFrameOrder: superFrameOrder,
		PermitJoin:      permitJoin,
	}
}

func writeNetwork(buf *bytes.Buffer, n *Network) {
	c := composer.NewWithRW(buf)
	c.Uint16le(n.NeighborPanID).Uint8(n.LogicalChannel)
	b := (byte(n.StackProfile) & 0x0F) | ((byte(n.ZigbeeVersion << 4)) & 0xF0)
	c.Uint8(b)
	b = (byte(n.BeaconOrder) & 0x0F) | ((byte(n.SuperFrameOrder << 4)) & 0xF0)
	c.Uint8(b).Uint8(n.PermitJoin)
	c.Flush()
}

func readCapabilities(buf *bytes.Buffer) *Capabilities {
	c := composer.NewWithRW(buf)
	b, _ := c.ReadUint16le()
	return &Capabilities{
		Sys:   uint8(b & 0x0001),
		Mac:   uint8((b & 0x0002) >> 1),
		Nwk:   uint8((b & 0x0004) >> 2),
		Af:    uint8((b & 0x0008) >> 3),
		Zdo:   uint8((b & 0x0010) >> 4),
		Sapi:  uint8((b & 0x0020) >> 5),
		Util:  uint8((b & 0x0040) >> 6),
		Debug: uint8((b & 0x0080) >> 7),
		App:   uint8((b & 0x0100) >> 8),
		Zoad:  uint8((b & 0x1000) >> 9),
	}
}

func writeCapabilities(buf *bytes.Buffer, n *Capabilities) {
	c := composer.NewWithRW(buf)
	b := uint16(n.Sys&0x0001) | (uint16(n.Mac<<1) & 0x0002) | (uint16(n.Nwk<<2) & 0x0004) |
		(uint16(n.Af<<3) & 0x0008) | (uint16(n.Zdo<<4) & 0x0010) | (uint16(n.Sapi<<5) & 0x0020) |
		(uint16(n.Util<<6) & 0x0040) | (uint16(n.Debug<<7) & 0x0080) | ((uint16(n.App) << 8) & 0x0100) |
		((uint16(n.Zoad) << 9) & 0x1000)
	c.Uint16le(b)
	c.Flush()
}
