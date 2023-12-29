package handofcultome

import (
	"log"

	"google.golang.org/protobuf/proto"
)

type RemoteManager struct {
	client *Client
}

func CreateRemoteManager(client *Client) *RemoteManager {
	return &RemoteManager{client: client}
}

func (m *RemoteManager) Configure() bool {
	m.remoteConfigure()
	m.remoteSetActive()

	return true
}

func (m *RemoteManager) remoteConfigure() *RemoteMessage {
	payload := m.remoteConfigurePayload()
	rawResponse := m.client.makeRequest(payload)
	response := m.ProcessRemoteResponse(rawResponse)

	return response
}

func (m *RemoteManager) remoteConfigurePayload() []byte {
	req := &RemoteMessage{
		RemoteConfigure: &RemoteConfigure{
			Code1: 622,
			DeviceInfo: &RemoteDeviceInfo{
				Model:       Params.Model,
				Vendor:      Params.Vendor,
				Unknown1:    Params.Unknown1,
				Unknown2:    Params.Unknown2,
				PackageName: Params.PackageName,
				AppVersion:  Params.AppVersion,
			},
		},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		log.Panicf("[remote] Problem marshalling RemoteConfigure: %+v", err)
	}

	return data
}

func (m *RemoteManager) remoteSetActive() *RemoteMessage {
	payload := m.remoteSetActivePayload()
	rawResponse := m.client.makeRequest(payload)
	response := m.ProcessRemoteResponse(rawResponse)

	return response
}

func (m *RemoteManager) remoteSetActivePayload() []byte {
	req := &RemoteMessage{
		RemoteSetActive: &RemoteSetActive{
			Active: 622,
		},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		log.Panicf("[remote] Problem marshalling RemoteSetActive: %+v", err)
	}

	return data
}

func (m *RemoteManager) RespondPing() *RemoteMessage {
	payload := m.remotePingResponsePayload()
	rawResponse := m.client.makeRequest(payload)
	response := m.ProcessRemoteResponse(rawResponse)

	return response
}

func (m *RemoteManager) remotePingResponsePayload() []byte {
	req := &RemoteMessage{
		RemotePingResponse: &RemotePingResponse{
			Val1: 25,
		},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		log.Panicf("[remote] Problem marshalling RemotePingResponse: %+v", err)
	}

	return data
}

func (m *RemoteManager) VolumeUp() *RemoteMessage {
	payload := m.remoteAdjustVolumeLevelPayload()
	rawResponse := m.client.makeRequest(payload)
	response := m.ProcessRemoteResponse(rawResponse)
	return response
}

func (m *RemoteManager) remoteAdjustVolumeLevelPayload() []byte {
	req := &RemoteMessage{
		RemoteKeyInject: &RemoteKeyInject{
			KeyCode:   RemoteKeyCode_KEYCODE_HOME,
			Direction: RemoteDirection_SHORT,
		},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		log.Panicf("[remote] Problem marshalling RemoteKeyInject: %+v", err)
	}

	return data
}

func (m *RemoteManager) ProcessRemoteResponse(response []byte) *RemoteMessage {
	msg := RemoteMessage{}
	proto.Unmarshal(response, &msg)

	log.Printf("[*] bytes: %v\n", response)
	log.Printf("[*] struct: %+v\n", msg)

	return &msg
}
