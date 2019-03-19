package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/dyrkin/unp-go"
	"github.com/dyrkin/znp-go"
	"go.bug.st/serial.v1"
)

func main() {
	spew.Config.DisableCapacities = true
	spew.Config.DisablePointerAddresses = true
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open("/dev/tty.usbmodem14101", mode)
	if err != nil {
		log.Fatalf("Can't open port. Reason: %s", err)
	}
	port.SetRTS(true)

	u := unp.New(1, port)
	z := znp.New(u)
	z.Start()

	go func() {
		for {
			select {
			case err := <-z.Errors():
				fmt.Printf("Error received: %s\n", err)
			case async := <-z.AsyncInbound():
				fmt.Printf("Async received: %s\n", spew.Sdump(async))
			}
		}
	}()

	go func() {
		for {
			select {
			case frame := <-z.OutFramesLog():
				fmt.Printf("Frame sent: %s\n", spew.Sdump(frame))
			case frame := <-z.InFramesLog():
				fmt.Printf("Frame received: %s\n", spew.Sdump(frame))
			}
		}
	}()

	// z.SysResetReq(1)

	var res interface{}

	res, err = z.SysPing()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysVersion()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysSetExtAddr("0x00124b00019c2ee9")
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilCallbackSubCmd(znp.SubsystemIdZdo, znp.ActionEnable)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysGetExtAddr()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SapiZbStartRequest()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SapiZbPermitJoiningRequest("0xFF00", 200)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SapiZbReadConfiguration(1)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SapiZbFindDeviceRequest("0x00124b00019c2ee9")
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysOsalStartTimer(1, 3000)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	t := time.Now()

	res, err = z.SysSetTime(0, uint8(t.Hour()), uint8(t.Minute()), uint8(t.Second()), uint8(t.Month()), uint8(t.Day()), uint16(t.Year()))
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.SysGetTime()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilGetDeviceInfo()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilGetNvInfo()
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilLoopback([]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9})
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilAssocFindDevice(1)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.UtilAssocGetWithAddr("0x0000000000000000", "0x25cc")
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	err = z.UtilSyncReq()
	if err != nil {
		log.Fatal(err)
	}

	res, err = z.ZdoNwkAddrReq("0x00124b00019c2ee9", znp.ReqTypeAssociatedDevicesResponse, 0)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.ZdoIeeeAddrReq("0x25cc", znp.ReqTypeAssociatedDevicesResponse, 0)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.ZdoUserDescSet("0x0000", "0x25cc", "hello")
	if err != nil {
		log.Fatal(err)
	}

	res, err = z.ZdoUserDescReq("0x0000", "0xe065")
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	res, err = z.ZdoServerDiscReq(&znp.ServerMask{PrimTrustCenter: 1})
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoMgmtNwkDiskReq("0x0000", &znp.Channels{Channel11: 1, Channel12: 1, Channel13: 1, Channel14: 1}, 1, 0)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoMgmtLqiReq("0x0000", 0)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoMgmtRtgReq("0x0000", 0)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	// res, err = z.ZdoBindReq("0x0000", "0x00124b00019c2ee9", 1, 30, znp.AddrModeAddr64Bit, "0x0000000000003000", 2)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// PrintStruct(res)

	res, err = z.ZdoMgmtPermitJoinReq(znp.AddrModeAddr16Bit, "0x0000", 255, 0)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoMgmtNwkUpdateReq("0x25cc", znp.AddrModeAddr16Bit, &znp.Channels{Channel11: 1, Channel12: 1, Channel13: 1, Channel14: 1}, 1)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoMsgCbRegister(1588)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoStartupFromApp(1)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoNwkDiscoveryReq(&znp.Channels{Channel11: 1, Channel12: 1, Channel13: 1, Channel14: 1}, 1)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoExtAddGroup(1, 5, "asdfghjklzxcvbn")
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoExtFindGroup(1, 5)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoExtFindAllGroupsEndpoint(1, []uint16{5})
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoExtCountAllGroups()
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoExtSwitchNwkKey("0x25cc", 0)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.UtilGetDeviceInfo()
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoNodeDescReq("0x0000", "0x0000")
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoActiveEpReq("0x0000", "0x0000")
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoPowerDescReq("0x0000", "0x0000")
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoComplexDescReq("0x0000", "0x0000")
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoUserDescReq("0x0000", "0x0000")
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoMgmtLqiReq("0x0000", 0)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.ZdoMgmtRtgReq("0x0000", 0)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	err = z.DebugMsg("hello")
	if err != nil {
		log.Fatal(err)
	}

	res, err = z.ZdoStartupFromApp(30)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	z.AfRegister(1, 0x0104, 1, 1, znp.LatencyNoLatency, []uint16{}, []uint16{})

	res, err = z.ZdoMgmtBindReq("0x0000", 1)
	if err != nil {
		log.Fatal(err)
	}

	PrintStruct(res)

	res, err = z.UtilLedControl(1, znp.ModeOFF)
	if err != nil {
		log.Fatal(err)
	}
	PrintStruct(res)

	time.Sleep(200 * time.Second)
}

func PrintStruct(v interface{}) {
	spew.Dump(v)
}
