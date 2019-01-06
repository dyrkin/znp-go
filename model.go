package znp

type Network struct {
	NeighborPanID   uint16
	LogicalChannel  uint8
	StackProfile    uint8
	ZigbeeVersion   uint8
	BeaconOrder     uint8
	SuperFrameOrder uint8
	PermitJoin      uint8
}

//Capabilities represents the interfaces that this device can handle (compiled into the device)
type Capabilities struct {
	Sys   uint8
	Mac   uint8
	Nwk   uint8
	Af    uint8
	Zdo   uint8
	Sapi  uint8
	Util  uint8
	Debug uint8
	App   uint8
	Zoad  uint8
}

//ResetRequest is sent by the tester to reset the target device
type ResetRequest struct {
	//This command will reset the device by using a hardware reset (i.e.
	//watchdog reset) if ‘Type’ is zero. Otherwise a soft reset (i.e. a jump to the
	//reset vector) is done. This is especially useful in the CC2531, for
	//instance, so that the USB host does not have to contend with the USB
	//H/W resetting (and thus causing the USB host to re-enumerate the device
	//which can cause an open virtual serial port to hang.)
	ResetType byte
}

type PingResponse struct {
	Capabilities *Capabilities
}

type VersionResponse struct {
	TransportRev uint8 //Transport protocol revision
	Product      uint8 //Product Id
	MajorRel     uint8 //Software major release number
	MinorRel     uint8 //Software minor release number
	MaintRel     uint8 //Software maintenance release number
}

type StatusResponse struct {
	Status uint8 //Failure 1, Success 0
}

//SetExtAddrReqest is used to set the extended address of the device
type SetExtAddrReqest struct {
	ExtAddress string `subtype:"longaddr"` //The device’s extended address.
}

//GetExtAddrResponse is used to get the extended address of the device.
type GetExtAddrResponse struct {
	ExtAddress string `subtype:"longaddr"` //The device’s extended address.
}

//RamReadRequest is used by the tester to read a single memory location in the target RAM. The
//command accepts an address value and returns the memory value present in the target RAM at that
//address.
type RamReadRequest struct {
	Address uint16 //Address of the memory that will be read.
	Len     uint8  //The number of bytes that will be read from the target RAM.
}

type RamReadResponse struct {
	Status uint8   //Status is either Success (0) or Failure (1).
	Value  []uint8 `len:"uint8"` //The value read from the target RAM.
}

//RamWriteRequest is used by the tester to write to a particular location in the target RAM. The
//command accepts an address location and a memory value. The memory value is written to the
//address location in the target RAM.
type RamWriteRequest struct {
	Address uint16  //Address of the memory that will be written.
	Value   []uint8 `len:"uint8"` //The value written to the target RAM.
}

type LedControlRequest struct {
	LedID uint8
	Mode  uint8
}
