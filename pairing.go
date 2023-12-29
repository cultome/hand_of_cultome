package handofcultome

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"
)

type PairingManager struct {
	client *Client
}

func CreatePairingManager(client *Client) *PairingManager {
	return &PairingManager{client: client}
}

func (m *PairingManager) Pair() bool {
	r := m.pairingRequest()

	if r.Status == PairingMessage_STATUS_OK {
		r = m.pairingOption()

		if r.Status == PairingMessage_STATUS_OK {
			r = m.pairingConfig()

			if r.Status == PairingMessage_STATUS_OK {
				r = m.pairingSecret()

				if r.Status == PairingMessage_STATUS_OK {
					return true
				}
			}
		}
	}

	return false
}

func (m *PairingManager) pairingRequest() *PairingMessage {
	payload := m.pairingRequestPayload()
	rawResponse := m.client.makeRequest(payload)
	response := m.processPairingResponse(rawResponse)

	return response
}

func (m *PairingManager) pairingOption() *PairingMessage {
	payload := m.pairingOptionPayload()
	rawResponse := m.client.makeRequest(payload)
	response := m.processPairingResponse(rawResponse)

	return response
}

func (m *PairingManager) pairingConfig() *PairingMessage {
	payload := m.pairingConfigPayload()
	rawResponse := m.client.makeRequest(payload)
	response := m.processPairingResponse(rawResponse)

	return response
}

func (m *PairingManager) pairingSecret() *PairingMessage {
	var input string

	fmt.Printf("TV Code: ")
	_, err := fmt.Scan(&input)
	if err != nil {
		log.Panicf("[pairing] Error reading user input: %+v\n", err)
	}

	secret := m.client.makeSecret(input)

	payload := m.pairingSecretPayload(secret)
	rawResponse := m.client.makeRequest(payload)
	response := m.processPairingResponse(rawResponse)

	return response
}

func (m *PairingManager) pairingRequestPayload() []byte {
	req := &PairingMessage{
		PairingRequest: &PairingRequest{
			ServiceName: Params.ServiceName,
			ClientName:  Params.ClientName,
		},
		ProtocolVersion: 2,
		Status:          PairingMessage_STATUS_OK,
	}

	data, err := proto.Marshal(req)
	if err != nil {
		log.Panicf("[pairing] Problem marshalling PairingRequest: %+v", err)
	}

	return data
}

func (m *PairingManager) pairingOptionPayload() []byte {
	req := &PairingMessage{
		PairingOption: &PairingOption{
			PreferredRole: RoleType_ROLE_TYPE_INPUT,
			InputEncodings: []*PairingEncoding{
				{
					Type:         PairingEncoding_ENCODING_TYPE_HEXADECIMAL,
					SymbolLength: 6,
				},
			},
		},
		ProtocolVersion: 2,
		Status:          PairingMessage_STATUS_OK,
	}

	data, err := proto.Marshal(req)
	if err != nil {
		log.Panicf("[pairing] Problem marshalling PairingOption: %+v", err)
	}

	return data
}

func (m *PairingManager) pairingConfigPayload() []byte {
	req := &PairingMessage{
		PairingConfiguration: &PairingConfiguration{
			Encoding: &PairingEncoding{
				Type:         PairingEncoding_ENCODING_TYPE_HEXADECIMAL,
				SymbolLength: 6,
			},
			ClientRole: RoleType_ROLE_TYPE_INPUT,
		},
		ProtocolVersion: 2,
		Status:          PairingMessage_STATUS_OK,
	}

	data, err := proto.Marshal(req)
	if err != nil {
		log.Panicf("[pairing] Problem marshalling PairingConfiguration: %+v", err)
	}

	return data
}

func (m *PairingManager) pairingSecretPayload(secret []byte) []byte {
	req := &PairingMessage{
		PairingSecret: &PairingSecret{
			Secret: secret,
		},
		ProtocolVersion: 2,
		Status:          PairingMessage_STATUS_OK,
	}

	data, err := proto.Marshal(req)
	if err != nil {
		log.Panicf("[pairing] Problem marshalling PairingSecret: %+v", err)
	}

	return data
}

func (m *PairingManager) processPairingResponse(response []byte) *PairingMessage {
	msg := PairingMessage{}
	proto.Unmarshal(response, &msg)

	return &msg
}
