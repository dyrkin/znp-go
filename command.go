package znp

import unpi "github.com/dyrkin/unpi-go"

var AsyncCommandRegistry = make(map[registryKey]interface{})

type LatencyReq uint8

const (
	NoLatency LatencyReq = iota
	FastBeacons
	SlowBeacons
)

type StartupFromAppStatus uint8

const (
	RestoredNetworkState StartupFromAppStatus = 0x00
	NewNetworkState      StartupFromAppStatus = 0x01
	LeaveAndNotStarted   StartupFromAppStatus = 0x02
)

type Status uint8

const (
	Success          Status = 0x00
	Failure          Status = 0x01
	InvalidParameter Status = 0x02

	ItemCreatedAndInitialized Status = 0x09
	InitializationFailed      Status = 0x0a
	BadLength                 Status = 0x0c

	// ZStack status values must start at 0x10, after the generic status values (defined in comdef.h)
	MemError        Status = 0x10
	BufferFull      Status = 0x11
	UnsupportedMode Status = 0x12
	MacMemError     Status = 0x13

	SapiInProgress Status = 0x20
	SapiTimeout    Status = 0x21
	SapiInit       Status = 0x22

	NotAuthorized Status = 0x7E

	MalformedCmd    Status = 0x80
	UnsupClusterCmd Status = 0x81

	ZdpInvalidEp         Status = 0x82 // Invalid endpoint value
	ZdpNotActive         Status = 0x83 // Endpoint not described by a simple desc.
	ZdpNotSupported      Status = 0x84 // Optional feature not supported
	ZdpTimeout           Status = 0x85 // Operation has timed out
	ZdpNoMatch           Status = 0x86 // No match for end device bind
	ZdpNoEntry           Status = 0x88 // Unbind request failed, no entry
	ZdpNoDescriptor      Status = 0x89 // Child descriptor not available
	ZdpInsufficientSpace Status = 0x8a // Insufficient space to support operation
	ZdpNotPermitted      Status = 0x8b // Not in proper state to support operation
	ZdpTableFull         Status = 0x8c // No table space to support operation
	ZdpNotAuthorized     Status = 0x8d // Permissions indicate request not authorized
	ZdpBindingTableFull  Status = 0x8e // No binding table space to support operation

	// OTA Status values
	OtaAbort            Status = 0x95
	OtaImageInvalid     Status = 0x96
	OtaWaitForData      Status = 0x97
	OtaNoImageAvailable Status = 0x98
	OtaRequireMoreImage Status = 0x99

	// APS status values
	ApsFail              Status = 0xb1
	ApsTableFull         Status = 0xb2
	ApsIllegalRequest    Status = 0xb3
	ApsInvalidBinding    Status = 0xb4
	ApsUnsupportedAttrib Status = 0xb5
	ApsNotSupported      Status = 0xb6
	ApsNoAck             Status = 0xb7
	ApsDuplicateEntry    Status = 0xb8
	ApsNoBoundDevice     Status = 0xb9
	ApsNotAllowed        Status = 0xba
	ApsNotAuthenticated  Status = 0xbb

	// Security status values
	SecNoKey       Status = 0xa1
	SecOldFrmCount Status = 0xa2
	SecMaxFrmCount Status = 0xa3
	SecCcmFail     Status = 0xa4
	SecFailure     Status = 0xad

	// NWK status values
	NwkInvalidParam         Status = 0xc1
	NwkInvalidRequest       Status = 0xc2
	NwkNotPermitted         Status = 0xc3
	NwkStartupFailure       Status = 0xc4
	NwkAlreadyPresent       Status = 0xc5
	NwkSyncFailure          Status = 0xc6
	NwkTableFull            Status = 0xc7
	NwkUnknownDevice        Status = 0xc8
	NwkUnsupportedAttribute Status = 0xc9
	NwkNoNetworks           Status = 0xca
	NwkLeaveUnconfirmed     Status = 0xcb
	NwkNoAck                Status = 0xcc // not in spec
	NwkNoRoute              Status = 0xcd

	// MAC status values
	// ZMacSuccess              Status = 0x00
	MacBeaconLoss           Status = 0xe0
	MacChannelAccessFailure Status = 0xe1
	MacDenied               Status = 0xe2
	MacDisableTrxFailure    Status = 0xe3
	MacFailedSecurityCheck  Status = 0xe4
	MacFrameTooLong         Status = 0xe5
	MacInvalidGTS           Status = 0xe6
	MacInvalidHandle        Status = 0xe7
	MacInvalidParameter     Status = 0xe8
	MacNoACK                Status = 0xe9
	MacNoBeacon             Status = 0xea
	MacNoData               Status = 0xeb
	MacNoShortAddr          Status = 0xec
	MacOutOfCap             Status = 0xed
	MacPANIDConflict        Status = 0xee
	MacRealignment          Status = 0xef
	MacTransactionExpired   Status = 0xf0
	MacTransactionOverFlow  Status = 0xf1
	MacTxActive             Status = 0xf2
	MacUnAvailableKey       Status = 0xf3
	MacUnsupportedAttribute Status = 0xf4
	MacUnsupported          Status = 0xf5
	MacSrcMatchInvalidIndex Status = 0xff
)

type AddrMode uint8

const (
	AddrNotPresent AddrMode = iota
	AddrGroup
	Addr16Bit
	Addr64Bit
	AddrBroadcast AddrMode = 15 //or 0xFF??????
)

type InterPanCommand uint8

const (
	InterPanClr InterPanCommand = iota
	InterPanSet
	InterPanReg
	InterPanChk
)

type Channel uint8

const (
	AIN0 Channel = iota
	AIN1
	AIN2
	AIN3
	AIN4
	AIN5
	AIN6
	AIN7
	TemperatureSensor Channel = 0x0E + iota
	VoltageReading
)

type Resolution uint8

const (
	Bit8 Resolution = iota
	Bit10
	Bit12
	Bit14
)

type Operation uint8

const (
	SetDirection Operation = iota
	SetInputMode
	Set
	Clear
	Toggle
	Read
)

type Reason uint8

const (
	PowerUp Reason = iota
	External
	WatchDog
)

type DeviceState uint8

const (
	InitializedNotStartedAutomatically DeviceState = iota
	InitializedNotConnectedToAnything
	DiscoveringPANsToJoin
	JoiningPAN
	RejoiningPAN
	JoinedButNotAuthenticated
	StartedAsDeviceAfterAuthentication
	DeviceJoinedAuthenticatedAndIsRouter
	StartingAsZigBeeCoordinator
	StartedAsZigBeeCoordinator
	DeviceHasLostInformationAboutItsParent
	DeviceSendingKeepAliveToParent
	DeviceWaitingBeforeRejoin
	ReJoiningPANInSecureModeScanningAllChannels
	ReJoiningPANInTrustCenterModeScanningCurrentChannel
	ReJoiningPANInTrustCenterModeScanningAllChannels
)

type SubsystemId uint16

const (
	Sys           SubsystemId = 0x0100
	Mac           SubsystemId = 0x0200
	Nwk           SubsystemId = 0x0300
	Af            SubsystemId = 0x0400
	Zdo           SubsystemId = 0x0500
	Sapi          SubsystemId = 0x0600
	Util          SubsystemId = 0x0700
	Debug         SubsystemId = 0x0800
	App           SubsystemId = 0x0900
	AllSubsystems SubsystemId = 0xFFFF
)

type Action uint8

const (
	Disable Action = 0
	Enable  Action = 1
)

type Shift uint8

const (
	NoShift  Shift = 0
	YesShift Shift = 1
)

type Mode uint8

const (
	OFF Mode = 0
	ON  Mode = 1
)

type Relation uint8

const (
	Parent Relation = iota
	ChildRfd
	ChildRfdRxIdle
	ChildFfd
	ChildFfdRxIdle
	Neighbor
	Other
)

type ReqType uint8

const (
	SingleDeviceResponse      ReqType = 0x00
	AssociatedDevicesResponse ReqType = 0x01
)

type RouteStatus uint8

const (
	Active            RouteStatus = 0x00
	DiscoveryUnderway RouteStatus = 0x01
	DiscoveryFailed   RouteStatus = 0x02
	Inactive          RouteStatus = 0x03
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
	AppInClusterList  []uint16 `size:"1"`
	AppOutClusterList []uint16 `size:"1"`
}

func (znp *Znp) AfRegister(endPoint uint8, appProfID uint16, appDeviceID uint16, addDevVer uint8,
	latencyReq LatencyReq, appInClusterList []uint16, appOutClusterList []uint16) (rsp *StatusResponse, err error) {
	req := &AfRegister{EndPoint: endPoint, AppProfID: appProfID, AppDeviceID: appDeviceID,
		AddDevVer: addDevVer, LatencyReq: latencyReq, AppInClusterList: appInClusterList, AppOutClusterList: appOutClusterList}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0, req, &rsp)
	return
}

type AfDataRequestOptions struct {
	WildcardProfileID uint8 `bits:"0b00000010" bitmask:"start" `
	APSAck            uint8 `bits:"0b00010000"`
	DiscoverRoute     uint8 `bits:"0b00100000"`
	APSSecurity       uint8 `bits:"0b01000000"`
	SkipRouting       uint8 `bits:"0b10000000" bitmask:"end" `
}

type AfDataRequest struct {
	DstAddr     string `hex:"2"`
	DstEndpoint uint8
	SrcEndpoint uint8
	ClusterID   uint16
	TransID     uint8
	Options     *AfDataRequestOptions
	Radius      uint8
	Data        []uint8 `size:"1"`
}

func (znp *Znp) AfDataRequest(dstAddr string, dstEndpoint uint8, srcEndpoint uint8, clusterId uint16,
	transId uint8, options *AfDataRequestOptions, radius uint8, data []uint8) (rsp *StatusResponse, err error) {
	req := &AfDataRequest{DstAddr: dstAddr, DstEndpoint: dstEndpoint, SrcEndpoint: srcEndpoint,
		ClusterID: clusterId, TransID: transId, Options: options, Radius: radius, Data: data}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x01, req, &rsp)
	return
}

type AfDataRequestExt struct {
	DstAddrMode AddrMode
	DstAddr     string `hex:"8"`
	DstEndpoint uint8
	DstPanID    uint16 //PAN - personal area networks
	SrcEndpoint uint8
	ClusterID   uint16
	TransID     uint8
	Options     *AfDataRequestOptions
	Radius      uint8
	Data        []uint8 `size:"2"`
}

func (znp *Znp) AfDataRequestExt(dstAddrMode AddrMode, dstAddr string, dstEndpoint uint8, dstPanId uint16,
	srcEndpoint uint8, clusterId uint16, transId uint8, options *AfDataRequestOptions, radius uint8,
	data []uint8) (rsp *StatusResponse, err error) {
	req := &AfDataRequestExt{DstAddrMode: dstAddrMode, DstAddr: dstAddr, DstEndpoint: dstEndpoint, DstPanID: dstPanId, SrcEndpoint: srcEndpoint,
		ClusterID: clusterId, TransID: transId, Options: options, Radius: radius, Data: data}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x02, req, &rsp)
	return
}

type AfDataRequestSrcRtgOptions struct {
	APSAck      uint8 `bits:"0b00000001" bitmask:"start`
	APSSecurity uint8 `bits:"0b00000100"`
	SkipRouting uint8 `bits:"0b00001000" bitmask:"end" `
}

type AfDataRequestSrcRtg struct {
	DstAddr     string `hex:"2"`
	DstEndpoint uint8
	SrcEndpoint uint8
	ClusterID   uint16
	TransID     uint8
	Options     *AfDataRequestSrcRtgOptions
	Radius      uint8
	RelayList   []string `size:"1" hex:"2"`
	Data        []uint8  `size:"1"`
}

