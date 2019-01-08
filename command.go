package znp

import unpi "github.com/dyrkin/unpi-go"

type LatencyReq uint8

const (
	NoLatency LatencyReq = iota
	FastBeacons
	SlowBeacons
)

type Status uint8

const (
	Success Status = iota
	Failure
)

type InterPanCtlStatus uint8

const (
	InterPanCtlSuccess InterPanCtlStatus = iota
	InterPanCtlFailure
	InterPanCtlInvalidParameter
	InterPanCtlZApsNotAllowed = 0xBA
)

type DataStatus uint8

const (
	DataSuccess DataStatus = iota
	DataFailure
	DataInvalidParameter
	DataMemFail = 0x10
)

type AddrMode uint8

const (
	AddrNotPresent AddrMode = iota
	Addr16Bit
	Addr64Bit
	AddrGroup
	AddrBroadcast = 15
)

type InterPanCommand uint8

const (
	InterPanClr InterPanCommand = iota
	InterPanSet
	InterPanReg
	InterPanChk
)

type StatusResponse struct {
	Status Status
}

type InterPanCtlStatusResponse struct {
	Status InterPanCtlStatus
}

type DataStatusResponse struct {
	Status DataStatus
}

// =======AF=======

type AfRegister struct {
	EndPoint          uint8
	AppProfID         uint16
	AppDeviceID       uint16
	AddDevVer         uint8
	LatencyReq        LatencyReq
	AppInClusterList  []uint16 `len:"uint8"`
	AppOutClusterList []uint16 `len:"uint8"`
}

