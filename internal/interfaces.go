package internal

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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
	CaptureLength int // bytes
}

type CHost struct {
	IP string
	HostName string
	HostNames []string
	ASNumber uint
	ASName string
	Bytes int
}

type HostQueue struct {
	lock sync.Mutex
	queue *list.List
}

func NewHostQueue() *HostQueue {
    return &HostQueue{
        queue: list.New(),
    }
}

func (q *HostQueue) Dequeue() *CHost {
    q.lock.Lock()
    defer q.lock.Unlock()

    if q.queue.Len() == 0 {
        return nil
    }

    frontElement := q.queue.Front()
    host := frontElement.Value.(CHost)
    q.queue.Remove(frontElement)
    
    return &host
}

// I implemented a map with mutex, but in my observation it's not neeeded.
type HostMap struct {
    lock sync.Mutex
    hosts map[string]CHost
}

func NewHostMap() *HostMap {
    return &HostMap{
        hosts: make(map[string]CHost),
    }
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
	if !err {
		return
	}
	sourceMAC := ethHandler.SrcMAC
	destinationMAC := ethHandler.DstMAC


	var sourceIP net.IP
	var destinationIP net.IP
	var protocol layers.IPProtocol

	httpHandler, err := ipLayer.(*layers.IPv4)
	if !err {
		return
	}

	// Capturing IPv6 Packets but commented out temporary due to the below error:

	// 2024/03/30 11:17:42 IP passed to Lookup cannot be nil
	// panic: IP passed to Lookup cannot be nil

	// goroutine 6 [running]:
	// log.Panic({0xc0000cfbc8?, 0x0?, 0x76918f?})
	// 		/usr/local/go/src/log/log.go:432 +0x5a
	// github.com/roopeshsn/netmon/internal.mmdbASNReader({0x76918f, 0x5})
	// 		/home/roopesh/developer/projects/netmon/internal/mmdb.go:36 +0x1eb
	// github.com/roopeshsn/netmon/internal.GetASN(...)
	// 		/home/roopesh/developer/projects/netmon/internal/mmdb.go:13
	// github.com/roopeshsn/netmon/internal.rHost({0x76918f?, 0xc000112510?})
	// 		/home/roopesh/developer/projects/netmon/internal/interfaces.go:429 +0x19
	// github.com/roopeshsn/netmon/internal.ResolveHostInformation({0x7ffe3ae2c7c4, 0x9}, 0xc000146000, 0xc000134030, 0x0?, 0x1, 0xc0000341d0)
	// 		/home/roopesh/developer/projects/netmon/internal/interfaces.go:382 +0x5b7
	// created by github.com/roopeshsn/netmon/internal.ResolveHostsInformation in goroutine 19
	// 		/home/roopesh/developer/projects/netmon/internal/interfaces.go:346 +0x152

	// ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
	// ipv6Handler, _ := ipv6Layer.(*layers.IPv6)

	// if httpHandler != nil {
    //     sourceIP = httpHandler.SrcIP
    //     destinationIP = httpHandler.DstIP
    //     protocol = httpHandler.Protocol
    // } else if ipv6Handler != nil {
    //     sourceIP = ipv6Handler.SrcIP
    //     destinationIP = ipv6Handler.DstIP
    //     protocol = ipv6Handler.NextHeader
    // }

	sourceIP = httpHandler.SrcIP
	destinationIP = httpHandler.DstIP
	protocol = httpHandler.Protocol

	var ports TransportLayerPortInfo
	if tcpLayer != nil {
		tcpHandler, ok := tcpLayer.(*layers.TCP) // Type assertion for TCP layer
		if !ok {
			return
		}

		ports = parseTCPHeader(tcpHandler) // Call parseTCPHeader to get the ports
	} else if udpLayer != nil {
		udpHandler, ok := udpLayer.(*layers.UDP)
		if !ok {
			return
		}
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
	cPacket := CPacket{SrcAddress: srcAddress, DstAddress: dstAddress, Protocol: protocol.String(), Timestamp: packet.Metadata().CaptureInfo.Timestamp, CaptureLength: packet.Metadata().CaptureInfo.CaptureLength}
		
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

func ResolveHostsInformation(interfaceName string, c chan CPacket, m *HostMap, mc chan CHost, html bool) {
	hostQueue := NewHostQueue()

	var wg sync.WaitGroup
	wg.Add(2)

	// ResolveHostInformation will put CHost to hostQueue.queue
	go ResolveHostInformation(interfaceName, c, m, mc, html, hostQueue)

	// pushToMC will push the CHost to a channel (mc), if queue is non-empty
	go pushToMC(hostQueue, mc)

	wg.Wait()
}

func ResolveHostInformation(interfaceName string, c chan CPacket, m *HostMap, mc chan CHost, html bool, hostQueue *HostQueue) {
	for p := range c {
		// determine which ip address to lookup (src, or dst) based on incoming and outgoing traffic
		var interfaceAddresses []string = getInterfaceAddresses(interfaceName)
		var trafficDirection string = getTrafficDirection(p.SrcAddress.IP, p.DstAddress.IP, p.SrcAddress.PORT, p.DstAddress.PORT, interfaceAddresses)
		var captureLength int = p.CaptureLength

		var ipAddressToLookup string
		if trafficDirection == "Incoming" {
			ipAddressToLookup = p.SrcAddress.IP
		} else {
			ipAddressToLookup = p.DstAddress.IP
		}

		m.lock.Lock()

		oldHost, ok := m.hosts[ipAddressToLookup]
		if ok {
			oldHost.Bytes += captureLength
			m.hosts[ipAddressToLookup] = oldHost

			oldHostClone, ok := m.hosts[ipAddressToLookup]
			if ok {
				hostQueue.lock.Lock()
				hostQueue.queue.PushBack(oldHostClone)
				hostQueue.lock.Unlock()
			}
		} else {
			asNumber, asName := rHost(ipAddressToLookup)
			hostNames, err := rDNS(ipAddressToLookup)

			if err != nil {
				// fmt.Printf("%v, %v\n", ipAddressToLookup, err)
			}

			var hostName string
			if len(hostNames) >= 1 {
				hostName = hostNames[0]
			} else {
				hostName = ipAddressToLookup
			}
			host := CHost{IP: ipAddressToLookup, HostName: hostName, HostNames: hostNames, ASNumber: asNumber, ASName: asName, Bytes: captureLength}
			m.hosts[ipAddressToLookup] = host

			hostClone := host
			hostQueue.lock.Lock()
			hostQueue.queue.PushBack(hostClone)
			hostQueue.lock.Unlock()

			if !html {
				fmt.Println(m.hosts[ipAddressToLookup])
				// fmt.Printf("%v %v - %v %v \n", ipAddressToLookup, hostName, asNumber, asName)
			}
		}

		m.lock.Unlock()
	}
}

func pushToMC(hostQueue *HostQueue, mc chan CHost) {
	for {
		host := hostQueue.Dequeue()

		// If the queue is empty, wait for a short duration before trying again
		if host == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Send the dequeued host to the channel
		mc <- *host
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
