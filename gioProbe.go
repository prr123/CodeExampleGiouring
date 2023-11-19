// rewrite in giouring
//
package main

import (
//	"bytes"
	"fmt"
	"os"
	"log"
//	"strconv"
//	"unsafe"

	util "github.com/prr123/utility/utilLib"
    gio "github.com/pawelgaczynski/giouring"
)

//const bufferSize = 4096

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg"}

    // default file

    useStr := "./gioProbe [/dbg]"
    helpStr := "program that tests giouring probe"

	if numarg < 2 {
		fmt.Println("insufficient arguments!")
        fmt.Printf("usage: %s\n", useStr)
        os.Exit(1)
	}

    if numarg > len(flags) + 1 {
        fmt.Println("too many arguments in cl!")
        fmt.Printf("usage: %s\n", useStr)
        os.Exit(1)
    }

    if numarg > 1 && os.Args[1] == "help" {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s\n", useStr)
        os.Exit(0)
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

	dbg := false
    if _, ok := flagMap["dbg"]; ok {dbg = true}

/*
	numEntries := 1
	entryStr :=""
    val, ok := flagMap["entries"]
	if ok {
		if val.(string) == "none" {log.Fatalf("entries flag requires a namber value!")}
		entryStr = val.(string)
		tmp, err := strconv.Atoi(entryStr)
		if err != nil {log.Fatalf(" entries flag value is not an integer!")}
		numEntries = tmp
	}
*/
    if dbg {
        fmt.Println("********* cli parameters *********")
//		fmt.Printf("Num of Entries: %d\n", numEntries)
//        fmt.Printf("directio: %t\n", directio)
        fmt.Println("**********************************")
	}

	probe, err:= gio.GetProbe()
	if err != nil {log.Fatalf("error -- GetProbe: %v", err)}

	PrintProbe(probe)

	log.Printf("exit  -- success!\n")
}

func PrintProbe(probe *gio.Probe){

	fmt.Println("********** Probe ***********")

	fmt.Printf("Last Op: %d\n", probe.LastOp)
	fmt.Printf("OpsLen: %d\n", probe.OpsLen)
	fmt.Printf("Res:    %d\n", probe.Res)
	fmt.Printf("Ops [%d]: \n", len(probe.Ops))
	last:=0
	for i:=0; i< 255; i++ {
		nam, _ := GetOps(uint8(i))
		if nam == "OpLast" {last = i; break;}
	}

//	for i:=0; i< len(probe.Ops); i++ {
	for i:=0; i< last; i++ {
		itmp :=probe.Ops[i].Op
		nam, err := GetOps(itmp)
		if err != nil {
			fmt.Printf("  %d - Op error: %v\n", err)
			break
		}
		fmt.Printf(" %d - Op: %s\n",i, nam)
		if nam == "OpLast" {break}
	}
	fmt.Println("******** End Probe *********")

}

func GetOps(id uint8) (Opsnam string, err error) {

//	if (idx<0) || (idx> 255) {return "", fmt.Errorf("invalid index")}

//	id := uint8(idx)

	OpsList := [255]string{
	"OpNop",
	"OpReadv",
	"OpWritev",
	"OpFsync",
	"OpReadFixed",
	"OpWriteFixed",
	"OpPollAdd",
	"OpPollRemove",
	"OpSyncFileRange",
	"OpSendmsg",
	"OpRecvmsg",
	"OpTimeout",
	"OpTimeoutRemove",
	"OpAccept",
	"OpAsyncCancel",
	"OpLinkTimeout",
	"OpConnect",
	"OpFallocate",
	"OpOpenat",
	"OpClose",
	"OpFilesUpdate",
	"OpStatx",
	"OpRead",
	"OpWrite",
	"OpFadvise",
	"OpMadvise",
	"OpSend",
	"OpRecv",
	"OpOpenat2",
	"OpEpollCtl",
	"OpSplice",
	"OpProvideBuffers",
	"OpRemoveBuffers",
	"OpTee",
	"OpShutdown",
	"OpRenameat",
	"OpUnlinkat",
	"OpMkdirat",
	"OpSymlinkat",
	"OpLinkat",
	"OpMsgRing",
	"OpFsetxattr",
	"OpSetxattr",
	"OpFgetxattr",
	"OpGetxattr",
	"OpSocket",
	"OpUringCmd",
	"OpSendZC",
	"OpSendMsgZC",
	"OpLast",
	}

	return OpsList[id], nil
}


