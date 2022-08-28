package gmx

import (
	"sync"
	"testing"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var gmxClient MBeanClient

func TestMain(m *testing.M) {

	gmxClient = CreateClient()
	defer gmxClient.Close()

	m.Run()

}

func teardownTest() {
	c := gmxClient.(*client)

	c.mbeans.Range(func(key, value any) bool {
		c.mbeans.Delete(key)
		return true
	})

	c.numberOfConnections = 0
}

func TestCanMakeMultipleAccountsInParrallel(t *testing.T) {
	if testing.Short() {
		return
	}
	wg := sync.WaitGroup{}
	defer teardownTest()

	ids := make([]uuid.UUID, 0)
	lock := sync.Mutex{}

	totalConnections := 10

	for i := 0; i < totalConnections; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			id, err := gmxClient.RegisterConnection("127.0.0.1", 9001)

			assert.Nil(t, err)
			lock.Lock()
			ids = append(ids, *id)
			lock.Unlock()
		}(&wg)

	}

	wg.Wait()

	assert.Equal(t, totalConnections, len(ids))

	c := gmxClient.(*client)
	assert.Equal(t, uint(totalConnections), c.numberOfConnections)
}

func TestCanMakeMultipleAccountsInParralelAndRegisterHandlers(t *testing.T) {
	if testing.Short() {
		return
	}
	defer teardownTest()
	wg := sync.WaitGroup{}

	ids := make([]uuid.UUID, 0)
	lock := sync.Mutex{}

	totalConnections := 25

	for i := 0; i < totalConnections; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			id, err := gmxClient.RegisterConnection("127.0.0.1", 9001)
			gmxClient.RegisterClassHandler(handlers.BoolClasspath, &handlers.BoolHandler{})

			assert.Nil(t, err)
			lock.Lock()
			ids = append(ids, *id)
			lock.Unlock()
		}(&wg)

	}

	wg.Wait()

	assert.Equal(t, totalConnections, len(ids))

	c := gmxClient.(*client)
	assert.Equal(t, uint(totalConnections), c.numberOfConnections)
}

func TestCanExecuteGetOperation(t *testing.T) {
	if testing.Short() {
		return
	}

	id, err := gmxClient.RegisterConnection("127.0.0.1", 9001)

	assert.Nil(t, err)

	operator := gmxClient.GetOperator()
	res, errmap := operator.ExecuteAgainstAll("org.example", "game", "getLong", MBeanArgs{Value: "hello"})

	assert.Nil(t, errmap[*id])
	assert.Equal(t, "3141592", res[*id])
}

func TestCanExecutePutOperation(t *testing.T) {
	if testing.Short() {
		return
	}
	id, err := gmxClient.RegisterConnection("127.0.0.1", 9001)

	assert.Nil(t, err)

	operator := gmxClient.GetOperator()
	res, errmap := operator.ExecuteAgainstAll("org.example", "game", "putLong", MBeanArgs{Value: "hello"}, MBeanArgs{Value: "12345678"})

	assert.Nil(t, errmap[*id])
	assert.Equal(t, "", res[*id])

	res, errmap = operator.ExecuteAgainstAll("org.example", "game", "getLong", MBeanArgs{Value: "hello"})
	assert.Nil(t, errmap[*id])
	assert.Equal(t, "12345678", res[*id])

	operator.ExecuteAgainstAll("org.example", "game", "putLong", MBeanArgs{Value: "hello"}, MBeanArgs{Value: "3141592"})
}

func TestCanGetAttribute(t *testing.T) {
	if testing.Short() {
		return
	}

	id, err := gmxClient.RegisterConnection("127.0.0.1", 9001)

	assert.Nil(t, err)

	mng := gmxClient.GetAttributeManager()
	resmap, errmap := mng.Get("org.example", "game", "NestedAttribute", MBeanArgs{})

	assert.Nil(t, errmap[*id])

	assert.Equal(t, "[[1,2,3]]", resmap[*id])
}

func TestCanPutAttribute(t *testing.T) {
	if testing.Short() {
		return
	}

	id, err := gmxClient.RegisterConnection("127.0.0.1", 9001)

	assert.Nil(t, err)

	mng := gmxClient.GetAttributeManager()
	resmap, errmap := mng.Put("org.example", "game", "StringAttribute", MBeanArgs{Value: "Hello"})

	assert.Nil(t, errmap[*id])
	assert.Equal(t, "", resmap[*id])

	resmap, errmap = mng.Get("org.example", "game", "StringAttribute", MBeanArgs{})
	assert.Nil(t, errmap[*id])
	assert.Equal(t, "Hello", resmap[*id])
}