func (znp *Znp) AfDataRequestSrcRtg(dstAddr string, dstEndpoint uint8, srcEndpoint uint8, clusterId uint16,
	transId uint8, options *AfDataRequestSrcRtgOptions, radius uint8, relayList []string, data []uint8) (rsp *StatusResponse, err error) {
	req := &AfDataRequestSrcRtg{DstAddr: dstAddr, DstEndpoint: dstEndpoint, SrcEndpoint: srcEndpoint,
		ClusterID: clusterId, TransID: transId, Options: options, Radius: radius, RelayList: relayList, Data: data}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x03, req, &rsp)
	return
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

func (znp *Znp) AfInterPanCtl(command InterPanCommand, data AfInterPanCtlData) (rsp *StatusResponse, err error) {
	req := &AfInterPanCtl{Command: command, Data: data}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x10, req, &rsp)
	return
}

type AfDataStore struct {
	Index uint16
	Data  []uint8 `size:"1"`
}

func (znp *Znp) AfDataStore(index uint16, data []uint8) (rsp *StatusResponse, err error) {
	req := &AfDataStore{Index: index, Data: data}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x11, req, &rsp)
	return
}

type AfDataRetrieve struct {
	Timestamp uint32
	Index     uint16
	Length    uint8
}

type AfDataRetrieveResponse struct {
	Status StatusResponse
	Data   []uint8 `size:"1"`
}

func (znp *Znp) AfDataRetrieve(timestamp uint32, index uint16, length uint8) (rsp *AfDataRetrieveResponse, err error) {
	req := &AfDataRetrieve{Timestamp: timestamp, Index: index, Length: length}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x12, req, &rsp)
	return
}

type AfApsfConfigSet struct {
	Endpoint   uint8
	FrameDelay uint8
	WindowSize uint8
}

func (znp *Znp) AfApsfConfigSet(endpoint uint8, frameDelay uint8, windowSize uint8) (rsp *StatusResponse, err error) {
	req := &AfApsfConfigSet{Endpoint: endpoint, FrameDelay: frameDelay, WindowSize: windowSize}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x13, req, &rsp)
	return
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
	DstAddr     string `hex:"2"`
}

type AfIncomingMessage struct {
	GroupID        uint16
	ClusterID      uint16
	SrcAddr        string `hex:"2"`
	SrcEndpoint    uint8
	DstEndpoint    uint8
	WasBroadcast   uint8
	LinkQuality    uint8
	SecurityUse    uint8
	Timestamp      uint32
	TransSeqNumber uint8
	Data           []uint8 `size:"1"`
}

type AfIncomingMessageExt struct {
	GroupID        uint16
	ClusterID      uint16
	SrcAddrMode    AddrMode
	SrcAddr        string `hex:"8"`
	SrcEndpoint    uint8
	SrcPanID       uint16
	DstEndpoint    uint8
	WasBroadcast   uint8
	LinkQuality    uint8
	SecurityUse    uint8
	Timestamp      uint32
	TransSeqNumber uint8
	Data           []uint8 `size:"2"`
}

// =======APP=======

type AppMsg struct {
	AppEndpoint uint8
	DstAddr     string `hex:"2"`
	DstEndpoint uint8
	ClusterID   uint16
	Message     []uint8 `size:"1"`
}

func (znp *Znp) AppMsg(appEndpoint uint8, dstAddr string, dstEndpoint uint8, clusterID uint16,
	message []uint8) (rsp *StatusResponse, err error) {
	req := &AppMsg{AppEndpoint: appEndpoint, DstAddr: dstAddr, DstEndpoint: dstEndpoint,
		ClusterID: clusterID, Message: message}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_APP, 0x00, req, &rsp)
	return
}

type AppUserTest struct {
	SrcEndpoint uint8
	CommandID   uint16
	Parameter1  uint16
	Parameter2  uint16
}

func (znp *Znp) AppUserTest(srcEndpoint uint8, commandId uint16, parameter1 uint16, parameter2 uint16) (rsp *StatusResponse, err error) {
	req := &AppUserTest{SrcEndpoint: srcEndpoint, CommandID: commandId, Parameter1: parameter1, Parameter2: parameter2}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_APP, 0x01, req, &rsp)
	return
}

// =======DEBUG=======

type DebugSetThreshold struct {
	ComponentID uint8
	Threshold   uint8
}

func (znp *Znp) DebugSetThreshold(componentId uint8, threshold uint8) (rsp *StatusResponse, err error) {
	req := &DebugSetThreshold{ComponentID: componentId, Threshold: threshold}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_DEBUG, 0x00, req, &rsp)
	return
}

type DebugMsg struct {
	String string `size:"1"`
}

func (znp *Znp) DebugMsg(str string) error {
	req := &DebugMsg{String: str}
	return znp.ProcessRequest(unpi.C_AREQ, unpi.S_DEBUG, 0x00, req, nil)
}

// =======MAC======= is not supported on my device

func (znp *Znp) MacInit() (rsp *StatusResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_MAC, 0x02, nil, &rsp)
	return
}

// =======SAPI=======

func (znp *Znp) SapiZbSystemReset() error {
	return znp.ProcessRequest(unpi.C_AREQ, unpi.S_SAPI, 0x09, nil, nil)
}

type EmptyResponse struct{}

func (znp *Znp) SapiZbStartRequest() (rsp *EmptyResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x00, nil, &rsp)
	return
}

type SapiZbPermitJoiningRequest struct {
	Destination string `hex:"2"`
	Timeout     uint8
}

func (znp *Znp) SapiZbPermitJoiningRequest(destination string, timeout uint8) (rsp *StatusResponse, err error) {
	req := &SapiZbPermitJoiningRequest{Destination: destination, Timeout: timeout}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x08, req, &rsp)
	return
}

type SapiZbBindDevice struct {
	Create      uint8
	CommandID   uint16
	Destination string `hex:"8"`
}

func (znp *Znp) SapiZbBindDevice(create uint8, commandId uint16, destination string) (rsp *EmptyResponse, err error) {
	req := &SapiZbBindDevice{Create: create, CommandID: commandId, Destination: destination}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x01, req, &rsp)
	return
}

type SapiZbAllowBind struct {
	Timeout uint8
}

func (znp *Znp) SapiZbAllowBind(timeout uint8) (rsp *EmptyResponse, err error) {
	req := &SapiZbAllowBind{Timeout: timeout}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x02, req, &rsp)
	return
}

const (
	ZbBindingAddr   = "0xFFFE"
	ZbBroadcastAddr = "0xFFFF"
)

type SapiZbSendDataRequest struct {
	Destination string `hex:"2"`
	CommandID   uint16
	Handle      uint8
	Ack         uint8
	Radius      uint8
	Data        []uint8 `size:"1"`
}

func (znp *Znp) SapiZbSendDataRequest(destination string, commandID uint16, handle uint8,
	ack uint8, radius uint8, data []uint8) (rsp *EmptyResponse, err error) {
	req := &SapiZbSendDataRequest{Destination: destination, CommandID: commandID,
		Handle: handle, Ack: ack, Radius: radius, Data: data}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x03, req, &rsp)
	return
}

type SapiZbReadConfiguration struct {
	ConfigID uint8
}

type SapiZbReadConfigurationResponse struct {
	Status   Status
	ConfigID uint8
	Value    []uint8 `size:"1"`
}

func (znp *Znp) SapiZbReadConfiguration(configID uint8) (rsp *SapiZbReadConfigurationResponse, err error) {
	req := &SapiZbReadConfiguration{ConfigID: configID}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x04, req, &rsp)
	return
}

type SapiZbWriteConfiguration struct {
	ConfigID uint8
	Value    []uint8 `size:"1"`
}

func (znp *Znp) SapiZbWriteConfiguration(configID uint8, value []uint8) (rsp *StatusResponse, err error) {
	req := &SapiZbWriteConfiguration{ConfigID: configID, Value: value}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x05, req, &rsp)
	return
}

type SapiZbGetDeviceInfo struct {
	Param uint8
}

type SapiZbGetDeviceInfoResponse struct {
	Param uint8
	Value uint16
}

func (znp *Znp) SapiZbGetDeviceInfo(param uint8) (rsp *SapiZbGetDeviceInfoResponse, err error) {
	req := &SapiZbGetDeviceInfo{Param: param}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x06, req, &rsp)
	return
}

type SapiZbFindDeviceRequest struct {
	SearchKey string `hex:"8"`
}

func (znp *Znp) SapiZbFindDeviceRequest(searchKey string) (rsp *EmptyResponse, err error) {
	req := &SapiZbFindDeviceRequest{SearchKey: searchKey}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x07, req, &rsp)
	return
}

type SapiZbStartConfirm struct {
	Status Status
}

type SapiZbBindConfirm struct {
	CommandID uint16
	Status    Status
}

type SapiZbAllowBindConfirm struct {
	Source string `hex:"2"`
}

type SapiZbSendDataConfirm struct {
	Handle uint8
	Status Status
}

type SapiZbReceiveDataIndication struct {
	Source    string `hex:"2"`
	CommandID uint16
	Data      []uint8 `size:"1"`
}

type SapiZbFindDeviceConfirm struct {
	SearchType uint8
	Result     string `hex:"2"`
	SearchKey  string `hex:"8"`
}

// =======SYS=======

type SysResetReq struct {
	//This command will reset the device by using a hardware reset (i.e.
	//watchdog reset) if ‘Type’ is zero. Otherwise a soft reset (i.e. a jump to the
	//reset vector) is done. This is especially useful in the CC2531, for
	//instance, so that the USB host does not have to contend with the USB
	//H/W resetting (and thus causing the USB host to re-enumerate the device
	//which can cause an open virtual serial port to hang.)
	ResetType byte
}

//SysReset is sent by the tester to reset the target device
func (znp *Znp) SysResetReq(resetType byte) error {
	req := &SysResetReq{resetType}
	return znp.ProcessRequest(unpi.C_AREQ, unpi.S_SYS, 0x00, req, nil)
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

type SysPingResponse struct {
	Capabilities *Capabilities
}

//SysPing issues PING requests to verify if a device is active and check the capability of the device.
func (znp *Znp) SysPing() (rsp *SysPingResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x01, nil, &rsp)
	return
}

type SysVersionResponse struct {
	TransportRev uint8 //Transport protocol revision
	Product      uint8 //Product Id
	MajorRel     uint8 //Software major release number
	MinorRel     uint8 //Software minor release number
	MaintRel     uint8 //Software maintenance release number
}

func (znp *Znp) SysVersion() (rsp *SysVersionResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x02, nil, &rsp)
	return
}

type SysSetExtAddr struct {
	ExtAddress string `hex:"8"` //The device’s extended address.
}

//SysSetExtAddr is used to set the extended address of the device
func (znp *Znp) SysSetExtAddr(extAddr string) (rsp *StatusResponse, err error) {
	req := &SysSetExtAddr{ExtAddress: extAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x03, req, &rsp)
	return
}

type SysGetExtAddrResponse struct {
	ExtAddress string `hex:"8"` //The device’s extended address.
}

//SysGetExtAddr is used to get the extended address of the device
func (znp *Znp) SysGetExtAddr() (rsp *SysGetExtAddrResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x04, nil, &rsp)
	return
}

type SysRamRead struct {
	Address uint16 //Address of the memory that will be read.
	Len     uint8  //The number of bytes that will be read from the target RAM.
}

type SysRamReadResponse struct {
	Status uint8   //Status is either Success (0) or Failure (1).
	Value  []uint8 `size:"1"` //The value read from the target RAM.
}

//SysRamRead is used by the tester to read a single memory location in the target RAM. The
//command accepts an address value and returns the memory value present in the target RAM at that address.
func (znp *Znp) SysRamRead(address uint16, len uint8) (rsp *SysRamReadResponse, err error) {
	req := &SysRamRead{Address: address, Len: len}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x05, req, &rsp)
	return
}

type SysRamWrite struct {
	Address uint16  //Address of the memory that will be written.
	Value   []uint8 `size:"1"` //The value written to the target RAM.
}

//SysRamWrite is used by the tester to write to a particular location in the target RAM. The
//command accepts an address location and a memory value. The memory value is written to the
//address location in the target RAM.
func (znp *Znp) SysRamWrite(address uint16, value []uint8) (rsp *StatusResponse, err error) {
	req := &SysRamWrite{Address: address, Value: value}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x06, req, &rsp)
	return
}

