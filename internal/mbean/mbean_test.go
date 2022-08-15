package mbean

import (
	"testing"

	"github.com/arinn1204/gmx/internal/handlers"
	"github.com/arinn1204/gmx/pkg/extensions"
	"github.com/stretchr/testify/assert"
)

func TestRegisterClassHandler_SingleTypedTest(t *testing.T) {
	client := &Client{
		ClassHandlers: make(map[string]extensions.IHandler),
	}
	doubleHandler := handlers.DoubleHandler{}
	client.RegisterClassHandler(handlers.DoubleJniRepresentation, &doubleHandler)

	assert.Equal(t, &doubleHandler, client.ClassHandlers[handlers.DoubleJniRepresentation])
}

func TestRegisterClassHandler_MultipleTypedTest(t *testing.T) {
	client := &Client{
		ClassHandlers: make(map[string]extensions.IHandler),
	}
	doubleHandler := handlers.DoubleHandler{}
	floatHandler := handlers.FloatHandler{}
	client.RegisterClassHandler(handlers.DoubleJniRepresentation, &doubleHandler)
	client.RegisterClassHandler(handlers.FloatJniRepresentation, &floatHandler)

	assert.Equal(t, &doubleHandler, client.ClassHandlers[handlers.DoubleJniRepresentation])
	assert.Equal(t, &floatHandler, client.ClassHandlers[handlers.FloatJniRepresentation])
}
