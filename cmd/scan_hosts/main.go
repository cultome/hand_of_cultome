package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	myIP := localIP()
	ipRange := strings.Join(strings.Split(myIP, ".")[:3], ".")

	fmt.Printf("Scanning...")
	for i := 0; i < 255; i += 1 {
		fmt.Printf(".")
		currentIp := ipRange + strconv.Itoa(i)
		found := scanPort(currentIp)

		if found {
			fmt.Printf("\nServer address found at: %s\n", currentIp)
			break
		}
	}
}

func scanPort(hostname string) bool {
	address := hostname + ":6467"
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}

func localIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Panicf("Problems getting network interface: %+v\n", err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Panicf("Problems getting network addresses %+v\n", err)
		}

		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.IsPrivate() && !ip.IsLoopback() {
				if strings.HasPrefix(ip.String(), "192.168.") {
					fmt.Printf("Local IP in [%s]: %v\n", i.Name, ip)
					return ip.String()
				}
			}
		}
	}

	return ""
}