type SysOsalNvRead struct {
	ID     uint16
	Offset uint8
}

type SysOsalNvReadResponse struct {
	Status Status
	Value  []uint8 `size:"1"`
}

//SysOsalNvRead is used by the tester to read a single memory item from the target non-volatile
//memory. The command accepts an attribute Id value and data offset and returns the memory value
//present in the target for the specified attribute Id.
func (znp *Znp) SysOsalNvRead(id uint16, offset uint8) (rsp *StatusResponse, err error) {
	req := &SysOsalNvRead{ID: id, Offset: offset}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x08, req, &rsp)
	return
}

type SysOsalNvWrite struct {
	ID     uint16
	Offset uint8
	Value  []uint8 `size:"1"`
}

//SysOsalNvWrite is used by the tester to write to a particular item in non-volatile memory. The
//command accepts an attribute Id, data offset, data length, and attribute value. The attribute value is
//written to the location specified for the attribute Id in the target.
func (znp *Znp) SysOsalNvWrite(id uint16, offset uint8, value []uint8) (rsp *StatusResponse, err error) {
	req := &SysOsalNvWrite{ID: id, Offset: offset, Value: value}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x09, req, &rsp)
	return
}

type SysOsalNvItemInit struct {
	ID       uint16
	ItemLen  uint16
	InitData []uint8 `size:"1"`
}

//SysOsalNvItemInit is used by the tester to create and initialize an item in non-volatile memory. The
//NV item will be created if it does not already exist. The data for the new NV item will be left
//uninitialized if the InitLen parameter is zero. When InitLen is non-zero, the data for the NV item
//will be initialized (starting at offset of zero) with the values from InitData. Note that it is not
//necessary to initialize the entire NV item (InitLen < ItemLen). It is also possible to create an NV
//item that is larger than the maximum length InitData – use the SYS_OSAL_NV_WRITE
//command to finish the initialization.
func (znp *Znp) SysOsalNvItemInit(id uint16, itemLen uint16, initData []uint8) (rsp *StatusResponse, err error) {
	req := &SysOsalNvItemInit{ID: id, ItemLen: itemLen, InitData: initData}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x07, req, &rsp)
	return
}

type SysOsalNvDelete struct {
	ID      uint16
	ItemLen uint16
}

//SysOsalNvDelete is used by the tester to delete an item from the non-volatile memory. The ItemLen
//parameter must match the length of the NV item or the command will fail. Use this command with
//caution – deleted items cannot be recovered.
func (znp *Znp) SysOsalNvDelete(id uint16, itemLen uint16) (rsp *StatusResponse, err error) {
	req := &SysOsalNvDelete{ID: id, ItemLen: itemLen}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x12, req, &rsp)
	return
}

type SysOsalNvLength struct {
	ID uint16
}

type SysOsalNvLengthResponse struct {
	Length uint16
}

//SysOsalNvLength is used by the tester to get the length of an item in non-volatile memory. A
//returned length of zero indicates that the NV item does not exist.
func (znp *Znp) SysOsalNvLength(id uint16) (rsp *SysOsalNvLengthResponse, err error) {
	req := &SysOsalNvLength{ID: id}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x13, req, &rsp)
	return
}

type SysOsalStartTimer struct {
	ID      uint8
	Timeout uint16
}

//SysOsalStartTimer is used by the tester to start a timer event. The event will expired after the indicated
//amount of time and a notification will be sent back to the tester.
func (znp *Znp) SysOsalStartTimer(id uint8, timeout uint16) (rsp *StatusResponse, err error) {
	req := &SysOsalStartTimer{ID: id, Timeout: timeout}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x0A, req, &rsp)
	return
}

type SysOsalStopTimer struct {
	ID uint8
}

//SysOsalStopTimer is used by the tester to stop a timer event.
func (znp *Znp) SysOsalStopTimer(id uint8) (rsp *StatusResponse, err error) {
	req := &SysOsalStopTimer{ID: id}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x0B, req, &rsp)
	return
}

type SysRandomResponse struct {
	Value uint16
}

//SysRandom is used by the tester to get a random 16-bit number.
func (znp *Znp) SysRandom() (rsp *SysRandomResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x0C, nil, &rsp)
	return
}

type SysAdcRead struct {
	Channel    Channel
	Resolution Resolution
}

type SysAdcReadResponse struct {
	Value uint16
}

//SysAdcRead reads a value from the ADC based on specified channel and resolution.
func (znp *Znp) SysAdcRead(channel Channel, resolution Resolution) (rsp *SysAdcReadResponse, err error) {
	req := &SysAdcRead{Channel: channel, Resolution: resolution}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x0D, req, &rsp)
	return
}

type SysGpio struct {
	Operation Operation
	Value     uint8
}

type SysGpioResponse struct {
	Value uint8
}

//SysGpio is used by the tester to control the 4 GPIO pins on the CC2530-ZNP build.
func (znp *Znp) SysGpio(operation Operation, value uint8) (rsp *SysGpioResponse, err error) {
	req := &SysGpio{Operation: operation, Value: value}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x0E, req, &rsp)
	return
}

type SysTime struct {
	UTCTime uint32
	Hour    uint8
	Minute  uint8
	Second  uint8
	Month   uint8
	Day     uint8
	Year    uint16
}

//SysSetTime is used by the tester to set the target system date and time. The time can be
//specified in “seconds since 00:00:00 on January 1, 2000” or in parsed date/time components
func (znp *Znp) SysSetTime(utcTime uint32, hour uint8, minute uint8, second uint8,
	month uint8, day uint8, year uint16) (rsp *StatusResponse, err error) {
	req := &SysTime{UTCTime: utcTime, Hour: hour, Minute: minute, Second: second, Month: month, Day: day, Year: year}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x10, req, &rsp)
	return
}

//SysGetTime is used by the tester to get the target system date and time. The time is returned in
//seconds since 00:00:00 on January 1, 2000” and parsed date/time components.
func (znp *Znp) SysGetTime() (rsp *SysTime, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x11, nil, &rsp)
	return
}

type SysSetTxPower struct {
	TXPower uint8
}

type SysSetTxPowerResponse struct {
	TXPower uint8
}

//SysSetTxPower is used by the tester to set the target system radio transmit power. The returned TX
//power is the actual setting applied to the radio – nearest characterized value for the specific radio
func (znp *Znp) SysSetTxPower(txPower uint8) (rsp *SysSetTxPowerResponse, err error) {
	req := &SysSetTxPower{TXPower: txPower}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x14, req, &rsp)
	return
}

//SysZDiagsInitStats is used to initialize the statistics table in NV memory.
func (znp *Znp) SysZDiagsInitStats() (rsp *StatusResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x17, nil, &rsp)
	return
}

type SysZDiagsClearStats struct {
	ClearNV uint8
}

type SysZDiagsClearStatsResponse struct {
	SysClock uint32
}

//SysZDiagsClearStats is used to clear the statistics table. To clear data in NV (including the Boot
//Counter) the clearNV flag shall be set to TRUE.
func (znp *Znp) SysZDiagsClearStats(clearNV uint8) (rsp *SysZDiagsClearStatsResponse, err error) {
	req := &SysZDiagsClearStats{ClearNV: clearNV}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x18, req, &rsp)
	return
}

type SysZDiagsGetStats struct {
	AttributeID uint16
}

type SysZDiagsGetStatsResponse struct {
	AttributeValue uint32
}

//SysZDiagsGetStats is used to read a specific system (attribute) ID statistics and/or metrics value.
func (znp *Znp) SysZDiagsGetStats(attributeID uint16) (rsp *SysZDiagsGetStatsResponse, err error) {
	req := &SysZDiagsGetStats{AttributeID: attributeID}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x19, req, &rsp)
	return
}

//SysZDiagsRestoreStatsNv is used to restore the statistics table from NV into the RAM table.
func (znp *Znp) SysZDiagsRestoreStatsNv() (rsp *StatusResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x1A, nil, &rsp)
	return
}

type SysZDiagsSaveStatsToNvResponse struct {
	SysClock uint32
}

//SysZDiagsSaveStatsToNv is used to save the statistics table from RAM to NV.
func (znp *Znp) SysZDiagsSaveStatsToNv() (rsp *SysZDiagsSaveStatsToNvResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x1B, nil, &rsp)
	return
}

type SysNvCreate struct {
	SysID  uint8
	ItemID uint16
	SubID  uint16
	Length uint32
}

//SysNvCreate is used to attempt to create an item in non-volatile memory.
func (znp *Znp) SysNvCreate(sysID uint8, itemID uint16, subID uint16, length uint32) (rsp *StatusResponse, err error) {
	req := &SysNvCreate{SysID: sysID, ItemID: itemID, SubID: subID, Length: length}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x30, req, &rsp)
	return
}

type SysNvDelete struct {
	SysID  uint8
	ItemID uint16
	SubID  uint16
}

//SysNvDelete is used to attempt to delete an item in non-volatile memory.
func (znp *Znp) SysNvDelete(sysID uint8, itemID uint16, subID uint16) (rsp *StatusResponse, err error) {
	req := &SysNvDelete{SysID: sysID, ItemID: itemID, SubID: subID}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x31, req, &rsp)
	return
}

type SysNvLength struct {
	SysID  uint8
	ItemID uint16
	SubID  uint16
}

type SysNvLengthResponse struct {
	Length uint8
}

//SysNvLength is used to get the length of an item in non-volatile memory.
func (znp *Znp) SysNvLength(sysID uint8, itemID uint16, subID uint16) (rsp *SysNvLengthResponse, err error) {
	req := &SysNvLength{SysID: sysID, ItemID: itemID, SubID: subID}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x32, req, &rsp)
	return
}

type SysNvRead struct {
	SysID  uint8
	ItemID uint16
	SubID  uint16
	Offset uint16
	Length uint8
}

type SysNvReadResponse struct {
	Status Status
	Value  []uint8 `size:"1"`
}

//SysNvRead is used to read an item in non-volatile memory
func (znp *Znp) SysNvRead(sysID uint8, itemID uint16, subID uint16, offset uint16, length uint8) (rsp *SysNvReadResponse, err error) {
	req := &SysNvRead{SysID: sysID, ItemID: itemID, SubID: subID, Offset: offset, Length: length}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x33, req, &rsp)
	return
}

type SysNvWrite struct {
	SysID  uint8
	ItemID uint16
	SubID  uint16
	Offset uint16
	Value  []uint8 `size:"1"`
}

//SysNvWrite is used to write an item in non-volatile memory
func (znp *Znp) SysNvWrite(sysID uint8, itemID uint16, subID uint16, offset uint16, value []uint8) (rsp *StatusResponse, err error) {
	req := &SysNvWrite{SysID: sysID, ItemID: itemID, SubID: subID, Offset: offset, Value: value}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x34, req, &rsp)
	return
}

type SysNvUpdate struct {
	SysID  uint8
	ItemID uint16
	SubID  uint16
	Value  []uint8 `size:"1"`
}

//SysNvUpdate is used to update an item in non-volatile memory
func (znp *Znp) SysNvUpdate(sysID uint8, itemID uint16, subID uint16, value []uint8) (rsp *StatusResponse, err error) {
	req := &SysNvUpdate{SysID: sysID, ItemID: itemID, SubID: subID, Value: value}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x35, req, &rsp)
	return
}

type SysNvCompact struct {
	Threshold uint16
}

//SysNvCompact is used to compact the active page in non-volatile memory
func (znp *Znp) SysNvCompact(threshold uint16) (rsp *StatusResponse, err error) {
	req := &SysNvCompact{Threshold: threshold}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x36, req, &rsp)
	return
}

type SysNvReadExt struct {
	ID     uint16
	Offset uint16
}

