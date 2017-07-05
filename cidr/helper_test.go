package cidr

import (
	"bytes"
	"net"
	"strconv"
	"testing"
)

type testVerifyNetwork struct {
	CIDRBlock string
	CIDRList  []string
}

func TestIncDec(t *testing.T) {

	testCase := [][]string{
		[]string{"0.0.0.0", "0.0.0.1"},
		[]string{"10.0.0.0", "10.0.0.1"},
		[]string{"9.255.255.255", "10.0.0.0"},
		[]string{"255.255.255.255", "0.0.0.0"},
		[]string{"::", "::1"},
		[]string{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "::"},
		[]string{"2001:db8:c001:ba00::", "2001:db8:c001:ba00::1"},
	}

	for _, tc := range testCase {
		ip1 := net.ParseIP(tc[0])
		ip2 := net.ParseIP(tc[1])
		Inc(ip1)
		if !ip1.Equal(ip2) {
			t.Logf("%s should equal %s\n", tc[0], tc[1])
			t.Errorf("%v should equal %v\n", ip1, ip2)
		}
	}
	for _, tc := range testCase {
		ip1 := net.ParseIP(tc[0])
		ip2 := net.ParseIP(tc[1])
		Dec(ip2)
		if !ip1.Equal(ip2) {
			t.Logf("%s should equal %s\n", tc[0], tc[1])
			t.Errorf("%v should equal %v\n", ip1, ip2)
		}
	}
}

func TestPreviousSubnet(t *testing.T) {

	testCases := [][]string{
		[]string{"10.0.0.0/24", "9.255.255.0/24", "false"},
		[]string{"100.0.0.0/26", "99.255.255.192/26", "false"},
		[]string{"0.0.0.0/26", "255.255.255.192/26", "true"},
		[]string{"2001:db8:e000::/36", "2001:db8:d000::/36", "false"},
		[]string{"::/64", "ffff:ffff:ffff:ffff::/64", "true"},
	}
	for _, tc := range testCases {
		_, c1, _ := net.ParseCIDR(tc[0])
		_, c2, _ := net.ParseCIDR(tc[1])
		mask, _ := c1.Mask.Size()
		p1, rollback := PreviousSubnet(c1, mask)
		if !p1.IP.Equal(c2.IP) {
			t.Errorf("IP expected %v, got %v\n", c2.IP, p1.IP)
		}
		if !bytes.Equal(p1.Mask, c2.Mask) {
			t.Errorf("Mask expected %v, got %v\n", c2.Mask, p1.Mask)
		}
		if p1.String() != c2.String() {
			t.Errorf("%s should have been equal %s\n", p1.String(), c2.String())
		}
		if check, _ := strconv.ParseBool(tc[2]); rollback != check {
			t.Errorf("%s to %s  should have rolled\n", tc[0], tc[1])
		}
	}
	for _, tc := range testCases {
		_, c1, _ := net.ParseCIDR(tc[0])
		_, c2, _ := net.ParseCIDR(tc[1])
		mask, _ := c1.Mask.Size()
		n1, rollover := NextSubnet(c2, mask)
		if !n1.IP.Equal(c1.IP) {
			t.Errorf("IP expected %v, got %v\n", c1.IP, n1.IP)
		}
		if !bytes.Equal(n1.Mask, c1.Mask) {
			t.Errorf("Mask expected %v, got %v\n", c1.Mask, n1.Mask)
		}
		if n1.String() != c1.String() {
			t.Errorf("%s should have been equal %s\n", n1.String(), c1.String())
		}
		if check, _ := strconv.ParseBool(tc[2]); rollover != check {
			t.Errorf("%s to %s  should have rolled\n", tc[0], tc[1])
		}
	}
}

func TestVerifyNetowrk(t *testing.T) {

	testCases := []*testVerifyNetwork{
		&testVerifyNetwork{
			CIDRBlock: "192.168.8.0/21",
			CIDRList: []string{
				"192.168.8.0/24",
				"192.168.9.0/24",
				"192.168.10.0/24",
				"192.168.11.0/25",
				"192.168.11.128/25",
				"192.168.12.0/25",
				"192.168.12.128/26",
				"192.168.12.192/26",
				"192.168.13.0/26",
				"192.168.13.64/27",
				"192.168.13.96/27",
				"192.168.13.128/27",
			},
		},
	}
	failCases := []*testVerifyNetwork{
		&testVerifyNetwork{
			CIDRBlock: "192.168.8.0/21",
			CIDRList: []string{
				"192.168.8.0/24",
				"192.168.9.0/24",
				"192.168.10.0/24",
				"192.168.11.0/25",
				"192.168.11.128/25",
				"192.168.12.0/25",
				"192.168.12.64/26",
				"192.168.12.128/26",
			},
		},
		&testVerifyNetwork{
			CIDRBlock: "192.168.8.0/21",
			CIDRList: []string{
				"192.168.7.0/24",
				"192.168.9.0/24",
				"192.168.10.0/24",
				"192.168.11.0/25",
				"192.168.11.128/25",
				"192.168.12.0/25",
				"192.168.12.64/26",
				"192.168.12.128/26",
			},
		},
	}

	for _, tc := range testCases {
		subnets := make([]*net.IPNet, len(tc.CIDRList))
		for i, s := range tc.CIDRList {
			_, n, err := net.ParseCIDR(s)
			if err != nil {
				t.Errorf("Bad test data %s\n", s)
			}
			subnets[i] = n
		}
		_, CIDRBlock, perr := net.ParseCIDR(tc.CIDRBlock)
		if perr != nil {
			t.Errorf("Bad test data %s\n", tc.CIDRBlock)
		}
		test := VerifyNetwork(subnets, CIDRBlock)
		if test != nil {
			t.Errorf("Failed test with %v\n", test)
		}
	}
	for _, tc := range failCases {
		subnets := make([]*net.IPNet, len(tc.CIDRList))
		for i, s := range tc.CIDRList {
			_, n, err := net.ParseCIDR(s)
			if err != nil {
				t.Errorf("Bad test data %s\n", s)
			}
			subnets[i] = n
		}
		_, CIDRBlock, perr := net.ParseCIDR(tc.CIDRBlock)
		if perr != nil {
			t.Errorf("Bad test data %s\n", tc.CIDRBlock)
		}
		test := VerifyNetwork(subnets, CIDRBlock)
		if test == nil {
			t.Errorf("Test should have failed with CIDR %s\n", tc.CIDRBlock)
		}
	}
}
