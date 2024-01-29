package internal

import (
	"net"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
	// "github.com/jedib0t/go-pretty/v6/table"
)

type Address struct {
	MAC string
	IP   string
	PORT uint16
}

type CPacket struct {
	SrcAddress Address
	DstAddress Address
	Protocol string
	Timestamp time.Time
}

type CHost struct {
	HostName string
	HostNames []string
	ASNumber uint
	ASName string
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

func WatchInterface(name string, c chan CPacket) {
	client, err := InitClient(name)
	if err != nil {
		fmt.Printf("%s", err)
	}
	packetSource := gopacket.NewPacketSource(client.handler, client.handler.LinkType())
	parsePackets(packetSource, c)
}

func parsePackets(packetSource *gopacket.PacketSource, c chan CPacket) {
	for packet := range packetSource.Packets() {
		readPacket(packet, c)
	}
}

func readPacket(packet gopacket.Packet,  c chan CPacket) {
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)

	ethHandler, err := ethLayer.(*layers.Ethernet)
	if err != true {
		return
	}
	sourceMAC := ethHandler.SrcMAC
	destinationMAC := ethHandler.DstMAC

	httpHandler, err := ipLayer.(*layers.IPv4)
	if err != true {
		return
	}
	sourceIP := httpHandler.SrcIP
	destinationIP := httpHandler.DstIP
	protocol := httpHandler.Protocol

	var ports TransportLayerPortInfo
	if tcpLayer != nil {
		// Type assertion for TCP layer
		tcpHandler, ok := tcpLayer.(*layers.TCP)
		if !ok {
			return
		}
		// Call parseTCPHeader to get the ports
		ports = parseTCPHeader(tcpHandler)
	} else if udpLayer != nil {
		// Type assertion for UDP layer
		udpHandler, ok := udpLayer.(*layers.UDP)
		if !ok {
			return
		}
		// Call parseUDPHeader to get the ports
		ports = parseUDPHeader(udpHandler)
	}

	// Extract source and destination ports
	srcPort := ports.SrcPort
	dstPort := ports.DstPort

	// Debugging
	// if tcpLayer != nil {
	// 	fmt.Println("TCP Layer")
	// 	srcPort, dstPort = parseUDPHeader(tcpLayer)
	// } else if udpLayer != nil {
	// 	fmt.Println("UDP Layer")
	// 	srcPort, dstPort = parseUDPHeader(udpLayer)
	// }

	srcAddress := Address{MAC: sourceMAC.String(), IP: sourceIP.String(), PORT: srcPort}
	dstAddress := Address{MAC: destinationMAC.String(), IP: destinationIP.String(), PORT: dstPort}
	cPacket := CPacket{SrcAddress: srcAddress, DstAddress: dstAddress, Protocol: protocol.String(), Timestamp: packet.Metadata().CaptureInfo.Timestamp}
		
	c <- cPacket
}

type TransportLayerPortInfo struct {
	SrcPort uint16
	DstPort uint16
}

func parseTCPHeader(tcpHandler *layers.TCP) TransportLayerPortInfo {
	var srcPort uint16 = uint16(tcpHandler.SrcPort)
	var dstPort uint16 = uint16(tcpHandler.DstPort)
	return TransportLayerPortInfo{SrcPort: srcPort, DstPort: dstPort}
}

func parseUDPHeader(udpHandler *layers.UDP) TransportLayerPortInfo {
	var srcPort uint16 = uint16(udpHandler.SrcPort)
	var dstPort uint16 = uint16(udpHandler.DstPort)
	return TransportLayerPortInfo{SrcPort: srcPort, DstPort: dstPort}
}

func PrintPacket(c chan CPacket) {
	fmt.Println("Timestamp                             SrcMAC            DestMAC           SrcIP             DestIP        SrcPort  DestPort   Protocol")
	// t := table.NewWriter()
	// t.AppendHeader(table.Row{"Timestamp", "SrcMac", "DstMac", "SrcIP", "DstIP", "SrcPort", "DstPort", "Protocol"})
	for {
		p := <- c
		fmt.Printf("%v %v %v %v -> %v %d %d %s\n", p.Timestamp, p.SrcAddress.MAC, p.DstAddress.MAC, p.SrcAddress.IP, p.DstAddress.IP, p.SrcAddress.PORT, p.DstAddress.PORT, p.Protocol)
	}
}

func getInterfaceAddresses(name string) []string{
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	var addresses []string
	for _, device := range devices {
		if name == device.Name {
			for _, address := range device.Addresses {
				addresses = append(addresses, address.IP.String())
			}
		}
	}

	return addresses
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func isLoopback(ipAddress string) bool {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}

func getTrafficDirection(srcIP string, dstIP string, srcPort uint16, dstPort uint16, interfaceAddresses[] string) string {
	if isLoopback(srcIP) && isLoopback(dstIP) {
		if srcPort > dstPort {
			return "Outgoing"
		} else {
			return "Incoming"
		}
	}

	if contains(interfaceAddresses, srcIP) {
		return "Outgoing"
	} else if srcIP != "0.0.0.0" {
		return "Incoming"
	} else if !contains(interfaceAddresses, dstIP) {
		return "Outgoing"
	} 

	return "Incoming"
}

func ResolveHostsInformation(interfaceName string, c chan CPacket) {
	m := make(map[string]CHost)

	for {
		p := <- c
		// determine which ip address to lookup (src, or dst) based on incoming and outgoing traffic
		var interfaceAddresses []string = getInterfaceAddresses(interfaceName)
		var trafficDirection string = getTrafficDirection(p.SrcAddress.IP, p.DstAddress.IP, p.SrcAddress.PORT, p.DstAddress.PORT, interfaceAddresses)

		var ipAddressToLookup string
		if trafficDirection == "Incoming" {
			ipAddressToLookup = p.SrcAddress.IP
		} else {
			ipAddressToLookup = p.DstAddress.IP
		}

		_, ok := m[ipAddressToLookup] 
		if ok {
			continue
		}

		asNumber, asName := rHost(ipAddressToLookup)
		hostNames, err := rDNS(ipAddressToLookup)

		if err != nil {
			fmt.Printf("%v, %v\n", ipAddressToLookup, err)
		}

		var hostName string
		if len(hostNames) >= 1 {
			hostName = hostNames[0]
		} else {
			hostName = ipAddressToLookup
		}
		host := CHost{HostName: hostName, HostNames: hostNames, ASNumber: asNumber, ASName: asName}
		m[ipAddressToLookup] = host
		fmt.Printf("%v %v - %v %v \n", ipAddressToLookup, hostName, asNumber, asName)
	}
}

func rHost(ipAddress string) (uint, string) {
	ASNumber, ASName := GetASN(ipAddress)
	return ASNumber, ASName
}

func rDNS(ipAddress string) ([]string, error) {
	domainNames, err := net.LookupAddr(ipAddress)
	if err != nil {
		return nil, err
	}
	return domainNames, nil
}
