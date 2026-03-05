package addr

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

func reverseIP(address string) (string, error) {
	splitAddress := strings.Split(address, ".")
	for i, j := 0, len(splitAddress)-1; i < len(splitAddress)/2; i, j = i+1, j-1 {
		splitAddress[i], splitAddress[j] = splitAddress[j], splitAddress[i]
	}
	return strings.Join(splitAddress, "."), nil
}

func GetPTR(address string) (string, error) {
	dnsClient := new(dns.Client)
	message := new(dns.Msg)
	reverseAddress, _ := reverseIP(address)
	reversed := fmt.Sprintf("%s.in-addr.arpa.", reverseAddress)
	message.SetQuestion(reversed, dns.TypePTR)
	message.RecursionDesired = true
	return doDNSLookup(dnsClient, message)
}

func doDNSLookup(client *dns.Client, message *dns.Msg) (string, error) {
	var err error
	response, _, cliErr := client.Exchange(message, "8.8.8.8:53")
	if cliErr != nil {
		log.Printf("Error on DNS lookup: %s", cliErr)
		return "", cliErr
	}
	if response.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf("DNS call not successful.  Response code: %d", response.Rcode)
		log.Print(err.Error())
	} else {
		for _, answer := range response.Answer {
			if ptrRecord, ok := answer.(*dns.PTR); ok {
				return ptrRecord.Ptr, nil
			}
		}
	}
	return "", err
}

func GetIpv4CIDR(addr string, cidr string) (IPv4CIDR, error) {
	search := fmt.Sprintf("%s/%s", addr, cidr)
	ip, network, err := net.ParseCIDR(search)

	v4cidr := IPv4CIDR{}

	if err != nil {
		return v4cidr, err
	}

	v4network, _ := CalculateV4CIDR(cidr)

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

	if host == "127.0.0.1" {
		host = r.Header.Get("X-Forwarded-For")
	}

	remote := net.ParseIP(host)
	if remote != nil {
		remoteHost = &remote
	}
	return remoteHost
}

func CalculateV4CIDR(cidr string) (CIDR, error) {
	cidrAsInt, err := strconv.Atoi(cidr)

	v4cidr := CIDR{}
	if err != nil || cidrAsInt < 0 || cidrAsInt > 32 {
		return v4cidr, err
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
