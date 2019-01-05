package znp

import (
	"bytes"

	"github.com/dyrkin/composer"
)

type Network struct {
	NeighborPanID   uint16
	LogicalChannel  uint8
	StackProfile    uint8
	ZigbeeVersion   uint8
	BeaconOrder     uint8
	SuperFrameOrder uint8
	PermitJoin      uint8
}

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
