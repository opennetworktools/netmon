package main

import (
	"os"
	"fmt"

	"github.com/roopeshsn/netmon/internal"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("No command entered!")
	} else if args[1] == "show" || args[1] == "sh" {
		if args[2] == "asn" {
			internal.GetASN("81.2.69.142")
		} else if args[2] == "country" || args[2] == "cc"  {
			internal.GetCountry("81.2.69.142")
		} else if args[2] == "interface" || args[2] == "int" {
			if len(args) >= 5 {
				if args[3] == "describe" || args[3] == "des" {
					internal.FindInterfaceDescribe(args[4])
				}
			}
		} else if args[2] == "interfaces" || args[2] == "ints" {
			if len(args) >= 4 {
				if args[3] == "describe" || args[3] == "des" {
					internal.FindAllInterfacesDescribe()
				}
				return
			}
			internal.FindAllInterfaces()
		} else if args[2] == "ip" {
			internal.GetLocalIP()
		}
	} else if args[1] == "watch" {
		internal.WatchInterface(args[2])
	} else {
		fmt.Printf("Command %v not found!\n", args[2])
	}
}