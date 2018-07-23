package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	geoip "github.com/elvis-yan/geoip/pkg"
)

var (
	addr = flag.String("addr", "", "geoip http address")
	dir  = flag.String("dir", "", "geoip data directory")
)

func main() {
	flag.Parse()

	if *dir == "" {
		log.Println("specify the data directory")
		flag.Usage()
	}

	if !strings.HasSuffix(*dir, "/") {
		*dir += "/"
	}

	if *addr == "" {
		shell(*dir)
	} else {
		log.Printf("geoip server listen on %s", *addr)
		server(*addr, *dir)
	}
}

func shell(dir string) {
	fmt.Println("Loading...")
	web := geoip.NewWeb(dir)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("IP: ")
	for scanner.Scan() {
		ip := scanner.Text()
		if err := geoip.ValidateIP(ip); err != nil {
			fmt.Printf("Bad IP: %s\n\n", ip)
			fmt.Print("IP: ")
			continue
		}

		if geoip.IsIntranet(ip) {
			fmt.Println("Intranet IP")
			fmt.Print("IP: ")
			continue
		}

		nw := web.City(ip)
		fmt.Printf("(%s, %s)\n", nw.City, nw.Country)
		fmt.Print("IP: ")
	}
}

func server(addr, dir string) {
	log.Println("Loading...")
	web := geoip.NewWeb(dir)

	handler := func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Path[7:]
		log.Println(ip)
		if err := geoip.ValidateIP(ip); err != nil {
			w.Write([]byte("Bad IP"))
			return
		}

		if geoip.IsIntranet(ip) {
			w.Write([]byte("Intranet IP"))
			return
		}

		nw := web.City(ip)
		fmt.Fprintf(w, "(%s, %s)", nw.City, nw.Country)
	}

	http.HandleFunc("/geoip/", handler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
