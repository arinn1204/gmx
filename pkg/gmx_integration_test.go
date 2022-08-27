package gmx

import (
	"sync"
	"testing"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCanMakeMultipleAccountsInParralel(t *testing.T) {
	if testing.Short() {
		return
	}
	gmxClient := CreateClient()
	wg := sync.WaitGroup{}

	ids := make([]uuid.UUID, 0)

	totalConnections := 10

	for i := 0; i < totalConnections; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			id, err := gmxClient.Connect("127.0.0.1", 9001)

			assert.Nil(t, err)
			ids = append(ids, *id)
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
	gmxClient := CreateClient()
	wg := sync.WaitGroup{}

	ids := make([]uuid.UUID, 0)

	totalConnections := 25

	for i := 0; i < totalConnections; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			id, err := gmxClient.Connect("127.0.0.1", 9001)
			gmxClient.RegisterClassHandler(handlers.BoolClasspath, &handlers.BoolHandler{})

			assert.Nil(t, err)
			ids = append(ids, *id)
		}(&wg)

	}

	wg.Wait()

	assert.Equal(t, totalConnections, len(ids))

	c := gmxClient.(*client)
	assert.Equal(t, uint(totalConnections), c.numberOfConnections)
}