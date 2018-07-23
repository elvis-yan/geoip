package geoip

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type Web struct {
	Networks Networks
}

type Networks []*Network

func (nws Networks) Len() int           { return len(nws) }
func (nws Networks) Less(i, j int) bool { return nws[i].Number < nws[j].Number }
func (nws Networks) Swap(i, j int)      { nws[i], nws[j] = nws[j], nws[i] }

func NewWeb(dir string) *Web {
	w := &Web{Networks: make(Networks, 0)}
	if err := w.load(dir + "network.gob"); err != nil {
		if err := w.genNetworks(dir + "data.csv"); err != nil {
			log.Fatalf("NewWeb(): %v", err)
		}
		sort.Sort(w.Networks)
		w.dump(dir + "network.gob")
	}
	return w
}

func (w *Web) City(ip string) *Network {
	n := str2int(ip)
	lo, hi := 0, len(w.Networks)
	for lo < hi {
		mid := (lo + hi) / 2
		nt := w.Networks[mid]
		if n >= nt.Number && n <= nt.Broadcast {
			return nt
		} else if n < nt.Number {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return &Network{}
}

func (w *Web) genNetworks(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Web.genNetwork(%q) openfile: %v", path, err)
	}
	r := csv.NewReader(f)
	// pass the csv header (first row)
	if _, err := r.Read(); err != nil {
		return fmt.Errorf("Web.genNetwork(%q) read csv: %v", path, err)
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Web.genNetwork(%q) read csv: %v", path, err)
		}
		if len(record) != 3 {
			return fmt.Errorf("Web.genNetwork(%q) bad record %q: %v", path, record, err)
		}

		nw, err := NewNetwork(record[0])

		if err != nil {
			return fmt.Errorf("Web.genNetwork(%q) record=%q , could not create Network: %v", path, record, err)
		}
		nw.City = record[1]
		nw.Country = record[2]

		w.Networks = append(w.Networks, nw)
	}

	return nil
}

func (w *Web) load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Web.load(%q) openfile: %v", path, err)
	}
	if err := gob.NewDecoder(f).Decode(&w.Networks); err != nil {
		return fmt.Errorf("Web.load(%q) decode: %v", path, err)
	}
	return nil
}

func (w *Web) dump(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("Web.dump(%q) openfile: %v", path, err)
	}
	if err := gob.NewEncoder(f).Encode(w.Networks); err != nil {
		return fmt.Errorf("Web.dump(%q) encode: %v", path, err)
	}
	return nil
}