//SysNvReadExt is used by the tester to read a single memory item from the target non-volatile
//memory. The command accepts an attribute Id value and data offset and returns the memory value
//present in the target for the specified attribute Id.
func (znp *Znp) SysNvReadExt(id uint16, offset uint16) (rsp *SysNvReadResponse, err error) {
	req := &SysNvReadExt{ID: id, Offset: offset}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x08, req, &rsp)
	return
}

type SysNvWriteExt struct {
	ID     uint16
	Offset uint16
	Value  []uint8 `size:"1"`
}

//SysNvWrite is used to write an item in non-volatile memory
func (znp *Znp) SysNvWriteExt(id uint16, offset uint16, value []uint8) (rsp *StatusResponse, err error) {
	req := &SysNvWriteExt{ID: id, Offset: offset, Value: value}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x09, req, &rsp)
	return
}

type SysResetInd struct {
	Reason       Reason
	TransportRev uint8
	Product      uint8
	MinorRel     uint8
	HwRev        uint8
}

type SysOsalTimerExpired struct {
	ID uint8
}

// =======UTIL=======

type DeviceType struct {
	Coordinator uint8 `bits:"0x01" bitmask:"start"`
	Router      uint8 `bits:"0x02"`
	EndDevice   uint8 `bits:"0x04" bitmask:"end"`
}

type UtilGetDeviceInfoResponse struct {
	Status           Status
	IEEEAddr         string `hex:"8"`
	ShortAddr        string `hex:"2"`
	DeviceType       *DeviceType
	DeviceState      DeviceState
	AssocDevicesList []string `size:"1" hex:"2"`
}

//UtilGetDeviceInfo is sent by the tester to retrieve the device info.
func (znp *Znp) UtilGetDeviceInfo() (rsp *UtilGetDeviceInfoResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x00, nil, &rsp)
	return
}

type NvInfoStatus struct {
	IEEEAddress   Status `bits:"0b00000001" bitmask:"start"`
	ScanChannels  Status `bits:"0b00000010"`
	PanID         Status `bits:"0b00000100"`
	SecurityLevel Status `bits:"0b00001000"`
	PreConfigKey  Status `bits:"0b00010000" bitmask:"end"`
}

type UtilGetNvInfoResponse struct {
	Status        *NvInfoStatus
	IEEEAddr      string `hex:"8"`
	ScanChannels  uint32
	PanID         uint16
	SecurityLevel uint8
	PreConfigKey  [16]uint8
}

//UtilGetNvInfo is used by the tester to read a block of parameters from non-volatile storage of the
//target device.
func (znp *Znp) UtilGetNvInfo() (rsp *UtilGetNvInfoResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x01, nil, &rsp)
	return
}

type UtilSetPanId struct {
	PanID uint16
}

//UtilSetPanId stores a PanId value into non-volatile memory to be used the next time the target device resets.
func (znp *Znp) UtilSetPanId(panId uint16) (rsp *StatusResponse, err error) {
	req := &UtilSetPanId{PanID: panId}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x02, req, &rsp)
	return
}

type UtilSetChannels struct {
	Channels uint32
}

//UtilSetChannels is used to store a channel select bit-mask into non-volatile memory to be used the
//next time the target device resets.
func (znp *Znp) UtilSetChannels(channels uint32) (rsp *StatusResponse, err error) {
	req := &UtilSetChannels{Channels: channels}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x03, req, &rsp)
	return
}

type UtilSetSecLevel struct {
	SecLevel uint8
}

//UtilSetSecLevel is used to store a security level value into non-volatile memory to be used the next time the target device
//resets.
func (znp *Znp) UtilSetSecLevel(secLevel uint8) (rsp *StatusResponse, err error) {
	req := &UtilSetSecLevel{SecLevel: secLevel}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x04, req, &rsp)
	return
}

type UtilSetPreCfgKey struct {
	PreCfgKey [16]uint8
}

//UtilSetPreCfgKey is used to store a pre-configured key array into non-volatile memory to be used the
//next time the target device resets.
func (znp *Znp) UtilSetPreCfgKey(preCfgKey [16]uint8) (rsp *StatusResponse, err error) {
	req := &UtilSetPreCfgKey{PreCfgKey: preCfgKey}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x05, req, &rsp)
	return
}

type UtilCallbackSubCmd struct {
	SubsystemID SubsystemId
	Action      Action
}

//UtilCallbackSubCmd subscribes/unsubscribes to layer callbacks. For particular subsystem callbacks to
//work, the software must be compiled with a special flag that is unique to that subsystem to enable
//the callback mechanism. For example to enable ZDO callbacks, MT_ZDO_CB_FUNC flag must
//be compiled when the software is built. For complete list of callback compile flags, check section
//1.2 or “Z-Stack Compile Options” document.
func (znp *Znp) UtilCallbackSubCmd(subsystemID SubsystemId, action Action) (rsp *StatusResponse, err error) {
	req := &UtilCallbackSubCmd{SubsystemID: subsystemID, Action: action}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x06, req, &rsp)
	return
}

type Keys struct {
	Key1 uint8 `bits:"0x01" bitmask:"start"`
	Key2 uint8 `bits:"0x02"`
	Key3 uint8 `bits:"0x04"`
	Key4 uint8 `bits:"0x08"`
	Key5 uint8 `bits:"0x10"`
	Key6 uint8 `bits:"0x20"`
	Key7 uint8 `bits:"0x40"`
	Key8 uint8 `bits:"0x80" bitmask:"end"`
}

type UtilKeyEvent struct {
	Keys  *Keys
	Shift Shift
}

//UtilKeyEvent sends key and shift codes to the application that registered for key events. The keys parameter is a
//bit mask, allowing for multiple keys in a single command. The return status indicates success if
//the command is processed by a registered key handler, not whether the key code was used. Not all
//applications support all key or shift codes but there is no indication when a key code is dropped.
func (znp *Znp) UtilKeyEvent(keys *Keys, shift Shift) (rsp *StatusResponse, err error) {
	req := &UtilKeyEvent{Keys: keys, Shift: shift}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x07, req, &rsp)
	return
}

type UtilTimeAliveResponse struct {
	Seconds uint32
}

//UtilTimeAlive is used by the tester to get the board’s time alive
func (znp *Znp) UtilTimeAlive() (rsp *UtilTimeAliveResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x09, nil, &rsp)
	return
}

type UtilLedControl struct {
	LedID uint8
	Mode  Mode
}

//UtilLedControl is used by the tester to control the LEDs on the board.
func (znp *Znp) UtilLedControl(ledID uint8, mode Mode) (rsp *StatusResponse, err error) {
	req := &UtilLedControl{LedID: ledID, Mode: mode}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x0A, req, &rsp)
	return
}

type UtilLoopback struct {
	Data []uint8
}

//UtilLoopback is used by the tester to test data buffer loopback.
func (znp *Znp) UtilLoopback(data []uint8) (rsp *UtilLoopback, err error) {
	req := &UtilLoopback{Data: data}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x10, req, &rsp)
	return
}

type UtilDataReq struct {
	SecurityUse uint8
}

//UtilDataReq is used by the tester to effect a MAC MLME Poll Request
func (znp *Znp) UtilDataReq(securityUse uint8) (rsp *StatusResponse, err error) {
	req := &UtilDataReq{SecurityUse: securityUse}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x11, req, &rsp)
	return
}

//UtilSrcMatchEnable is used to enable AUTOPEND and source address matching.
func (znp *Znp) UtilSrcMatchEnable() (rsp *StatusResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x20, nil, &rsp)
	return
}

type UtilSrcMatchAddEntry struct {
	AddrMode AddrMode
	Address  string `hex:"8"`
	PanID    uint16
}

//UtilSrcMatchAddEntry is used to add a short or extended address to the source address table
func (znp *Znp) UtilSrcMatchAddEntry(addrMode AddrMode, address string, panId uint16) (rsp *StatusResponse, err error) {
	req := &UtilSrcMatchAddEntry{AddrMode: addrMode, Address: address, PanID: panId}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x21, req, &rsp)
	return
}

type UtilSrcMatchDelEntry struct {
	AddrMode AddrMode
	Address  string `hex:"8"`
	PanID    uint16
}

//UtilSrcMatchDelEntry is used to delete a short or extended address from the source address table.
func (znp *Znp) UtilSrcMatchDelEntry(addrMode AddrMode, address string, panId uint16) (rsp *StatusResponse, err error) {
	req := &UtilSrcMatchDelEntry{AddrMode: addrMode, Address: address, PanID: panId}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x22, req, &rsp)
	return
}

type UtilSrcMatchCheckSrcAddr struct {
	AddrMode AddrMode
	Address  string `hex:"8"`
	PanID    uint16
}

//UtilSrcMatchCheckSrcAddr is used to delete a short or extended address from the source address table.
func (znp *Znp) UtilSrcMatchCheckSrcAddr(addrMode AddrMode, address string, panId uint16) (rsp *StatusResponse, err error) {
	req := &UtilSrcMatchCheckSrcAddr{AddrMode: addrMode, Address: address, PanID: panId}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x23, req, &rsp)
	return
}

type UtilSrcMatchAckAllPending struct {
	Option Action
}

//UtilSrcMatchAckAllPending is used to enable/disable acknowledging all packets with pending bit set.
func (znp *Znp) UtilSrcMatchAckAllPending(option Action) (rsp *StatusResponse, err error) {
	req := &UtilSrcMatchAckAllPending{Option: option}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x24, req, &rsp)
	return
}

type UtilSrcMatchCheckAllPendingResponse struct {
	Status Status
	Value  uint8
}

//UtilSrcMatchCheckAllPending is used to check if acknowledging all packets with pending bit set is enabled.
func (znp *Znp) UtilSrcMatchCheckAllPending() (rsp *UtilSrcMatchCheckAllPendingResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x25, nil, &rsp)
	return
}

type UtilAddrMgrExtAddrLookup struct {
	ExtAddr string `hex:"8"`
}

type UtilAddrMgrExtAddrLookupResponse struct {
	NwkAddr string `hex:"2"`
}

//UtilAddrMgrExtAddrLookup is a proxy call to the AddrMgrEntryLookupExt() function.
func (znp *Znp) UtilAddrMgrExtAddrLookup(extAddr string) (rsp *UtilAddrMgrExtAddrLookupResponse, err error) {
	req := &UtilAddrMgrExtAddrLookup{ExtAddr: extAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x40, req, &rsp)
	return
}

type UtilAddrMgrAddrLookup struct {
	NwkAddr string `hex:"2"`
}

type UtilAddrMgrAddrLookupResponse struct {
	ExtAddr string `hex:"8"`
}

//UtilAddrMgrAddrLookup is a proxy call to the AddrMgrEntryLookupNwk() function.
func (znp *Znp) UtilAddrMgrAddrLookup(nwkAddr string) (rsp *UtilAddrMgrAddrLookupResponse, err error) {
	req := &UtilAddrMgrAddrLookup{NwkAddr: nwkAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x41, req, &rsp)
	return
}

type UtilApsmeLinkKeyDataGet struct {
	ExtAddr string `hex:"8"`
}

type UtilApsmeLinkKeyDataGetResponse struct {
	Status    Status
	SecKey    [16]uint8
	TxFrmCntr uint32
	RxFrmCntr uint32
}

//UtilApsmeLinkKeyDataGet retrieves APS link key data, Tx and Rx frame counters
func (znp *Znp) UtilApsmeLinkKeyDataGet(extAddr string) (rsp *UtilApsmeLinkKeyDataGetResponse, err error) {
	req := &UtilApsmeLinkKeyDataGet{ExtAddr: extAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x44, req, &rsp)
	return
}

type UtilApsmeLinkKeyNvIdGet struct {
	ExtAddr string `hex:"8"`
}

type UtilApsmeLinkKeyNvIdGetResponse struct {
	Status      Status
	LinkKeyNvId uint16
}

//UtilApsmeLinkKeyNvIdGet is a proxy call to the APSME_LinkKeyNvIdGet() function.
func (znp *Znp) UtilApsmeLinkKeyNvIdGet(extAddr string) (rsp *UtilApsmeLinkKeyNvIdGetResponse, err error) {
	req := &UtilApsmeLinkKeyNvIdGet{ExtAddr: extAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x45, req, &rsp)
	return
}

