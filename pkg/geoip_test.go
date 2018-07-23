package geoip

import (
	"os"
	"testing"
)

func TestWeb(t *testing.T) {
	os.Remove("../data/network.gob")
	defer os.Remove("../data/network.gob")

	web := NewWeb("../data/")
	var nw1, nw2 *Network
	for i := 0; i < len(web.Networks)-1; i++ {
		nw1 = web.Networks[i]
		nw2 = web.Networks[i+1]
		if nw1.Number < nw1.Broadcast && nw1.Broadcast < nw2.Number {
			continue
		} else {
			if nw1.Number == nw1.Broadcast && nw1.Mask == 1<<32-1 {
				continue
			}
			t.Error("NewWeb test failed")
		}
	}
	if nw2.Number >= nw2.Broadcast {
		t.Error("NewWeb test failed")
	}
}

func TestWebCity(t *testing.T) {
	web := NewWeb("../data/")
	table := map[string]struct{ city, country string }{
		"36.2.12.15":     {"Tokyo", "JP"},
		"123.2.15.98":    {"Coogee", "AU"},
		"16.58.68.216":   {"Palo Alto", "US"},
		"183.240.196.59": {"Beijing", "CN"},
		"18.19.20.21":    {"Cambridge", "US"},
	}
	for ip := range table {
		v := table[ip]
		nw := web.City(ip)
		if nw.City != v.city || nw.Country != v.country {
			t.Errorf("Web.City(%q) test failed", ip)
		}
	}
}
