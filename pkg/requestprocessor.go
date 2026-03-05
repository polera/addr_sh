package addr

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func GetIpv4CIDR(addr string, cidr string) (IPv4CIDR, error) {
	search := fmt.Sprintf("%s/%s", addr, cidr)
	ip, network, err := net.ParseCIDR(search)

	v4cidr := IPv4CIDR{}

	if err != nil {
		return v4cidr, err
	}

	v4network, err := CalculateV4CIDR(cidr)
	if err != nil {
		return v4cidr, err
	}

	// fantastic writeup here: https://networkengineering.stackexchange.com/a/53994
	// build host mask
	// host mask is the bitwise NOT of the network mask
	res := make(net.IP, 4)
	for i, b := range network.Mask {
		res[i] = ^b
	}

	// calculate broadcast addr
	// broadcast address is the sum of network address and the host mask
	broadcast := make(net.IP, 4)
	for i, o := range network.IP.To4() {
		broadcast[i] = o + res[i]
	}

	firstUsable := make(net.IP, 4)
	copy(firstUsable, network.IP)
	firstUsable[3] = network.IP[3] + 1

	// last usable = broadcast - 1
	lastUsable := make(net.IP, 4)
	copy(lastUsable, broadcast)
	lastUsable[3] = lastUsable[3] - 1

	// if this is a host only network, there's no broadcast
	netOnes, netBits := network.Mask.Size()
	if netOnes == netBits {
		lastUsable = broadcast
		broadcast = make(net.IP, 4)
	}

	v4cidr.FirstUsableAddress = firstUsable.String()
	v4cidr.LastUsableAddress = lastUsable.String()
	v4cidr.BroadcastAddress = broadcast.String()

	if v4network.NumAddresses == 1 {
		v4cidr.HostOnly = true
		v4cidr.FirstUsableAddress = ip.String()
		v4cidr.LastUsableAddress = ip.String()
		v4cidr.BroadcastAddress = broadcast.String()
	}

	v4cidr.Network = network.String()
	v4cidr.Netmask = net.IP(network.Mask)
	v4cidr.NumAddresses = v4network.NumAddresses

	return v4cidr, nil

}

func GetRemoteHost(r *http.Request) *net.IP {
	remoteHost := &net.IP{}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return remoteHost
	}

	if host == "127.0.0.1" || host == "::1" {
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			// Take the first (leftmost) IP from comma-separated list
			if idx := strings.IndexByte(xff, ','); idx != -1 {
				xff = xff[:idx]
			}
			host = strings.TrimSpace(xff)
		}
	}

	remote := net.ParseIP(host)
	if remote != nil {
		remoteHost = &remote
	}
	return remoteHost
}

func SplitCIDR(cidr string, count int) (CIDRSplit, error) {
	result := CIDRSplit{}

	if count < 2 {
		return result, fmt.Errorf("count must be at least 2")
	}
	if count&(count-1) != 0 {
		return result, fmt.Errorf("count must be a power of 2, got %d", count)
	}

	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return result, fmt.Errorf("invalid CIDR: %w", err)
	}

	ones, bits := network.Mask.Size()
	if bits != 32 {
		return result, fmt.Errorf("only IPv4 CIDRs are supported")
	}

	bitsNeeded := int(math.Log2(float64(count)))
	newPrefix := ones + bitsNeeded
	if newPrefix > 32 {
		return result, fmt.Errorf("cannot split /%d into %d subnets: would require /%d which exceeds /32", ones, count, newPrefix)
	}

	newMask := net.CIDRMask(newPrefix, 32)
	ip := network.IP.To4()
	ipInt := uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
	subnetSize := uint32(1) << uint(32-newPrefix)

	subnets := make([]string, count)
	for i := 0; i < count; i++ {
		subnetIP := make(net.IP, 4)
		addr := ipInt + uint32(i)*subnetSize
		subnetIP[0] = byte(addr >> 24)
		subnetIP[1] = byte(addr >> 16)
		subnetIP[2] = byte(addr >> 8)
		subnetIP[3] = byte(addr)
		subnet := net.IPNet{IP: subnetIP, Mask: newMask}
		subnets[i] = subnet.String()
	}

	result.Subnets = subnets
	result.Count = count
	return result, nil
}

func CalculateV4CIDR(cidr string) (CIDR, error) {
	cidrAsInt, err := strconv.Atoi(cidr)

	v4cidr := CIDR{}
	if err != nil {
		return v4cidr, err
	}
	if cidrAsInt < 0 || cidrAsInt > 32 {
		return v4cidr, fmt.Errorf("CIDR must be between 0 and 32, got %d", cidrAsInt)
	}

	numAddresses := math.Pow(float64(2), (float64(32) - float64(cidrAsInt)))
	v4cidr.NumAddresses = numAddresses
	if numAddresses > 1 {
		v4cidr.NumUsableAddresses = numAddresses - 2
	} else {
		v4cidr.NumUsableAddresses = 1
	}

	return v4cidr, nil

}
