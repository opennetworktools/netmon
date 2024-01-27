package internal

import (
	"net"
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
)

type Address struct {
	MAC string
	IP   string
	PORT string
}

func GetLocalIP() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}

	for _, iface := range interfaces {
		// Exclude loopback and non-up interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				fmt.Println(err)
				return
			}

			// Check if it's an IPv4 address
			if ip.To4() != nil {
				fmt.Println(ip.String())
				return
			}
		}
	}

	fmt.Println("local IP address not found")
}

func FindAllInterfaces() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Println(device.Name)
	}
}

func FindInterfaceDescribe(name string) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		if name == device.Name {
			fmt.Println(device.Name)
			for _, address := range device.Addresses {
				fmt.Println("- IP address: ", address.IP)
				fmt.Println("- Subnet mask: ", address.Netmask)
			}
			return
		}
	}
	fmt.Printf("Unable to find interface with the name %s\n", name)
}

func FindAllInterfacesDescribe() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Println("\n", device.Name)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
}

func WatchInterface(name string) {
	client, err := InitClient(name)
	if err != nil {
		fmt.Printf("%s", err)
	}
	packetSource := gopacket.NewPacketSource(client.handler, client.handler.LinkType())
	fmt.Println("SrcMAC            DestMAC           SrcIP             DestIP        SrcPort  DestPort")
	for packet := range packetSource.Packets() {
		readPacket(packet)
	}
}

func readPacket(packet gopacket.Packet) {
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)

	if tcpLayer != nil {	
		ethHandler, err := ethLayer.(*layers.Ethernet)
		if err != true {
			// fmt.Println("ethHandler")
			return
		}
		sourceMAC := ethHandler.SrcMAC
		destinationMAC := ethHandler.DstMAC

		tcpHandler, err := tcpLayer.(*layers.TCP)
		if err != true {
			// fmt.Println("tcpHandler")
			return
		}
		sourcePort := tcpHandler.SrcPort
		destinationPort := tcpHandler.DstPort
		
		httpHandler, err := ipLayer.(*layers.IPv4)
		if err != true {
			// fmt.Println("httpHandler")
			return
		}
		sourceIP := httpHandler.SrcIP
		destinationIP := httpHandler.DstIP
	
		sourceAddress := &Address{MAC: sourceMAC.String(), IP: sourceIP.String(), PORT: sourcePort.String()}
		destinationAddress := &Address{MAC: destinationMAC.String(), IP: destinationIP.String(), PORT: destinationPort.String()}
		fmt.Printf("%v %v %v -> %v %v %v\n", sourceAddress.MAC, destinationAddress.MAC, sourceAddress.IP, destinationAddress.IP, sourceAddress.PORT, destinationAddress.PORT)
	}
}

func rDNS() {
	domainNames, err := net.LookupAddr("13.232.193.230")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, domain := range domainNames {
		fmt.Println(domain)
	}
}
