package gmx

import (
	"errors"
	"fmt"
	"gmx/internal/jvm"
	"gmx/internal/mbean"
	"log"
	"sync"

	"github.com/google/uuid"
)

var java jvm.IJava
var lock *sync.Mutex

func init() {
	lock = &sync.Mutex{}
	java = &jvm.Java{}
}

// Initialize is the initial method to create a GMX client.
// This will initialize the JVM if necessary as well as setting up the object
func (client *Client) Initialize() error {
	startJvm()

	client.mbeans = make(map[uuid.UUID]*mbean.Client)

	return nil
}

func (client *Client) Connect(hostname string, port int) (uuid.UUID, error) {
	jmxUri := fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%d/jmxrmi", hostname, port)
	bean, err := java.CreateMBeanConnection(jmxUri)

	if err != nil {
		return uuid.UUID{}, errors.New("failed to create a connection::" + err.Error())
	}

	id := uuid.New()

	client.mbeans[id] = bean

	return id, nil
}

func (client *Client) Close() {
	for uri := range client.mbeans {
		client.mbeans[uri].Close()
		client.mbeans[uri] = nil
	}

	java.ShutdownJvm()
}

func startJvm() {
	if java.IsStarted() {
		return
	}

	lock.Lock()
	if java.IsStarted() {
		return
	}

	var err error
	java, err = java.CreateJVM()

	if err != nil {
		log.Fatalf("Failed to create the JVM::" + err.Error())
	}
}