func (znp *Znp) AfRegister(endPoint uint8, appProfID uint16, appDeviceID uint16, addDevVer uint8,
	latencyReq LatencyReq, appInClusterList []uint16, appOutClusterList []uint16) (*StatusResponse, error) {
	req := &AfRegister{EndPoint: endPoint, AppProfID: appProfID, AppDeviceID: appDeviceID,
		AddDevVer: addDevVer, LatencyReq: latencyReq, AppInClusterList: appInClusterList, AppOutClusterList: appOutClusterList}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type AfDataRequestOptions struct {
	WildcardProfileID uint8 `bits:"0b00000010" bitmask:"start" `
	APSAck            uint8 `bits:"0b00010000"`
	DiscoverRoute     uint8 `bits:"0b00100000"`
	APSSecurity       uint8 `bits:"0b01000000"`
	SkipRouting       uint8 `bits:"0b10000000" bitmask:"end" `
}

type AfDataRequest struct {
	DstAddr     string `hex:"uint16"`
	DstEndpoint uint8
	SrcEndpoint uint8
	ClusterID   uint16
	TransID     uint8
	Options     *AfDataRequestOptions
	Radius      uint8
	Data        []uint8 `len:"uint8"`
}

func (znp *Znp) AfDataRequest(dstAddr string, dstEndpoint uint8, srcEndpoint uint8, clusterId uint16,
	transId uint8, options *AfDataRequestOptions, radius uint8, data []uint8) (*StatusResponse, error) {
	req := &AfDataRequest{DstAddr: dstAddr, DstEndpoint: dstEndpoint, SrcEndpoint: srcEndpoint,
		ClusterID: clusterId, TransID: transId, Options: options, Radius: radius, Data: data}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x01, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type AfDataRequestExt struct {
	DstAddrMode AddrMode
	DstAddr     string `hex:"uint64"`
	DstEndpoint uint8
	DstPanID    uint16
	SrcEndpoint uint8
	ClusterID   uint16
	TransID     uint8
	Options     *AfDataRequestOptions
	Radius      uint8
	Data        []uint8 `len:"uint16"`
}

func (znp *Znp) AfDataRequestExt(dstAddrMode AddrMode, dstAddr string, dstEndpoint uint8, dstPanId uint16, srcEndpoint uint8, clusterId uint16,
	transId uint8, options *AfDataRequestOptions, radius uint8, data []uint8) (*StatusResponse, error) {
	req := &AfDataRequestExt{DstAddrMode: dstAddrMode, DstAddr: dstAddr, DstEndpoint: dstEndpoint, DstPanID: dstPanId, SrcEndpoint: srcEndpoint,
		ClusterID: clusterId, TransID: transId, Options: options, Radius: radius, Data: data}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x02, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type AfDataRequestSrcRtgOptions struct {
	APSAck      uint8 `bits:"0b00000001" bitmask:"start`
	APSSecurity uint8 `bits:"0b00000100"`
	SkipRouting uint8 `bits:"0b00001000" bitmask:"end" `
}

type AfDataRequestSrcRtg struct {
	DstAddr     string `hex:"uint16"`
	DstEndpoint uint8
	SrcEndpoint uint8
	ClusterID   uint16
	TransID     uint8
	Options     *AfDataRequestSrcRtgOptions
	Radius      uint8
	RelayList   []string `len:"uint8" hex:"uint16"`
	Data        []uint8  `len:"uint8"`
}

func (znp *Znp) AfDataRequestSrcRtg(dstAddr string, dstEndpoint uint8, srcEndpoint uint8, clusterId uint16,
	transId uint8, options *AfDataRequestSrcRtgOptions, radius uint8, relayList []string, data []uint8) (*StatusResponse, error) {
	req := &AfDataRequestSrcRtg{DstAddr: dstAddr, DstEndpoint: dstEndpoint, SrcEndpoint: srcEndpoint,
		ClusterID: clusterId, TransID: transId, Options: options, Radius: radius, RelayList: relayList, Data: data}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x03, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type AfInterPanCtlData interface {
	AfInterPanCtlData()
}

type AfInterPanClrData struct{}

func (a *AfInterPanClrData) AfInterPanCtlData() {}

type AfInterPanSetData struct {
	Channel uint8
}

func (a *AfInterPanSetData) AfInterPanCtlData() {}

type AfInterPanRegData struct {
	Endpoint uint8
}

func (a *AfInterPanRegData) AfInterPanCtlData() {}

type AfInterPanChkData struct {
	PanID    uint16
	Endpoint uint8
}

func (a *AfInterPanChkData) AfInterPanCtlData() {}

type AfInterPanCtl struct {
	Command InterPanCommand
	Data    AfInterPanCtlData
}

func (znp *Znp) AfInterPanCtl(command InterPanCommand, data AfInterPanCtlData) (*InterPanCtlStatusResponse, error) {
	req := &AfInterPanCtl{Command: command, Data: data}
	rsp := &InterPanCtlStatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x10, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type AfDataStore struct {
	Index uint16
	Data  []uint8 `len:"uint8"`
}

func (znp *Znp) AfDataStore(index uint16, data []uint8) (*DataStatusResponse, error) {
	req := &AfDataStore{Index: index, Data: data}
	rsp := &DataStatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x11, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type AfDataRetrieve struct {
	Timestamp uint32
	Index     uint16
	Length    uint8
}

type AfDataRetrieveResponse struct {
	Status DataStatus
	Data   []uint8 `len:"uint8"`
}

func (znp *Znp) AfDataRetrieve(timestamp uint32, index uint16, length uint8) (*AfDataRetrieveResponse, error) {
	req := &AfDataRetrieve{Timestamp: timestamp, Index: index, Length: length}
	rsp := &AfDataRetrieveResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x12, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// =======SYS=======

func (znp *Znp) Reset(resetType byte) {
	znp.ProcessRequest(unpi.C_AREQ, unpi.S_SYS, 0, &ResetRequest{resetType}, nil)
}

//This command issues PING requests to verify if a device is active and check the capability of the device.
func (znp *Znp) Ping() (*PingResponse, error) {
	rsp := &PingResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 1, nil, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (znp *Znp) Version() (*VersionResponse, error) {
	rsp := &VersionResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 2, nil, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (znp *Znp) SetExtAddr(extAddr string) (*StatusResponse, error) {
	req := &SetExtAddrReqest{ExtAddress: extAddr}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 3, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (znp *Znp) GetExtAddr() (*GetExtAddrResponse, error) {
	rsp := &GetExtAddrResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 4, nil, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (znp *Znp) RamRead(address uint16, len uint8) (*RamReadResponse, error) {
	req := &RamReadRequest{Address: address, Len: len}
	rsp := &RamReadResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 5, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (znp *Znp) RamWrite(address uint16, value []uint8) (*StatusResponse, error) {
	req := &RamWriteRequest{Address: address, Value: value}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 6, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (znp *Znp) LedControl(ledID uint8, mode uint8) (*StatusResponse, error) {
	req := &LedControlRequest{LedID: ledID, Mode: mode}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 10, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
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

//SetExtAddrReqest is used to set the extended address of the device
type SetExtAddrReqest struct {
	ExtAddress string `hex:"uint64"` //The device’s extended address.
}

//GetExtAddrResponse is used to get the extended address of the device.
type GetExtAddrResponse struct {
	ExtAddress string `hex:"uint64"` //The device’s extended address.
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

type Network struct {
	NeighborPanID   uint16
	LogicalChannel  uint8
	StackProfile    uint8 `bitmask:"start" bits:"0b00001111"`
	ZigbeeVersion   uint8 `bitmask:"end" bits:"0b11110000"`
	BeaconOrder     uint8 `bitmask:"start" bits:"0b00001111"`
	SuperFrameOrder uint8 `bitmask:"end" bits:"0b11110000"`
	PermitJoin      uint8
}

//Capabilities represents the interfaces that this device can handle (compiled into the device)
type Capabilities struct {
	Sys   uint16 `bitmask:"start" bits:"0x0001"`
	Mac   uint16 `bits:"0x0002"`
	Nwk   uint16 `bits:"0x0004"`
	Af    uint16 `bits:"0x0008"`
	Zdo   uint16 `bits:"0x0010"`
	Sapi  uint16 `bits:"0x0020"`
	Util  uint16 `bits:"0x0040"`
	Debug uint16 `bits:"0x0080"`
	App   uint16 `bits:"0x0100"`
	Zoad  uint16 `bitmask:"end" bits:"0x1000"`
}
