package znp

import unpi "github.com/dyrkin/unpi-go"

var AsyncCommandRegistry = make(map[registryKey]interface{})

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
	InvalidParameter
	MemError Status = iota + 0x10
	BufferFull
	UnsupportedMode
	MacMemError
	SapiInProgress Status = iota + 0x20
	SapiTimeout
	SapiInit
	NotAuthorized          = 0x7E
	MalformedCmd           = 0x80
	UnsupClusterCmd        = 0x81
	OtaAbort        Status = iota + 0x95
	OtaImageInvalid
	OtaWaitForData
	OtaNoImageAvailable
	OtaRequireMoreImage
	ApsFail Status = iota + 0xb1
	ApsTableFull
	ApsIllegalRequest
	ApsInvalidBinding
	ApsUnsupportedAttrib
	ApsNotSupported
	ApsNoAck
	ApsDuplicateEntry
	ApsNoBoundDevice
	ApsNotAllowed
	ApsNotAuthenticated
	SecNoKey Status = iota + 0xa1
	SecOldFrmCount
	SecMaxFrmCount
	SecCcmFail
	SecFailure             = 0xad
	NwkInvalidParam Status = iota + 0xc1
	NwkInvalidRequest
	NwkNotPermitted
	NwkStartupFailure
	NwkAlreadyPresent
	NwkSyncFailure
	NwkTableFull
	NwkUnknownDevice
	NwkUnsupportedAttribute
	NwkNoNetworks
	NwkLeaveUnconfirmed
	NwkNoAck
	NwkNoRoute
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
	DstPanID    uint16 //PAN - personal area networks
	SrcEndpoint uint8
	ClusterID   uint16
	TransID     uint8
	Options     *AfDataRequestOptions
	Radius      uint8
	Data        []uint8 `len:"uint16"`
}

