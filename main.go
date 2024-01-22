package main

import (
	"os"
	"github.com/roopeshsn/netmon/internal"
)

func main() {
	args := os.Args
	if args[1] == "show" || args[1] == "sh" {
		if args[2] == "interfaces" || args[2] == "int" {
			// if args[3] == "describe" {
			// 	internal.FindAllInterfacesDescribe()
			// }
			internal.FindAllInterfaces()
		}
	}
}