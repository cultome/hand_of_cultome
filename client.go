package handofcultome

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
)

type Client struct {
	Certificate       *x509.Certificate
	RemoteCertificate *x509.Certificate
	PublicKey         *rsa.PublicKey
	RemotePublicKey   *rsa.PublicKey
	Conn              *tls.Conn
}

func CreateClient(host, port string) *Client {
	cert, err := tls.LoadX509KeyPair("keys/cert.pem", "keys/private.pem")
	if err != nil {
		log.Panicf("[client] Error: %+v", err)
	}

	publicCertPem, err := os.ReadFile("keys/cert.pem")
	pubBlock, _ := pem.Decode(publicCertPem)
	var publicCert *x509.Certificate
	publicCert, _ = x509.ParseCertificate(pubBlock.Bytes)
	publicKey := publicCert.PublicKey.(*rsa.PublicKey)

	conf := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	log.Printf("[client] Connecting to %s:%s...\n", host, port)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), conf)
	if err != nil {
		log.Panicf("[client] Error: %+v", err)
	}

	log.Printf("[client] Connection successful!")

	serverCert := conn.ConnectionState().PeerCertificates[0]
	serverKey := serverCert.PublicKey.(*rsa.PublicKey)

	return &Client{
		Certificate:       publicCert,
		PublicKey:         publicKey,
		RemoteCertificate: serverCert,
		RemotePublicKey:   serverKey,
		Conn:              conn,
	}
}

func (client *Client) AddListener(onData func([]byte)) {
	log.Printf("[client] Setting listener\n")
	for {
		// log.Printf("[client] Waiting for data...\n")
		data := client.readPayload()
		onData(data)
	}
}

func (c *Client) makeRequest(payload []byte) []byte {
	c.writePayload(payload)
	return c.readPayload()
}

func (c *Client) readPayload() []byte {
	buf := make([]byte, 100)
	n, err := c.Conn.Read(buf)
	if err != nil {
		log.Panicf("[client] Problems reading from network: %+v\n", err)
	}

	bytesToRead := int(buf[0])

	resp := make([]byte, bytesToRead+100)
	bytesReaded := 0

	if n > 1 {
		copy(resp[:n-1], buf[1:n])
		bytesReaded = n - 1
	}

	for bytesReaded < bytesToRead {
		n, err = c.Conn.Read(buf)
		copy(resp[bytesReaded:bytesReaded+n], buf[:n])
		bytesReaded += n
	}

	return resp
}

func (c *Client) writePayload(payload []byte) int {
	reqSize := c.intToBytes(len(payload))

	// log.Printf("[client] Sending size: %v\n", reqSize)
	_, err := c.Conn.Write(reqSize)
	if err != nil {
		log.Panic("[client] Error: Unable to send request size")
	}

	// log.Printf("[client>] %v\n", payload)
	n, err := c.Conn.Write(payload)
	if err != nil {
		log.Panic("[client] Error: Unable to send payload")
	}

	return n
}

func (c *Client) makeSecret(code string) []byte {
	codeBytes := c.hexStringToBytes(code)

	h := sha256.New()

	v1 := c.PublicKey.N.Bytes()
	v2 := c.intToBytes(c.PublicKey.E)
	v3 := c.RemotePublicKey.N.Bytes()
	v4 := c.intToBytes(c.RemotePublicKey.E)
	v5 := codeBytes[1:]

	h.Write(v1)
	h.Write(v2)
	h.Write(v3)
	h.Write(v4)
	h.Write([]byte(v5))

	hash := h.Sum(nil)

	if codeBytes[0] != hash[0] {
		log.Panicf("[client] Error: Incorrect code [%v != %v]\n", codeBytes[0], hash[0])
	}

	return hash
}

func (c *Client) intToBytes(size int) []byte {
	return big.NewInt(int64(size)).Bytes()
}

func (c *Client) hexStringToBytes(value string) []byte {
	v1 := c.hexDigitToByte(value[0:2])
	v2 := c.hexDigitToByte(value[2:4])
	v3 := c.hexDigitToByte(value[4:6])

	return []byte{v1, v2, v3}
}

func (c *Client) hexDigitToByte(value string) byte {
	v, _ := strconv.ParseInt(value, 16, 16)

	if v > 127 {
		v = v - 256
	}

	return byte(v)
}
