// kernelVer
// testing giouring
// Author prr, azul software
// Date: 20/10/2023
// copyright (c) 2023 prr, azul software
//

package main

import (
	"fmt"
	"os"
	gio "github.com/pawelgaczynski/giouring"
	)


func errmsg(msg string, err error) {

	if err == nil {
		fmt.Printf("%s\n", msg)
	} else {
		fmt.Printf("%s: %v\n", msg, err)
	}
	os.Exit(-1)
}

func PrintKernel(kern *gio.KernelVersion) {

	fmt.Println("*** kernel version ***")
	fmt.Printf("Kernel: %d Maj: %d Min: %d Flavor: %s\n", (*kern).Kernel, (*kern).Major, (*kern).Minor, (*kern).Flavor)
}

func main() {

	numarg := len(os.Args)

	if numarg < 2 {
		errmsg("insufficient args -- try help",nil)
	}

	cmdStr := os.Args[1]

//	fmt.Printf("command: %s\n", cmdStr)

	switch cmdStr {
	case "ver","version":
		fmt.Printf("giouring library version: %d-%d\n", gio.MajorVersion(), gio.MinorVersion())

	case "kernel":
//		fmt.Println("kernel vers")
		kern, err := gio.GetKernelVersion()
		if err != nil {errmsg("gio.GetKernelVersion", err)}
		PrintKernel(kern)

	case "help":
		out := fmt.Sprintf("./kernelVer [ver/version] | [kernel]\n")
		errmsg(out, nil)

	default:
		out := fmt.Sprintf("invalid command: %s", cmdStr)
		errmsg(out, nil)
	}

	fmt.Println("success!")
}
