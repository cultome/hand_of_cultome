package main

import (
	"log"
	"os"
	"time"

	h "github.com/cultome/handofcultome"
)

func main() {
	ip := os.Args[1] // requires IP address for the Android TV

	pairingPort := "6467"
	remotePort := "6466"

	pairingClient := h.CreateClient(ip, pairingPort)
	defer pairingClient.Conn.Close()

	paring := h.CreatePairingManager(pairingClient)

	if paring.Pair() {
		log.Printf("[pairing] Pairing success!")
		pairingClient.Conn.Close()
	} else {
		log.Printf("[pairing] Pairing failed!")
		os.Exit(1)
	}

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
		msg := remote.ProcessRemoteResponse(b)

		if msg.RemotePingRequest != nil {
			remote.RespondPing(msg.RemotePingRequest.Val1)
		} else {
			log.Printf("[main] %+v\n", msg)
		}
	})

	// time.Sleep(10 * time.Second)
	// remote.VolumeUp()

	for {
		time.Sleep(1 * time.Second)
	}
}