type UtilApsmeRequestKeyCmd struct {
	PartnerAddr string `hex:"8"`
}

//UtilApsmeRequestKeyCmd is used to send a request key to the Trust Center from an originator device who
//wants to exchange messages with a partner device.
func (znp *Znp) UtilApsmeRequestKeyCmd(partnerAddr string) (rsp *StatusResponse, err error) {
	req := &UtilApsmeRequestKeyCmd{PartnerAddr: partnerAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x4B, req, &rsp)
	return
}

type UtilAssocCount struct {
	StartRelation Relation
	EndRelation   Relation
}

type UtilAssocCountResponse struct {
	Count uint16
}

//UtilAssocCount is a proxy call to the AssocCount() function
func (znp *Znp) UtilAssocCount(startRelation Relation, endRelation Relation) (rsp *UtilAssocCountResponse, err error) {
	req := &UtilAssocCount{StartRelation: startRelation, EndRelation: endRelation}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x48, req, &rsp)
	return
}

const (
	InvalidNodeAddr = "0xFFFE"
)

type LinkInfo struct {
	TxCounter uint8 // Counter of transmission success/failures
	TxCost    uint8 // Average of sending rssi values if link staus is enabled
	// i.e. NWK_LINK_STATUS_PERIOD is defined as non zero
	RxLqi uint8 // average of received rssi values
	// needs to be converted to link cost (1-7) before used
	InKeySeqNum uint8  // security key sequence number
	InFrmCntr   uint32 // security frame counter..
	TxFailure   uint16 // higher values indicate more failures
}

type AgingEndDevice struct {
	EndDevCfg     uint8
	DeviceTimeout uint32
}

type Device struct {
	ShortAddr      string `hex:"2"` // Short address of associated device, or invalid 0xfffe
	AddrIdx        uint16 // Index from the address manager
	NodeRelation   uint8
	DevStatus      uint8 // bitmap of various status values
	AssocCnt       uint8
	Age            uint8
	LinkInfo       *LinkInfo
	EndDev         *AgingEndDevice
	TimeoutCounter uint32
	KeepaliveRcv   uint8
}

type UtilAssocFindDevice struct {
	Number uint8
}

type UtilAssocFindDeviceResponse struct {
	Device *Device
}

//UtilAssocFindDevice is a proxy call to the AssocFindDevice() function.
func (znp *Znp) UtilAssocFindDevice(number uint8) (rsp *UtilAssocFindDeviceResponse, err error) {
	req := &UtilAssocFindDevice{Number: number}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x49, req, &rsp)
	return
}

type UtilAssocGetWithAddr struct {
	ExtAddr string `hex:"8"`
	NwkAddr string `hex:"2"`
}

type UtilAssocGetWithAddrResponse struct {
	Device *Device
}

//UtilAssocGetWithAddr is a proxy call to the AssocGetWithAddress() function.
func (znp *Znp) UtilAssocGetWithAddr(extAddr string, nwkAddr string) (rsp *UtilAssocGetWithAddrResponse, err error) {
	req := &UtilAssocGetWithAddr{ExtAddr: extAddr, NwkAddr: nwkAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x4A, req, &rsp)
	return
}

type UtilBindAddEntry struct {
	AddrMode    AddrMode
	DstAddr     string `hex:"8"`
	DstEndpoint uint8
	ClusterIDs  []uint16 `size:"1"`
}

type BindEntry struct {
	SrcEP         uint8
	DstGroupMode  uint8
	DstIdx        uint16
	DstEP         uint8
	ClusterIDList []uint16 `size:"1"`
}

type UtilBindAddEntryResponse struct {
	BindEntry *BindEntry
}

//UtilBindAddEntry is a proxy call to the bindAddEntry() function
func (znp *Znp) UtilBindAddEntry(addrMode AddrMode, dstAddr string, dstEndpoint uint8, clusterIds []uint16) (rsp *UtilBindAddEntryResponse, err error) {
	req := &UtilBindAddEntry{AddrMode: addrMode, DstAddr: dstAddr, DstEndpoint: dstEndpoint, ClusterIDs: clusterIds}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x4D, req, &rsp)
	return
}

type UtilZclKeyEstInitEst struct {
	TaskID   uint8
	SeqNum   uint8
	EndPoint uint8
	AddrMode AddrMode
	Addr     string `hex:"8"`
}

//UtilZclKeyEstInitEst is a proxy call to zclGeneral_KeyEstablish_InitiateKeyEstablishment().
func (znp *Znp) UtilZclKeyEstInitEst(taskId uint8, seqNum uint8, endPoint uint8, addrMode AddrMode, addr string) (rsp *StatusResponse, err error) {
	req := &UtilZclKeyEstInitEst{TaskID: taskId, SeqNum: seqNum, EndPoint: endPoint, AddrMode: addrMode, Addr: addr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x80, req, &rsp)
	return
}

type UtilZclKeyEstSign struct {
	Input []uint8 `size:"1"`
}

type UtilZclKeyEstSignResponse struct {
	Status Status
	Key    [42]uint8
}

//UtilZclKeyEstSign is a proxy call to zclGeneral_KeyEstablishment_ECDSASign().
func (znp *Znp) UtilZclKeyEstSign(input []uint8) (rsp *UtilZclKeyEstSignResponse, err error) {
	req := &UtilZclKeyEstSign{Input: input}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x81, req, &rsp)
	return
}

type UtilSrngGenResponse struct {
	SecureRandomNumbers [100]uint8
}

//UtilSrngGen is used to generate Secure Random Number. It generates 1,000,000 bits in sets of
//100 bytes. As in 100 bytes of secure random numbers are generated until 1,000,000 bits are
//generated. 100 bytes are generate
func (znp *Znp) UtilSrngGen() (rsp *UtilSrngGenResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x4C, nil, &rsp)
	return
}

type UtilSyncReq struct {
}

//UtilSyncReq is an asynchronous request/response handshake.
func (znp *Znp) UtilSyncReq() (err error) {
	err = znp.ProcessRequest(unpi.C_AREQ, unpi.S_UTIL, 0xE0, nil, nil)
	return
}

type UtilZclKeyEstablishInd struct {
	TaskId   uint8
	Event    uint8
	Status   uint8
	WaitTime uint8
	Suite    uint16
}

type ZdoNwkAddrReq struct {
	IEEEAddress string `hex:"8"`
	ReqType     ReqType
	StartIndex  uint8
}

//ZdoNwkAddrReq will request the device to send a “Network Address Request”. This message sends a
//broadcast message looking for a 16 bit address with a known 64 bit IEEE address. You must
//subscribe to “ZDO Network Address Response” to receive the response to this message. Check
//section 3.0.1.7 for more details on callback subscription. The response message listed below only
//indicates whether or not the message was received properly.
func (znp *Znp) ZdoNwkAddrReq(ieeeAddress string, reqType ReqType, startIndex uint8) (rsp *StatusResponse, err error) {
	req := &ZdoNwkAddrReq{IEEEAddress: ieeeAddress, ReqType: reqType, StartIndex: startIndex}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x00, req, &rsp)
	return
}

type ZdoIeeeAddrReq struct {
	ShortAddr  string `hex:"2"`
	ReqType    ReqType
	StartIndex uint8
}

//ZdoIeeeAddrReq will request a device’s IEEE 64-bit address. You must subscribe to “ZDO IEEE
//Address Response” to receive the data response to this message. The response message listed
//below only indicates whether or not the message was received properly.
func (znp *Znp) ZdoIeeeAddrReq(shortAddr string, reqType ReqType, startIndex uint8) (rsp *StatusResponse, err error) {
	req := &ZdoIeeeAddrReq{ShortAddr: shortAddr, ReqType: reqType, StartIndex: startIndex}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x01, req, &rsp)
	return
}

type ZdoNodeDescReq struct {
	DstAddr           string `hex:"2"`
	NWKAddrOfInterest string `hex:"2"`
}

//ZdoNodeDescReq is generated to inquire about the Node Descriptor information of the destination
//device.
func (znp *Znp) ZdoNodeDescReq(dstAddr string, nwkAddrOfInterest string) (rsp *StatusResponse, err error) {
	req := &ZdoNodeDescReq{DstAddr: dstAddr, NWKAddrOfInterest: nwkAddrOfInterest}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x02, req, &rsp)
	return
}

type ZdoPowerDescReq struct {
	DstAddr           string `hex:"2"`
	NWKAddrOfInterest string `hex:"2"`
}

//ZdoPowerDescReq is generated to inquire about the Power Descriptor information of the destination
//device.
func (znp *Znp) ZdoPowerDescReq(dstAddr string, nwkAddrOfInterest string) (rsp *StatusResponse, err error) {
	req := &ZdoPowerDescReq{DstAddr: dstAddr, NWKAddrOfInterest: nwkAddrOfInterest}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x03, req, &rsp)
	return
}

type ZdoSimpleDescReq struct {
	DstAddr           string `hex:"2"`
	NWKAddrOfInterest string `hex:"2"`
	Endpoint          uint8
}

//ZdoSimpleDescReq is generated to inquire as to the Simple Descriptor of the destination device’s
//Endpoint.
func (znp *Znp) ZdoSimpleDescReq(dstAddr string, nwkAddrOfInterest string, endpoint uint8) (rsp *StatusResponse, err error) {
	req := &ZdoSimpleDescReq{DstAddr: dstAddr, NWKAddrOfInterest: nwkAddrOfInterest, Endpoint: endpoint}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x04, req, &rsp)
	return
}

type ZdoActiveEpReq struct {
	DstAddr           string `hex:"2"`
	NWKAddrOfInterest string `hex:"2"`
}

//ZdoActiveEpReq is generated to request a list of active endpoint from the destination device
func (znp *Znp) ZdoActiveEpReq(dstAddr string, nwkAddrOfInterest string) (rsp *StatusResponse, err error) {
	req := &ZdoActiveEpReq{DstAddr: dstAddr, NWKAddrOfInterest: nwkAddrOfInterest}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x05, req, &rsp)
	return
}

type ZdoMatchDescReq struct {
	DstAddr           string `hex:"2"`
	NWKAddrOfInterest string `hex:"2"`
	ProfileID         uint16
	InClusterList     []uint16 `size:"1"`
	OutClusterList    []uint16 `size:"1"`
}

//ZdoMatchDescReq is generated to request the device match descriptor
func (znp *Znp) ZdoMatchDescReq(dstAddr string, nwkAddrOfInterest string, profileId uint16,
	inClusterList []uint16, outClusterList []uint16) (rsp *StatusResponse, err error) {
	req := &ZdoMatchDescReq{DstAddr: dstAddr, NWKAddrOfInterest: nwkAddrOfInterest, ProfileID: profileId,
		InClusterList: inClusterList, OutClusterList: outClusterList}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x06, req, &rsp)
	return
}

type ZdoComplexDescReq struct {
	DstAddr           string `hex:"2"`
	NWKAddrOfInterest string `hex:"2"`
}

//ZdoComplexDescReq is generated to request for the destination device’s complex descriptor.
func (znp *Znp) ZdoComplexDescReq(dstAddr string, nwkAddrOfInterest string) (rsp *StatusResponse, err error) {
	req := &ZdoComplexDescReq{DstAddr: dstAddr, NWKAddrOfInterest: nwkAddrOfInterest}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x07, req, &rsp)
	return
}

type ZdoUserDescReq struct {
	DstAddr           string `hex:"2"`
	NWKAddrOfInterest string `hex:"2"`
}

//ZdoUserDescReq is generated to request for the destination device’s user descriptor
func (znp *Znp) ZdoUserDescReq(dstAddr string, nwkAddrOfInterest string) (rsp *StatusResponse, err error) {
	req := &ZdoUserDescReq{DstAddr: dstAddr, NWKAddrOfInterest: nwkAddrOfInterest}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x08, req, &rsp)
	return
}

