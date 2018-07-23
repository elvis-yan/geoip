package geoip

import (
	"testing"
)

func TestNewNetwork(t *testing.T) {
	table := map[string]struct {
		mask      string
		number    string
		broadcast string
		ip0       string
		ip1       string
	}{
		"168.173.70.134/29": {
			"255.255.255.248",
			"168.173.70.128",
			"168.173.70.135",
			"168.173.70.129",
			"168.173.70.134",
		},
		"153.34.173.242/17": {
			"255.255.128.0",
			"153.34.128.0",
			"153.34.255.255",
			"153.34.128.1",
			"153.34.255.254",
		},
	}
	assertEqual := func(a, b uint32) {
		if a != b {
			t.Error("geoip.NewNetwork(_) test failed")
		}
	}

	for addr := range table {
		nw, err := NewNetwork(addr)
		if err != nil {
			t.Error("geoip.NewNetwork(_) test failed")
		}
		v := table[addr]
		assertEqual(nw.Mask, str2int(v.mask))
		assertEqual(nw.Number, str2int(v.number))
		assertEqual(nw.Broadcast, str2int(v.broadcast))
		assertEqual(nw.IP0, str2int(v.ip0))
		assertEqual(nw.IP1, str2int(v.ip1))
	}
}

func TestInt2str(t *testing.T) {
	if int2str(3232235620) != "192.168.0.100" {
		t.Error("geoip.int2str(_) test failed")
	}
}

func TestStr2int(t *testing.T) {
	if str2int("192.168.0.100") != 3232235620 {
		t.Error("geoip.str2int(_) test failed")
	}
}

func TestStrIntIPTransform(t *testing.T) {
	ips := randips(100000)
	for _, ip := range ips {
		n := str2int(ip)
		s := int2str(n)
		if s != ip {
			t.Error(ip, s, n)
		}
	}
}
