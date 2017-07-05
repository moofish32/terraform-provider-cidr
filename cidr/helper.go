package cidr

import (
	"fmt"
	"net"

	gcidr "github.com/apparentlymart/go-cidr/cidr"
)

// Consider moving all of these functions to go-cidr/cidr

//VerifyNetwork takes a list subnets and supernet (CIDRBlock) and verifies
//none of the subnets overlap and all subnets are in the supernet
//it returns an error if any of those conditions are not satisfied
func VerifyNetwork(subnets []*net.IPNet, CIDRBlock *net.IPNet) error {
	firstLastIP := make([][]net.IP, len(subnets))
	for i, s := range subnets {
		first, last := gcidr.AddressRange(s)
		firstLastIP[i] = []net.IP{first, last}
	}
	for i, s := range subnets {
		if !CIDRBlock.Contains(firstLastIP[i][0]) || !CIDRBlock.Contains(firstLastIP[i][1]) {
			return fmt.Errorf("%s does not fully contain %s", CIDRBlock.String(), s.String())
		}
		for j := i + 1; j < len(subnets); j++ {
			first := firstLastIP[j][0]
			last := firstLastIP[j][1]
			if s.Contains(first) || s.Contains(last) {
				return fmt.Errorf("%s overlaps with %s", subnets[j].String(), s.String())
			}
		}
	}
	return nil
}

// PreviousSubnet returns the subnet of the desired mask in the IP space
// just lower than the start of IPNet provided. If the IP space rolls over
// then the second return value is true
func PreviousSubnet(startNet *net.IPNet, mask int) (*net.IPNet, bool) {
	startIP := checkIPv4(startNet.IP)
	previousIP := make(net.IP, len(startIP))
	copy(previousIP, startIP)
	cMask := net.CIDRMask(mask, 8*len(previousIP))
	previousIP = Dec(previousIP)
	previous := &net.IPNet{IP: previousIP.Mask(cMask), Mask: cMask}
	if startIP.Equal(net.IPv4zero) || startIP.Equal(net.IPv6zero) {
		return previous, true
	}
	return previous, false
}

// NextSubnet returns the next available subnet of the desired mask size
// starting for the maximum IP of the offset subnet
// If the IP exceeds the maxium IP then the second return value is true
func NextSubnet(offset *net.IPNet, desiredMask int) (*net.IPNet, bool) {
	_, currentLast := gcidr.AddressRange(offset)
	mask := net.CIDRMask(desiredMask, 8*len(currentLast))
	currentSubnet := &net.IPNet{IP: currentLast.Mask(mask), Mask: mask}
	_, last := gcidr.AddressRange(currentSubnet)
	last = Inc(last)
	next := &net.IPNet{IP: last.Mask(mask), Mask: mask}
	if last.Equal(net.IPv4zero) || last.Equal(net.IPv6zero) {
		return next, true
	}
	return next, false
}

//Inc increases the IP by one
func Inc(ip net.IP) net.IP {
	ip = checkIPv4(ip)
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
	return ip
}

//Dec decreases the IP by one
func Dec(ip net.IP) net.IP {
	ip = checkIPv4(ip)
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]--
		if ip[j] < 255 {
			break
		}
	}
	return ip
}

func checkIPv4(ip net.IP) net.IP {
	// Go for some reason allocs IPv6len for IPv4 so we have to correct it
	if v4 := ip.To4(); v4 != nil {
		return v4
	}
	return ip
}
