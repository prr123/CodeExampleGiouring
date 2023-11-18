// rewrite in giouring
//
package main

import (
//	"bytes"
	"fmt"
	"os"
	"log"
	"strconv"
	"unsafe"
//	"sync"
//	"syscall"
//	"time"

	util "github.com/prr123/utility/utilLib"
    gio "github.com/pawelgaczynski/giouring"
//	"github.com/iceber/iouring-go"
)

const bufferSize = 4096

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg","file"}

    // default file

    useStr := "./gioSimpleRead /file=filenam [/dbg]"
    helpStr := "program that reads a file"

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

	directio := false
    if _, ok := flagMap["directio"]; ok {directio = true}

	rdFilnam := ""
    rdVal, ok := flagMap["file"]
    if !ok {
        log.Fatalf("error: no file flag provided!")
	} else {
        if rdVal.(string) == "none" {log.Fatalf("error: no output file name provided with /file flag!")}
        rdFilnam = rdVal.(string)
    }

	inFilnam := rdFilnam + ".dat"
	altFilnam := rdFilnam + ".alt"

    if dbg {
        fmt.Println("********* cli parameters *********")
        fmt.Printf("directio: %t\n", directio)
        fmt.Printf("file:     %s\n", inFilnam)
        fmt.Printf("alt file: %s\n", altFilnam)
        fmt.Println("**********************************")
	}

//	size := 1073741824 // 1GiB

	falt, err := os.Open(altFilnam)
	defer falt.Close()
	if err != nil {log.Fatalf("error -- Open Alt File: %v",err)}

	data, err := ReadUring(falt)
	if err != nil {log.Fatalf("error -- Read Alt File: %v",err)}
	log.Printf("Read File %s, Length %d\n", altFilnam, len(data))

/*
	fin, err := os.Open(inFilnam)
	defer fin.Close()
	if err != nil {log.Fatalf("error -- Open File: %v",err)}
*/
	data2, err := os.ReadFile(inFilnam)
	if err != nil {log.Fatalf("error -- os.ReadFile: %v", err)}
	log.Printf("Read File %s, Length %d\n", inFilnam, len(data))

	err =CmpData(data,data2)
	if err != nil {log.Fatalf("error -- cmpData: %v", err)}
	log.Printf("comparison was successful!")
	fmt.Println("success!")
}


func CmpData(data, data2 []byte) (err error) {

	dlen := len(data)
	dlen2 := len(data2)

	if dlen != dlen2 {return fmt.Errorf("unequal file sizes!")}

	for i:=0; i< dlen; i++ {
		if data[i] != data2[i] {
			return fmt.Errorf("data [%d] are not equal!", i)
		}
	}

	return nil
}
func ReadUring(fil *os.File) (data []byte, err error) {


//	numEntries := len(data)/bufferSize
//fmt.Printf("ReadUring: num of blocks: %d\n", numEntries)

	info, err := fil.Stat()
	if err != nil {return data, fmt.Errorf("fil.Stat: %v", err)}
	size := info.Size()

	data = make([]byte, int(size))


	// start iouring
	entries :=uint32(1)
	ring, err:= gio.CreateRing(entries)
	if err != nil {return data, fmt.Errorf("CreateRing: %v", err)}

	if err := ring.QueueInit(entries, 0); err != nil {return data, fmt.Errorf("QueueInit: %v", err)}

	// create sqe
	sqe := ring.GetSQE()
	DispSQE(sqe)


	// prep read
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

	return data, nil
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
	fmt.Printf("Opcode:   %d\n", sqe.OpCode)
	fmt.Printf("Flags:    %d\n", sqe.Flags)
	fmt.Printf("IOPrio:   %d\n", sqe.IoPrio)
	fmt.Printf("Fd:       %d\n", sqe.Fd)
	fmt.Printf("Off:      %d\n", sqe.Off)
	fmt.Printf("Addr:     %d\n", sqe.Addr)
	fmt.Printf("Len:      %d\n", sqe.Len)
	fmt.Printf("OpFlags:  %d\n", sqe.OpcodeFlags)
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

	rfil, err := os.Open(fn)
	if err != nil {
		return data, fmt.Errorf("os.Open %v!", err)
	}
	defer rfil.Close()

	log.Printf("rfil opened!")

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
