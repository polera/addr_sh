package addr

import (
	"net"
)

type Addr struct {
	AboutRoute string            `json:"about"`
	IP         *net.IP           `json:"ip"`
	Tools      map[string]string `json:"tools"`
}

type About struct {
	Text   string `json:"about"`
	Email  string `json:"email"`
	GitHub string `json:"github"`
}

type IPv4CIDR struct {
	NumAddresses       float64 `json:"num_addresses"`
	FirstUsableAddress string  `json:"first_usable_address"`
	LastUsableAddress  string  `json:"last_usable_address"`
	Network            string  `json:"network"`
	Netmask            net.IP  `json:"netmask"`
	BroadcastAddress   string  `json:"broadcast_address"`
	HostOnly           bool    `json:"host_only"`
}

type CIDR struct {
	NumAddresses       float64 `json:"num_addresses"`
	NumUsableAddresses float64 `json:"num_usable_addresses"`
}

type CIDRSplit struct {
	Subnets []IPv4CIDR `json:"subnets"`
	Count   int        `json:"count"`
}
