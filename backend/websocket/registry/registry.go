package registry

import (
	"sync"

	"github.com/bseto/arcade/backend/log"
	"github.com/bseto/arcade/backend/websocket/identifier"
)

// Registry defines an interface in which `Registry`'s should provide
type Registry interface {
	Register(send chan []byte, clientID identifier.Client)
	Unregister(clientID identifier.Client)

	CheckIfHubExists(hubName identifier.HubNameStruct) bool

	SendToSameHub(clientID identifier.Client, message []byte)
	SendToCaller(clientID identifier.Client, message []byte)
}

// RegistryProvider will provide the actual registry functionality
type RegistryProvider struct {
	lookupLock sync.RWMutex
	lookupMap  map[identifier.HubNameStruct]map[identifier.ClientUUIDStruct](chan []byte)
}

func GetRegistryProvider() *RegistryProvider {
	return &RegistryProvider{
		lookupMap: make(
			map[identifier.HubNameStruct]map[identifier.ClientUUIDStruct](chan []byte),
		),
	}
}

// Register should take the send chan and fill in the lookupMap
// This function should be threadsafe
func (r *RegistryProvider) Register(
	send chan []byte,
	clientID identifier.Client,
) {
	r.lookupLock.Lock()
	defer r.lookupLock.Unlock()

	_, ok := r.lookupMap[clientID.HubName]

	if ok != true {
		// we did not find a clientMap under this hubName
		clientMap := make(map[identifier.ClientUUIDStruct](chan []byte))
		r.lookupMap[clientID.HubName] = clientMap
	}

	r.lookupMap[clientID.HubName][clientID.ClientUUID] = send
}

func (r *RegistryProvider) Unregister(
	clientID identifier.Client,
) {
	r.lookupLock.Lock()
	defer r.lookupLock.Unlock()

	_, ok := r.lookupMap[clientID.HubName][clientID.ClientUUID]
	if !ok {
		log.Errorf("could not find client to unregister: %v", clientID)
		return
	}

	delete(r.lookupMap[clientID.HubName], clientID.ClientUUID)
	if len(r.lookupMap[clientID.HubName]) == 0 {
		delete(r.lookupMap, clientID.HubName)
	}

	return
}

func (r *RegistryProvider) SendToSameHub(
	clientID identifier.Client,
	message []byte,
) {
	r.lookupLock.Lock()
	defer r.lookupLock.Unlock()

	sendHubChannel, ok := r.lookupMap[clientID.HubName]
	if ok != true {
		log.Errorf("cannot find channel for %v", clientID)
		return
	}
	for _, clientChannel := range sendHubChannel {
		clientChannel <- message
	}
	return

}

func (r *RegistryProvider) SendToCaller(
	clientID identifier.Client,
	message []byte,
) {
	r.lookupLock.Lock()
	defer r.lookupLock.Unlock()

	sendChannel, ok := r.lookupMap[clientID.HubName][clientID.ClientUUID]
	if ok != true {
		log.Errorf("could not find channel for ID: %v", clientID)
		return
	}

	sendChannel <- message
}

// CheckIfHubExists will return whether or not the hub exists within this
// registry
func (r *RegistryProvider) CheckIfHubExists(
	hubName identifier.HubNameStruct,
) bool {
	r.lookupLock.RLock()
	defer r.lookupLock.RUnlock()

	_, ok := r.lookupMap[hubName]
	return ok
}