/*
func UringSetupTest(numEntries int) (err error) {

	// start iouring
	entries :=uint32(numEntries)
	ring, err:= gio.CreateRing(entries)
	if err != nil {return fmt.Errorf("CreateRing: %v", err)}

	if err := ring.QueueInit(entries, 0); err != nil {return fmt.Errorf("QueueInit: %v", err)}

	// create sqe
	sqe := ring.GetSQE()
	DispSQE(sqe)

	// prep read "Operation
	lenDat := uint32(len(data))
	datPtr := uintptr(unsafe.Pointer(&data[0]))
	offset := uint64(0)
	fd := int(fil.Fd())
	sqe.PrepareRead(fd, datPtr, lenDat, offset)

	numSQE, err := ring.Submit()
	if err != nil {return nil, fmt.Errorf("ringSubmit: %v", err)}
	fmt.Printf("NumSqes: %d\n", numSQE)

	cqe, err := ring.WaitCQE()
	if err != nil {return nil, fmt.Errorf("ringWait: %v", err)}
	DispCQE(cqe)

	ring.CQESeen(cqe)
	ring.QueueExit()

	return nil
}

// CQE completion queue event
func DispCQE(cqe *gio.CompletionQueueEvent) {

	fmt.Println("************* cqe **************")
	fmt.Printf("UserData: %d\n", cqe.UserData)
	fmt.Printf("Res:      %d\n", cqe.Res)
	fmt.Printf("Flags:    %d\n", cqe.Flags)
	fmt.Println("*********** end cqe ************")
	return
}

func DispSQE(sqe *gio.SubmissionQueueEntry) {

	fmt.Println("************* sqe **************")
	fmt.Printf(""Opcode:   %d\n", sqe."OpCode)
	fmt.Printf("Flags:    %d\n", sqe.Flags)
	fmt.Printf("I"Oprio:   %d\n", sqe.I"Oprio)
	fmt.Printf("Fd:       %d\n", sqe.Fd)
	fmt.Printf("Off:      %d\n", sqe.Off)
	fmt.Printf("Addr:     %d\n", sqe.Addr)
	fmt.Printf("Len:      %d\n", sqe.Len)
	fmt.Printf(""OpFlags:  %d\n", sqe."OpcodeFlags)
	fmt.Printf("UserData: %d\n", sqe.UserData)
	fmt.Printf("BufIG:    %d\n", sqe.BufIG)
	fmt.Printf("Person:   %d\n", sqe.Personality)
	fmt.Printf("SpliceFd: %d\n", sqe.SpliceFdIn)
	fmt.Printf("Addr3:    %d\n", sqe.Addr3)
	fmt.Println("*********** end sqe ************")
	return
}


// redo into blocks
func readNBytes(fn string, n int) (data []byte, err error) {

	rfil, err := os."Open(fn)
	if err != nil {
		return data, fmt.Errorf("os."Open %v!", err)
	}
	defer rfil.Close()

	log.Printf("rfil "Opened!")

	data = make([]byte, 0, n)

	var buffer = make([]byte, bufferSize)

	for len(data) < n {
		nByt, err := rfil.Read(buffer)
		if err != nil {
			return data, fmt.Errorf("rfil.Read: %v", err)
		}
//	log.Printf("read %d bytes\n", nByt)
		data = append(data, buffer[:nByt]...)
//	log.Printf("data len: %d\n", len(data))

	}

	return data, nil
}


func CvtSize(sizeStr string, two bool) (siz uint64, err error) {

    // check last letter of size
    let := sizeStr[len(sizeStr) -1]
//    fmt.Printf("last letter: %q ", let)

    // if last letter is a letter, convert the rest into a number
    var mult uint64 = 1
    switch let {
    case 'K':
        mult = 1000

    case 'M':
        mult = 1000000

    case 'G':
        mult = 1000000000

    default:
        if !util.IsNumeric(let) {
            return 0 , fmt.Errorf("let is not alphaNumeric!")
		}
	}
    intStr:=""
    if mult>1 {
        intStr = string(sizeStr[:len(sizeStr)-1])
    } else {
        intStr = string(sizeStr[:len(sizeStr)])
    }

    inum, err := strconv.Atoi(intStr)
    if err !=nil {return 0, fmt.Errorf("error -- cannot convert intStr: %s: %v", intStr, err)}
    num := uint64(inum)*uint64(mult)

//    fmt.Printf("res: %d\n", inum)

    if !two {return num, nil}
    num--
    num = num | num>>1
    num = num | num>>2
    num = num | num>>4
    num = num | num>>8
    num = num | num>>16
    num = num | num>>32
    num++

    return num, nil
}
*/