type CapInfo struct {
	AlternatePANCoordinator uint8 `bits:"0b00000001" bitmask:"start"`
	Router                  uint8 `bits:"0b00000010"`
	MainPowered             uint8 `bits:"0b00000100"`
	ReceiverOnWhenIdle      uint8 `bits:"0b00001000"`
	Reserved1               uint8 `bits:"0b00010000"`
	Reserved2               uint8 `bits:"0b00100000"`
	Security                uint8 `bits:"0b01000000"`
	AllocAddr               uint8 `bits:"0b10000000" bitmask:"end"`
}

type ZdoEndDeviceAnnce struct {
	NwkAddr      string `hex:"2"`
	IEEEAddr     string `hex:"8"`
	Capabilities *CapInfo
}

//ZdoEndDeviceAnnce will cause the device to issue an “End device announce” broadcast packet to the
//network. This is typically used by an end-device to announce itself to the network.
func (znp *Znp) ZdoEndDeviceAnnce(nwkAddr string, ieeeAddr string, capabilities *CapInfo) (rsp *StatusResponse, err error) {
	req := &ZdoEndDeviceAnnce{NwkAddr: nwkAddr, IEEEAddr: ieeeAddr, Capabilities: capabilities}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x0A, req, &rsp)
	return
}

type ZdoUserDescSet struct {
	DstAddr           string `hex:"2"`
	NWKAddrOfInterest string `hex:"2"`
	UserDescriptor    string `size:"1"`
}

//ZdoUserDescSet is generated to write a User Descriptor value to the targeted device.
func (znp *Znp) ZdoUserDescSet(dstAddr string, nwkAddrOfInterest string, userDescriptor string) (rsp *StatusResponse, err error) {
	req := &ZdoUserDescSet{DstAddr: dstAddr, NWKAddrOfInterest: nwkAddrOfInterest, UserDescriptor: userDescriptor}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x0B, req, &rsp)
	return
}

type ServerMask struct {
	PrimTrustCenter uint16 `bits:"0x01" bitmask:"start"`
	BkupTrustCenter uint16 `bits:"0x02"`
	PrimBindTable   uint16 `bits:"0x04"`
	BkupBindTable   uint16 `bits:"0x08"`
	PrimDiscTable   uint16 `bits:"0x10"`
	BkupDiscTable   uint16 `bits:"0x20"`
	NetworkManager  uint16 `bits:"0x40" bitmask:"end"`
}

type ZdoServerDiscReq struct {
	ServerMask *ServerMask
}

//ZdoServerDiscReq is used for local device to discover the location of a particular system server or
//servers as indicated by the ServerMask parameter. The destination addressing on this request is
//‘broadcast to all RxOnWhenIdle devices’.
func (znp *Znp) ZdoServerDiscReq(serverMask *ServerMask) (rsp *StatusResponse, err error) {
	req := &ZdoServerDiscReq{ServerMask: serverMask}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x0C, req, &rsp)
	return
}

type ZdoEndDeviceBindReq struct {
	DstAddr              string `hex:"2"`
	LocalCoordinatorAddr string `hex:"2"`
	IEEEAddr             string `hex:"8"`
	Endpoint             uint8
	ProfileID            uint16
	InClusterList        []uint16 `size:"1"`
	OutClusterList       []uint16 `size:"1"`
}

//ZdoEndDeviceBindReq is generated to request an End Device Bind with the destination device.
func (znp *Znp) ZdoEndDeviceBindReq(dstAddr string, localCoordinatorAddr string, ieeeAddr string, endpoint uint8,
	profileId uint16, inClusterList []uint16, outClusterList []uint16) (rsp *StatusResponse, err error) {
	req := &ZdoEndDeviceBindReq{DstAddr: dstAddr, LocalCoordinatorAddr: localCoordinatorAddr, IEEEAddr: ieeeAddr,
		Endpoint: endpoint, ProfileID: profileId, InClusterList: inClusterList, OutClusterList: outClusterList}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x20, req, &rsp)
	return
}

type ZdoBindUnbindReq struct {
	DstAddr     string `hex:"2"`
	SrcAddress  string `hex:"8"`
	SrcEndpoint uint8
	ClusterID   uint16
	DstAddrMode AddrMode
	DstAddress  string `hex:"8"`
	DstEndpoint uint8
}

//ZdoBindReq is generated to request an End Device Bind with the destination device.
func (znp *Znp) ZdoBindReq(dstAddr string, srcAddress string, srcEndpoint uint8, clusterId uint16,
	dstAddrMode AddrMode, dstAddress string, dstEndpoint uint8) (rsp *StatusResponse, err error) {
	req := &ZdoBindUnbindReq{DstAddr: dstAddr, SrcAddress: srcAddress, SrcEndpoint: srcEndpoint, ClusterID: clusterId,
		DstAddrMode: dstAddrMode, DstAddress: dstAddress, DstEndpoint: dstEndpoint}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x21, req, &rsp)
	return
}

//ZdoUnbindReq is generated to request a un-bind.
func (znp *Znp) ZdoUnbindReq(dstAddr string, srcAddress string, srcEndpoint uint8, clusterId uint16,
	dstAddrMode AddrMode, dstAddress string, dstEndpoint uint8) (rsp *StatusResponse, err error) {
	req := &ZdoBindUnbindReq{DstAddr: dstAddr, SrcAddress: srcAddress, SrcEndpoint: srcEndpoint, ClusterID: clusterId,
		DstAddrMode: dstAddrMode, DstAddress: dstAddress, DstEndpoint: dstEndpoint}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x22, req, &rsp)
	return
}

type Channels struct {
	Channel11 uint32 `bits:"0x00000800" bitmask:"start"`
	Channel12 uint32 `bits:"0x00001000"`
	Channel13 uint32 `bits:"0x00002000"`
	Channel14 uint32 `bits:"0x00004000"`
	Channel15 uint32 `bits:"0x00008000"`
	Channel16 uint32 `bits:"0x00010000"`
	Channel17 uint32 `bits:"0x00020000"`
	Channel18 uint32 `bits:"0x00040000"`
	Channel19 uint32 `bits:"0x00080000"`
	Channel20 uint32 `bits:"0x00100000"`
	Channel21 uint32 `bits:"0x00200000"`
	Channel22 uint32 `bits:"0x00400000"`
	Channel23 uint32 `bits:"0x00800000"`
	Channel24 uint32 `bits:"0x01000000"`
	Channel25 uint32 `bits:"0x02000000"`
	Channel26 uint32 `bits:"0b04000000" bitmask:"end"`
}

type ZdoMgmtNwkDiskReq struct {
	DstAddr      string `hex:"2"`
	ScanChannels *Channels
	ScanDuration uint8
	StartIndex   uint8
}

//ZdoMgmtNwkDiskReq is generated to request the destination device to perform a network discovery
func (znp *Znp) ZdoMgmtNwkDiskReq(dstAddr string, scanChannels *Channels, scanDuration uint8, startIndex uint8) (rsp *StatusResponse, err error) {
	req := &ZdoMgmtNwkDiskReq{DstAddr: dstAddr, ScanChannels: scanChannels, ScanDuration: scanDuration, StartIndex: startIndex}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x30, req, &rsp)
	return
}

type ZdoMgmtLqiReq struct {
	DstAddr    string `hex:"2"`
	StartIndex uint8
}

//ZdoMgmtLqiReq is generated to request the destination device to perform a LQI query of other
//devices in the network.
func (znp *Znp) ZdoMgmtLqiReq(dstAddr string, startIndex uint8) (rsp *StatusResponse, err error) {
	req := &ZdoMgmtLqiReq{DstAddr: dstAddr, StartIndex: startIndex}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x31, req, &rsp)
	return
}

type ZdoMgmtRtgReq struct {
	DstAddr    string `hex:"2"`
	StartIndex uint8
}

//ZdoMgmtRtgReq is generated to request the Routing Table of the destination device
func (znp *Znp) ZdoMgmtRtgReq(dstAddr string, startIndex uint8) (rsp *StatusResponse, err error) {
	req := &ZdoMgmtRtgReq{DstAddr: dstAddr, StartIndex: startIndex}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x32, req, &rsp)
	return
}

type ZdoMgmtBindReq struct {
	DstAddr    string `hex:"2"`
	StartIndex uint8
}

//ZdoMgmtBindReq is generated to request the Binding Table of the destination device.
func (znp *Znp) ZdoMgmtBindReq(dstAddr string, startIndex uint8) (rsp *StatusResponse, err error) {
	req := &ZdoMgmtBindReq{DstAddr: dstAddr, StartIndex: startIndex}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x33, req, &rsp)
	return
}

type RemoveChildrenRejoin struct {
	Rejoin         uint8 `bits:"0b00000001" bitmask:"start"`
	RemoveChildren uint8 `bits:"0b00000010" bitmask:"end"`
}

type ZdoMgmtLeaveReq struct {
	DstAddr              string `hex:"2"`
	DeviceAddr           string `hex:"8"`
	RemoveChildrenRejoin *RemoveChildrenRejoin
}

//ZdoMgmtLeaveReq is generated to request a Management Leave Request for the target device
func (znp *Znp) ZdoMgmtLeaveReq(dstAddr string, deviceAddr string, removeChildrenRejoin *RemoveChildrenRejoin) (rsp *StatusResponse, err error) {
	req := &ZdoMgmtLeaveReq{DstAddr: dstAddr, DeviceAddr: deviceAddr, RemoveChildrenRejoin: removeChildrenRejoin}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x34, req, &rsp)
	return
}

type ZdoMgmtDirectJoinReq struct {
	DstAddr    string `hex:"2"`
	DeviceAddr string `hex:"8"`
	CapInfo    *CapInfo
}

//ZdoMgmtDirectJoinReq is generated to request the Management Direct Join Request of a designated
//device.
func (znp *Znp) ZdoMgmtDirectJoinReq(dstAddr string, deviceAddr string, capInfo *CapInfo) (rsp *StatusResponse, err error) {
	req := &ZdoMgmtDirectJoinReq{DstAddr: dstAddr, DeviceAddr: deviceAddr, CapInfo: capInfo}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x35, req, &rsp)
	return
}

type ZdoMgmtPermitJoinReq struct {
	AddrMode       AddrMode
	DstAddr        string `hex:"2"`
	Duration       uint8
	TCSignificance uint8
}

//ZdoMgmtPermitJoinReq is generated to set the Permit Join for the destination device.
func (znp *Znp) ZdoMgmtPermitJoinReq(addrMode AddrMode, dstAddr string, duration uint8, tcSignificance uint8) (rsp *StatusResponse, err error) {
	req := &ZdoMgmtPermitJoinReq{AddrMode: addrMode, DstAddr: dstAddr, Duration: duration, TCSignificance: tcSignificance}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x36, req, &rsp)
	return
}

type ZdoMgmtNwkUpdateReq struct {
	DstAddr      string `hex:"2"`
	DstAddrMode  AddrMode
	ChannelMask  *Channels
	ScanDuration uint8
}

//ZdoMgmtNwkUpdateReq is provided to allow updating of network configuration parameters or to request
//information from devices on network conditions in the local operating environment.
func (znp *Znp) ZdoMgmtNwkUpdateReq(dstAddr string, dstAddrMode AddrMode, channelMask *Channels, scanDuration uint8) (rsp *StatusResponse, err error) {
	req := &ZdoMgmtNwkUpdateReq{DstAddr: dstAddr, DstAddrMode: dstAddrMode, ChannelMask: channelMask, ScanDuration: scanDuration}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x37, req, &rsp)
	return
}

type ZdoMsgCbRegister struct {
	ClusterID uint16
}

//ZdoMsgCbRegister registers for a ZDO callback (see reference [3], “6. ZDO Message Requests” for
//example usage).
func (znp *Znp) ZdoMsgCbRegister(clusterId uint16) (rsp *StatusResponse, err error) {
	req := &ZdoMsgCbRegister{ClusterID: clusterId}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x3E, req, &rsp)
	return
}

type ZdoMsgCbRemove struct {
	ClusterID uint16
}

