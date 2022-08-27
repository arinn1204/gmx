package mbean

import (
	"sync"
	"testing"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestRegisterClassHandler_SingleTypedTest(t *testing.T) {
	client := &Client{
		ClassHandlers: sync.Map{},
	}
	doubleHandler := handlers.DoubleHandler{}
	client.RegisterClassHandler(handlers.DoubleJniRepresentation, &doubleHandler)

	hnd, ok := client.ClassHandlers.Load(handlers.DoubleJniRepresentation)
	assert.True(t, ok)
	assert.Equal(t, &doubleHandler, hnd)
}

func TestRegisterClassHandler_MultipleTypedTest(t *testing.T) {
	client := &Client{
		ClassHandlers: sync.Map{},
	}
	doubleHandler := handlers.DoubleHandler{}
	floatHandler := handlers.FloatHandler{}
	client.RegisterClassHandler(handlers.DoubleJniRepresentation, &doubleHandler)
	client.RegisterClassHandler(handlers.FloatJniRepresentation, &floatHandler)

	dhnd, ok := client.ClassHandlers.Load(handlers.DoubleJniRepresentation)
	assert.True(t, ok)
	fhnd, ok := client.ClassHandlers.Load(handlers.FloatJniRepresentation)
	assert.True(t, ok)
	assert.Equal(t, &doubleHandler, dhnd)
	assert.Equal(t, &floatHandler, fhnd)
}
