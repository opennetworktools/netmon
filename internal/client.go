package internal

import (
	"github.com/google/gopacket/pcap"
)

type Client struct {
	handler *pcap.Handle
}

func InitClient(name string) (*Client, error) {
	c := new(Client)
	handler, err := pcap.OpenLive(name, 65536, false, 1000)
	if err != nil {
		return nil, err
	}
	c.handler = handler
	return c, nil
}