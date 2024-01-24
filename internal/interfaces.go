package internal

import (
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

func FindAllInterfaces() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Println(device.Name)
	}
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
	fmt.Println("SrcMAC            DestMAC           SrcIP          DestIP        SrcPort  DestPort")
	for packet := range packetSource.Packets() {
		readPacket(packet)
	}
}

func readPacket(packet gopacket.Packet) {
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)

	if tcpLayer != nil {	
		ethHandler, _ := ethLayer.(*layers.Ethernet)
		sourceMAC := ethHandler.SrcMAC
		destinationMAC := ethHandler.DstMAC

		tcpHandler, _ := tcpLayer.(*layers.TCP)
		sourcePort := tcpHandler.SrcPort
		destinationPort := tcpHandler.DstPort
		
		httpHandler, _ := ipLayer.(*layers.IPv4)
		sourceIP := httpHandler.SrcIP
		destinationIP := httpHandler.DstIP
	
		sourceAddress := &Address{MAC: sourceMAC.String(), IP: sourceIP.String(), PORT: sourcePort.String()}
		destinationAddress := &Address{MAC: destinationMAC.String(), IP: destinationIP.String(), PORT: destinationPort.String()}
		fmt.Printf("%v %v %v %v %v %v\n", sourceAddress.MAC, destinationAddress.MAC, sourceAddress.IP, destinationAddress.IP, sourceAddress.PORT, destinationAddress.PORT)
	}
}