//ZdoMsgCbRemove removes a registration for a ZDO callback (see reference [3], “6. ZDO Message
//Requests” for example usage).
func (znp *Znp) ZdoMsgCbRemove(clusterId uint16) (rsp *StatusResponse, err error) {
	req := &ZdoMsgCbRemove{ClusterID: clusterId}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x3F, req, &rsp)
	return
}

type ZdoStartupFromApp struct {
	StartDelay uint16
}

type ZdoStartupFromAppResponse struct {
	Status StartupFromAppStatus
}

//ZdoStartupFromApp starts the device in the network.
func (znp *Znp) ZdoStartupFromApp(startDelay uint16) (rsp *ZdoStartupFromAppResponse, err error) {
	req := &ZdoStartupFromApp{StartDelay: startDelay}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x40, req, &rsp)
	return
}

type ZdoSetLinkKey struct {
	ShortAddr   string `hex:"2"`
	IEEEAddr    string `hex:"8"`
	LinkKeyData [16]uint8
}

//ZdoSetLinkKey starts the device in the network.
func (znp *Znp) ZdoSetLinkKey(shortAddr string, ieeeAddr string, linkKeyData [16]uint8) (rsp *StatusResponse, err error) {
	req := &ZdoSetLinkKey{ShortAddr: shortAddr, IEEEAddr: ieeeAddr, LinkKeyData: linkKeyData}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x23, req, &rsp)
	return
}

type ZdoRemoveLinkKey struct {
	IEEEAddr string `hex:"8"`
}

//ZdoRemoveLinkKey removes the application link key of a given device.
func (znp *Znp) ZdoRemoveLinkKey(ieeeAddr string) (rsp *StatusResponse, err error) {
	req := &ZdoRemoveLinkKey{IEEEAddr: ieeeAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x24, req, &rsp)
	return
}

type ZdoGetLinkKey struct {
	IEEEAddr string `hex:"8"`
}

type ZdoGetLinkKeyResponse struct {
	Status      Status
	IEEEAddr    string `hex:"8"`
	LinkKeyData [16]uint8
}

//ZdoGetLinkKey retrieves the application link key of a given device.
func (znp *Znp) ZdoGetLinkKey(ieeeAddr string) (rsp *ZdoGetLinkKeyResponse, err error) {
	req := &ZdoGetLinkKey{IEEEAddr: ieeeAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x25, req, &rsp)
	return
}

type ZdoNwkDiscoveryReq struct {
	ScanChannels *Channels
	ScanDuration uint8
}

//ZdoNwkDiscoveryReq is used to initiate a network discovery (active scan).
//Strange response SecOldFrmCount(0xa1)
func (znp *Znp) ZdoNwkDiscoveryReq(scanChannels *Channels, scanDuration uint8) (rsp *StatusResponse, err error) {
	req := &ZdoNwkDiscoveryReq{ScanChannels: scanChannels, ScanDuration: scanDuration}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x26, req, &rsp)
	return
}

type ZdoJoinReq struct {
	LogicalChannel uint8
	PanID          uint16
	ExtendedPanID  uint64 //64-bit extended PAN ID (ver. 1.1 only). If not v1.1 or don't care, use all 0xFF
	ChosenParent   string `hex:"2"`
	ParentDepth    uint8
	StackProfile   uint8
}

//ZdoJoinReq is used to request the device to join itself to a parent device on a network.
func (znp *Znp) ZdoJoinReq(logicalChannel uint8, panId uint16, extendedPanId uint64,
	chosenParent string, parentDepth uint8, stackProfile uint8) (rsp *StatusResponse, err error) {
	req := &ZdoJoinReq{LogicalChannel: logicalChannel, PanID: panId, ExtendedPanID: extendedPanId,
		ChosenParent: chosenParent, ParentDepth: parentDepth, StackProfile: stackProfile}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x27, req, &rsp)
	return
}

type ZdoSetRejoinParameters struct {
	BackoffDuration uint32
	ScanDuration    uint32
}

//ZdoSetRejoinParameters is used to set rejoin backoff duration and rejoin scan duration for an end device
func (znp *Znp) ZdoSetRejoinParameters(backoffDuration uint32, scanDuration uint32) (rsp *StatusResponse, err error) {
	req := &ZdoSetRejoinParameters{BackoffDuration: backoffDuration, ScanDuration: scanDuration}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0xCC, req, &rsp)
	return
}

type ZdoSecAddLinkKey struct {
	ShortAddress    string `hex:"2"`
	ExtendedAddress string `hex:"8"`
	Key             [16]uint8
}

//ZdoSecAddLinkKey handles the ZDO security add link key extension message.
func (znp *Znp) ZdoSecAddLinkKey(shortAddress string, extendedAddress string, key [16]uint8) (rsp *StatusResponse, err error) {
	req := &ZdoSecAddLinkKey{ShortAddress: shortAddress, ExtendedAddress: extendedAddress, Key: key}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x42, req, &rsp)
	return
}

type ZdoSecEntryLookupExt struct {
	ExtendedAddress string `hex:"8"`
	Entry           [5]uint8
}

type ZdoSecEntryLookupExtResponse struct {
	AMI                  uint16
	KeyNVID              uint16
	AuthenticationOption uint8
}

//ZdoSecEntryLookupExt handles the ZDO security entry lookup extended extension message
func (znp *Znp) ZdoSecEntryLookupExt(extendedAddress string, entry [5]uint8) (rsp *ZdoSecEntryLookupExtResponse, err error) {
	req := &ZdoSecEntryLookupExt{ExtendedAddress: extendedAddress, Entry: entry}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x43, req, &rsp)
	return
}

type ZdoSecDeviceRemove struct {
	ExtendedAddress string `hex:"8"`
}

//ZdoSecDeviceRemove handles the ZDO security remove device extended extension message.
func (znp *Znp) ZdoSecDeviceRemove(extendedAddress string) (rsp *StatusResponse, err error) {
	req := &ZdoSecDeviceRemove{ExtendedAddress: extendedAddress}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x44, req, &rsp)
	return
}

type ZdoExtRouteDisc struct {
	DestinationAddress string `hex:"2"`
	Options            uint8
	Radius             uint8
}

//ZdoExtRouteDisc handles the ZDO route discovery extension message.
func (znp *Znp) ZdoExtRouteDisc(destinationAddress string, options uint8, radius uint8) (rsp *StatusResponse, err error) {
	req := &ZdoExtRouteDisc{DestinationAddress: destinationAddress, Options: options, Radius: radius}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x45, req, &rsp)
	return
}

type ZdoExtRouteCheck struct {
	DestinationAddress string `hex:"2"`
	RTStatus           uint8
	Options            uint8
}

//ZdoExtRouteCheck handles the ZDO route check extension message.
func (znp *Znp) ZdoExtRouteCheck(destinationAddress string, rtStatus uint8, options uint8) (rsp *StatusResponse, err error) {
	req := &ZdoExtRouteCheck{DestinationAddress: destinationAddress, RTStatus: rtStatus, Options: options}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x46, req, &rsp)
	return
}

type ZdoExtRemoveGroup struct {
	Endpoint uint8
	GroupID  uint16
}

//ZdoExtRemoveGroup handles the ZDO extended remove group extension message.
func (znp *Znp) ZdoExtRemoveGroup(endpoint uint8, groupId uint16) (rsp *StatusResponse, err error) {
	req := &ZdoExtRemoveGroup{Endpoint: endpoint, GroupID: groupId}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x47, req, &rsp)
	return
}

type ZdoExtRemoveAllGroup struct {
	Endpoint uint8
}

//ZdoExtRemoveAllGroup handles the ZDO extended remove all group extension message.
func (znp *Znp) ZdoExtRemoveAllGroup(endpoint uint8) (rsp *StatusResponse, err error) {
	req := &ZdoExtRemoveAllGroup{Endpoint: endpoint}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x48, req, &rsp)
	return
}

type ZdoExtFindAllGroupsEndpoint struct {
	Endpoint  uint8
	GroupList []uint16 `size:"1"`
}

type ZdoExtFindAllGroupsEndpointResponse struct {
	Groups []uint16 `size:"1"`
}

//ZdoExtFindAllGroupsEndpoint handles the ZDO extension find all groups for endpoint message
func (znp *Znp) ZdoExtFindAllGroupsEndpoint(endpoint uint8, groupList []uint16) (rsp *ZdoExtFindAllGroupsEndpointResponse, err error) {
	req := &ZdoExtFindAllGroupsEndpoint{Endpoint: endpoint, GroupList: groupList}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x49, req, &rsp)
	return
}

type ZdoExtFindGroup struct {
	Endpoint uint8
	GroupID  uint16
}

type ZdoExtFindGroupResponse struct {
	Status  Status
	GroupID uint16
	Name    string `size:"1"`
}

//ZdoExtFindGroup handles the ZDO extension find all groups for endpoint message
func (znp *Znp) ZdoExtFindGroup(endpoint uint8, groupID uint16) (rsp *ZdoExtFindGroupResponse, err error) {
	req := &ZdoExtFindGroup{Endpoint: endpoint, GroupID: groupID}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x4A, req, &rsp)
	return
}

type ZdoExtAddGroup struct {
	Endpoint  uint8
	GroupID   uint16
	GroupName string `size:"1"`
}

//ZdoExtAddGroup handles the ZDO extension add group message.
func (znp *Znp) ZdoExtAddGroup(endpoint uint8, groupID uint16, groupName string) (rsp *StatusResponse, err error) {
	req := &ZdoExtAddGroup{Endpoint: endpoint, GroupID: groupID, GroupName: groupName}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x4B, req, &rsp)
	return
}

type ZdoExtCountAllGroupsResponse struct {
	Count uint8
}

//ZdoExtCountAllGroups handles the ZDO extension count all groups message.
func (znp *Znp) ZdoExtCountAllGroups() (rsp *ZdoExtCountAllGroupsResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x4C, nil, &rsp)
	return
}

type ZdoExtRxIdle struct {
	SetFlag  uint8
	SetValue uint8
}

//ZdoExtRxIdle handles the ZDO extension Get/Set RxOnIdle to ZMac message
func (znp *Znp) ZdoExtRxIdle(setFlag uint8, setValue uint8) (rsp *StatusResponse, err error) { //very unclear from the docs and the code
	req := &ZdoExtRxIdle{SetFlag: setFlag, SetValue: setValue}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x4D, req, &rsp)
	return
}

type ZdoExtUpdateNwkKey struct {
	DestinationAddress string `hex:"2"`
	KeySeqNum          uint8
	Key                [128]uint8
}

//ZdoExtUpdateNwkKey handles the ZDO security update network key extension message.
func (znp *Znp) ZdoExtUpdateNwkKey(destinationAddress string, keySeqNum uint8, key [128]uint8) (rsp *StatusResponse, err error) {
	req := &ZdoExtUpdateNwkKey{DestinationAddress: destinationAddress, KeySeqNum: keySeqNum, Key: key}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x4E, req, &rsp)
	return
}

type ZdoExtSwitchNwkKey struct {
	DestinationAddress string `hex:"2"`
	KeySeqNum          uint8
}

//ZdoExtSwitchNwkKey handles the ZDO security switch network key extension message.
func (znp *Znp) ZdoExtSwitchNwkKey(destinationAddress string, keySeqNum uint8) (rsp *StatusResponse, err error) {
	req := &ZdoExtSwitchNwkKey{DestinationAddress: destinationAddress, KeySeqNum: keySeqNum}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x4F, req, &rsp)
	return
}

type ZdoExtNwkInfoResponse struct {
	ShortAddress          string `hex:"2"`
	PanID                 uint16
	ParentAddress         string `hex:"2"`
	ExtendedPanID         uint64
	ExtendedParentAddress string `hex:"8"`
	Channel               uint16 //uint16 or uint8?????
}

//ZdoExtNwkInfo handles the ZDO extension network message.
func (znp *Znp) ZdoExtNwkInfo() (rsp *ZdoExtNwkInfoResponse, err error) {
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x50, nil, &rsp)
	return
}