func (znp *Znp) AfDataRequestExt(dstAddrMode AddrMode, dstAddr string, dstEndpoint uint8, dstPanId uint16,
	srcEndpoint uint8, clusterId uint16, transId uint8, options *AfDataRequestOptions, radius uint8,
	data []uint8) (*StatusResponse, error) {
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

func (znp *Znp) AfInterPanCtl(command InterPanCommand, data AfInterPanCtlData) (*StatusResponse, error) {
	req := &AfInterPanCtl{Command: command, Data: data}
	rsp := &StatusResponse{}
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

func (znp *Znp) AfDataStore(index uint16, data []uint8) (*StatusResponse, error) {
	req := &AfDataStore{Index: index, Data: data}
	rsp := &StatusResponse{}
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
	Status StatusResponse
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

type AfApsfConfigSet struct {
	Endpoint   uint8
	FrameDelay uint8
	WindowSize uint8
}

func (znp *Znp) AfApsfConfigSet(endpoint uint8, frameDelay uint8, windowSize uint8) (*StatusResponse, error) {
	req := &AfApsfConfigSet{Endpoint: endpoint, FrameDelay: frameDelay, WindowSize: windowSize}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x13, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type AfDataConfirm struct {
	Status   Status
	Endpoint uint8
	TransID  uint8
}

type AfReflectError struct {
	Status      Status
	Endpoint    uint8
	TransID     uint8
	DstAddrMode AddrMode
	DstAddr     string `hex:"uint16"`
}

type AfIncomingMessage struct {
	GroupID        uint16
	ClusterID      uint16
	SrcAddr        string `hex:"uint16"`
	SrcEndpoint    uint8
	DstEndpoint    uint8
	WasBroadcast   uint8
	LinkQuality    uint8
	SecurityUse    uint8
	Timestamp      uint32
	TransSeqNumber uint8
	Data           []uint8 `len:"uint8"`
}

type AfIncomingMessageExt struct {
	GroupID        uint16
	ClusterID      uint16
	SrcAddrMode    AddrMode
	SrcAddr        string `hex:"uint64"`
	SrcEndpoint    uint8
	SrcPanID       uint16
	DstEndpoint    uint8
	WasBroadcast   uint8
	LinkQuality    uint8
	SecurityUse    uint8
	Timestamp      uint32
	TransSeqNumber uint8
	Data           []uint8 `len:"uint16"`
}

// =======APP=======

type AppMsg struct {
	AppEndpoint uint8
	DstAddr     string `hex:"uint16"`
	DstEndpoint uint8
	ClusterID   uint16
	Message     []uint8 `len:"uint8"`
}

func (znp *Znp) AppMsg(appEndpoint uint8, dstAddr string, dstEndpoint uint8, clusterID uint16,
	message []uint8) (*StatusResponse, error) {
	req := &AppMsg{AppEndpoint: appEndpoint, DstAddr: dstAddr, DstEndpoint: dstEndpoint,
		ClusterID: clusterID, Message: message}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_APP, 0x00, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type AppUserTest struct {
	SrcEndpoint uint8
	CommandID   uint16
	Parameter1  uint16
	Parameter2  uint16
}

func (znp *Znp) AppUserTest(srcEndpoint uint8, commandId uint16, parameter1 uint16, parameter2 uint16) (*StatusResponse, error) {
	req := &AppUserTest{SrcEndpoint: srcEndpoint, CommandID: commandId, Parameter1: parameter1, Parameter2: parameter2}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_APP, 0x01, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// =======DEBUG=======

type DebugSetThreshold struct {
	ComponentID uint8
	Threshold   uint8
}

func (znp *Znp) DebugSetThreshold(componentId uint8, threshold uint8) (*StatusResponse, error) {
	req := &DebugSetThreshold{ComponentID: componentId, Threshold: threshold}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_DEBUG, 0x00, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type DebugMsg struct {
	String string `len:"uint8"`
}

func (znp *Znp) DebugMsg(str string) error {
	req := &DebugMsg{String: str}
	return znp.ProcessRequest(unpi.C_AREQ, unpi.S_DEBUG, 0x00, req, nil)
}

// =======MAC======= is not supported on my device

func (znp *Znp) MacInit() (*StatusResponse, error) {
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_MAC, 0x02, nil, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// =======SAPI=======

func (znp *Znp) SapiZbSystemReset() error {
	return znp.ProcessRequest(unpi.C_AREQ, unpi.S_SAPI, 0x09, nil, nil)
}

type EmptyResponse struct{}

func (znp *Znp) SapiZbStartRequest() (*EmptyResponse, error) {
	rsp := &EmptyResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x00, nil, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type SapiZbPermitJoiningRequest struct {
	Destination string `hex:"uint16"`
	Timeout     uint8
}

func (znp *Znp) SapiZbPermitJoiningRequest(destination string, timeout uint8) (*StatusResponse, error) {
	req := &SapiZbPermitJoiningRequest{Destination: destination, Timeout: timeout}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x08, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type SapiZbBindDevice struct {
	Create      uint8
	CommandID   uint16
	Destination string `hex:"uint64"`
}

func (znp *Znp) SapiZbBindDevice(create uint8, commandId uint16, destination string) (*EmptyResponse, error) {
	req := &SapiZbBindDevice{Create: create, CommandID: commandId, Destination: destination}
	rsp := &EmptyResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x01, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type SapiZbAllowBind struct {
	Timeout uint8
}

func (znp *Znp) SapiZbAllowBind(timeout uint8) (*EmptyResponse, error) {
	req := &SapiZbAllowBind{Timeout: timeout}
	rsp := &EmptyResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x02, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

const (
	ZbBindingAddr   = "0xFFFE"
	ZbBroadcastAddr = "0xFFFF"
)

type SapiZbSendDataRequest struct {
	Destination string `hex:"uint16"`
	CommandID   uint16
	Handle      uint8
	Ack         uint8
	Radius      uint8
	Data        []uint8 `len:"uint8"`
}

func (znp *Znp) SapiZbSendDataRequest(destination string, commandID uint16, handle uint8,
	ack uint8, radius uint8, data []uint8) (*EmptyResponse, error) {
	req := &SapiZbSendDataRequest{Destination: destination, CommandID: commandID,
		Handle: handle, Ack: ack, Radius: radius, Data: data}
	rsp := &EmptyResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x03, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type SapiZbReadConfiguration struct {
	ConfigID uint8
}

type SapiZbReadConfigurationResponse struct {
	Status   Status
	ConfigID uint8
	Value    []uint8 `len:"uint8"`
}

func (znp *Znp) SapiZbReadConfiguration(configID uint8) (*SapiZbReadConfigurationResponse, error) {
	req := &SapiZbReadConfiguration{ConfigID: configID}
	rsp := &SapiZbReadConfigurationResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x04, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type SapiZbWriteConfiguration struct {
	ConfigID uint8
	Value    []uint8 `len:"uint8"`
}

func (znp *Znp) SapiZbWriteConfiguration(configID uint8, value []uint8) (*StatusResponse, error) {
	req := &SapiZbWriteConfiguration{ConfigID: configID, Value: value}
	rsp := &StatusResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x05, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type SapiZbGetDeviceInfo struct {
	Param uint8
}

type SapiZbGetDeviceInfoResponse struct {
	Param uint8
	Value uint16
}

func (znp *Znp) SapiZbGetDeviceInfo(param uint8) (*SapiZbGetDeviceInfoResponse, error) {
	req := &SapiZbGetDeviceInfo{Param: param}
	rsp := &SapiZbGetDeviceInfoResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x06, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type SapiZbFindDeviceRequest struct {
	SearchKey string `hex:"uint64"`
}

func (znp *Znp) SapiZbFindDeviceRequest(searchKey string) (*EmptyResponse, error) {
	req := &SapiZbFindDeviceRequest{SearchKey: searchKey}
	rsp := &EmptyResponse{}
	err := znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x07, req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type SapiZbStartConfirm struct {
	Status Status
}

type SapiZbBindConfirm struct {
	CommandID uint16
	Status    Status
}

type SapiZbAllowBindConfirm struct {
	Source string `hex:"uint16"`
}

type SapiZbSendDataConfirm struct {
	Handle uint8
	Status Status
}

type SapiZbReceiveDataIndication struct {
	Source    string `hex:"uint16"`
	CommandID uint16
	Data      []uint8 `len:"uint8"`
}

type SapiZbFindDeviceConfirm struct {
	SearchType uint8
	Result     string `hex:"uint16"`
	SearchKey  string `hex:"uint64"`
}

func init() {
	//AF
	AsyncCommandRegistry[registryKey{unpi.S_AF, 0x80}] = &AfDataConfirm{}
	AsyncCommandRegistry[registryKey{unpi.S_AF, 0x83}] = &AfReflectError{}
	AsyncCommandRegistry[registryKey{unpi.S_AF, 0x81}] = &AfIncomingMessage{}
	AsyncCommandRegistry[registryKey{unpi.S_AF, 0x82}] = &AfIncomingMessageExt{}

	//DEBUG
	AsyncCommandRegistry[registryKey{unpi.S_DEBUG, 0x00}] = &DebugMsg{}

	//SAPI
	AsyncCommandRegistry[registryKey{unpi.S_SAPI, 0x80}] = &SapiZbStartConfirm{}
	AsyncCommandRegistry[registryKey{unpi.S_SAPI, 0x81}] = &SapiZbBindConfirm{}
	AsyncCommandRegistry[registryKey{unpi.S_SAPI, 0x82}] = &SapiZbAllowBindConfirm{}
	AsyncCommandRegistry[registryKey{unpi.S_SAPI, 0x83}] = &SapiZbSendDataConfirm{}
	AsyncCommandRegistry[registryKey{unpi.S_SAPI, 0x87}] = &SapiZbReceiveDataIndication{}
	AsyncCommandRegistry[registryKey{unpi.S_SAPI, 0x85}] = &SapiZbFindDeviceConfirm{}
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
