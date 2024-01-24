package main

import (
	"os"
	"github.com/roopeshsn/netmon/internal"
)

func main() {
	args := os.Args
	if args[1] == "show" || args[1] == "sh" {
		if args[2] == "interfaces" || args[2] == "int" {
			if len(args) >= 4 {
				if args[3] == "describe" || args[3] == "des" {
					internal.FindAllInterfacesDescribe()
				}
				return
			}
			internal.FindAllInterfaces()
		}
	}
	if args[1] == "watch" {
		internal.WatchInterface(args[2])
	}
}