type ZdoExtSeqApsRemoveReq struct {
	NwkAddress      string `hex:"2"`
	ExtendedAddress string `hex:"8"`
	ParentAddress   string `hex:"2"`
}

//ZdoExtSeqApsRemoveReq handles the ZDO extension Security Manager APS Remove Request message.
func (znp *Znp) ZdoExtSeqApsRemoveReq(nwkAddress string, extendedAddress string, parentAddress string) (rsp *StatusResponse, err error) {
	req := &ZdoExtSeqApsRemoveReq{NwkAddress: nwkAddress, ExtendedAddress: extendedAddress, ParentAddress: parentAddress}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x51, req, &rsp)
	return
}

//ZdoForceConcentratorChange forces a network concentrator change by resetting zgConcentratorEnable and
//zgConcentratorDiscoveryTime from NV and set nwk event.
func (znp *Znp) ZdoForceConcentratorChange() error {
	return znp.ProcessRequest(unpi.C_AREQ, unpi.S_ZDO, 0x52, nil, nil)
}

type ZdoExtSetParams struct {
	UseMulticast uint8
}

//ZdoExtSeqApsRemoveReq set parameters not settable through NV.
func (znp *Znp) ZdoExtSetParams(useMulticast uint8) (rsp *StatusResponse, err error) {
	req := &ZdoExtSetParams{UseMulticast: useMulticast}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x53, req, &rsp)
	return
}

type ZdoNwkAddrOfInterestReq struct {
	DestAddr          string `hex:"2"`
	NwkAddrOfInterest string `hex:"2"`
	Cmd               uint8
}

//ZdoNwkAddrOfInterestReq handles ZDO network address of interest request.
func (znp *Znp) ZdoNwkAddrOfInterestReq(destAddr string, nwkAddrOfInterest string, cmd uint8) (rsp *StatusResponse, err error) {
	req := &ZdoNwkAddrOfInterestReq{DestAddr: destAddr, NwkAddrOfInterest: nwkAddrOfInterest, Cmd: cmd}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_ZDO, 0x29, req, &rsp)
	return
}

type ZdoNwkAddrRsp struct {
	Status       Status
	IEEEAddr     string `hex:"8"`
	NwkAddr      string `hex:"2"`
	StartIndex   uint8
	AssocDevList []string `size:"1" hex:"2"`
}

type ZdoIEEEAddrRsp struct {
	Status       Status
	IEEEAddr     string `hex:"8"`
	NwkAddr      string `hex:"2"`
	StartIndex   uint8
	AssocDevList []string `size:"1" hex:"2"`
}

type LogicalType uint8

func (t LogicalType) Coordinator() bool {
	return t&0x1 > 0
}

func (t LogicalType) Router() bool {
	return t&0x2 > 0
}

func (t LogicalType) EndDevice() bool {
	return t&0x4 > 0
}

type ZdoNodeDescRsp struct {
	SrcAddr                    string `hex:"2"`
	Status                     Status
	NWKAddrOfInterest          string      `hex:"2"`
	LogicalType                LogicalType `bits:"0b00000011" bitmask:"start"`
	ComplexDescriptorAvailable uint8       `bits:"0b00001000"`
	UserDescriptorAvailable    uint8       `bits:"0b00010000"  bitmask:"end"`
	APSFlags                   uint8       `bits:"0b00011111" bitmask:"start"`
	FrequencyBand              uint8       `bits:"0b11100000" bitmask:"end"`
	MacCapabilitiesFlags       *CapInfo
	ManufacturerCode           uint16
	MaxBufferSize              uint8
	MaxInTransferSize          uint16
	ServerMask                 *ServerMask
	MaxOutTransferSize         uint16
	DescriptorCapabilities     uint8
}

type ZdoPowerDescRsp struct {
	SrcAddr                 string `hex:"2"`
	Status                  Status
	NWKAddr                 string `hex:"2"`
	CurrentPowerMode        uint8  `bits:"0b00001111" bitmask:"start"`
	AvailablePowerSources   uint8  `bits:"0b11110000"  bitmask:"end"`
	CurrentPowerSource      uint8  `bits:"0b00001111" bitmask:"start"`
	CurrentPowerSourceLevel uint8  `bits:"0b11110000"  bitmask:"end"`
}

type ZdoSimpleDescRsp struct {
	SrcAddr        string `hex:"2"`
	Status         Status
	NWKAddr        string `hex:"2"`
	Len            uint8
	Endpoint       uint8
	ProfileID      uint16
	DeviceID       uint16
	DeviceVersion  uint8
	InClusterList  []uint16 `size:"1"`
	OutClusterList []uint16 `size:"1"`
}

type ZdoActiveEpRsp struct {
	SrcAddr      string `hex:"2"`
	Status       Status
	NWKAddr      string  `hex:"2"`
	ActiveEPList []uint8 `size:"1"`
}

type ZdoMatchDescRsp struct {
	SrcAddr   string `hex:"2"`
	Status    Status
	NWKAddr   string  `hex:"2"`
	MatchList []uint8 `size:"1"`
}

type ZdoComplexDescRsp struct {
	SrcAddr           string `hex:"2"`
	Status            Status
	NWKAddr           string `hex:"2"`
	ComplexDescriptor string `size:"1"`
}

type ZdoUserDescRsp struct {
	SrcAddr        string `hex:"2"`
	Status         Status
	NWKAddr        string `hex:"2"`
	UserDescriptor string `size:"1"`
}

type ZdoUserDescConf struct {
	SrcAddr string `hex:"2"`
	Status  Status
	NWKAddr string `hex:"2"`
}

type ZdoServerDiscRsp struct {
	SrcAddr    string `hex:"2"`
	Status     Status
	ServerMask *ServerMask
}

type ZdoEndDeviceBindRsp struct {
	SrcAddr string `hex:"2"`
	Status  Status
}

type ZdoBindRsp struct {
	SrcAddr string `hex:"2"`
	Status  Status
}

type ZdoUnbindRsp struct {
	SrcAddr string `hex:"2"`
	Status  Status
}

type Network struct {
	PanID           uint16 `bound:"8"`
	LogicalChannel  uint8
	StackProfile    uint8 `bits:"0b00001111" bitmask:"start"`
	ZigbeeVersion   uint8 `bits:"0b11110000" bitmask:"end"`
	BeaconOrder     uint8 `bits:"0b00001111" bitmask:"start"`
	SuperFrameOrder uint8 `bits:"0b11110000" bitmask:"end"`
	PermitJoin      uint8
}

type ZdoMgmtNwkDiscRsp struct {
	SrcAddr      string `hex:"2"`
	Status       Status
	NetworkCount uint8
	StartIndex   uint8
	NetworkList  []*Network `size:"1"`
}

type LqiDeviceType uint8

const (
	Coordinator LqiDeviceType = 0x00
	Router      LqiDeviceType = 0x01
	EndDevice   LqiDeviceType = 0x02
)

type NeighborLqi struct {
	ExtendedPanID   uint64
	ExtendedAddress string        `hex:"8"`
	NetworkAddress  string        `hex:"2"`
	DeviceType      LqiDeviceType `bits:"0b00000011" bitmask:"start"`
	RxOnWhenIdle    uint8         `bits:"0b00001100"`
	Relationship    uint8         `bits:"0b00110000"`
	PermitJoining   uint8
	Depth           uint8
	LQI             uint8
}

type ZdoMgmtLqiRsp struct {
	SrcAddr              string `hex:"2"`
	Status               Status
	NeighborTableEntries uint8
	StartIndex           uint8
	NeighborLqiList      []*NeighborLqi `size:"1"`
}

type Route struct {
	DestinationAddress string `hex:"2"`
	Status             RouteStatus
	NextHop            string `hex:"2"`
}

type ZdoMgmtRtgRsp struct {
	SrcAddr             string `hex:"2"`
	Status              Status
	RoutingTableEntries uint8
	StartIndex          uint8
	RoutingTable        []*Route `size:"1"`
}

type Addr struct {
	AddrMode     AddrMode
	ShortAddr    string `hex:"2" cond:"uint:AddrMode!=3"`
	ExtendedAddr string `hex:"8" cond:"uint:AddrMode==3"`
	DstEndpoint  uint8  `cond:"uint:AddrMode==3"`
}

type Binding struct {
	SrcAddr     string `hex:"8"`
	SrcEndpoint uint8
	ClusterID   uint16
	DstAddr     *Addr
}

type ZdoMgmtBindRsp struct {
	SrcAddr          string `hex:"2"`
	Status           Status
	BindTableEntries uint8
	StartIndex       uint8
	BindTable        []*Binding `size:"1"`
}

type ZdoMgmtLeaveRsp struct {
	SrcAddr string `hex:"2"`
	Status  Status
}

type ZdoMgmtDirectJoinRsp struct {
	SrcAddr string `hex:"2"`
	Status  Status
}

type ZdoMgmtPermitJoinRsp struct {
	SrcAddr string `hex:"2"`
	Status  Status
}

type ZdoStateChangeInd struct {
	State DeviceState
}

type ZdoEndDeviceAnnceInd struct {
	SrcAddr      string `hex:"2"`
	NwkAddr      string `hex:"2"`
	IEEEAddr     string `hex:"8"`
	Capabilities *CapInfo
}

type ZdoMatchDescRpsSent struct {
	NwkAddr        string   `hex:"2"`
	InClusterList  []uint16 `size:"1"`
	OutClusterList []uint16 `size:"1"`
}

type ZdoStatusErrorRsp struct {
	SrcAddr string `hex:"2"`
	Status  Status
}

type ZdoSrcRtgInd struct {
	DstAddr   string   `hex:"2"`
	RelayList []string `size:"1" hex:"2"`
}

type Beacon struct {
	SrcAddr         string `hex:"2"`
	PanID           uint16
	LogicalChannel  uint8
	PermitJoining   uint8
	RouterCapacity  uint8
	DeviceCapacity  uint8
	ProtocolVersion uint8
	StackProfile    uint8
	LQI             uint8
	Depth           uint8
	UpdateID        uint8
	ExtendedPanID   uint64
}

type ZdoBeaconNotifyInd struct {
	BeaconList []*Beacon `size:"1"`
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

	//SYS
	AsyncCommandRegistry[registryKey{unpi.S_SYS, 0x80}] = &SysResetInd{}
	AsyncCommandRegistry[registryKey{unpi.S_SYS, 0x81}] = &SysOsalTimerExpired{}

	//UTIL
	AsyncCommandRegistry[registryKey{unpi.S_UTIL, 0xE0}] = &UtilSyncReq{}
	AsyncCommandRegistry[registryKey{unpi.S_UTIL, 0xE1}] = &UtilZclKeyEstablishInd{}

	//ZDO
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x80}] = &ZdoNwkAddrRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x81}] = &ZdoIEEEAddrRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x82}] = &ZdoNodeDescRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x83}] = &ZdoPowerDescRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x84}] = &ZdoSimpleDescRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x85}] = &ZdoActiveEpRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x86}] = &ZdoMatchDescRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x87}] = &ZdoComplexDescRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x88}] = &ZdoUserDescRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x89}] = &ZdoUserDescConf{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0x8A}] = &ZdoServerDiscRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xA0}] = &ZdoEndDeviceBindRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xA1}] = &ZdoBindRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xA2}] = &ZdoUnbindRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xB0}] = &ZdoMgmtNwkDiscRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xB1}] = &ZdoMgmtLqiRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xB2}] = &ZdoMgmtRtgRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xB3}] = &ZdoMgmtBindRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xB4}] = &ZdoMgmtLeaveRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xB5}] = &ZdoMgmtDirectJoinRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xB6}] = &ZdoMgmtPermitJoinRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xC0}] = &ZdoStateChangeInd{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xC1}] = &ZdoEndDeviceAnnceInd{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xC2}] = &ZdoMatchDescRpsSent{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xC3}] = &ZdoStatusErrorRsp{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xC4}] = &ZdoSrcRtgInd{}
	AsyncCommandRegistry[registryKey{unpi.S_ZDO, 0xC5}] = &ZdoBeaconNotifyInd{}
}
