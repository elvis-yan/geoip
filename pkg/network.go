package geoip

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Network struct {
	Addr      string
	Number    uint32
	Mask      uint32
	IP0       uint32
	IP1       uint32
	Broadcast uint32
	City      string
	Country   string
}

func NewNetwork(addr string) (*Network, error) {
	nw := &Network{Addr: addr}
	ip, mask, err := nw.analyze()

	if err != nil {
		return nil, err
	}
	nw.Mask = mask
	nw.Number = ip & mask
	nw.Broadcast = ip | ^mask
	if nw.Broadcast-nw.Number > 2 {
		nw.IP0 = nw.Number + 1
		nw.IP1 = nw.Broadcast - 1
	}

	return nw, nil
}

func (nw *Network) analyze() (ip, mask uint32, err error) {
	ss := strings.Split(nw.Addr, "/")
	if len(ss) != 2 {
		return 0, 0, fmt.Errorf("Bad Addr: %s", nw.Addr)
	}
	v, err := strconv.Atoi(ss[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Bad Addr: %s", nw.Addr)
	}
	if v < 0 || v > 32 {
		return 0, 0, fmt.Errorf("Bad Addr: %s", nw.Addr)
	}

	t := 1<<32 - 1<<(32-uint8(v))
	mask = uint32(t)

	if err := ValidateIP(ss[0]); err != nil {
		return 0, 0, fmt.Errorf("Bad Addr: %s", nw.Addr)
	}
	ip = str2int(ss[0])

	return ip, mask, nil
}

func ValidateIP(s string) error {
	ss := strings.Split(s, ".")
	if len(ss) != 4 {
		return fmt.Errorf("Bad IP: %s", s)
	}
	for _, s := range ss {
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Bad IP: %s", s)
		}
		if v < 0 || v > 255 {
			return fmt.Errorf("Bad IP: %s", s)
		}
	}
	return nil
}

func IsIntranet(ip string) bool {
	in := func(nw *Network) bool {
		n := str2int(ip)
		return n >= nw.Number && n <= nw.Broadcast
	}
	nw_a, _ := NewNetwork("10.0.0.0/8")
	nw_b, _ := NewNetwork("172.16.0.0/12")
	nw_c, _ := NewNetwork("192.168.0.0/16")

	return in(nw_a) || in(nw_b) || in(nw_c)
}

func int2str(n uint32) (ip string) {
	nums := make([]string, 4)
	i := 3
	var t uint32
	for n != 0 {
		n, t = n/256, n%256
		nums[i] = strconv.Itoa(int(t))
		i--
	}
	return strings.Join(nums, ".")
}

// input `ip` should be validated before.
func str2int(ip string) (n uint32) {
	ss := strings.Split(ip, ".")
	nums := []uint32{1 << 24, 1 << 16, 1 << 8, 1}
	for i := range ss {
		v, err := strconv.Atoi(ss[i])
		if err != nil {
			panic("IP should be validated before")
		}
		n += uint32(v) * nums[i]
	}
	return n
}

func randips(n int) []string {
	rand.Seed(time.Now().UnixNano())
	ips := make([]string, n)
	ip := func() string {
		ns := make([]int, 4)
		for i := range ns {
			ns[i] = rand.Intn(256)
		loop:
			if i == 0 && ns[i] == 0 {
				ns[i] = rand.Intn(256)
				goto loop
			}
		}
		return fmt.Sprintf("%d.%d.%d.%d", ns[0], ns[1], ns[2], ns[3])
	}
	for i := range ips {
		ips[i] = ip()
	}
	return ips
}
