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
	ItemCreatedAndInitialized Status = iota + 0x09
	InitializationFailed
	BadLength Status = iota + 0x0C
	MemError  Status = iota + 0x10
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
	transId uint8, options *AfDataRequestOptions, radius uint8, data []uint8) (rsp *StatusResponse, err error) {
	req := &AfDataRequest{DstAddr: dstAddr, DstEndpoint: dstEndpoint, SrcEndpoint: srcEndpoint,
		ClusterID: clusterId, TransID: transId, Options: options, Radius: radius, Data: data}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_AF, 0x01, req, &rsp)
	return
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
	Data  []uint8 `len:"uint8"`
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
	Data   []uint8 `len:"uint8"`
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
	String string `len:"uint8"`
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
	Destination string `hex:"uint16"`
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
	Destination string `hex:"uint64"`
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
	Destination string `hex:"uint16"`
	CommandID   uint16
	Handle      uint8
	Ack         uint8
	Radius      uint8
	Data        []uint8 `len:"uint8"`
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
	Value    []uint8 `len:"uint8"`
}

func (znp *Znp) SapiZbReadConfiguration(configID uint8) (rsp *SapiZbReadConfigurationResponse, err error) {
	req := &SapiZbReadConfiguration{ConfigID: configID}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SAPI, 0x04, req, &rsp)
	return
}

type SapiZbWriteConfiguration struct {
	ConfigID uint8
	Value    []uint8 `len:"uint8"`
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
	SearchKey string `hex:"uint64"`
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
	ExtAddress string `hex:"uint64"` //The device’s extended address.
}

//SysSetExtAddr is used to set the extended address of the device
func (znp *Znp) SysSetExtAddr(extAddr string) (rsp *StatusResponse, err error) {
	req := &SysSetExtAddr{ExtAddress: extAddr}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_SYS, 0x03, req, &rsp)
	return
}

type SysGetExtAddrResponse struct {
	ExtAddress string `hex:"uint64"` //The device’s extended address.
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
	Value  []uint8 `len:"uint8"` //The value read from the target RAM.
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
	Value   []uint8 `len:"uint8"` //The value written to the target RAM.
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
	Value  []uint8 `len:"uint8"`
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
	Value  []uint8 `len:"uint8"`
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
	InitData []uint8 `len:"uint8"`
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
	Value  []uint8 `len:"uint8"`
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
	Value  []uint8 `len:"uint8"`
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
	Value  []uint8 `len:"uint8"`
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
	Value  []uint8 `len:"uint8"`
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
	IEEEAddr         string `hex:"uint64"`
	ShortAddr        string `hex:"uint16"`
	DeviceType       *DeviceType
	DeviceState      DeviceState
	AssocDevicesList []string `len:"uint8" hex:"uint16"`
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
	IEEEAddr      string `hex:"uint64"`
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
	Address  string `hex:"uint64"`
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
	Address  string `hex:"uint64"`
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
	Address  string `hex:"uint64"`
	PanID    uint16
}

//UtilSrcMatchCheckSrcAddr is used to delete a short or extended address from the source address table.
func (znp *Znp) UtilSrcMatchCheckSrcAddr(addrMode AddrMode, address string, panId uint16) (rsp *StatusResponse, err error) {
	req := &UtilSrcMatchCheckSrcAddr{AddrMode: addrMode, Address: address, PanID: panId}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x23, req, &rsp)
	return
}

type UtilSrcMatchAckAllPending struct {
	Action Action
}

//UtilSrcMatchAckAllPending is used to enable/disable acknowledging all packets with pending bit set.
func (znp *Znp) UtilSrcMatchAckAllPending(action Action) (rsp *StatusResponse, err error) {
	req := &UtilSrcMatchAckAllPending{Action: action}
	err = znp.ProcessRequest(unpi.C_SREQ, unpi.S_UTIL, 0x24, req, &rsp)
	return
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
