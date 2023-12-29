package main

import (
	"log"
	"os"
	"time"

	h "github.com/cultome/handofcultome"
)

func main() {
	ip := "192.168.0.6"
	pairingPort := "6467"
	remotePort := "6466"

	pairingClient := h.CreateClient(ip, pairingPort)
	defer pairingClient.Conn.Close()

	paring := h.CreatePairingManager(pairingClient)

	if paring.Pair() {
		log.Printf("[pairing] Pairing success!")
	} else {
		log.Printf("[pairing] Pairing failed!")
		os.Exit(1)
	}

	go pairingClient.AddListener(func(b []byte) {
		log.Printf("[<pairing] %v\n", b)
	})

	remoteClient := h.CreateClient(ip, remotePort)
	defer remoteClient.Conn.Close()

	remote := h.CreateRemoteManager(remoteClient)

	time.Sleep(1)

	if remote.Configure() {
		log.Printf("[remote] Configuration correct!")
	} else {
		log.Printf("[remote] Configuration failed!")
		os.Exit(2)
	}

	go remoteClient.AddListener(func(b []byte) {
		log.Printf("[<remote] %v\n", b)

		remote.ProcessRemoteResponse(b)

		if b[1] == 66 && b[2] == 8 {
			remote.RespondPing()
		}
	})

	time.Sleep(10 * time.Second)
	remote.VolumeUp()

	for {
		time.Sleep(1 * time.Second)
	}
}
